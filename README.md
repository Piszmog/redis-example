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