// Package main provides ...
package gotestwithdocker

import (
	"archive/tar"
	"fmt"
	"os"
	"testing"
)

func TestBuildDockerFileTar(t *testing.T) {
	dockerFileTar, err := BuildDockerFileTar("./test_directory")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	tarFile, err := os.OpenFile(dockerFileTar, 0, 0666)
	tarReader := tar.NewReader(tarFile)
	for {
		header, _ := tarReader.Next()
		if header == nil {
			break
		}
		fmt.Println(header.Name)
	}
}
