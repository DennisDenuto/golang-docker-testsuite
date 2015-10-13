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

type DockerSuite struct {
	ConfigYaml string
}

var docker *dockerclient.DockerClient
var testConfig *DockerTestConfig

func (s *DockerSuite) SetUpSuite(c *C) {
	var err error = nil
	if s.SetUpSuite == nil {
		fmt.Println("Missing docker config file")
		c.FailNow()
	}
	testConfig, err = NewConfig(s.ConfigYaml)
	failOnError(err, c)

	err = initDockerClient()
	failOnError(err, c)

	if testConfig.HasBuildConfig() {
		buildDockerImage()
	}

	containerId, err := createTestContainer(getContainerName())
	failOnError(err, c)

	exposedPorts, _ := testConfig.GetExposePorts()
	fmt.Println(exposedPorts)
	if err := docker.StartContainer(containerId, &dockerclient.HostConfig{PortBindings: exposedPorts}); err != nil {
		fmt.Printf("error starting container: %s", err)
		c.FailNow()
	}
	c.Assert(waitForContainerToStartup(), Equals, true)
}

func (s *DockerSuite) TearDownSuite(c *C) {
	containerName := getContainerName()
	docker.KillContainer(containerName, "9")
	docker.RemoveContainer(containerName, false, true)
}

func initDockerClient() error {
	var err error
	var tlsc tls.Config
	var cert tls.Certificate

	cert, err = tls.LoadX509KeyPair(getX509KeyPairConfig())
	if err != nil {
		return err
	}
	tlsc.Certificates = append(tlsc.Certificates, cert)
	tlsc.InsecureSkipVerify = true
	docker, err = dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), &tlsc)
	return err
}

func getX509KeyPairConfig() (certFile, keyFile string) {
	return os.Getenv("DOCKER_CERT_PATH") + "/cert.pem", os.Getenv("DOCKER_CERT_PATH") + "/key.pem"
}

func getContainerName() string {
	return testConfig.GetContainerName()
}

func waitForContainerToStartup() bool {
	logMessageIndicatingTestHasStarted, _ := testConfig.GetWaitLogMessage()
	reader, err := docker.ContainerLogs(
		getContainerName(),
		&dockerclient.LogOptions{Stdout: true, Stderr: true, Follow: true},
	)
	if err != nil {
		fmt.Printf("error reading container logs: %s", err)
		return false
	}
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
	imageName, err := testConfig.GetImageName()
	if err != nil {
		return "", err
	}
	containerConfig := &dockerclient.ContainerConfig{
		Image:       imageName + ":latest",
		AttachStdin: false,
		Tty:         false}
	return docker.CreateContainer(containerConfig, containerName)
}

func failOnError(err error, c *C) {
	if err != nil {
		fmt.Printf("error: %s", err)
		c.FailNow()
	}
}

func findIp(input string) string {
	numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	regexPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock

	regEx := regexp.MustCompile(regexPattern)
	return regEx.FindString(input)
}
