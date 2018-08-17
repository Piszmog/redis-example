package main

import (
	"github.com/Piszmog/cfservices"
	"github.com/Piszmog/redis-example/cache"
	"github.com/Piszmog/redis-example/webrouter"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

const (
	Key          = "movie"
	CloudService = "rediscloud"
)

func main() {
	environment := cfservices.LoadFromEnvironment()
	var cacheClient *cache.RedisClient
	if len(environment) == 0 {
		cacheClient = cache.CreateLocalRedisClient(Key)
	} else {
		credentials, err := cfservices.GetServiceCredentials(CloudService, environment)
		if err != nil {
			log.Fatal(err)
		}
		credential := credentials.Credentials[0]
		cacheClient = cache.CreateRedisClient(credential.Hostname, credential.Port, credential.Password, Key)
	}
	defer cacheClient.Close()
	router := httprouter.New()
	webrouter.SetupMovieRoutes(router, cacheClient)
	log.Fatal(http.ListenAndServe(":"+"8080", router))
}
