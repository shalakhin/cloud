package cloudfiles

import (
	"errors"
	"fmt"
	"github.com/OShalakhin/cloud/storage"
	"github.com/ncw/swift/rs"
	"net/url"
	"time"
)

type (
	// Storage is an implementation of Rackspace CloudFiles cloud storage
	Storage struct {
		Provider  storage.Provider
		Container storage.Container
		Info      *Info
		Conn      rs.RsConnection
	}
	// Info holds storage info
	Info struct {
		URL *url.URL
	}
)

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
	s.Conn = rs.RsConnection{}
	s.Conn.UserName = s.Provider.Name
	s.Conn.ApiKey = s.Provider.Key
	s.Conn.AuthUrl = s.GetAuthURL()
	s.Conn.Region = s.Provider.AuthURL

	if err := s.Conn.Authenticate(); err != nil {
		panic(err)
	}
	// Set larger data timeout to reduce number of failed transfer
	s.Conn.Timeout = time.Duration(90) * time.Second
	// Set Info.URL empty as it must be initialized by GetURL
	s.Info.URL = new(url.URL)
	return nil
}

// Create file
func (s *Storage) Create(filename string, data []byte) error {
	if !s.Conn.Authenticated() {
		return fmt.Errorf("not authenticated")
	}
	if err := s.Conn.ObjectPutBytes(s.Container.Name, filename, data, ""); err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

// Read file
func (s *Storage) Read(filename string) ([]byte, error) {
	data, err := s.Conn.ObjectGetBytes(s.Container.Name, filename)
	if err != nil {
		return []byte(""), err
	}
	return data, nil
}

// Update file
func (s *Storage) Update(filename string, data []byte) error {
	if !s.Conn.Authenticated() {
		return fmt.Errorf("not authenticated")
	}
	// delete
	if err := s.Delete(filename); err != nil {
		return err
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
	if !s.Conn.Authenticated() {
		return fmt.Errorf("not authenticated")
	}
	if err := s.Conn.ObjectDelete(s.Container.Name, filename); err != nil {
		return err
	}
	return nil
}

// GetURL returns url to the storage to use i.e. in templates
func (s *Storage) GetURL() *url.URL {
	if !s.Conn.Authenticated() {
		err := s.Conn.Authenticate()
		if err != nil {
			panic(err)
		}
	}

	u := s.Info.URL
	// generate it if not exist
	if u.RequestURI() == "/" {
		h, err := s.Conn.ContainerCDNMeta(s.Container.Name)
		if err != nil {
			panic(err)
		}
		u, err = u.Parse(h["X-Cdn-Ssl-Uri"])
		if err != nil {
			panic(err)
		}
		s.Info.URL = u
	}
	return s.Info.URL
}
