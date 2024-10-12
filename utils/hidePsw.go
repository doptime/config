package utils

import (
	"encoding/json"
	"regexp"
	"strings"
)

func ToHidePswdString(obj interface{}) (jsonStr string) {
	var (
		err       error
		jsonBytes []byte
	)

	jsonBytes, err = json.Marshal(obj)
	if err != nil {
		return ""
	}
	jsonStr = string(jsonBytes)
	// Replace the password, keeping the first character and masking the rest
	re := regexp.MustCompile(`("pswd[0-9]?":\s*")([^"]+)(")`)
	jsonStr = re.ReplaceAllStringFunc(jsonStr, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) == 4 {
			password := submatches[2]
			if len(password) > 0 {
				maskedPassword := string(password[0]) + strings.Repeat("*", len(password)-1)
				return submatches[1] + maskedPassword + submatches[3]
			}
		}
		return match
	})

	if jsonStr == "" {
		return "error: failed to marshal config to json string"
	}
	var c1 = map[string]interface{}{}

	err = json.Unmarshal([]byte(jsonStr), &c1)
	if err != nil {
		return "error: " + err.Error() + " when unmarshaling json string"
	}
	jsonBytes, _ = json.Marshal(c1)
	return string(jsonBytes)
}
