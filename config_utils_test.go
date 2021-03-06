package gotestwithdocker

import (
	"fmt"
	"testing"
)

func TestGetExposePorts(t *testing.T) {
	yamlContent := `
ports:
   - "8080:8081"
`
	config, err := NewConfigContent(yamlContent)
	ports, err := config.GetExposePorts()

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	portMapping, ok := ports["8080/tcp"]
	if !ok {
		fmt.Println("port parsing failed")
		t.Fail()
	}

	if portMapping[0].HostPort != "8081" {
		fmt.Println("port parsing failed")
		t.Fail()
	}
}

func TestGetDockerBuildDirectory(t *testing.T) {
	yamlContent := `
build:
   dockerfile_dir: ./some-directory/
`
	config, err := NewConfigContent(yamlContent)
	buildDirectory, err := config.GetBuildContextDirectory()

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if buildDirectory != "./some-directory/" {
		fmt.Printf("yaml parsing failed, expected :%s, GOT: %s", "./some-directory/", buildDirectory)
		t.Fail()
	}
}
