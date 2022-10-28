# go-follow-microservice

microservice in golang to handle follow/unfollow social functionality for games etc. 

windows and bash exectuable included

run `golang-follow-microservice.exe` to start a server on `localhost:9090`

interact via `curl` requests

The following functionality is provided:
- Get list of users:      `curl localhost:9090`
- Get user by integer ID: `curl localhost:9090/[ID]` 
- Create new user:        `curl -X POST -d "{\"id\": [ID]}" localhost:9090`
- Follow:                 `curl -X PATCH localhost:9090/follow/[ID_1]/[ID_2]`
- Unfollow:               `curl -X PATCH localhost:9090/unfollow/[ID_1]/[ID_2]`

GET requests can be piped into `jq` for extra readability
