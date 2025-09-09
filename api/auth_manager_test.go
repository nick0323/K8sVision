package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAuthManager(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	am := NewAuthManager(logger)
	defer am.Close()

	t.Run("记录登录失败", func(t *testing.T) {
		key := "testuser|127.0.0.1"
		maxFail := 3
		lockDuration := 1 * time.Minute

		// 记录失败
		am.RecordLoginFailure(key, maxFail, lockDuration)

		// 检查是否被锁定（应该还没达到最大失败次数）
		locked := am.IsLocked(key, maxFail, lockDuration)
		assert.False(t, locked)

		// 检查剩余尝试次数
		remaining := am.GetRemainingAttempts(key, maxFail)
		assert.Equal(t, 2, remaining)
	})

	t.Run("达到最大失败次数后锁定", func(t *testing.T) {
		key := "testuser2|127.0.0.1"
		maxFail := 2
		lockDuration := 1 * time.Minute

		// 记录失败直到达到最大次数
		am.RecordLoginFailure(key, maxFail, lockDuration)
		am.RecordLoginFailure(key, maxFail, lockDuration)

		// 检查是否被锁定
		locked := am.IsLocked(key, maxFail, lockDuration)
		assert.True(t, locked)

		// 检查剩余尝试次数
		remaining := am.GetRemainingAttempts(key, maxFail)
		assert.Equal(t, 0, remaining)
	})

	t.Run("清除登录失败记录", func(t *testing.T) {
		key := "testuser3|127.0.0.1"
		maxFail := 2
		lockDuration := 1 * time.Minute

		// 记录失败
		am.RecordLoginFailure(key, maxFail, lockDuration)

		// 清除记录
		am.ClearLoginFailures(key)

		// 检查是否被锁定（应该不被锁定）
		locked := am.IsLocked(key, maxFail, lockDuration)
		assert.False(t, locked)

		// 检查剩余尝试次数（应该重置）
		remaining := am.GetRemainingAttempts(key, maxFail)
		assert.Equal(t, maxFail, remaining)
	})

	t.Run("获取统计信息", func(t *testing.T) {
		key := "testuser4|127.0.0.1"
		maxFail := 2
		lockDuration := 1 * time.Minute

		// 记录一些失败
		am.RecordLoginFailure(key, maxFail, lockDuration)

		stats := am.GetStats()
		assert.Contains(t, stats, "totalRecords")
		assert.Contains(t, stats, "lockedCount")
		assert.GreaterOrEqual(t, stats["totalRecords"], 1)
	})

	t.Run("锁定时间过期后自动重置", func(t *testing.T) {
		key := "testuser5|127.0.0.1"
		maxFail := 2
		lockDuration := 100 * time.Millisecond // 很短的锁定时间

		// 记录失败直到锁定
		am.RecordLoginFailure(key, maxFail, lockDuration)
		am.RecordLoginFailure(key, maxFail, lockDuration)

		// 确认被锁定
		locked := am.IsLocked(key, maxFail, lockDuration)
		assert.True(t, locked)

		// 等待锁定时间过期
		time.Sleep(200 * time.Millisecond)

		// 检查是否自动解锁
		locked = am.IsLocked(key, maxFail, lockDuration)
		assert.False(t, locked)
	})
}
