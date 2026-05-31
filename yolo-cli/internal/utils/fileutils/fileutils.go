package fileutils

import (
	"fmt"
	"os"
)

func Write(path, content string) {
	os.WriteFile(path, []byte(content), 0644)
}

func Read(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func BuildDocPath(workspace string, projectID int64, title string) string {
	dir := fmt.Sprintf("%s/documents/%d", workspace, projectID)
	os.MkdirAll(dir, 0755)
	return dir + "/" + title + ".md"
}
