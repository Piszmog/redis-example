package webrouter

import (
	"encoding/json"
	"fmt"
	"github.com/Piszmog/redis-example/cache"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	KeyContentType       = "Content-Type"
	ValueApplicationJson = "application/json"
	Id                   = "id"
)

type Movie struct {
	Id          string `bson:"_id" json:"_id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

func (s Movie) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s Movie) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, s)
}

type ResponseId struct {
	Id string `json:"id"`
}

var cacheClient cache.Client

func SetupMovieRoutes(router *httprouter.Router, client cache.Client) {
	router.GET("/movies", GetAllMovies)
	router.GET("/movies/:id", FindMovie)
	router.POST("/movies", CreateMovie)
	router.PUT("/movies/:id", UpdateMovie)
	router.DELETE("/movies/:id", DeleteMovie)
	router.DELETE("/movies", DeleteMovies)
	cacheClient = client
}

func GetAllMovies(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	movies, err := cacheClient.GetAll()
	movieList := make([]Movie, len(movies))
	if err != nil {
		fmt.Printf("failure %+v", err)
		WriteResponse(writer, http.StatusInternalServerError, nil)
		return
	}
	i := 0
	for _, movie := range movies {
		var movieStruct Movie
		err := json.Unmarshal([]byte(movie), &movieStruct)
		if err != nil {
			fmt.Printf("failure %+v", err)
			WriteResponse(writer, http.StatusInternalServerError, nil)
			return
		}
		movieList[i] = movieStruct
		i++
	}
	WriteOkResponse(writer, movieList)
}

// Finds the movie matching the provided id
func FindMovie(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var movie Movie
	err := cacheClient.Get(params.ByName(Id), &movie)
	if err != nil {
		fmt.Printf("failure %+v", err)
		WriteResponse(writer, http.StatusNotFound, nil)
		return
	}
	WriteOkResponse(writer, movie)
}

func CreateMovie(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	defer request.Body.Close()
	var movie Movie
	if err := json.NewDecoder(request.Body).Decode(&movie); err != nil {
		fmt.Printf("failure %+v", err)
		WriteResponse(writer, http.StatusInternalServerError, nil)
		return
	}
	movie.Id = uuid.New().String()
	if err := cacheClient.Insert(movie.Id, movie); err != nil {
		fmt.Printf("failure %+v", err)
		WriteResponse(writer, http.StatusInternalServerError, nil)
		return
	}
	WriteOkResponse(writer, ResponseId{Id: movie.Id})
}

func UpdateMovie(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	defer request.Body.Close()
	movieId := params.ByName(Id)
	var movie Movie
	if err := json.NewDecoder(request.Body).Decode(&movie); err != nil {
		fmt.Printf("failure %+v", err)
		WriteResponse(writer, http.StatusInternalServerError, nil)
		return
	}
	movie.Id = movieId
	err := cacheClient.Insert(movieId, movie)
	if err != nil {
		fmt.Printf("failure %+v", err)
		WriteResponse(writer, http.StatusNotFound, nil)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func DeleteMovie(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	movieId := params.ByName(Id)
	err := cacheClient.Remove(movieId)
	if err != nil {
		fmt.Printf("failure %+v", err)
		WriteResponse(writer, http.StatusNotFound, nil)
		return
	}
	WriteResponse(writer, http.StatusOK, nil)
}

func DeleteMovies(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	err := cacheClient.RemoveAll()
	if err != nil {
		fmt.Printf("failure %+v", err)
		WriteResponse(writer, http.StatusNotFound, nil)
		return
	}
	WriteResponse(writer, http.StatusOK, nil)
}

func WriteOkResponse(writer http.ResponseWriter, payload interface{}) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("failure %+v", err)
	}
	WriteResponse(writer, http.StatusOK, bytes)
}

func WriteResponse(writer http.ResponseWriter, httpStatus int, bytes []byte) {
	writer.Header().Set(KeyContentType, ValueApplicationJson)
	writer.WriteHeader(httpStatus)
	if bytes != nil {
		writer.Write(bytes)
	}
}
