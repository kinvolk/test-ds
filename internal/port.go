package internal

import (
	"fmt"
)

func ValidatePort(port int, desc string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("Invalid %s port %d", desc, port)
	}
	return nil
}
