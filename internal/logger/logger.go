package logger

import "log"

func LogInfo(value string) {
	log.Printf("[INFO]:  %s", value)
}

func LogDebug(value string) {
	log.Printf("[DEBUG]: %s", value)
}

func LogError(value string) {
	log.Printf("[ERROR]: %s", value)
}

func LogWarning(value string) {
	log.Printf("[WARN]:  %s", value)
}
