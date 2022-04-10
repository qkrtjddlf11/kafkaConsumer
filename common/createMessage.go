package common

import "fmt"

// CreateMessage creates message in Measurement
func CreateMessage(typeOf, level, value string) (message string) {
	// message := fmt.Sprintf("[%s] MEM => INFO: %s | USED PERCENT : %.2f%%", level, "Common", memUsedPercent)
	switch typeOf {
	case "mem":
		message = fmt.Sprintf("[%s] MEM => INFO: %s | USED PERCENT : %s", level, "COMMON", value)
	}

	return message
}
