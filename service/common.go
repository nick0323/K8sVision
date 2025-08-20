package service

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
)

// GenericResourceLister 通用资源列表获取函数
func GenericResourceLister[T any](
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
	listFunc func(string) ([]T, error),
) ([]T, error) {
	return listFunc(namespace)
}

// FormatResourceUsage 格式化资源使用量
func FormatResourceUsage(cpu, mem int64) (cpuStr, memStr string) {
	if cpu > 0 {
		cpuStr = fmt.Sprintf("%.2f mCPU", float64(cpu))
	} else {
		cpuStr = "-"
	}
	if mem > 0 {
		memStr = fmt.Sprintf("%.2f MiB", float64(mem)/(1024*1024))
	} else {
		memStr = "-"
	}
	return cpuStr, memStr
}

// GetResourceStatus 通用状态判断函数
func GetResourceStatus(ready, desired int32) string {
	if ready == desired && desired > 0 {
		return "Ready"
	} else if ready > 0 {
		return "PartialAvailable"
	} else if desired == 0 {
		return "Scaled to zero"
	} else {
		return "Not Ready"
	}
}

// GetWorkloadStatus 工作负载状态判断函数
func GetWorkloadStatus(ready, desired int32) string {
	if ready == desired && desired > 0 {
		return "Healthy"
	} else if ready > 0 {
		return "PartialAvailable"
	} else {
		return "Abnormal"
	}
}

// GetJobStatus Job状态判断函数
func GetJobStatus(succeeded, failed, active int32) string {
	if succeeded > 0 {
		return "Succeeded"
	} else if failed > 0 {
		return "Failed"
	} else if active > 0 {
		return "Running"
	} else {
		return "Pending"
	}
}

// GetCronJobStatus CronJob状态判断函数
func GetCronJobStatus(activeCount int, lastSuccessfulTime interface{}) string {
	if activeCount > 0 {
		return "Running"
	} else if lastSuccessfulTime != nil {
		return "Succeeded"
	} else {
		return "Pending"
	}
}

// ExtractKeys 从map中提取keys
func ExtractKeys[T any](data map[string]T) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	return keys
}

// SafeInt32Ptr 安全获取int32指针值
func SafeInt32Ptr(ptr *int32, defaultValue int32) int32 {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

// SafeBoolPtr 安全获取bool指针值
func SafeBoolPtr(ptr *bool, defaultValue bool) bool {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}
