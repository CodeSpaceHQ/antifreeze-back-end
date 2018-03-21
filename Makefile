build:
	env GOOS=linux go build -i common/*.go
	env GOOS=linux go build -i common/db/*.go
	env GOOS=linux go build -i ws/mux/*.go
	env GOOS=linux go build -i rest/*.go
	env GOOS=linux go build -i ws/*.go
	env GOOS=linux go build -o bin/antifreeze-back-end main.go
