package gotestwithdocker

import (
	"os"

	"github.com/olebedev/config"
)

type DockerTestConfig struct {
	config *config.Config
}

func NewConfig() (*DockerTestConfig, error) {
	file, err := os.OpenFile("docker/config.yaml", os.O_RDONLY, 0666)
	if err != nil {
		return &DockerTestConfig{}, err
	}
	config, err := config.ParseYamlFile(file.Name())
	return &DockerTestConfig{config}, err
}

func (testconfig *DockerTestConfig) GetContainerName() (string, error) {
	return testconfig.config.String("container_name")
}

func (testconfig *DockerTestConfig) GetImageName() (string, error) {
	return testconfig.config.String("image")
}

func (testconfig *DockerTestConfig) GetWaitLogMessage() (string, error) {
	return testconfig.config.String("wait.log")
}

func (testconfig *DockerTestConfig) GetWaitTimeout() int {
	return testconfig.config.UInt("wait.timeout", 180)
}

func (testconfig *DockerTestConfig) HasBuildConfig() bool {
	_, err := testconfig.config.Get("build.dockerfile_dir")
	return err == nil
}
