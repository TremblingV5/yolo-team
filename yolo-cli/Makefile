.PHONY: gen swag build build-frontend build-backend

gen:
	gorm gen -i internal/model -o internal/query

swag:
	swag init -g cmd/yolo/main.go -o docs --parseDependency --parseInternal

build-frontend:
	cd ../yolo-client && npm run build
	touch internal/frontend/client-dist/.gitkeep

build-backend:
	go build -o yolo.exe ./cmd/yolo

build: build-frontend build-backend
