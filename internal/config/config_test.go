package config

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	appConfig := LoadConfig()
	golang := appConfig.Image["golang"]
	fmt.Printf(golang)
}
