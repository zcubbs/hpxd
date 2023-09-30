#!/bin/bash

# Variables
ARCH=$(uname -m) # Detects architecture of the system
TARBALL_URL="https://github.com/zcubbs/hpxd/releases/latest/download/Hpxd_Linux_$ARCH.tar.gz"
INSTALL_DIR="/opt/hpxd"
LOG_DIR="$INSTALL_DIR/logs"
SERVICE_PATH="/etc/systemd/system/hpxd.service"
UNINSTALL_PATH="$INSTALL_DIR/uninstall.sh"

# Configuration Variables
REPO_URL=""
BRANCH="main"
HAPROXY_CONFIG_PATH=""
POLLING_INTERVAL="5s"
ENABLE_PROMETHEUS=true
PROMETHEUS_PORT=9100

# Ensure the script is run as root
if [ "$(id -u)" -ne 0 ]; then
    echo "Please run this script as root."
    exit 1
fi

# Parse arguments
while [[ $# -gt 0 ]]
do
    key="$1"

    case $key in
        --repo-url)
        REPO_URL="$2"
        shift
        shift
        ;;
        --branch)
        BRANCH="$2"
        shift
        shift
        ;;
        --haproxy-config-path)
        HAPROXY_CONFIG_PATH="$2"
        shift
        shift
        ;;
        --polling-interval)
        POLLING_INTERVAL="$2"
        shift
        shift
        ;;
        --enable-prometheus)
        ENABLE_PROMETHEUS="$2"
        shift
        shift
        ;;
        --prometheus-port)
        PROMETHEUS_PORT="$2"
        shift
        shift
        ;;
        *)    # unknown option
        shift # past argument
        ;;
    esac
done

# Check if mandatory parameters are provided
if [[ -z "$REPO_URL" || -z "$HAPROXY_CONFIG_PATH" ]]; then
    echo "Mandatory parameters --repo-url and --haproxy-config-path are missing."
    exit 1
fi

# Download, extract and install binary
echo "Downloading and installing hpxd..."
mkdir -p $INSTALL_DIR
curl -L -o "$INSTALL_DIR/hpxd.tar.gz" $TARBALL_URL
tar -xzf "$INSTALL_DIR/hpxd.tar.gz" -C $INSTALL_DIR
rm "$INSTALL_DIR/hpxd.tar.gz"

# Create logs directory
echo "Creating logs directory at $LOG_DIR..."
mkdir -p $LOG_DIR

# Create config directory
echo "Creating config directory at $INSTALL_DIR/config..."
mkdir -p $INSTALL_DIR/config

# Create and pre-populate the config file
echo "Creating and pre-populating config file at $INSTALL_DIR/config/hpxd.yaml..."
cat <<EOL > $INSTALL_DIR/config/hpxd.yaml
repoURL: $REPO_URL
branch: $BRANCH
haproxyConfigPath: $HAPROXY_CONFIG_PATH
pollingInterval: $POLLING_INTERVAL
enablePrometheus: $ENABLE_PROMETHEUS
prometheusPort: $PROMETHEUS_PORT
EOL

# Configure systemd service
echo "Configuring systemd service..."
cat <<EOL > $SERVICE_PATH
[Unit]
Description=HPXD Service
After=network.target

[Service]
ExecStart=$INSTALL_DIR/hpxd -config $INSTALL_DIR/config/hpxd.yaml
Restart=always
User=nobody
Group=nogroup
Environment=PATH=/usr/bin:/usr/local/bin
WorkingDirectory=$INSTALL_DIR
StandardOutput=append:$LOG_DIR/output.log
StandardError=append:$LOG_DIR/error.log

[Install]
WantedBy=multi-user.target
EOL

# Generate uninstall script
echo "Generating uninstall script..."
cat <<EOL > $UNINSTALL_PATH
#!/bin/bash

# Stop and disable service
systemctl stop hpxd
systemctl disable hpxd

# Remove systemd service
rm -f $SERVICE_PATH

# Remove installation directory
rm -rf $INSTALL_DIR

echo "Uninstallation complete."

EOL
chmod +x $UNINSTALL_PATH

# Reload systemd, enable and start service
systemctl daemon-reload
systemctl enable hpxd
systemctl start hpxd

echo "Installation complete. HPXD is now running as a systemd service."
echo "To uninstall, run $UNINSTALL_PATH"
