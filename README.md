# hpxd

`hpxd` is a daemon that runs on a node and manages the state of HAProxy. It listens for changes from a specified Git repository and updates the HAProxy configuration accordingly.

![hpxd logo](docs/assets/logo.png)

## Features

- **Git Integration**: Automatically detects changes in the HAProxy configuration from the Git repository.
- **Auto-Reload**: Gracefully reloads HAProxy whenever a valid configuration change is detected.
- **Validation**: Ensures that HAProxy configurations are valid before attempting a reload.

## Getting Started

### Prerequisites

- Go (version 1.17 or newer)
- Git
- HAProxy installed on the node

### Installation

1. Clone the repository:

```bash
git clone https://github.com/zcubbs/hpxd.git
```
2. Navigate to the project directory:

```bash
cd hpxd
```
3. Build the project:

```bash
go build
```
4. Run the project:

```bash
./hpxd
```

### Configuration

Create a configuration file named `hpxd.yaml` inside the `configs` directory.

Sample configuration:

```yaml
repoURL: https://github.com/user/haproxy-configs.git
branch: main
haproxyConfigPath: /path/to/haproxy/config
pollingInterval: 60  # in seconds
```

### Docker

1. Build the Docker image:

```bash
docker build -t hpxd .
```

2. Run the Docker image:

```bash
docker run -d --name hpxd hpxd
```

### Contributing

Contributions are welcome! Please read the contribution guidelines before submitting a pull request.

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

### Support and Feedback

If you need support or have any feedback, please open an issue [here]()


