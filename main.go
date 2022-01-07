package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// DockerConfigJSON represents a local docker auth config file
// for pulling images.
type DockerConfigJSON struct {
	Auths DockerConfig `json:"auths" datapolicy:"token"`
	// +optional
	HttpHeaders map[string]string `json:"HttpHeaders,omitempty" datapolicy:"token"`
}

// DockerConfig represents the config file used by the docker CLI.
// This config that represents the credentials that should be used
// when pulling images from specific image repositories.
type DockerConfig map[string]DockerConfigEntry

// DockerConfigEntry holds the user information that grant the access to docker registry
type DockerConfigEntry struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty" datapolicy:"password"`
	Email    string `json:"email,omitempty"`
	Auth     string `json:"auth,omitempty" datapolicy:"token"`
}

// handleDockerCfgJSONContent serializes a ~/.docker/config.json file
func handleDockerCfgJSONContent(username, password, email, server string) ([]byte, error) {
	dockerConfigAuth := DockerConfigEntry{
		Username: username,
		Password: password,
		Email:    email,
		Auth:     encodeDockerConfigFieldAuth(username, password),
	}
	dockerConfigJSON := DockerConfigJSON{
		Auths: map[string]DockerConfigEntry{server: dockerConfigAuth},
	}

	return json.Marshal(dockerConfigJSON)
}

// encodeDockerConfigFieldAuth returns base64 encoding of the username and password string
func encodeDockerConfigFieldAuth(username, password string) string {
	fieldValue := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(fieldValue))
}

func main() {
	username := ""
	password := ""
	email := ""
	server := ""

	flag.StringVar(&username, "username", "", "Registry Username")
	flag.StringVar(&password, "password", "", "Registry Password")
	flag.StringVar(&email, "email", "", "Registry email")
	flag.StringVar(&server, "server", "", "Registry Url")

	flag.Parse()

	if len(email) == 0 || len(username) == 0 || len(password) == 0 || len(server) == 0 {
		fmt.Println("Inputs Missing ! Usage: gen-dockercfg")
		flag.PrintDefaults()
		os.Exit(1)
	}

	dockerConfigJSONContent, err := handleDockerCfgJSONContent(username, password, email, server)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Print(string(dockerConfigJSONContent))
	os.Exit(0)
}
