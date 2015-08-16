package gotestwithdocker

import (
	"bufio"
	"fmt"
	"os"

	"github.com/samalba/dockerclient"
)

func buildTestDockerImage() (string, error) {
	fmt.Println("Building jenkins image with test data")

	dockerFile, err := os.Open("./docker/Dockerfile.tar")
	if err != nil {
		return "", err
	}
	defer dockerFile.Close()

	imageName, _ := testConfig.GetImageName()
	reader, err := docker.BuildImage(&dockerclient.BuildImage{
		Context:        dockerFile,
		RepoName:       imageName,
		SuppressOutput: false,
		Remove:         true,
	})
	defer reader.Close()

	if err != nil {
		fmt.Printf("error building image: %s", err)
		return "", err
	}
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() == true {
		fmt.Println(scanner.Text())
	}

	return imageName, nil
}
