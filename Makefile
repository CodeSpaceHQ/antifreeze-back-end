build:
	env GOOS=linux go install github.com/NilsG-S/antifreeze-back-end/common/env
	env GOOS=linux go install github.com/NilsG-S/antifreeze-back-end/common/user
	env GOOS=linux go install github.com/NilsG-S/antifreeze-back-end/common
	env GOOS=linux go install github.com/NilsG-S/antifreeze-back-end/ws/mux
	env GOOS=linux go install github.com/NilsG-S/antifreeze-back-end/rest/user
	env GOOS=linux go install github.com/NilsG-S/antifreeze-back-end/rest
	env GOOS=linux go install github.com/NilsG-S/antifreeze-back-end/ws
	env GOOS=linux go build -o bin/antifreeze-back-end main.go
