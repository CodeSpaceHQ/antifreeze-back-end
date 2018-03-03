# antifreeze-back-end
Repository for the back-end of Freeze-B-Gone

## Development

This section contains instructions for setting up and using the development environment.

### Environment

- Node.js 8.9.4
- NPM 5.6.0
- Serverless framework 1.26 (global install)
  -- Type "npm install -g serverless" into your command line
  -- Check by running "serverless --version"
- Go 1.9.3
  -- Follow: https://golang.org/dl/ and test your installation
- Go Dep 0.4.1
  -- Follow: https://github.com/golang/dep/releases
- JRE 6.x or newer
- [DynamoDB Local]
  -- Follow: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html
  -- Download the .jar file for US West
  -- Extract the contents to where you want them
  -- Open cmd prompt, navigate to where you stored the file and run the following command to start it up:
      java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb
  -- Note: It uses port 8000 by default and processes incoming requests until you stop it with Ctrl+C
- Postman
- Git

### Set Up

1. Install the required tools
2. Clone the repository
3. Use Dep to install Go dependencies
4. Add a `.env` file to the top level directory with the following contents (for development purposes):

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

Instructions pending...

# WebSocket

Requires a mapping from devices to clients, i.e. `map[deviceId]*Client`.

## Connection:

- email
- isAuthed
- permissions
- subsriptions

## Message

```json
{
    "subsription": "string",
    "operation": "ADD/REMOVE/UPDATE"
}
```

## Path

### /user/devices

Make sure to update the front-end before doing server logic (to avoid situations where the front-end hasn't been setup to handle temp/alarm updates for a specific device). But then the server would have to remove the new device if something goes wrong.

- deviceId

#### ADD
#### REMOVE

### /device/alarm

- deviceId

#### UPDATE

Allow for `nil` values

- temp

### /device/history

- deviceId

#### ADD

- temp
- time

# REST

## Endpoint

Many of these should send updates to certain WebSockets

### /user

For creating users

#### POST

### /user/devices

For managing a user's devices.

#### GET

Get all devices (for setup).

#### POST

- Adding a device to a user.
- Requires a WebSocket push to `/user/device` subscribers.

#### DELETE

- Delete a device.
- Requires a WebSocket push to `/user/device` subscribers.

### /device

#### POST

- For adding devices to the database.

#### PUT

- For updating the name in the database.

```json
{
    "name": "string"
}
```

### /device/alarm

#### PUT

- Allow for `nil` values.
- Requires a WebSocket push to `/device/alarm` subscribers.

### /device/history

#### GET

- deviceId

This endpoint gets all history for a single device

#### POST

- Requires a WebSocket push to `/device/history` subscribers.
