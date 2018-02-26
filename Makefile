build:
	dep ensure
	# env GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go
