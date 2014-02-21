package cloudfiles

import (
	"github.com/OShalakhin/cloud/storage"
	"github.com/ncw/swift"
)

// Storage is an implementation of Rackspace CloudFiles cloud storage
type Storage struct {
	storage.Provider
	Connection swift.Connection
}

// GetAuthURL gets auth URL based on storage provider AuthURL (ORD, DFW, HKG, LON, IAD, SYD)
func (s *Storage) GetAuthURL() string {
	switch {
	case s.Provider.AuthURL == "ORD" || s.Provider.AuthURL == "DFW" || s.Provider.AuthURL == "HKG" || s.Provider.AuthURL == "IAD" || s.Provider.AuthURL == "SYD":
		return authUS
	case s.Provider.AuthURL == "LON":
		return authLON
	default:
		return authUS
	}
}

// Authenticate CloudFiles storage
func (s *Storage) Authenticate() error {
	s.Connection = swift.Connection{
		UserName: s.Provider.Name,
		ApiKey:   s.Provider.Key,
		AuthUrl:  s.GetAuthURL(),
	}
	if err := s.Connection.Authenticate(); err != nil {
		panic(err)
	}
	return nil
}

// Create file
func (s *Storage) Create(filename string, data []byte) error {
	return nil
}

// Read file
func (s *Storage) Read(filename string) ([]byte, error) {
	return []byte("OK"), nil
}

// Update file
func (s *Storage) Update(filename string, data []byte) error {
	return nil
}

// Upsert file
func (s *Storage) Upsert(filename string, data []byte) error {
	return nil
}

// Delete file
func (s *Storage) Delete(filename string) error {
	return nil
}
