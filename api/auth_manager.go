package api

import (
	"fmt"
	"sync"
	"time"

	"github.com/nick0323/K8sVision/config"
	"go.uber.org/zap"
)

type LoginAttempt struct {
	FailCount int       `json:"failCount"`
	LockUntil time.Time `json:"lockUntil"`
	LastFail  time.Time `json:"lastFail"`
}

type AuthManager struct {
	attempts map[string]*LoginAttempt
	mutex    sync.RWMutex
	logger   *zap.Logger
	config   *config.Manager
	stopCh   chan struct{}
}

func NewAuthManager(logger *zap.Logger, configMgr *config.Manager) *AuthManager {
	am := &AuthManager{
		attempts: make(map[string]*LoginAttempt),
		logger:   logger,
		config:   configMgr,
		stopCh:   make(chan struct{}),
	}

	go am.startCleanup()
	return am
}

func (am *AuthManager) IsLocked(username, ip string) bool {
	key := am.makeKey(username, ip)

	am.mutex.RLock()
	defer am.mutex.RUnlock()

	attempt, exists := am.attempts[key]
	if !exists {
		return false
	}

	if time.Now().After(attempt.LockUntil) {
		go am.clearAttempt(key)
		return false
	}

	authConfig := am.config.GetAuthConfig()
	return attempt.FailCount >= authConfig.MaxLoginFail
}

func (am *AuthManager) RecordFailure(username, ip string) {
	key := am.makeKey(username, ip)
	authConfig := am.config.GetAuthConfig()

	am.mutex.Lock()
	defer am.mutex.Unlock()

	attempt, exists := am.attempts[key]
	if !exists {
		attempt = &LoginAttempt{}
		am.attempts[key] = attempt
	}

	attempt.FailCount++
	attempt.LastFail = time.Now()

	if attempt.FailCount >= authConfig.MaxLoginFail {
		attempt.LockUntil = time.Now().Add(authConfig.LockDuration)
		am.logger.Warn("用户已被锁定",
			zap.String("username", username),
			zap.String("ip", ip),
			zap.Int("failCount", attempt.FailCount),
			zap.Time("lockUntil", attempt.LockUntil),
		)
	}
}

func (am *AuthManager) RecordSuccess(username, ip string) {
	key := am.makeKey(username, ip)
	am.clearAttempt(key)
}

func (am *AuthManager) GetRemainingAttempts(username, ip string) int {
	key := am.makeKey(username, ip)
	authConfig := am.config.GetAuthConfig()

	am.mutex.RLock()
	defer am.mutex.RUnlock()

	attempt, exists := am.attempts[key]
	if !exists {
		return authConfig.MaxLoginFail
	}

	remaining := authConfig.MaxLoginFail - attempt.FailCount
	if remaining < 0 {
		remaining = 0
	}
	return remaining
}

func (am *AuthManager) GetLockTime(username, ip string) time.Duration {
	key := am.makeKey(username, ip)

	am.mutex.RLock()
	defer am.mutex.RUnlock()

	attempt, exists := am.attempts[key]
	if !exists {
		return 0
	}

	if time.Now().After(attempt.LockUntil) {
		return 0
	}

	return time.Until(attempt.LockUntil)
}

func (am *AuthManager) GetStats() map[string]interface{} {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	now := time.Now()
	lockedCount := 0

	for _, attempt := range am.attempts {
		if now.Before(attempt.LockUntil) {
			authConfig := am.config.GetAuthConfig()
			if attempt.FailCount >= authConfig.MaxLoginFail {
				lockedCount++
			}
		}
	}

	authConfig := am.config.GetAuthConfig()
	return map[string]interface{}{
		"totalAttempts": len(am.attempts),
		"lockedUsers":   lockedCount,
		"maxFailCount":  authConfig.MaxLoginFail,
		"lockDuration":  authConfig.LockDuration.String(),
	}
}

func (am *AuthManager) Close() {
	close(am.stopCh)
	am.logger.Info("认证管理器已关闭")
}

func (am *AuthManager) makeKey(username, ip string) string {
	return fmt.Sprintf("%s|%s", username, ip)
}

func (am *AuthManager) clearAttempt(key string) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	delete(am.attempts, key)
}

func (am *AuthManager) startCleanup() {
	ticker := time.NewTicker(10 * time.Minute)
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

func (am *AuthManager) cleanup() {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	now := time.Now()
	cleaned := 0

	for key, attempt := range am.attempts {
		if now.After(attempt.LockUntil) && now.Sub(attempt.LastFail) > time.Hour {
			delete(am.attempts, key)
			cleaned++
		}
	}

	if cleaned > 0 {
		am.logger.Debug("清理过期登录记录", zap.Int("count", cleaned))
	}
}
