package gui

import "fmt"

func ColorRed(value string) string {
	return fmt.Sprintf("\033[31;1m%s\033[0m", value)
}

func ColorGreen(value string) string {
	return fmt.Sprintf("\033[1;32m%s\033[0m", value)
}

func ColorMagenta(value string) string {
	return fmt.Sprintf("\033[1;35m%s\033[0m", value)
}

func ColorBlue(value string) string {
	return fmt.Sprintf("\033[1;34m%s\033[0m", value)
}

func ColorYellow(value string) string {
	return fmt.Sprintf("\033[0;33m%s\033[0m", value)
}

func ColorCyan(value string) string {
	return fmt.Sprintf("\033[0;36m%s\033[0m", value)
}

func ColorLightCyan(value string) string {
	return fmt.Sprintf("\033[1;36m%s\033[0m", value)
}
