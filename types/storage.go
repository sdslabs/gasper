package types

import "sync"

// RecordStorage is the data structure for storing records
type RecordStorage struct {
	sync.Mutex
	// Records are stored in the form of Key : Value pairs
	Holder map[string]string
}

// Get retrieves a record from the storage
func (rs *RecordStorage) Get(key string) (string, bool) {
	value, success := rs.Holder[key]
	return value, success
}

// Set adds/updates a single record to the storage
func (rs *RecordStorage) Set(key, value string) {
	rs.Lock()
	defer rs.Unlock()
	rs.Holder[key] = value
}

// SetBulk adds/updates multiple records to the storage
func (rs *RecordStorage) SetBulk(data map[string]string) {
	rs.Lock()
	defer rs.Unlock()
	for key, value := range data {
		rs.Holder[key] = value
	}
}

// Replace replaces the records in the storage with new records
func (rs *RecordStorage) Replace(replacement map[string]string) {
	rs.Lock()
	defer rs.Unlock()
	rs.Holder = replacement
}

// NewRecordStorage returns a new instance of RecordStorage data structure
func NewRecordStorage() *RecordStorage {
	return &RecordStorage{
		Holder: make(map[string]string),
	}
}
