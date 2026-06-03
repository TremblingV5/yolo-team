package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func trunc(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n-2]) + ".."
	}
	return s
}

func fmtStr(v interface{}) string {
	s := fmt.Sprintf("%v", v)
	if s == "<nil>" {
		return "-"
	}
	return s
}

func fmtTime(v interface{}) string {
	s := fmt.Sprintf("%v", v)
	s = strings.TrimSuffix(s, " +0000 UTC")
	if len(s) > 19 {
		return s[:19]
	}
	return s
}
