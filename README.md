# TOTP Generator for CLI with Clipboard Support

A fast, secure command-line TOTP (Time-based One-Time Password) generator that automatically copies codes to your clipboard. Built in Go for maximum performance and zero dependencies.

## 🎯 Key Features

✅ **Auto-clipboard copy** - Automatically copies TOTP to clipboard for instant pasting
✅ **Single executable** - Golang Binary with no dependencies
✅ **Multiple output modes** - Print + copy, quiet mode, or print-only
✅ **Case insensitive** - `totp USER_1` works the same as `totp user_1`
✅ **Cross-platform** - Works on macOS, Linux, Windows
✅ **Lightning fast** - Pure Go implementation, sub-millisecond generation
✅ **Secure** - Industry-standard TOTP algorithm with proper time-based hashing
✅ **Tiny binary** - Usually under 2MB, no external dependencies

## 🚀 Quick Installation

### Option 1: Pre-built Binary (Recommended)

1. **Download the right binary for your Mac:**

   ```bash
   # Intel Mac
   cp totp-macos-intel /usr/local/bin/totp

   # Apple Silicon Mac (M1/M2/M3)
   cp totp-macos-arm64 /usr/local/bin/totp

   # Make executable
   sudo chmod +x /usr/local/bin/totp
   ```

2. **Setup your config:**

   ```bash
   cp sample_config.json ~/.totp_config.json
   nano ~/.totp_config.json  # Add your real TOTP secrets
   ```

3. **Ready to use:**
   ```bash
   totp user_1
   # Output: 987654
   # ✅ Code automatically copied to clipboard!
   ```

### Option 2: Build from Source

1. **Install Go:**

   ```bash
   brew install go
   ```

2. **Build:**

   ```bash
   chmod +x build.sh
   ./build.sh
   ```

3. **Install:**
   ```bash
   sudo cp totp-macos-arm64 /usr/local/bin/totp  # or totp-macos-intel
   sudo chmod +x /usr/local/bin/totp
   ```

## 💻 Usage Examples

### Basic Usage (Default Mode)

```bash
totp user_1
# Output: 987654
# ✅ Also copies "987654" to clipboard automatically!
# Just Command+V to paste anywhere
```

### Quiet Mode (Silent Clipboard Copy)

```bash
totp github --quiet
# No terminal output
# ✅ Code silently copied to clipboard
# Perfect for scripts and clean workflows
```

### Print-Only Mode (No Clipboard)

```bash
totp aws --no-copy
# Output: 123456
# ❌ Does NOT copy to clipboard
# Useful when you want to see the code but not copy it
```

### Case Insensitive Examples

```bash
totp USER_1              # Works perfectly
totp GitHub --quiet      # Silent copy for GitHub
totp AWS_PROD           # Production AWS account
totp work-vpn --no-copy  # VPN code without copying
```

## 📋 Clipboard Magic

The tool automatically detects your operating system and uses the right clipboard command:

- **macOS**: `pbcopy` (built-in) ✅
- **Linux**: `xclip` or `xsel` (install via package manager)
- **Windows**: `clip` (built-in) ✅

If clipboard copy fails, you'll get a warning but the program continues normally.

## ⚡ Perfect Workflows

### Super Fast Login Flow

```bash
# 1. Generate and copy TOTP silently
totp github --quiet

# 2. Switch to browser/app and paste
# Command+V (macOS) or Ctrl+V (Windows/Linux)

# Done! No manual copying needed.
```

### Visual Confirmation Flow

```bash
# 1. See the code AND copy to clipboard
totp work_vpn
# Output: 456789

# 2. Verify the code visually, then paste
# Command+V to paste the same code
```

### Script Integration

```bash
#!/bin/bash
# Generate TOTP and verify it's copied
totp production_server --quiet
echo "TOTP copied to clipboard. Paste with Cmd+V"
```

## 📁 Configuration

### Config File Location

```
~/.totp_config.json
```

### Config File Format

```json
{
  "user_1": "JBSWY3DPEHPK3PXP",
  "user_2": "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ",
  "github": "YOUR_GITHUB_TOTP_SECRET",
  "aws_prod": "YOUR_AWS_PRODUCTION_SECRET",
  "aws_dev": "YOUR_AWS_DEVELOPMENT_SECRET",
  "work_vpn": "YOUR_VPN_TOTP_SECRET",
  "personal_bank": "YOUR_BANK_TOTP_SECRET"
}
```

### Real-World Config Example

```json
{
  "github_personal": "KBSWY3DPEHPK3PXP",
  "github_work": "NBSWY3DPEHPK3PXP",
  "aws_production": "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ",
  "aws_staging": "MXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ",
  "google_workspace": "LBSWY3DPEHPK3PXP",
  "office365": "PBSWY3DPEHPK3PXP",
  "vpn": "QBSWY3DPEHPK3PXP",
  "chase_bank": "RBSWY3DPEHPK3PXP",
  "coinbase": "SBSWY3DPEHPK3PXP"
}
```

## 🔒 Security Best Practices

