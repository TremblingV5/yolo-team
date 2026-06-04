package keyutils

import "fmt"

func Project(id int64) string {
	return fmt.Sprintf("YOLO-PROJECT-%d", id)
}

func Issue(id int64) string {
	return fmt.Sprintf("YOLO-ISSUE-%d", id)
}

func Document(id int64) string {
	return fmt.Sprintf("YOLO-DOC-%d", id)
}

func Task(id int64) string {
	return fmt.Sprintf("YOLO-TASK-%d", id)
}
