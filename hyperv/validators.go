package hyperv

import (
	"fmt"
	"strings"
)

func validateMacAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if strings.Contains(value, "-") || strings.Contains(value, ":") {
		errors = append(errors, fmt.Errorf("MAC address cannot contain special characters. Example value: '7824AF34D9B9'"))
	}
	return
}
