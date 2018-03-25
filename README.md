# antifreeze-back-end
Repository for the back-end of Freeze-B-Gone

## Development

This section contains instructions for setting up and using the development environment.

### Environment

- Postman
- Git
- Docker
- Docker-Compose (comes with Windows and Mac installations of Docker)

### Set Up

1. Install the required tools
2. Clone the repository

### Test Changes

1. Change directories to `$GOPATH/src/github.com/NilsG-S/antifreeze-back-end`
2. If necessary, run `dep ensure` to update dependency files
3. Run dev environment: `docker-compose -f dev-compose.yml up --build`
    - Has to be done every time a change is made to the server
    - First run and runs after dependency changes take way longer
4. Target `0.0.0.0:8081` when using Postman to make requests

### Deploying

1. Run `terraform refresh` to update state
2. Run `terraform apply` to deploy the new infrastructure

### Notes

- You'll slowly accumulate dangling Docker images that consume gigabytes of space.
 These can be removed by running `docker image prune`.
- Running the server as root can result in an inability to access environment variables.
