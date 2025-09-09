package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatResourceUsage(t *testing.T) {
	t.Run("正常资源使用量", func(t *testing.T) {
		cpuStr, memStr := FormatResourceUsage(1000, 1024*1024*1024) // 1 CPU, 1GB

		assert.Equal(t, "1.00 mCPU", cpuStr)
		assert.Equal(t, "1.00 MiB", memStr)
	})

	t.Run("零资源使用量", func(t *testing.T) {
		cpuStr, memStr := FormatResourceUsage(0, 0)

		assert.Equal(t, "-", cpuStr)
		assert.Equal(t, "-", memStr)
	})

	t.Run("负资源使用量", func(t *testing.T) {
		cpuStr, memStr := FormatResourceUsage(-100, -1024)

		assert.Equal(t, "-", cpuStr)
		assert.Equal(t, "-", memStr)
	})
}

func TestGetResourceStatus(t *testing.T) {
	t.Run("完全就绪", func(t *testing.T) {
		status := GetResourceStatus(5, 5)
		assert.Equal(t, "Ready", status)
	})

	t.Run("部分可用", func(t *testing.T) {
		status := GetResourceStatus(3, 5)
		assert.Equal(t, "PartialAvailable", status)
	})

	t.Run("缩放到零", func(t *testing.T) {
		status := GetResourceStatus(0, 0)
		assert.Equal(t, "Scaled to zero", status)
	})

	t.Run("未就绪", func(t *testing.T) {
		status := GetResourceStatus(0, 5)
		assert.Equal(t, "Not Ready", status)
	})
}

func TestGetWorkloadStatus(t *testing.T) {
	t.Run("健康状态", func(t *testing.T) {
		status := GetWorkloadStatus(5, 5)
		assert.Equal(t, "Healthy", status)
	})

	t.Run("部分可用", func(t *testing.T) {
		status := GetWorkloadStatus(3, 5)
		assert.Equal(t, "PartialAvailable", status)
	})

	t.Run("异常状态", func(t *testing.T) {
		status := GetWorkloadStatus(0, 5)
		assert.Equal(t, "Abnormal", status)
	})
}

func TestGetJobStatus(t *testing.T) {
	t.Run("成功状态", func(t *testing.T) {
		status := GetJobStatus(5, 0, 0)
		assert.Equal(t, "Succeeded", status)
	})

	t.Run("失败状态", func(t *testing.T) {
		status := GetJobStatus(0, 3, 0)
		assert.Equal(t, "Failed", status)
	})

	t.Run("运行状态", func(t *testing.T) {
		status := GetJobStatus(0, 0, 2)
		assert.Equal(t, "Running", status)
	})

	t.Run("等待状态", func(t *testing.T) {
		status := GetJobStatus(0, 0, 0)
		assert.Equal(t, "Pending", status)
	})
}

func TestGetCronJobStatus(t *testing.T) {
	t.Run("运行状态", func(t *testing.T) {
		status := GetCronJobStatus(2, "2023-01-01T00:00:00Z")
		assert.Equal(t, "Running", status)
	})

	t.Run("成功状态", func(t *testing.T) {
		status := GetCronJobStatus(0, "2023-01-01T00:00:00Z")
		assert.Equal(t, "Succeeded", status)
	})

	t.Run("等待状态", func(t *testing.T) {
		status := GetCronJobStatus(0, nil)
		assert.Equal(t, "Pending", status)
	})
}

func TestExtractKeys(t *testing.T) {
	t.Run("正常map", func(t *testing.T) {
		data := map[string]int{
			"key1": 1,
			"key2": 2,
			"key3": 3,
		}

		keys := ExtractKeys(data)
		assert.Len(t, keys, 3)
		assert.Contains(t, keys, "key1")
		assert.Contains(t, keys, "key2")
		assert.Contains(t, keys, "key3")
	})

	t.Run("空map", func(t *testing.T) {
		data := map[string]int{}

		keys := ExtractKeys(data)
		assert.Len(t, keys, 0)
	})
}

func TestSafeInt32Ptr(t *testing.T) {
	t.Run("有效指针", func(t *testing.T) {
		val := int32(42)
		result := SafeInt32Ptr(&val, 0)
		assert.Equal(t, int32(42), result)
	})

	t.Run("空指针", func(t *testing.T) {
		result := SafeInt32Ptr(nil, 10)
		assert.Equal(t, int32(10), result)
	})
}

func TestSafeBoolPtr(t *testing.T) {
	t.Run("有效指针", func(t *testing.T) {
		val := true
		result := SafeBoolPtr(&val, false)
		assert.True(t, result)
	})

	t.Run("空指针", func(t *testing.T) {
		result := SafeBoolPtr(nil, true)
		assert.True(t, result)
	})
}
