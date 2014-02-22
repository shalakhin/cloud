package cloudfiles

import (
	"errors"
	"fmt"
	"github.com/OShalakhin/cloud/storage"
	"github.com/ncw/swift"
	"time"
)

// Storage is an implementation of Rackspace CloudFiles cloud storage
type Storage struct {
	Provider   storage.Provider
	Container  storage.Container
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

// GetContainer from cloudfiles.Storage
func (s *Storage) GetContainer() (*storage.Container, error) {
	if s.Container.Name == "" || s.Container.Provider == "" {
		return nil, errors.New("empty container")
	}
	return &s.Container, nil
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
	// Set larger data timeout to reduce number of failed transfer
	s.Connection.Timeout = time.Duration(90) * time.Second
	return nil
}

// Create file
func (s *Storage) Create(filename string, data []byte) error {
	// s.Authenticate()
	if err := s.Connection.ObjectPutBytes(s.Container.Name, filename, data, ""); err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

// Read file
func (s *Storage) Read(filename string) ([]byte, error) {
	data, err := s.Connection.ObjectGetBytes(s.Container.Name, filename)
	if err != nil {
		return []byte(""), err
	}
	return data, nil
}

// Update file
func (s *Storage) Update(filename string, data []byte) error {
	// delete
	if err := s.Delete(filename); err != nil {
		panic(err)
	}
	// create new
	return s.Create(filename, data)
}

// Upsert file
// func (s *Storage) Upsert(filename string, data []byte) error {
// 	return nil
// }

// Delete file
func (s *Storage) Delete(filename string) error {
	if err := s.Connection.ObjectDelete(s.Container.Name, filename); err != nil {
		return err
	}
	return nil
}
