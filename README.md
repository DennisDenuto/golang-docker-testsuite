# gotestwithdocker

A testing library (for golang projects) to assist in writing integration tests with docker!.

You specify the docker container you want running for your integration tests.

The library will spin it up for your tests to use and tear it down at the end.

## Usage

Inside your test file:

        import "github.com/DennisDenuto/golang-docker-testsuite"

        func Test(t *testing.T) { TestingT(t) }

        type LocalDockerSuite struct {
        	gotestwithdocker.DockerSuite
        }

        var _ = Suite(&LocalDockerSuite{gotestwithdocker.DockerSuite{ConfigYaml: "./docker-config.yaml"}})

minimalistic docker-config.yaml

        container_name: some-name
        image: a-docker-image

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

