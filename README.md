# canvas

This repository is used in the course [Build Cloud Apps in Go](https://www.golang.dk/courses/build-cloud-apps-in-go).


## Installation

### Pre-Requiste

1. Install `air` for developing locally (_this program provides hot reloding for any code changes you make_)
    ```shell
    curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
    sudo mv ./bin/air /usr/local/bin/air
    ```


### General

1. Install dependencies
    ```shell
    make
    ```

1. Run locally:
    ```shell
    # start docker db
    make db-start

    # run app
    make run
    ```

1. Run tests:
    ```shell
    make test
    ```

1. Run integration tests:
    ```shell
    make test-integration
    ```

1. Deploy to ACI (Azure Container Instances):
    ```shell
    make aci-deploy
    ```
