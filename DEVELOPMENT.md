# Combo Development Guidelines
Thank you for seeking to contribute to Combo! Any contribution is welcome and considered. 

Want to get involved? Please check the [Github project](https://github.com/operator-framework/combo/projects/1) for Combo. 

Have an issue you would like fix? Feel free to raise one [here](https://github.com/operator-framework/combo/issues/new/choose)!

## Requirements

### Tooling

| Requirement | Purpose               | macOS                |
|-------------|-----------------------|----------------------|
| Go          | Compiler              | brew install go      |
| Docker      | Packaging             | [Docker for Mac]     |
| kubebuilder | Code generation       | [kubebuilder docs]   |
| Ginkgo      | Testing               | [Ginkgo docs]        |

[Docker for Mac]: https://store.docker.com/editions/community/docker-ce-desktop-mac
[kubebuilder docs]: https://book.kubebuilder.io/quick-start.html#installation
[Ginkgo docs]: https://onsi.github.io/ginkgo/

### E2E test environments

| Requirement | install docs         |
|-------------|----------------------|
| Minikube    | [Minikube docs]      |
| Kind        | [Kind docs]          |

[Minikube docs]: https://minikube.sigs.k8s.io/docs/start
[Kind docs]: https://kind.sigs.k8s.io/docs/user/quick-start

## Usage

### Commit Messages
Commits made in this repo follow Conventional Commit standards. For more information on this please check out [the offical documentation](https://www.conventionalcommits.org/en/v1.0.0/).

### Local targets
When building the Docker container for this repository, it can often take a long amount of time. The majority of this time is spent building the binary as it downloads the dependencies. To get around this, local targets will build the binary on your local machine and then copy it onto the Docker container.

This speeds build and test times up significantly when developing locally. All of these targets have production ready versions by simply remove the `local` portion of the target.

### Testing

This project uses the built-in testing support for golang. To run all of the unit tests, run:
```bash
$ make test-unit
```

For this project, the E2E test suite is written using Ginkgo. Prior to running them, you will need to make sure you have deployed Combo's manifests onto a local cluster. Luckily, there is a wrapper to make this process doable in one command.
```bash
$ make run-e2e-local
```

**NOTE:** The `run-e2e-local` target supports Minikube and Kind environments. If you want to run the e2e tests on Minikube, you will need to make sure Minikube is deployed in your local environment. Additionally, you will need to specify the command that will load the image onto the cluster. 

```bash
$ make IMAGE_LOAD_COMMAND="minikube image load" run-e2e-local
```

If you want to run the e2e tests on Kind, you will need to make sure Kind is deployed in your local environment and switch the kubeconfig to an existing Kind cluster. This is done by default when running`run-e2e-local`.

### Building

Ensure your version of go is up to date; check that you're running the same version as in go.mod with the
commands:
```bash
$ head go.mod
$ go version
```

To build the go binary, run:
```bash
$ make build-cli
```

To build the docker container, run:
```bash
$ make build-container
```

### Loading

To load the controller, run:
```bash
$ make load-image
```

To load the bundle, run:
```bash
$ make load-image IMAGE=quay.io/operator-framework/combo-bundle
```

### Style
This project utilizes various validators that can be run with a single aggregate make target.

```bash
$ make verify
```

When running `make verify` a number of things happen.
1. `make tidy` - Ensure all dependencies are accounted for (or removed if not used).
2. `make generate` - Use Go's built in generation tools to look in the code for generation targets
3. `make format` - Ensure code meets Go's built in formatter
4. `make lint` - Use `golangci-lint` tool to run a collection of various linters. Check [here](https://github.com/golangci/golangci-lint) for more info.

This makes the most sense to run ***before pushing your code up***. If you do not, the CI/CD will catch any issues and fail a Github action as a result.
