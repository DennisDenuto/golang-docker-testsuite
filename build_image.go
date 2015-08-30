package gotestwithdocker

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/samalba/dockerclient"
)

func BuildDockerFileTar(directoryPath string) (string, error) {
	dockerFileTar, err := os.Create(os.TempDir() + "/Dockerfile.tar")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	directories, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	tarWriter := tar.NewWriter(dockerFileTar)
	defer tarWriter.Close()

	for _, v := range directories {
		fileContent, _ := ioutil.ReadFile(directoryPath + "/" + v.Name())
		header, err := tar.FileInfoHeader(v, "")
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		tarWriter.WriteHeader(header)
		tarWriter.Write(fileContent)
	}

	tarWriter.Flush()
	return dockerFileTar.Name(), nil
}

func buildDockerImage() (string, error) {
	fmt.Println("Building docker test image")

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
