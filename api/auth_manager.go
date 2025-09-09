package api

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// LoginAttempt 登录尝试记录
type LoginAttempt struct {
	Count     int
	LastFail  time.Time
	LockUntil time.Time
}

// AuthManager 认证管理器
type AuthManager struct {
	loginFailMap map[string]*LoginAttempt
	mutex        sync.RWMutex
	logger       *zap.Logger
	stopCh       chan struct{}

	// 新增：安全增强字段
	blacklist    map[string]time.Time // IP黑名单
	whitelist    map[string]bool      // IP白名单
	maxFailCount int                  // 最大失败次数
	lockDuration time.Duration        // 锁定时间
}

// NewAuthManager 创建认证管理器
func NewAuthManager(logger *zap.Logger) *AuthManager {
	am := &AuthManager{
		loginFailMap: make(map[string]*LoginAttempt),
		blacklist:    make(map[string]time.Time),
		whitelist:    make(map[string]bool),
		logger:       logger,
		stopCh:       make(chan struct{}),
		maxFailCount: 5,
		lockDuration: 10 * time.Minute,
	}

	// 启动清理协程
	go am.cleanupWorker()

	return am
}

// RecordLoginFailure 记录登录失败
func (am *AuthManager) RecordLoginFailure(key string, maxFail int, lockDuration time.Duration) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	attempt, exists := am.loginFailMap[key]
	if !exists {
		attempt = &LoginAttempt{}
		am.loginFailMap[key] = attempt
	}

	attempt.Count++
	attempt.LastFail = time.Now()

	// 如果达到最大失败次数，设置锁定时间
	if attempt.Count >= maxFail {
		attempt.LockUntil = time.Now().Add(lockDuration)
		am.logger.Warn("用户登录被锁定",
			zap.String("key", key),
			zap.Int("failCount", attempt.Count),
			zap.Time("lockUntil", attempt.LockUntil),
		)
	}
}

// IsLocked 检查是否被锁定
func (am *AuthManager) IsLocked(key string, maxFail int, lockDuration time.Duration) bool {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	attempt, exists := am.loginFailMap[key]
	if !exists {
		return false
	}

	// 检查是否超过锁定时间
	if time.Now().After(attempt.LockUntil) {
		// 锁定时间已过，重置计数
		am.mutex.RUnlock()
		am.mutex.Lock()
		attempt.Count = 0
		attempt.LockUntil = time.Time{}
		am.mutex.Unlock()
		am.mutex.RLock()
		return false
	}

	return attempt.Count >= maxFail
}

// GetRemainingAttempts 获取剩余尝试次数
func (am *AuthManager) GetRemainingAttempts(key string, maxFail int) int {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	attempt, exists := am.loginFailMap[key]
	if !exists {
		return maxFail
	}

	remaining := maxFail - attempt.Count
	if remaining < 0 {
		remaining = 0
	}

	return remaining
}

// ClearLoginFailures 清除登录失败记录
func (am *AuthManager) ClearLoginFailures(key string) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	delete(am.loginFailMap, key)
	am.logger.Debug("清除登录失败记录", zap.String("key", key))
}

// cleanupWorker 清理工作协程
func (am *AuthManager) cleanupWorker() {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟清理一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			am.cleanup()
		case <-am.stopCh:
			return
		}
	}
}

// cleanup 清理过期的登录失败记录
func (am *AuthManager) cleanup() {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	now := time.Now()
	cleanedCount := 0

	for key, attempt := range am.loginFailMap {
		// 清理超过24小时的记录
		if now.Sub(attempt.LastFail) > 24*time.Hour {
			delete(am.loginFailMap, key)
			cleanedCount++
		}
	}

	if cleanedCount > 0 {
		am.logger.Info("清理过期登录失败记录", zap.Int("count", cleanedCount))
	}
}

// GetStats 获取统计信息
func (am *AuthManager) GetStats() map[string]interface{} {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	totalRecords := len(am.loginFailMap)
	lockedCount := 0

	for _, attempt := range am.loginFailMap {
		if time.Now().Before(attempt.LockUntil) {
			lockedCount++
		}
	}

	return map[string]interface{}{
		"totalRecords": totalRecords,
		"lockedCount":  lockedCount,
	}
}

// Close 关闭认证管理器
func (am *AuthManager) Close() {
	close(am.stopCh)
	am.logger.Info("认证管理器已关闭")
}

// IsIPBlacklisted 检查IP是否在黑名单中
func (am *AuthManager) IsIPBlacklisted(ip string) bool {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	if blacklistTime, exists := am.blacklist[ip]; exists {
		if time.Now().Before(blacklistTime) {
			return true
		}
		// 黑名单已过期，删除
		delete(am.blacklist, ip)
	}
	return false
}

// IsIPWhitelisted 检查IP是否在白名单中
func (am *AuthManager) IsIPWhitelisted(ip string) bool {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	return am.whitelist[ip]
}

// AddToBlacklist 添加IP到黑名单
func (am *AuthManager) AddToBlacklist(ip string, duration time.Duration) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	am.blacklist[ip] = time.Now().Add(duration)
	am.logger.Warn("IP已添加到黑名单",
		zap.String("ip", ip),
		zap.Duration("duration", duration))
}

// AddToWhitelist 添加IP到白名单
func (am *AuthManager) AddToWhitelist(ip string) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	am.whitelist[ip] = true
	am.logger.Info("IP已添加到白名单", zap.String("ip", ip))
}

// RemoveFromBlacklist 从黑名单移除IP
func (am *AuthManager) RemoveFromBlacklist(ip string) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	delete(am.blacklist, ip)
	am.logger.Info("IP已从黑名单移除", zap.String("ip", ip))
}

// RemoveFromWhitelist 从白名单移除IP
func (am *AuthManager) RemoveFromWhitelist(ip string) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	delete(am.whitelist, ip)
	am.logger.Info("IP已从白名单移除", zap.String("ip", ip))
}

// GetSecurityStats 获取安全统计信息
func (am *AuthManager) GetSecurityStats() map[string]interface{} {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	// 统计当前锁定的用户
	lockedCount := 0
	for _, attempt := range am.loginFailMap {
		if time.Now().Before(attempt.LockUntil) {
			lockedCount++
		}
	}

	// 统计黑名单中的IP
	blacklistCount := 0
	for _, blacklistTime := range am.blacklist {
		if time.Now().Before(blacklistTime) {
			blacklistCount++
		}
	}

	return map[string]interface{}{
		"totalLoginAttempts": len(am.loginFailMap),
		"lockedUsers":        lockedCount,
		"blacklistedIPs":     blacklistCount,
		"whitelistedIPs":     len(am.whitelist),
		"maxFailCount":       am.maxFailCount,
		"lockDuration":       am.lockDuration.String(),
	}
}

// ValidateLoginAttempt 验证登录尝试（增强版）
func (am *AuthManager) ValidateLoginAttempt(ip, username string) error {
	// 检查IP黑名单
	if am.IsIPBlacklisted(ip) {
		return fmt.Errorf("IP %s 已被封禁", ip)
	}

	// 白名单IP直接通过
	if am.IsIPWhitelisted(ip) {
		return nil
	}

	// 检查用户锁定状态
	key := username + "|" + ip
	if am.IsLocked(key, am.maxFailCount, am.lockDuration) {
		remainingAttempts := am.GetRemainingAttempts(key, am.maxFailCount)
		return fmt.Errorf("用户 %s 已被锁定，剩余尝试次数: %d", username, remainingAttempts)
	}

	return nil
}
