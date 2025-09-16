package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Config represents the TOTP configuration
type Config map[string]string

// loadConfig loads the TOTP secrets from the config file
func loadConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get home directory: %v", err)
	}

	configPath := filepath.Join(homeDir, ".totp_config.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s\nCreate a JSON file with format: {\"user_1\": \"totp_secret_1\", \"user_2\": \"totp_secret_2\"}", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid JSON in config file: %v", err)
	}

	return config, nil
}

// generateTOTP generates a TOTP code from a base32 secret
func generateTOTP(secret string) (string, error) {
	// Remove any whitespace and convert to uppercase
	secret = strings.ToUpper(strings.ReplaceAll(secret, " ", ""))

	// Decode base32 secret
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("invalid base32 secret: %v", err)
	}

	// Get current time step (30-second intervals)
	timeStep := time.Now().Unix() / 30

	// Convert time step to bytes
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(timeStep))

	// Create HMAC-SHA1 hash
	h := hmac.New(sha1.New, key)
	h.Write(timeBytes)
	hash := h.Sum(nil)

	// Dynamic truncation
	offset := hash[len(hash)-1] & 0x0F
	truncatedHash := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF

	// Generate 6-digit code
	code := truncatedHash % 1000000

	return fmt.Sprintf("%06d", code), nil
}

// copyToClipboard copies text to the system clipboard
func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("pbcopy")
	case "linux":
		// Try xclip first, then xsel
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("no clipboard utility found (install xclip or xsel)")
		}
	case "windows":
		cmd = exec.Command("cmd", "/c", "clip")
	default:
		return fmt.Errorf("clipboard not supported on %s", runtime.GOOS)
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// createCaseInsensitiveMap creates a map with lowercase keys for case-insensitive lookup
func createCaseInsensitiveMap(config Config) map[string]string {
	caseInsensitiveMap := make(map[string]string)
	for key, value := range config {
		caseInsensitiveMap[strings.ToLower(key)] = value
	}
	return caseInsensitiveMap
}

// printUsage prints the usage information
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <user_id> [options]\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	fmt.Fprintf(os.Stderr, "  --no-copy    Don't copy to clipboard\n")
	fmt.Fprintf(os.Stderr, "  --quiet      Only copy to clipboard, don't print to stdout\n")
	fmt.Fprintf(os.Stderr, "  --help       Show this help message\n")
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s user_1              # Print code and copy to clipboard\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "  %s user_1 --quiet      # Only copy to clipboard (silent)\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "  %s user_1 --no-copy    # Only print, don't copy\n", filepath.Base(os.Args[0]))
}

func main() {
	// Parse arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Check for help flag
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage()
		os.Exit(0)
	}

	userID := strings.ToLower(os.Args[1]) // Case insensitive
	var copyToClip = true
	var quietMode = false

	// Parse flags
	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--no-copy":
			copyToClip = false
		case "--quiet":
			quietMode = true
		default:
			fmt.Fprintf(os.Stderr, "⚠️ Unknown option: %s\n", os.Args[i])
			printUsage()
			os.Exit(1)
		}
	}

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️ Error: %v\n", err)
		os.Exit(1)
	}

	// Create case-insensitive lookup
	caseInsensitiveConfig := createCaseInsensitiveMap(config)

	// Find the secret for the user
	secret, exists := caseInsensitiveConfig[userID]
	if !exists {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️ Error: %v\n", err)
			os.Exit(1)
		}
		configPath := filepath.Join(homeDir, ".totp_config.json")
		fmt.Fprintf(os.Stderr, "⚠️ User '%s' not found in config file at %s\n", os.Args[1], configPath)

		// Show available users
		var users []string
		for key := range config {
			users = append(users, key)
		}
		if len(users) > 0 {
			fmt.Fprintf(os.Stderr, "⚠️ Available users: %s\n", strings.Join(users, ", "))
		}
		os.Exit(1)
	}

	// Generate TOTP code
	code, err := generateTOTP(secret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️ Error generating TOTP: %v\n", err)
		fmt.Fprintf(os.Stderr, "⚠️ Make sure the secret is a valid base32 string\n")
		os.Exit(1)
	}

	// Copy to clipboard (unless disabled)
	if copyToClip {
		if err := copyToClipboard(code); err != nil {
			// Don't fail the program if clipboard copy fails, just warn
			if !quietMode {
				fmt.Fprintf(os.Stderr, "⚠️ Warning: Could not copy to clipboard: %v\n", err)
			}
		}
	}

	// Output the code (unless in quiet mode)
	if !quietMode {
		fmt.Println("👤 User		: ", userID)
		fmt.Println("🔑 TOTP Code	: ", code)
		fmt.Println("📋 Copied to clipboard")
	}

}
