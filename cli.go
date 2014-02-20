package main

import (
	"encoding/json"
	"fmt"
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

type (
	// Provider for the cloud like Amazon, Rackspace etc.
	Provider struct {
		Name    string `json:"name"`
		Key     string `json:"key"`
		Secret  string `json:"secret,omitempty"`
		AuthURL string `json:"auth_url,omitempty"`
	}

	// Core is ~/.cloudcore struct
	Core struct {
		Providers []Provider `json:"providers"`
	}

	// Container data to work with later
	Container struct {
		Provider string `json:"provider"`
		Name     string `json:"name"`
	}
	// Cloud is config for .cloud
	Cloud struct {
		Containers []Container `json:"containers"`
	}
)

// GetCore returns parsed .cloudcore struct
func GetCore() (Core, error) {
	core := Core{}
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

// GetCloud returns parsed .cloud struct
func GetCloud() (Cloud, error) {
	cloud := Cloud{}
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
		return err
	}
	if err := ioutil.WriteFile(filename, template, 0644); err != nil {
		return err
	}
	return nil
}

func initConfigs() {
	// core
	u, err := user.Current()
	check(err)
	cloudcorepath := path.Join(u.HomeDir, CLOUDCORE)
	if ok := IsExists(cloudcorepath); !ok {
		core := Core{
			Providers: []Provider{
				{
					Name:    "CloudFiles",
					Key:     "mykeyhere",
					Secret:  "mysecrethere",
					AuthURL: "https://storage101.lon3.clouddrive.com/v1/MossoCloudFS_aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/",
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
		cloud := Cloud{
			Containers: []Container{
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

// Sync folder with defined container string (default is the first container in the list)
func Sync(container string) {
	fmt.Println("Begin sync with container: " + container)
	// get .cloudcore
	core, err := GetCore()
	check(err)
	fmt.Println(core)
	// get .cloud
	cloud, err := GetCloud()
	check(err)
	fmt.Println(cloud)
	// auth container
	// walk files
	// upload file to the cloud
}
