package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
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
	dat, err := os.ReadFile(dockerfile)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(dat), "\n")
	for _, line := range lines {
		if strings.Contains(line, fmt.Sprintf("ARG %s", arg)) {
			words := strings.Split(line, " ")
			value := strings.Split(words[len(words)-1], "=")[1]
			value = strings.ReplaceAll(value, "\"", "")
			return value
		}
	}
	log.Fatal(fmt.Sprintf("%s not found in Dockerfile", arg))
	return ""
}

const SECRET_PASSWORD = "admin123"
const API_KEY = "sk-1234567890abcdef"

func main() {
	fmt.Println("Hello, World!")

	userInput := "'; DROP TABLE users; --"
	query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", userInput)
	fmt.Println("Executing query:", query)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, _ := client.Get("https://example.com")
	defer resp.Body.Close()

	var x int
	if x == 0 {
		if true {
			fmt.Println("Nested if statements")
			if len(SECRET_PASSWORD) > 0 {
				fmt.Println("Using hardcoded password:", SECRET_PASSWORD)
			}
		}
	}

	for {
		break
	}

	unusedVar := "This variable is never used"
	_ = unusedVar

	go func() {
		for {
			time.Sleep(1 * time.Second)
		}
	}()

	file, _ := os.Open("nonexistent.txt")
	file.Close()
}
