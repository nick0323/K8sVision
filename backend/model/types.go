package model

// 资源状态结构体

type OverviewStatus struct {
	NodeCount      int           `json:"nodeCount"`
	NodeReady      int           `json:"nodeReady"`
	PodCount       int           `json:"podCount"`
	PodNotReady    int           `json:"podNotReady"`
	NamespaceCount int           `json:"namespaceCount"`
	ServiceCount   int           `json:"serviceCount"`
	CPUCapacity    float64       `json:"cpuCapacity"`
	CPURequests    float64       `json:"cpuRequests"`
	CPULimits      float64       `json:"cpuLimits"`
	MemoryCapacity float64       `json:"memoryCapacity"`
	MemoryRequests float64       `json:"memoryRequests"`
	MemoryLimits   float64       `json:"memoryLimits"`
	Events         []EventStatus `json:"events"`
}

type NodeStatus struct {
	Name         string   `json:"name"`
	IP           string   `json:"ip"`
	Status       string   `json:"status"`
	CPUUsage     float64  `json:"cpuUsage"`
	MemoryUsage  float64  `json:"memoryUsage"`
	Role         []string `json:"role"`
	PodsUsed     int      `json:"podsUsed"`
	PodsCapacity int      `json:"podsCapacity"`
}

type PodStatus struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	CPUUsage    string `json:"cpuUsage"`
	MemoryUsage string `json:"memoryUsage"`
	PodIP       string `json:"podIP"`
	NodeName    string `json:"nodeName"`
}

type DeploymentStatus struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Available int32  `json:"availableReplicas"`
	Desired   int32  `json:"desiredReplicas"`
	Status    string `json:"status"`
}

type StatefulSetStatus struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Available int32  `json:"availableReplicas"`
	Desired   int32  `json:"desiredReplicas"`
	Status    string `json:"status"`
}

type DaemonSetStatus struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Available int32  `json:"availableReplicas"`
	Desired   int32  `json:"desiredReplicas"`
	Status    string `json:"status"`
}

type ServiceStatus struct {
	Namespace string   `json:"namespace"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	ClusterIP string   `json:"clusterIP"`
	Ports     []string `json:"ports"`
}

type EventStatus struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	Count     int32  `json:"count"`
	FirstSeen string `json:"firstSeen"`
	LastSeen  string `json:"lastSeen"`
	Duration  string `json:"duration"`
}

type CronJobStatus struct {
	Namespace        string `json:"namespace"`
	Name             string `json:"name"`
	Schedule         string `json:"schedule"`
	Suspend          bool   `json:"suspend"`
	Active           int    `json:"active"`
	LastScheduleTime string `json:"lastScheduleTime"`
	Status           string `json:"status"`
}

type JobStatus struct {
	Namespace      string `json:"namespace"`
	Name           string `json:"name"`
	Completions    int32  `json:"completions"`
	Succeeded      int32  `json:"succeeded"`
	Failed         int32  `json:"failed"`
	StartTime      string `json:"startTime"`
	CompletionTime string `json:"completionTime"`
	Status         string `json:"status"`
}

type IngressStatus struct {
	Namespace     string   `json:"namespace"`
	Name          string   `json:"name"`
	Hosts         []string `json:"hosts"`
	Address       string   `json:"address"`
	Ports         []string `json:"ports"`
	Class         string   `json:"class"`
	Status        string   `json:"status"`
	Path          []string `json:"path"`
	TargetService []string `json:"targetService"`
}

// 统一 API 响应结构体

// APIResponse 统一 API 响应结构
// code: 0 表示成功，非0为错误码
// message: 响应消息
// data: 业务数据
// traceId: 链路追踪ID
// timestamp: 响应时间戳
// page: 分页元数据（可选）
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	TraceID   string      `json:"traceId,omitempty"`
	Timestamp int64       `json:"timestamp"`
	Page      *PageMeta   `json:"page,omitempty"`
}

// PageMeta 分页元数据
// total: 总条数
// limit: 每页数量
// offset: 偏移量
type PageMeta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// 统一指标类型，避免匿名结构体类型不一致导致的编译问题

// PodMetrics 统一 Pod 指标类型
// CPU: 使用量（m核）
// Mem: 使用量（字节）
type PodMetrics struct {
	CPU int64
	Mem int64
}
type PodMetricsMap map[string]PodMetrics

// NodeMetrics 统一 Node 指标类型
// CPU: 使用量字符串（如 "123m"）
// Mem: 使用量字符串（如 "512Mi"）
type NodeMetrics struct {
	CPU string
	Mem string
}
type NodeMetricsMap map[string]NodeMetrics

// PodDetail 提供给前端的 Pod 详情结构体
// 可根据实际需求补充字段
type PodDetail struct {
	Namespace   string            `json:"namespace"`
	Name        string            `json:"name"`
	Status      string            `json:"status"`
	PodIP       string            `json:"podIP"`
	NodeName    string            `json:"nodeName"`
	StartTime   string            `json:"startTime"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Containers  []string          `json:"containers"`
}

// LoginRequest 登录参数
// @Description 登录参数
// @name LoginRequest
// @Param username body string true "用户名"
// @Param password body string true "密码"
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
