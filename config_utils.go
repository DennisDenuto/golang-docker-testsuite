package gotestwithdocker

import (
	"os"
	"strings"

	"github.com/olebedev/config"
	"github.com/samalba/dockerclient"
)

type DockerTestConfig struct {
	config *config.Config
}

func NewConfigContent(configContent string) (*DockerTestConfig, error) {
	config, err := config.ParseYaml(configContent)
	return &DockerTestConfig{config}, err
}

func NewConfig(configYamlFile string) (*DockerTestConfig, error) {
	file, err := os.OpenFile(configYamlFile, os.O_RDONLY, 0666)
	if err != nil {
		return &DockerTestConfig{}, err
	}
	config, err := config.ParseYamlFile(file.Name())
	return &DockerTestConfig{config}, err
}

func (testconfig *DockerTestConfig) GetContainerName() string {
	return testconfig.config.UString("container_name", "test-container-name")
}

func (testconfig *DockerTestConfig) GetImageName() (string, error) {
	return testconfig.config.String("image")
}

func (testconfig *DockerTestConfig) GetExposePorts() (map[string][]dockerclient.PortBinding, error) {
	var exposePorts = make(map[string][]dockerclient.PortBinding, 0)
	ports, err := testconfig.config.List("ports")
	for _, v := range ports {
		yamlPortValue, _ := v.(string)
		portSplit := strings.Split(yamlPortValue, ":")
		containerPort := portSplit[0] + "/tcp"
		hostPort := portSplit[1]
		exposePorts[containerPort] = []dockerclient.PortBinding{{HostPort: hostPort}}
	}
	return exposePorts, err
}

func (testconfig *DockerTestConfig) GetWaitLogMessage() (string, error) {
	return testconfig.config.String("wait.log")
}

func (testconfig *DockerTestConfig) GetWaitTimeout() int {
	return testconfig.config.UInt("wait.timeout", 180)
}

func (testconfig *DockerTestConfig) GetBuildContextDirectory() (string, error) {
	return testconfig.config.String("build.dockerfile_dir")
}

func (testconfig *DockerTestConfig) HasBuildConfig() bool {
	_, err := testconfig.config.Get("build.dockerfile_dir")
	return err == nil
}
