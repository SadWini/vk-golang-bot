package bot

import (
	"log"
)

func isDebugging() bool {
	return API.DEBUG
}

func DebugPrint(format string, values ...interface{}) {
	if isDebugging() {
		log.Printf("[VKBOT-DEBUG] "+format, values...)
	}
}
