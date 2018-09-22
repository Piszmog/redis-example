# Example of Redis with Go
Example of using Redis with a Go application that exposes REST endpoints to interact with the cache.

## Client
The base client library being used is `github.com/go-redis/redis`. This library is wrapped in `cache/cache.go`.

### Persistence
All persistence being used is hash sets.

## Local
If `VCAP_SERVICES` is not an environment variable, then the application assumes a locally Redis is running and will attempt 
to connect to it via the default configurations.

## Cloud Foundry
When deployed to Cloud Foundry and a Redis service is bounded to the application, `VCAP_SERVICES` will be set as an 
environment variable. This variable contains the credential information for connecting to the Redis instance.

### Deploying a Linux Binary
To deploy the app as a linux binary instead of the Go source code, first build the code to linux

* `go build`
* `GOOS=linux GOARCH=amd64 go build` -- if building from a Windows machine

Update the manifest to the following,
```yaml
applications:
- name: piszmog-redis-demo
  buildpack: binary_buildpack
  memory: 5M
  instances: 1
  target: redis-example
  services:
  - redis
```

Then push the application with the command `cf push -c './redis-example'`

#### Start command
There are two ways to start a binary.

1. `-c './${binary name}'`
2. Procfile at the root of the application with `web: ./${binary name}` in the file