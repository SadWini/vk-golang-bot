package utils

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// HasPrefix tests case insensitive whether the string s begins with prefix.
func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.ToLower(s[0:len(prefix)]) == strings.ToLower(prefix)
}

// TrimPrefix returns s without the provided leading case insensitive prefix string.
// If s doesn't start with prefix, s is returned unchanged.
func TrimPrefix(s, prefix string) string {
	if HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

func ReadJSON(fn string, v interface{}) {
	file, _ := os.Open(fn)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(v)
	if err != nil {
		log.Println("error:", err)
	}
}
func ParseArg(str, arg1, arg2 string) string {
	return str[strings.Index(str, arg1)+len(arg1) : strings.Index(str, arg2)-1]
}
func ReadText(fn string) string {
	content, err := os.ReadFile(fn)
	if err != nil {
		log.Println("error:", err)
	}
	return string(content)
}
