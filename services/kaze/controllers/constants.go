package controllers

const (
	// WorkerNode is the reference to the worker nodes
	WorkerNode = "workers"
	// MasterNode is the reference to the master nodes
	MasterNode = "master"
)

// timeConversionMap holds various units of time and their conversion
// factor to seconds
var timeConversionMap = map[string]int64{
	"seconds": 1,
	"minutes": 60,
	"hours":   3600,
	"days":    24 * 3600,
	"weeks":   7 * 24 * 3600,
	"months":  30 * 24 * 3600,
	"years":   365 * 24 * 3600,
	"decades": 10 * 365 * 24 * 3600,
}
