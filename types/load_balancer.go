package types

import (
	"sync"
)

// LoadBalancer is a data structure for load balancing requests among
// multiple instances using round-robin scheduling algorithm
type LoadBalancer struct {
	sync.Mutex
	// Instances stores the reverse-proxy containers among which network load is balanced
	Instances []*ProxyInfo
	// Counter stores the index of the server instance for directing the next request to
	Counter int
}

// Get returns an instance from the LoadBalancer
func (lb *LoadBalancer) Get() (*ProxyInfo, bool) {
	instances := lb.Instances
	numInstances := len(instances)
	if numInstances == 0 {
		return nil, false
	}
	instance := instances[lb.Counter%numInstances]
	lb.Counter = (lb.Counter + 1) % numInstances
	return instance, true
}

// Update updates the LoadBalancer instances
func (lb *LoadBalancer) Update(newInstances []string) {
	lb.Lock()
	defer lb.Unlock()
	newProxyInstances := make([]*ProxyInfo, 0)
	for _, instance := range newInstances {
		newProxyInstances = append(newProxyInstances, NewProxyInfo(instance))
	}
	lb.Instances = newProxyInstances
}

// NewLoadBalancer returns a new LoadBalancer instance
func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		Instances: make([]*ProxyInfo, 0),
		Counter:   0,
	}
}
