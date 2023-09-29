# hpxd

`hpxd` is a daemon that runs on a node and manages the state of HAProxy. It listens for changes from a specified Git repository and updates the HAProxy configuration accordingly.

![hpxd logo](docs/assets/logo.png)

## Features

- **Dynamic Configuration Updates**: Polls a Git repository for changes in HAProxy configuration and applies them dynamically.
- **Prometheus Metrics**: Provides metrics on Git pull successes/failures, HAProxy reloads, and configuration validation.
- **Cross-Platform**: Builds available for Linux (`amd64` and `arm64`).

## Installation

### From Binary

You can download the latest release from [here](https://github.com/zcubbs/hpxd/releases).

If you prefer to install using a script:

```bash
curl -sL https://github.com/zcubbs/hpxd/scripts/install.sh | sudo bash
```

### From Source

```bash
go get github.com/zcubbs/hpxd
cd $GOPATH/src/github.com/zcubbs/hpxd
go install ./...
```

## Usage

You'll need a configuration file (by default, the tool looks for ./configs/hpxd.yaml).

```yaml
repoURL: https://github.com/yourusername/haproxy-configs.git
branch: main
haproxyConfigPath: /path/to/haproxy.cfg
pollingInterval: 5s
enablePrometheus: true
prometheusPort: 9100
```

Then, run the tool:

```bash
hpxd -c /path/to/config.yaml
```

## Development

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


