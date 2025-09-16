package service

import (
	"testing"

	"github.com/nick0323/K8sVision/model"
)

// TestGetResourceStatus 测试资源状态判断
func TestGetResourceStatus(t *testing.T) {
	tests := []struct {
		name     string
		ready    int32
		desired  int32
		expected string
	}{
		{
			name:     "Ready状态",
			ready:    3,
			desired:  3,
			expected: model.StatusReady,
		},
		{
			name:     "部分可用状态",
			ready:    2,
			desired:  3,
			expected: model.StatusPartial,
		},
		{
			name:     "缩放到零状态",
			ready:    0,
			desired:  0,
			expected: model.StatusScaledToZero,
		},
		{
			name:     "未就绪状态",
			ready:    0,
			desired:  3,
			expected: model.StatusNotReady,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetResourceStatus(tt.ready, tt.desired)
			if result != tt.expected {
				t.Errorf("GetResourceStatus(%d, %d) = %v, want %v", tt.ready, tt.desired, result, tt.expected)
			}
		})
	}
}

// TestGetWorkloadStatus 测试工作负载状态判断
func TestGetWorkloadStatus(t *testing.T) {
	tests := []struct {
		name     string
		ready    int32
		desired  int32
		expected string
	}{
		{
			name:     "健康状态",
			ready:    3,
			desired:  3,
			expected: model.StatusHealthy,
		},
		{
			name:     "部分可用状态",
			ready:    2,
			desired:  3,
			expected: model.StatusPartial,
		},
		{
			name:     "异常状态",
			ready:    0,
			desired:  3,
			expected: model.StatusAbnormal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetWorkloadStatus(tt.ready, tt.desired)
			if result != tt.expected {
				t.Errorf("GetWorkloadStatus(%d, %d) = %v, want %v", tt.ready, tt.desired, result, tt.expected)
			}
		})
	}
}

// TestGetJobStatus 测试Job状态判断
func TestGetJobStatus(t *testing.T) {
	tests := []struct {
		name      string
		succeeded int32
		failed    int32
		active    int32
		expected  string
	}{
		{
			name:      "成功状态",
			succeeded: 1,
			failed:    0,
			active:    0,
			expected:  model.StatusSucceeded,
		},
		{
			name:      "失败状态",
			succeeded: 0,
			failed:    1,
			active:    0,
			expected:  model.StatusFailed,
		},
		{
			name:      "运行状态",
			succeeded: 0,
			failed:    0,
			active:    1,
			expected:  model.StatusRunning,
		},
		{
			name:      "等待状态",
			succeeded: 0,
			failed:    0,
			active:    0,
			expected:  model.StatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetJobStatus(tt.succeeded, tt.failed, tt.active)
			if result != tt.expected {
				t.Errorf("GetJobStatus(%d, %d, %d) = %v, want %v", tt.succeeded, tt.failed, tt.active, result, tt.expected)
			}
		})
	}
}

// TestSafeInt32Ptr 测试安全获取int32指针值
func TestSafeInt32Ptr(t *testing.T) {
	value := int32(42)
	tests := []struct {
		name         string
		ptr          *int32
		defaultValue int32
		expected     int32
	}{
		{
			name:         "非空指针",
			ptr:          &value,
			defaultValue: 0,
			expected:     42,
		},
		{
			name:         "空指针",
			ptr:          nil,
			defaultValue: 10,
			expected:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeInt32Ptr(tt.ptr, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("SafeInt32Ptr() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSafeBoolPtr 测试安全获取bool指针值
func TestSafeBoolPtr(t *testing.T) {
	value := true
	tests := []struct {
		name         string
		ptr          *bool
		defaultValue bool
		expected     bool
	}{
		{
			name:         "非空指针",
			ptr:          &value,
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "空指针",
			ptr:          nil,
			defaultValue: false,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeBoolPtr(tt.ptr, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("SafeBoolPtr() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsEmptyString 测试字符串是否为空
func TestIsEmptyString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "空字符串",
			input:    "",
			expected: true,
		},
		{
			name:     "非空字符串",
			input:    "hello",
			expected: false,
		},
		{
			name:     "空格字符串",
			input:    " ",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmptyString(tt.input)
			if result != tt.expected {
				t.Errorf("IsEmptyString(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestTruncateString 测试字符串截断
func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "短字符串",
			input:    "hello",
			maxLen:   10,
			expected: "hello",
		},
		{
			name:     "长字符串",
			input:    "hello world",
			maxLen:   5,
			expected: "hello...",
		},
		{
			name:     "空字符串",
			input:    "",
			maxLen:   5,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("TruncateString(%q, %d) = %v, want %v", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

// TestContainsString 测试字符串切片是否包含指定字符串
func TestContainsString(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}
	tests := []struct {
		name     string
		item     string
		expected bool
	}{
		{
			name:     "包含的字符串",
			item:     "banana",
			expected: true,
		},
		{
			name:     "不包含的字符串",
			item:     "orange",
			expected: false,
		},
		{
			name:     "空字符串",
			item:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsString(slice, tt.item)
			if result != tt.expected {
				t.Errorf("ContainsString(%v, %q) = %v, want %v", slice, tt.item, result, tt.expected)
			}
		})
	}
}

// TestRemoveEmptyStrings 测试移除空字符串
func TestRemoveEmptyStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "包含空字符串",
			input:    []string{"apple", "", "banana", "", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "无空字符串",
			input:    []string{"apple", "banana", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "全为空字符串",
			input:    []string{"", "", ""},
			expected: []string{},
		},
		{
			name:     "空切片",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveEmptyStrings(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("RemoveEmptyStrings(%v) length = %d, want %d", tt.input, len(result), len(tt.expected))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("RemoveEmptyStrings(%v)[%d] = %q, want %q", tt.input, i, v, tt.expected[i])
				}
			}
		})
	}
}

// TestMergeStringMaps 测试合并字符串map
func TestMergeStringMaps(t *testing.T) {
	m1 := map[string]string{"a": "1", "b": "2"}
	m2 := map[string]string{"b": "3", "c": "4"}
	expected := map[string]string{"a": "1", "b": "3", "c": "4"}

	result := MergeStringMaps(m1, m2)

	if len(result) != len(expected) {
		t.Errorf("MergeStringMaps() length = %d, want %d", len(result), len(expected))
	}

	for k, v := range expected {
		if result[k] != v {
			t.Errorf("MergeStringMaps()[%q] = %q, want %q", k, result[k], v)
		}
	}
}
