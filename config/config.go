package config

import "os"

var DebugMode = false

func init() {
	if os.Getenv("DEBUG") != "" {
		DebugMode = true
	}
}
