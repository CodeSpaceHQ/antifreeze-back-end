build:
	# dep ensure
	env GOOS=linux go build common/*.go
	env GOOS=linux go build ws/mux/*.go
	env GOOS=linux go build rest/*.go
	env GOOS=linux go build ws/*.go
	env GOOS=linux go build -o bin/antifreeze-back-end main.go
