# MCP Argo Server

An MCP-compliant server for running Argo Workflows written in Golang.

## Overview

MCP Argo Server is a lightweight CLI tool that wraps Argo Workflows using JSON-RPC over STDIN/STDOUT. It leverages Foxy Contexts for RPC handling and client-go for interacting with Kubernetes and Argo Workflow resources. The project provides tools for launching workflows, checking workflow status, and retrieving results.

## Installation

This project is configured to run inside a development container. Simply open the repository in your dev container-enabled editor (e.g., VS Code Remote - Containers) and all dependencies are pre-installed.  
If you prefer to run it locally, clone the repository and run:
   ```
   go mod tidy
   ```

## Usage

Open the project in the dev container. 

Run `make cluster` which will install the k3d cluster and set up Argo. 

You can check that's worked by typing `kubectl cluster-info`. 

You can run a test workflow by typing `argo submit -n argo --watch ./kube/argo-hello-world.yaml`.

You can see the Argo interface at [https://localhost:2746/workflows/argo/](https://localhost:2746/workflows/argo/)

You can check that the app is building and the MCP is working by typing `make run`.

## Testing with Python

The project includes a Python test client that demonstrates how to interact with the MCP Argo server. The test client is located in `python/test_with_autogen.py` and showcases:

- Submitting Argo workflows
- Checking workflow status
- Waiting for workflow completion
- Retrieving workflow results

To run the Python test:

1. Ensure you have Python dependencies installed:
   ```bash
   cd python
   make install
   ```

2. Run the test script:
   ```bash
   python test_with_autogen.py
   ```

or... just debug it to step through. 

The script will:
- Connect to the MCP Argo server
- Iterate the tools and print them out
- Submit a sample workflow from `kube/argo-hello-world.yaml`
- Monitor the workflow status until completion
- Display the workflow results

## Contributing

Contributions are welcome! Please open issues and submit pull requests. Before submitting changes, ensure that you follow the project's coding guidelines and that all tests pass.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

## Additional Resources

- [Argo Workflows](https://argoproj.github.io/argo-workflows/)
- [Foxy Contexts](https://github.com/strowk/foxy-contexts)
