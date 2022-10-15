# canvas

This repository is used in the course [Build Cloud Apps in Go](https://www.golang.dk/courses/build-cloud-apps-in-go).


## Installation

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
