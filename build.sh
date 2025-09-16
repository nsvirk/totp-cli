#!/bin/bash

# Build script for TOTP generator

echo "Building CLI TOTP generator..."

# Build for current platform
go build -o totp main.go

echo "Build complete!"


# Install globally
echo "Installing globally..."
sudo cp totp /usr/local/bin/totp

echo "Install complete!"

rm totp

echo "Done!"