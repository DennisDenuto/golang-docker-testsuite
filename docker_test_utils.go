package gotestwithdocker

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/samalba/dockerclient"
	. "gopkg.in/check.v1"
)

type DockerSuite struct{}

var docker *dockerclient.DockerClient
var testConfig *DockerTestConfig

func (s *DockerSuite) SetUpSuite(c *C) {
	var err error = nil
	testConfig, err = NewConfig()
	if err != nil {
		fmt.Printf("error parsing config yaml file: %s", err)
		c.FailNow()
	}

	var tlsc tls.Config
	cert, _ := tls.LoadX509KeyPair(os.Getenv("DOCKER_CERT_PATH")+"/cert.pem", os.Getenv("DOCKER_CERT_PATH")+"/key.pem")
	tlsc.Certificates = append(tlsc.Certificates, cert)
	tlsc.InsecureSkipVerify = true
	docker, _ = dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), &tlsc)

	if testConfig.HasBuildConfig() {
		buildTestDockerImage()
	}

	containerName, _ := testConfig.GetContainerName()
	containerId, err := createTestContainer(containerName)
	if err != nil {
		fmt.Printf("error creating container: %s", err)
		c.FailNow()
	}

	if err := docker.StartContainer(containerId, &dockerclient.HostConfig{PublishAllPorts: true}); err != nil {
		fmt.Printf("error starting container: %s", err)
		c.FailNow()
	}
	c.Assert(waitForContainerToStartup(), Equals, true)
}

func (s *DockerSuite) TearDownSuite(c *C) {
	containerName, _ := testConfig.GetContainerName()
	docker.KillContainer(containerName, "9")
	docker.RemoveContainer(containerName, false, true)
}

func waitForContainerToStartup() bool {
	containerName, _ := testConfig.GetContainerName()
	logMessageIndicatingTestHasStarted, _ := testConfig.GetWaitLogMessage()
	reader, _ := docker.ContainerLogs(containerName, &dockerclient.LogOptions{Stdout: true, Stderr: true, Follow: true})
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	success := make(chan bool)

	go func() {
		for scanner.Scan() == true {
			if strings.Contains(scanner.Text(), logMessageIndicatingTestHasStarted) {
				success <- true
				break
			}
		}
	}()

	go func() {
		time.Sleep(time.Duration(testConfig.GetWaitTimeout()) * time.Second)
		success <- false
	}()

	return <-success
}

func createTestContainer(containerName string) (string, error) {
	imageName, _ := testConfig.GetImageName()
	containerConfig := &dockerclient.ContainerConfig{
		Image:       imageName + ":latest",
		AttachStdin: false,
		Tty:         false}
	return docker.CreateContainer(containerConfig, containerName)
}

func findIp(input string) string {
	numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	regexPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock

	regEx := regexp.MustCompile(regexPattern)
	return regEx.FindString(input)
}
