// Package local implements local storage
package local

import (
	"github.com/OShalakhin/cloud/storage"
	"io/ioutil"
	"net/url"
	"os"
)

type (
	// Storage to hold files locally
	Storage struct {
		Provider  storage.Provider
		Container storage.Container
	}
	// Info holds storage information
	Info struct {
		URL url.URL
		Dir string
	}
)

// Create file
func (s *Storage) Create(filename string, data []byte) error {
	if err := ioutil.WriteFile(s.Info.Dir+filename, data, 0644); err != nil {
		return err
	}
	return nil
}

// Read file
func (s *Storage) Read(filename string) ([]byte, error) {
	return ioutil.ReadFile(s.Info.Dir + filename)
}

// Update file
func (s *Storage) Update(filename string, data []byte) error {
	var err error
	// Remove
	if err = s.Delete(filename); err != nil {
		return err
	}
	// Create
	if err = s.Create(filename); err != nil {
		return err
	}
	return nil
}

// Delete file
func (s *Storage) Delete(filename string) error {
	return os.Remove(s.Info.Dir + filename)
}

// Authenticate storage. Local storage doesn't need authentication as storage
// happens locally in the filesystem.
func (s *Storage) Authenticate() error {
	return nil
}

// GetContainer where to store
func (s *Storage) GetContainer() (*storage.Container, error) {
	return new(storage.Container), nil
}

// GetURL where it must be accessed
func (s *Storage) GetURL() *url.URL {
	return new(url.URL)
}
