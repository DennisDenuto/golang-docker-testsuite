package gotestwithdocker

import (
	"archive/tar"
	"bufio"
	"io/ioutil"
	"os"

	"github.com/samalba/dockerclient"
)

func BuildDockerFileTar(directoryPath string) (string, error) {
	dockerFileTar, err := os.Create(os.TempDir() + "/Dockerfile.tar")
	if err != nil {
		log.Error("%s", err)
		return "", err
	}
	directories, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		log.Error("%s", err)
		return "", err
	}
	tarWriter := tar.NewWriter(dockerFileTar)
	defer tarWriter.Close()

	for _, v := range directories {
		fileContent, _ := ioutil.ReadFile(directoryPath + "/" + v.Name())
		header, err := tar.FileInfoHeader(v, "")
		if err != nil {
			log.Error("%s", err)
			return "", err
		}
		tarWriter.WriteHeader(header)
		tarWriter.Write(fileContent)
	}

	tarWriter.Flush()
	return dockerFileTar.Name(), nil
}

func pullDockerImage() (string, error) {
	log.Info("Pulling docker test image")
	imageName, _ := testConfig.GetImageName()
	log.Info(imageName)
	err := docker.PullImage(imageName, &dockerclient.AuthConfig{})

	if err != nil {
		log.Error("error building image: %s", err)
		return "", err
	}

	return imageName, nil
}

func buildDockerImage() (string, error) {
	log.Notice("Building docker test image")
	buildContextDir, err := testConfig.GetBuildContextDirectory()
	dockerFileTar, err := BuildDockerFileTar(buildContextDir)
	if err != nil {
		return "", err
	}
	dockerFile, err := os.Open(dockerFileTar)
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
		log.Error("error building image: %s", err)
		return "", err
	}
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() == true {
		log.Notice(scanner.Text())
	}

	return imageName, nil
}