### File Permissions

```bash
# Restrict access to your config file
chmod 600 ~/.totp_config.json

# Verify permissions
ls -la ~/.totp_config.json
# Should show: -rw------- (only you can read/write)
```

### Security Notes

- ✅ Binary contains **no secrets** - all secrets stored in config file
- ✅ Uses industry-standard HMAC-SHA1 TOTP algorithm (RFC 6238)
- ✅ Secrets never logged or cached
- ✅ Clipboard is managed by OS (auto-clears after timeout)
- ✅ Config file is local only (never transmitted)

## 🛠 Advanced Options

### Help Command

```bash
totp --help
# Shows all available options and examples
```

### All Available Flags

```bash
totp <user_id>              # Default: print + copy to clipboard
totp <user_id> --quiet      # Only copy to clipboard (silent)
totp <user_id> --no-copy    # Only print to terminal
totp --help                 # Show help message
```

### Error Handling

```bash
# User not found
totp nonexistent_user
# Error: User 'nonexistent_user' not found in config
# Available users: github, aws_prod, vpn

# Invalid secret
# Error: invalid base32 secret
# Make sure the secret is a valid base32 string

# Clipboard unavailable
totp user_1
# Output: 123456
# Warning: Could not copy to clipboard: xclip not found
```

## 📥 Getting TOTP Secrets

### Common Sources

1. **QR Codes**: Scan with any QR reader to extract the secret
2. **Setup URLs**: Look for `otpauth://` URLs and extract the `secret` parameter
3. **Base32 Strings**: Directly provided as strings like "JBSWY3DPEHPK3PXP"
4. **App Settings**: Many apps show the secret in account settings

### QR Code Example

```
otpauth://totp/GitHub:username?secret=JBSWY3DPEHPK3PXP&issuer=GitHub
                                      ^^^^^^^^^^^^^^^^
                                      This is your secret
```

### Manual Entry

When setting up 2FA, most services offer both QR code and manual entry options. Choose manual entry to get the base32 secret directly.

## 🚀 Why This Tool Rocks

### Before (Manual Process)

```bash
# Old way with Google Authenticator on phone:
# 1. Find phone
# 2. Open authenticator app
# 3. Find the right account
# 4. Read 6-digit code
# 5. Type it manually (hope you don't make mistakes)
# 6. Code expires while typing...
```

### After (This Tool)

```bash
# New way:
totp github --quiet    # ✅ Code instantly in clipboard
# Command+V to paste   # ✅ Zero typing errors
# Takes 2 seconds total
```

### Performance Comparison

- **Phone apps**: 10-30 seconds (find phone, unlock, navigate, read, type)
- **This tool**: 1-2 seconds (type command, paste)
- **Error rate**: Near zero vs. common typos with manual entry

## 🔧 Troubleshooting

### Common Issues

**1. "command not found: totp"**

```bash
# Make sure it's in your PATH
which totp
# If nothing, reinstall to /usr/local/bin/
```

**2. "Config file not found"**

```bash
# Create the config file
cp sample_config.json ~/.totp_config.json
nano ~/.totp_config.json
```

**3. "Could not copy to clipboard"**

```bash
# macOS: Should work out of the box
# Linux: Install clipboard utility
sudo apt install xclip     # Ubuntu/Debian
sudo yum install xclip     # CentOS/RHEL
```

**4. "Invalid base32 secret"**

```bash
# Check your secret format:
# ✅ Correct: "JBSWY3DPEHPK3PXP"
# ❌ Wrong:   "jbswy3dpehpk3pxp" (lowercase)
# ❌ Wrong:   "JBSWY3DP-EHPK-3PXP" (dashes)
```

### Verify Installation

```bash
# Test the tool
totp --help

# Test config file
cat ~/.totp_config.json

# Test permissions
ls -la ~/.totp_config.json

# Test clipboard (macOS)
totp test_user --quiet && pbpaste
```

## 📦 What's Included

- **`main_enhanced.go`** - Full-featured source code
- **`main.go`** - Basic version source code
- **`build.sh`** - Cross-platform build script
- **`sample_config.json`** - Example configuration
- **`go.mod`** - Go module definition
- **Pre-built binaries** for Intel and Apple Silicon Macs

## 🎉 Pro Tips

1. **Use descriptive names in config:**

   ```json
   {
     "github_work": "...",
     "github_personal": "...",
     "aws_production": "...",
     "aws_development": "..."
   }
   ```

2. **Create shell aliases for common accounts:**

   ```bash
   # Add to ~/.zshrc or ~/.bashrc
   alias gh-work='totp github_work --quiet'
   alias aws-prod='totp aws_production --quiet'
   ```

3. **Integrate with password managers:**

   ```bash
   # Generate TOTP and then open password manager
   totp important_account && open -a "1Password"
   ```

4. **Quick verification:**
   ```bash
   # See and copy in one command
   totp bank_account && echo "✅ Code copied to clipboard"
   ```

---

**Built with ❤️ in Go | Zero dependencies | Maximum security | Instant clipboard**
