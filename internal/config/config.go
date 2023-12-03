package config

import "os"

// tokenPath is path of token
const tokenPath string = "configs/token.txt"

var token string

// GetBotToken returns loaded token.
//
// Note: It is better to call 'runtime.GC()' after using token to clear
// this variable from memory.
func GetBotToken() string { return token }

// LoadConfig loads all configurations and tokens for running bot service.
func LoadConfig() error {
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return err
	}
	token = string(data)
	return nil
}
