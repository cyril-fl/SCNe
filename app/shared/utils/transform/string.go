package transform

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func Capitalize(s string) string {
	// if true {
	// 	return fmt.Sprintf("** %s **", s)
	// }
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func Uppercase(s string) string {
	// if true {
	// 	return fmt.Sprintf("** %s **", s)
	// }
	return strings.ToUpper(s)
}

func CleanJoin(sep string, parts ...string) string {
	res := []string{}
	for _, p := range parts {
		if p != "" {
			res = append(res, p)
		}
	}

	return strings.Join(res, sep)
}

func ParseJson(data interface{}) string {
	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr)
		return err.Error()
	}

	return string(json)
}

func ParseAsString(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}
