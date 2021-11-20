package config

import "os"

var DebugMode = false
var DynamoTableName = ""

func init() {
	if os.Getenv("DEBUG") != "" {
		DebugMode = true
	}
	DynamoTableName = os.Getenv("ROOMS_NAME")
}
