package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OShalakhin/cloud/storage"
	"github.com/OShalakhin/cloud/storage/cloudfiles"
	"github.com/ncw/swift"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

const (
	// CLOUDCORE cloud core config name
	CLOUDCORE = ".cloudcore"
	// CLOUD cloud container config name
	CLOUD = ".cloud"
	// CLOUDIGNORE cloud ignore config name
	CLOUDIGNORE = ".cloudignore"
)

// GetCore returns parsed .cloudcore struct
func GetCore() (storage.Core, error) {
	core := storage.Core{}
	// open file
	u, err := user.Current()
	if err != nil {
		return core, err
	}
	corepath := path.Join(u.HomeDir, CLOUDCORE)
	data, err := ioutil.ReadFile(corepath)
	if err != nil {
		return core, err
	}
	// parse it
	if err = json.Unmarshal(data, &core); err != nil {
		return core, err
	}
	return core, nil
}

// GetProvider from container string from the core config
func GetProvider(container *storage.Container, core *storage.Core) (storage.Provider, error) {
	for _, v := range core.Providers {
		if v.Provider == container.Provider {
			fmt.Println("Provider:", v.Provider)
			return v, nil
		}
	}
	return storage.Provider{}, errors.New("no provider with such name found")
}

// GetCloud returns parsed .cloud struct
func GetCloud() (storage.Cloud, error) {
	cloud := storage.Cloud{}
	// open file
	data, err := ioutil.ReadFile(CLOUD)
	if err != nil {
		return cloud, err
	}
	// parse
	if err = json.Unmarshal(data, &cloud); err != nil {
		return cloud, err
	}
	return cloud, nil
}

// GetContainer by name
func GetContainer(name string, cloud *storage.Cloud) (storage.Container, error) {
	if len(cloud.Containers) == 0 {
		return storage.Container{}, errors.New(".cloud has no containers")
	}
	if name == "" {
		fmt.Println("First container to be used:", cloud.Containers[0].Name)
		return cloud.Containers[0], nil
	}
	for _, v := range cloud.Containers {
		if v.Name == name {
			return v, nil
		}
	}
	return storage.Container{}, errors.New("no container found")
}

// IsExists checks if config exists
func IsExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Create config
func Create(filename string, v interface{}) error {
	template, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(filename, template, 0644); err != nil {
		panic(err)
	}
	return nil
}

func initConfigs() {
	// core
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	cloudcorepath := path.Join(u.HomeDir, CLOUDCORE)
	if ok := IsExists(cloudcorepath); !ok {
		core := storage.Core{
			Providers: []storage.Provider{
				{
					Provider: "CloudFiles",
					Name:     "myaccountname",
					Key:      "mykeyhere",
					Secret:   "mysecrethere",
					AuthURL:  "LON",
				},
			},
		}
		if err := Create(cloudcorepath, core); err != nil {
			panic(err)
		}
		fmt.Println("Initializing file:\t", cloudcorepath)
	}
	// cloud
	if ok := IsExists(CLOUD); !ok {
		cloud := storage.Cloud{
			Containers: []storage.Container{
				{
					Provider: "CloudFiles",
					Name:     "containername",
				},
			},
		}
		if err := Create(CLOUD, cloud); err != nil {
			panic(err)
		}
		fmt.Println("Initializing file:\t", CLOUD)
	}
	// TODO cloudignore
}

// Sync folder with defined container name string (default is the first container in the list in local .cloud)
func Sync(name string) {
	fmt.Printf("Beginning sync with container: \"%s\"\n", name)
	// get .cloudcore
	core, err := GetCore()
	if err != nil {
		panic(err)
	}
	fmt.Println("Found\t.cloudcore")
	// get .cloud
	cloud, err := GetCloud()
	if err != nil {
		panic(err)
	}
	fmt.Println("Found\t.cloud")
	// auth container
	// get container by name
	c, err := GetContainer(name, &cloud)
	if err != nil {
		panic(err)
	}
	// get container provider
	var s storage.Storage
	switch {
	case c.Provider == storage.CLOUDFILES:
		p, err := GetProvider(&c, &core)
		if err != nil {
			panic(err)
		}
		fmt.Println("Container found:", c.Name)
		if err != nil {
			panic(err)
		}
		s = &cloudfiles.Storage{p, swift.Connection{}}
	default:
		// TODO what would be better to write here
		fmt.Println("Something went wrong!")
		return
	}
	// Authenticate
	if err = s.Authenticate(); err != nil {
		panic(err)
	}
	fmt.Println("Authenticated")
	return
	// walk files upload file to the cloud
}
