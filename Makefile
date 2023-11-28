up:
	docker compose up --build
subscriber:
	go run cmd\subscriber\main.go
publisher:
	go run cmd\publisher\main.go