package types

import "sync"

// ProxyStorage maps the application name to its appropriate reverse-proxy container
type ProxyStorage struct {
	sync.Mutex
	Holder map[string]*ProxyInfo
}

// Get returns an application's reverse-proxy container along with a success message
func (ps *ProxyStorage) Get(key string) (*ProxyInfo, bool) {
	value, success := ps.Holder[key]
	return value, success
}

// Update updates the application information in the ProxyStorage container
func (ps *ProxyStorage) Update(body map[string]string) {
	ps.Lock()
	defer ps.Unlock()
	for name, host := range body {
		if ps.Holder[name] == nil {
			ps.Holder[name] = NewProxyInfo(host)
			continue
		}
		if ps.Holder[name].host == host {
			continue
		}
		ps.Holder[name].UpdateDirector(host)
	}
}

// NewProxyStorage returns a new ProxyStorage container
func NewProxyStorage() *ProxyStorage {
	return &ProxyStorage{
		Holder: make(map[string]*ProxyInfo),
	}
}
