package configs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"iDevopzAgent/models"
	"iDevopzAgent/security"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/denisbrodbeck/machineid"
)

func getConfigPath() string {
	var baseDir string
	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("APPDATA")
	default:
		baseDir = "/"
	}
	return filepath.Join(baseDir, "metrics-agent", "config.json")
}

func LoadUserID() (string, string, error) {
	path := getConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
	}

	var cfg models.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", "", err
	}
	fmt.Printf("Raw JSON: %s\n", string(data))
	fmt.Printf("Parsed struct: %+v\n", cfg)
	// Decrypt both UserID and MachineID before returning
	decUserID, err := security.Decrypt(cfg.UserID)
	if err != nil {
		return "", "", fmt.Errorf("failed to decrypt UserID: %w", err)
	}

	decMachineID, err := security.Decrypt(cfg.MachineID)
	if err != nil {
		return "", "", fmt.Errorf("failed to decrypt MachineID: %w", err)
	}

	return decUserID, decMachineID, nil
}

func PromptAndSaveUserID() (string, string) {
	configPath := getConfigPath()

	// 1. If config already exists, load it
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err == nil {
			var cfg models.Config
			if json.Unmarshal(data, &cfg) == nil && cfg.UserID != "" && cfg.MachineID != "" {
				//  Decrypt stored values
				decryptedUserID, _ := security.Decrypt(cfg.UserID)
				decryptedMachineID, _ := security.Decrypt(cfg.MachineID)

				fmt.Println("✔ Using stored UserID:", decryptedUserID)
				fmt.Println("✔ Using stored MachineID:", decryptedMachineID)

				return cfg.UserID, cfg.MachineID // still return encrypted values
			}
		}
	}

	// 2. Ask user for UserID
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your User ID (or Device Key): ")
	userID, _ := reader.ReadString('\n')
	userID = string(bytes.TrimSpace([]byte(userID)))

	// 3. Get MachineID
	machineID, err := machineid.ID()
	if err != nil {
		log.Fatalf("Failed to get machine ID: %v", err)
	}

	// 4. Encrypt both
	encryptedUserID, err := security.Encrypt(userID)
	if err != nil {
		log.Fatalf("Failed to encrypt UserID: %v", err)
	}

	encryptedMachineID, err := security.Encrypt(machineID)
	if err != nil {
		log.Fatalf("Failed to encrypt MachineID: %v", err)
	}

	// 5. Save config
	config := models.Config{
		UserID:    encryptedUserID,
		MachineID: encryptedMachineID,
	}
	data, _ := json.MarshalIndent(config, "", "  ")

	os.MkdirAll(filepath.Dir(configPath), 0700)
	_ = os.WriteFile(configPath, data, 0600)

	fmt.Println("✔ User ID stored (encrypted):", encryptedUserID)
	fmt.Println("✔ Machine ID stored (encrypted):", encryptedMachineID)

	return encryptedUserID, encryptedMachineID
}

func SaveUserID(userID string) error {
	configPath := getConfigPath()
	
	// Get MachineID
	machineID, err := machineid.ID()
	if err != nil {
		return fmt.Errorf("failed to get machine ID: %w", err)
	}
	
	// Encrypt both
	encryptedUserID, err := security.Encrypt(userID)
	if err != nil {
		return fmt.Errorf("failed to encrypt UserID: %w", err)
	}
	
	encryptedMachineID, err := security.Encrypt(machineID)
	if err != nil {
		return fmt.Errorf("failed to encrypt MachineID: %w", err)
	}
	
	// Save config
	config := models.Config{
		UserID:    encryptedUserID,
		MachineID: encryptedMachineID,
	}
	data, _ := json.MarshalIndent(config, "", "  ")
	
	os.MkdirAll(filepath.Dir(configPath), 0700)
	return os.WriteFile(configPath, data, 0600)
}
