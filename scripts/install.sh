#!/bin/bash

# Variables
ARCH=$(uname -m) # Detects architecture of the system
TARBALL_URL="https://github.com/zcubbs/hpxd/releases/latest/download/Hpxd_Linux_$ARCH.tar.gz"
INSTALL_DIR="/opt/hpxd"
LOG_DIR="$INSTALL_DIR/logs"
SERVICE_PATH="/etc/systemd/system/hpxd.service"
UNINSTALL_PATH="$INSTALL_DIR/uninstall.sh"
HPXD_USER="hpxd" # User that will run the hpxd service

# Configuration Variables
REPO_URL=""
BRANCH="main"
REPO_FILE_PATH=""
HAPROXY_CONFIG_PATH=""
POLLING_INTERVAL="5s"
ENABLE_PROMETHEUS=true
PROMETHEUS_PORT=9100
GIT_USERNAME=""
GIT_PASSWORD=""
LOG_LEVEL="info" # default log level is 'info'

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
        --path)
        REPO_FILE_PATH="$2"
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
        --git-username)
        GIT_USERNAME="$2"
        shift
        shift
        ;;
        --git-password)
        GIT_PASSWORD="$2"
        shift
        shift
        ;;
        --log-level)
        LOG_LEVEL="$2"
        shift
        shift
              ;;
        *)    # unknown option
        shift # past argument
        ;;
    esac
done

# Check if mandatory parameters are provided
if [[ -z "$REPO_URL" || -z "$REPO_FILE_PATH" || -z "$HAPROXY_CONFIG_PATH" ]]; then
    echo "Mandatory parameters --repo-url, --path and --haproxy-config-path are missing."
    exit 1
fi

# Download, extract and install binary
echo "Downloading and installing hpxd..."
mkdir -p $INSTALL_DIR
curl -L -o "$INSTALL_DIR/hpxd.tar.gz" "$TARBALL_URL"
tar -xzf "$INSTALL_DIR/hpxd.tar.gz" -C $INSTALL_DIR
rm "$INSTALL_DIR/hpxd.tar.gz"

# Create logs directory
echo "Creating logs directory at $LOG_DIR..."
mkdir -p $LOG_DIR

# Create config directory
echo "Creating config directory at $INSTALL_DIR/config..."
mkdir -p $INSTALL_DIR/config

# Create a dedicated user for hpxd and grant it necessary permissions
if ! id "$HPXD_USER" &>/dev/null; then
    useradd -r -s /sbin/nologin $HPXD_USER
fi
chown $HPXD_USER: $INSTALL_DIR -R
chown $HPXD_USER: "$HAPROXY_CONFIG_PATH"

# Allow hpxd user permissions for haproxy
cat <<EOL | sudo tee /etc/sudoers.d/hpxd_permissions
$HPXD_USER ALL=NOPASSWD: /bin/systemctl reload haproxy
$HPXD_USER ALL=NOPASSWD: /bin/systemctl status haproxy
EOL

# Create and pre-populate the config file
echo "Creating and pre-populating config file at $INSTALL_DIR/config/hpxd.yaml..."
cat <<EOL > $INSTALL_DIR/config/hpxd.yaml
repoURL: $REPO_URL
branch: $BRANCH
path: $REPO_FILE_PATH
haproxyConfigPath: $HAPROXY_CONFIG_PATH
pollingInterval: $POLLING_INTERVAL
enablePrometheus: $ENABLE_PROMETHEUS
prometheusPort: $PROMETHEUS_PORT
logLevel: $LOG_LEVEL
EOL

# Create an environment file for hpxd
ENV_FILE="$INSTALL_DIR/.hpxd_vars"

echo "Creating environment file at $ENV_FILE..."
touch $ENV_FILE

# If GIT_USERNAME is passed as an argument, write to the environment file
if [[ -n "$GIT_USERNAME" ]]; then
    echo "HPXD_GIT_USERNAME=$GIT_USERNAME" >> $ENV_FILE
fi

# If GIT_PASSWORD is passed as an argument, write to the environment file
if [[ -n "$GIT_PASSWORD" ]]; then
    echo "HPXD_GIT_PASSWORD=$GIT_PASSWORD" >> $ENV_FILE
fi

# Set permissions for the environment file so only hpxd (or root) can read it
chown $HPXD_USER: $ENV_FILE
chmod 600 $ENV_FILE

# Configure systemd service
echo "Configuring systemd service..."
cat <<EOL > $SERVICE_PATH
[Unit]
Description=HPXD Service
After=network.target

[Service]
ExecStart=$INSTALL_DIR/hpxd -config $INSTALL_DIR/config
EnvironmentFile=$ENV_FILE
Restart=always
User=$HPXD_USER
Group=nogroup
Environment=PATH=/usr/bin:/usr/local/bin:/usr/sbin
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

# Remove sudoers permission for hpxd user
rm -f /etc/sudoers.d/hpxd_permissions

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
