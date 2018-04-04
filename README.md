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
3. Set the following environment variables:

```
ANTIFREEZE_SECRET=example
```

### Test Changes

1. Change directories to `$GOPATH/src/github.com/NilsG-S/antifreeze-back-end`
2. If necessary, run `dep ensure` to update dependency files
3. Run dev environment: `docker-compose -f dev-compose.yml up --build`
    - Has to be done every time a change is made to the server
    - First run and runs after dependency changes take way longer
4. Target `0.0.0.0:8081` when using Postman to make requests

Note:
- The error "The system cannot find the file specified." can indicate that docker hasn't been started.

### Deploying

1. Change directories to `${project_path}/deploy/terraform`
2. Create `secret.tfvars` with the following contents:

```
master_username = "example"
master_password = "example"
```

3. Run `terraform init` to get required resources
4. Run `terraform apply --var-file="secret.tfvars"` to deploy the new infrastructure

Notes:

- Running terraform in a script requires `-auto-approve` flag
- When making major infrastructure changes, it's better to run `terraform destory` before applying the new plan.
 Otherwise resources have a way of being orphaned.
 Basically, destroying resources as a result of updates.
- The Terraform config for this project requires the following APIs be enabled:
    1. Identity and Access Management API
    2. Cloud Resource Manager API
- Make sure to keep docker file version tag updated, both in Terraform and Travis.
 This is the only way to deploy new container versions to the cloud.
 The files are `deploy.sh` and `deploy/terraform/prod.tf`.

### Notes

- You'll slowly accumulate dangling Docker images that consume gigabytes of space.
 These can be removed by running `docker image prune`.
- Running the server as root can result in an inability to access environment variables.
