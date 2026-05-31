.PHONY: gen swag

gen:
	gorm gen -i internal/model -o internal/query

swag:
	swag init -g cmd/yolo/main.go -o docs --parseDependency --parseInternal
