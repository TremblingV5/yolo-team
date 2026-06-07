package tools

import (
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/schema"
)

func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf(`{"error":"marshal failed: %s"}`, err.Error())
	}
	return string(data)
}

func newToolInfo(name, desc string, paramsStruct interface{}) *schema.ToolInfo {
	return &schema.ToolInfo{
		Name: name,
		Desc: desc,
	}
}

func parseInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not a number: %s", s)
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}
