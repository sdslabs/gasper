package types

// MemoryStats defines a struct for storing a container's memory statistics
type MemoryStats struct {
	Usage    float64 `json:"usage"`
	MaxUsage float64 `json:"max_usage"`
	Limit    float64 `json:"limit"`
}

// CPUUsageStats defines a struct for storing a container's CPU usage statistics
type CPUUsageStats struct {
	TotalUsage float64 `json:"total_usage"`
}

// CPUStats defines a struct for storing a container's CPU statistics
type CPUStats struct {
	OnlineCPUs float64       `json:"online_cpus"`
	CPUUsage   CPUUsageStats `json:"cpu_usage"`
}

// Stats defines a struct for storing container statistics
type Stats struct {
	Memory MemoryStats `json:"memory_stats"`
	CPU    CPUStats    `json:"cpu_stats"`
}

// Metrics defines a struct for storing container metrics
type Metrics struct {
	Name           string  `json:"name" bson:"name"`
	CPUUsage       float64 `json:"cpu_usage" bson:"cpu_usage"`
	OnlineCPUs     float64 `json:"online_cpus" bson:"online_cpus"`
	MemoryUsage    float64 `json:"memory_usage" bson:"memory_usage"`
	MaxMemoryUsage float64 `json:"max_memory_usage" bson:"max_memory_usage"`
	MemoryLimit    float64 `json:"memory_limit" bson:"memory_limit"`
	ReadTime       int64   `json:"timestamp" bson:"timestamp"`
	Alive          bool    `json:"alive" bson:"alive"`
	HostIP         string  `json:"host_ip" bson:"host_ip"`
	Logs           string  `json:"logs" bson:"logs"`
}
