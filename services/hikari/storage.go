package hikari

import (
	"sync"
)

// RecordStorage is the data structure for storing DNS records
type RecordStorage struct {
	sync.Mutex
	// Holds the DNS A records in the form of Key : Value pairs
	// with Domain Name as the key and the IPv4 Address as the value
	Holder map[string]string
}

// Get retrieves an A record for a given domain name
func (rs *RecordStorage) Get(name string) (string, bool) {
	answer, success := rs.Holder[name]
	return answer, success
}

// SetBulk adds/updates multiple A records to the storage
func (rs *RecordStorage) SetBulk(data map[string]string) {
	rs.Lock()
	defer rs.Unlock()
	for domainName, IP := range data {
		rs.Holder[domainName] = IP
	}
}

// Set adds/updates a single A record to the storage
func (rs *RecordStorage) Set(domainName, IP string) {
	rs.SetBulk(map[string]string{
		domainName: IP,
	})
}

// NewRecordStorage returns a new instance of RecordStorage data structure
func NewRecordStorage() *RecordStorage {
	return &RecordStorage{
		Holder: make(map[string]string),
	}
}

// storage is the RecordStorage instance being used in this package
var storage = NewRecordStorage()
