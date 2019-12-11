package types

// Stats defines a struct for storing container stats information
type Stats struct {
	Name           string `json:"name"`
	CPUUsage       string `json:"cpu_usage"`
	OnlineCPUS     string `json:"online_cpus"`
	MemoryUsage    string `json:"memory_usage"`
	MaxMemoryUsage string `json:"max_memory_usage"`
	ReadTime       int64  `json:"read_time"`
	Alive          bool   `json:"alive"`
}
