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

// 存储资源状态结构体

type PVCStatus struct {
	Namespace    string `json:"namespace"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Capacity     string `json:"capacity"`
	AccessMode   string `json:"accessMode"`
	StorageClass string `json:"storageClass"`
	VolumeName   string `json:"volumeName"`
}

type PVStatus struct {
	Name          string `json:"name"`
	Status        string `json:"status"`
	Capacity      string `json:"capacity"`
	AccessMode    string `json:"accessMode"`
	StorageClass  string `json:"storageClass"`
	ClaimRef      string `json:"claimRef"`
	ReclaimPolicy string `json:"reclaimPolicy"`
}

type StorageClassStatus struct {
	Name              string `json:"name"`
	Provisioner       string `json:"provisioner"`
	ReclaimPolicy     string `json:"reclaimPolicy"`
	VolumeBindingMode string `json:"volumeBindingMode"`
	IsDefault         bool   `json:"isDefault"`
}

// 配置资源状态结构体

type ConfigMapStatus struct {
	Namespace string   `json:"namespace"`
	Name      string   `json:"name"`
	DataCount int      `json:"dataCount"`
	Keys      []string `json:"keys"`
}

type SecretStatus struct {
	Namespace string   `json:"namespace"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	DataCount int      `json:"dataCount"`
	Keys      []string `json:"keys"`
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

// APIError 自定义API错误结构
type APIError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error 实现error接口
func (e *APIError) Error() string {
	return e.Message
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

// NodeDetail 提供给前端的 Node 详情结构体
type NodeDetail struct {
	Name         string            `json:"name"`
	IP           string            `json:"ip"`
	Status       string            `json:"status"`
	CPUUsage     float64           `json:"cpuUsage"`
	MemoryUsage  float64           `json:"memoryUsage"`
	Role         []string          `json:"role"`
	PodsUsed     int               `json:"podsUsed"`
	PodsCapacity int               `json:"podsCapacity"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
}

// ServiceDetail 提供给前端的 Service 详情结构体
type ServiceDetail struct {
	Namespace   string            `json:"namespace"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	ClusterIP   string            `json:"clusterIP"`
	Ports       []string          `json:"ports"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Selector    map[string]string `json:"selector"`
}

// 工作负载资源详情结构体
type DeploymentDetail struct {
	Namespace   string            `json:"namespace"`
	Name        string            `json:"name"`
	Replicas    int32             `json:"replicas"`
	Available   int32             `json:"available"`
	Desired     int32             `json:"desired"`
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Selector    map[string]string `json:"selector"`
	Strategy    string            `json:"strategy"`
	Image       string            `json:"image"`
}

type StatefulSetDetail struct {
	Namespace   string            `json:"namespace"`
	Name        string            `json:"name"`
	Replicas    int32             `json:"replicas"`
	Available   int32             `json:"available"`
	Desired     int32             `json:"desired"`
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Selector    map[string]string `json:"selector"`
	ServiceName string            `json:"serviceName"`
	Image       string            `json:"image"`
}

type DaemonSetDetail struct {
	Namespace   string            `json:"namespace"`
	Name        string            `json:"name"`
	Available   int32             `json:"available"`
	Desired     int32             `json:"desired"`
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Selector    map[string]string `json:"selector"`
	Image       string            `json:"image"`
}

type JobDetail struct {
	Namespace      string            `json:"namespace"`
	Name           string            `json:"name"`
	Completions    int32             `json:"completions"`
	Succeeded      int32             `json:"succeeded"`
	Failed         int32             `json:"failed"`
	StartTime      string            `json:"startTime"`
	CompletionTime string            `json:"completionTime"`
	Status         string            `json:"status"`
	Labels         map[string]string `json:"labels"`
	Annotations    map[string]string `json:"annotations"`
	Image          string            `json:"image"`
}

type CronJobDetail struct {
	Namespace        string            `json:"namespace"`
	Name             string            `json:"name"`
	Schedule         string            `json:"schedule"`
	Suspend          bool              `json:"suspend"`
	Active           int               `json:"active"`
	LastScheduleTime string            `json:"lastScheduleTime"`
	Status           string            `json:"status"`
	Labels           map[string]string `json:"labels"`
	Annotations      map[string]string `json:"annotations"`
	Image            string            `json:"image"`
}

// 网络资源详情结构体
type IngressDetail struct {
	Namespace     string            `json:"namespace"`
	Name          string            `json:"name"`
	Hosts         []string          `json:"hosts"`
	Address       string            `json:"address"`
	Ports         []string          `json:"ports"`
	Class         string            `json:"class"`
	Status        string            `json:"status"`
	Path          []string          `json:"path"`
	TargetService []string          `json:"targetService"`
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
}

// 其他资源详情结构体
type NamespaceDetail struct {
	Name        string            `json:"name"`
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type EventDetail struct {
	Namespace   string            `json:"namespace"`
	Name        string            `json:"name"`
	Reason      string            `json:"reason"`
	Message     string            `json:"message"`
	Type        string            `json:"type"`
	Count       int32             `json:"count"`
	FirstSeen   string            `json:"firstSeen"`
	LastSeen    string            `json:"lastSeen"`
	Duration    string            `json:"duration"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

// 存储资源详情结构体
type PVCDetail struct {
	Namespace    string            `json:"namespace"`
	Name         string            `json:"name"`
	Status       string            `json:"status"`
	Capacity     string            `json:"capacity"`
	AccessMode   []string          `json:"accessMode"`
	StorageClass string            `json:"storageClass"`
	VolumeName   string            `json:"volumeName"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	Data         map[string]string `json:"data"`
}

type PVDetail struct {
	Name          string            `json:"name"`
	Status        string            `json:"status"`
	Capacity      string            `json:"capacity"`
	AccessMode    []string          `json:"accessMode"`
	StorageClass  string            `json:"storageClass"`
	ClaimRef      string            `json:"claimRef"`
	ReclaimPolicy string            `json:"reclaimPolicy"`
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
}

type StorageClassDetail struct {
	Name              string            `json:"name"`
	Provisioner       string            `json:"provisioner"`
	ReclaimPolicy     string            `json:"reclaimPolicy"`
	VolumeBindingMode string            `json:"volumeBindingMode"`
	IsDefault         bool              `json:"isDefault"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
	Parameters        map[string]string `json:"parameters"`
}

// 配置资源详情结构体
type ConfigMapDetail struct {
	Namespace   string            `json:"namespace"`
	Name        string            `json:"name"`
	DataCount   int               `json:"dataCount"`
	Keys        []string          `json:"keys"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Data        map[string]string `json:"data"`
}

type SecretDetail struct {
	Namespace   string            `json:"namespace"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	DataCount   int               `json:"dataCount"`
	Keys        []string          `json:"keys"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Data        map[string]string `json:"data"`
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

// 为各种状态结构体实现SearchableItem接口

// GetSearchableFields 实现SearchableItem接口
func (p PodStatus) GetSearchableFields() map[string]string {
	return map[string]string{
		"Name":      p.Name,
		"Namespace": p.Namespace,
		"Status":    p.Status,
		"PodIP":     p.PodIP,
		"NodeName":  p.NodeName,
	}
}

// GetSearchableFields 实现SearchableItem接口
func (d DeploymentStatus) GetSearchableFields() map[string]string {
	return map[string]string{
		"Name":      d.Name,
		"Namespace": d.Namespace,
		"Status":    d.Status,
	}
}

// GetSearchableFields 实现SearchableItem接口
func (s StatefulSetStatus) GetSearchableFields() map[string]string {
	return map[string]string{
		"Name":      s.Name,
		"Namespace": s.Namespace,
		"Status":    s.Status,
	}
}

// GetSearchableFields 实现SearchableItem接口
func (d DaemonSetStatus) GetSearchableFields() map[string]string {
	return map[string]string{
		"Name":      d.Name,
		"Namespace": d.Namespace,
		"Status":    d.Status,
	}
}

// GetSearchableFields 实现SearchableItem接口
func (s ServiceStatus) GetSearchableFields() map[string]string {
	return map[string]string{
		"Name":      s.Name,
		"Namespace": s.Namespace,
		"Type":      s.Type,
		"ClusterIP": s.ClusterIP,
	}
}

// GetSearchableFields 实现SearchableItem接口
func (n NodeStatus) GetSearchableFields() map[string]string {
	return map[string]string{
		"Name":   n.Name,
		"IP":     n.IP,
		"Status": n.Status,
	}
}
