#!/bin/bash

# Variables
ARCH=$(uname -m) # Detects architecture of the system
TARBALL_URL="https://github.com/zcubbs/hpxd/releases/latest/download/Hpxd_Linux_$ARCH.tar.gz"
INSTALL_DIR="/opt/hpxd"
LOG_DIR="$INSTALL_DIR/logs"
SERVICE_PATH="/etc/systemd/system/hpxd.service"
UNINSTALL_PATH="$INSTALL_DIR/uninstall.sh"

# Ensure the script is run as root
if [ "$(id -u)" -ne 0 ]; then
    echo "Please run this script as root."
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

# Create config file
echo "Creating config file at $INSTALL_DIR/config/hpxd.yaml..."
touch $INSTALL_DIR/config/hpxd.yaml

# Configure systemd service
echo "Configuring systemd service..."
cat <<EOL > $SERVICE_PATH
[Unit]
Description=HPXD Service
After=network.target

[Service]
ExecStart=$INSTALL_DIR/hpxd -config $INSTALL_DIR/config
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
