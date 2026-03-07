#!/usr/bin/env bash
set -e

# detect OS & ARCH
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
  ARCH=amd64
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
  ARCH=arm64
fi

# latest release
TAG=$(curl -s https://api.github.com/repos/nxrmqlly/silo/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# binary URL
URL="https://github.com/nxrmqlly/silo/releases/download/${TAG}/silo-${OS}-${ARCH}"
if [[ "$OS" == "windows" ]]; then
  URL="${URL}.exe"
fi

echo $URL
echo "Downloading Silo $TAG..."
curl -L -o silo "$URL"
chmod +x silo
sudo mv silo /usr/local/bin/silo

echo "Installed Silo $TAG!"