package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/docker"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/shell"
	"golang.org/x/exp/slices"
)

func BuildAndRunContainer(t *testing.T, runOptions *docker.RunOptions, buildOptions *docker.BuildOptions) string {
	fmt.Println("--- Building Dockerfile")

	docker.Build(t, "./", buildOptions)

	fmt.Println("--- Running image")

	return docker.RunAndGetID(t, buildOptions.Tags[0], runOptions)
}

func BuildContainer(t *testing.T, buildOptions *docker.BuildOptions) {
	fmt.Println("--- Building Dockerfile")

	docker.Build(t, "./", buildOptions)
}

func RemoveContainer(t *testing.T, id string) {
	cmd := shell.Command{
		Command: "docker",
		Args:    []string{"container", "rm", "--force", id},
	}

	shell.RunCommand(t, cmd)
}

func GetHostPort(t *testing.T, containerId string, port uint16) uint16 {
	c := docker.Inspect(t, containerId)
	idx := slices.IndexFunc(c.Ports, func(p docker.Port) bool { return p.ContainerPort == port })

	return c.Ports[idx].HostPort
}

func GetHostPWD() string {
	pwd, _ := os.Getwd()

	return fmt.Sprint(strings.Replace(pwd, "/app", os.Getenv("HOST_PWD"), 1))
}

func HttpGet(t *testing.T, url string, validate func(int, string) bool) {
	tlsConfig := tls.Config{}
	maxRetries := 6
	timeBetweenRetries := 5 * time.Second

	http_helper.HttpGetWithRetryWithCustomValidation(t, url, &tlsConfig, maxRetries, timeBetweenRetries, validate)
}

func GetArgFromDockerfile(t *testing.T, arg string, dockerfile string) string {
	// read the arg from Dockerfile ARG
	dat, err := os.ReadFile(dockerfile)
	if err != nil {
		log.Fatal(err)
	}
	// convert []byte to string and split by new line
	lines := strings.Split(string(dat), "\n")
	// iterate over lines and find the line with ARG %arg
	for _, line := range lines {
		if strings.Contains(line, fmt.Sprintf("ARG %s", arg)) {
			words := strings.Split(line, " ")
			value := strings.Split(words[len(words)-1], "=")[1]
			// remove the quotes
			value = strings.ReplaceAll(value, "\"", "")
			return value
		}
	}
	log.Fatal(fmt.Sprintf("%s not found in Dockerfile", arg))
	return ""
}

func main() {
	fmt.Println("Hello, World!")
}
