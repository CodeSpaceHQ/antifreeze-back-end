# antifreeze-back-end
Repository for the back-end of Freeze-B-Gone

## Development

This section contains instructions for setting up and using the development environment.

### Environment

- Postman
- Git
- Docker

### Set Up

1. Install the required tools
2. Clone the repository
3. Add a `.env` file to the top level directory with the following contents (for development purposes):

```
PORT=3000
SECRET=thisisoursecret
SALT=10

PROD_DB_URI=pending
DEV_DB_URI=pending
TEST_DB_URI=pending

USERNAME=test@ttu.edu
PASSWORD=test
```

### Test Changes

#### Docker Setup

1. Setup Cloud Datastore emulator `docker build -f emulator.Dockerfile -t emulator .` (only has to be done once)
2. Setup server container `docker build -f dev.Dockerfile -t antifreeze .` (has to be done every time a change is made to the server)
3. Run `emulator` `docker run -dit emulator`
4. Run `antifreeze` `docker run -dit antifreeze`

Note: these are primarily for users of the back-end. Developer instructions are pending.
