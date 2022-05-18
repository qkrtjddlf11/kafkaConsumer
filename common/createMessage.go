package common

import "fmt"

// CreateMessage creates message in Measurement
func CreateMessage(typeOf, level, value, path string) (message string) {
	// message := fmt.Sprintf("[%s] MEM => INFO: %s | USED PERCENT : %.2f%%", level, "Common", memUsedPercent)
	switch typeOf {
	case "mem":
		message = fmt.Sprintf("[%s] MEM => INFO: %s | USED PERCENT : %s", level, "COMMON", value)

	case "cpu":
		message = fmt.Sprintf("[%s] CPU => INFO: %s | USED PERCENT : %s", level, "COMMON", value)

	case "disk":
		message = fmt.Sprintf("[%s] DISK => INFO: %s | PATH : %s | USED PERCENT : %s", level, "COMMON", path, value)

	case "swap":
		message = fmt.Sprintf("[%s] SWAP => INFO: %s | USED PERCENT : %s", level, "COMMON", value)

	case "load5":
		message = fmt.Sprintf("[%s] LOAD5 => INFO: %s | USED PERCENT : %s", level, "COMMON", value)
	}

	return message
}
