package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/waste3d/Hikari-Anime/metadata/proto"
	"google.golang.org/grpc"
)

const (
	gatewayPort            = ":8080"
	metadataServiceAddress = "localhost:50051"
)

func main() {
	grpcServer, err := grpc.Dial(metadataServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to metadata service: %v", err)
	}
	defer grpcServer.Close()

	metadataServiceClient := pb.NewMetadataServiceClient(grpcServer)

	router := gin.Default()

	router.GET("/api/v1/movies/popular", getPopularMoviesHandler(metadataServiceClient))
	router.GET("/api/v1/movies/search", searchMoviesHandler(metadataServiceClient))
	router.GET("/api/v1/movies/:id", movieByIDHandler(metadataServiceClient))
	router.GET("/api/v1/tv/search", searchTVShowsHandler(metadataServiceClient))

	err = router.Run(gatewayPort)
	if err != nil {
		log.Fatalf("failed to start gateway: %v", err)
	}
}

func getPopularMoviesHandler(metadataServiceClient pb.MetadataServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		page := c.Query("page")
		language := c.Query("language")

		if page == "" || language == "" {
			page = c.DefaultQuery("page", "1")
			language = c.DefaultQuery("language", "ru-RU")
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := metadataServiceClient.GetPopularMovies(ctx, &pb.GetPopularMoviesRequest{
			Page:     int32(pageInt),
			Language: language,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch popular movies"})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func searchMoviesHandler(metadataServiceClient pb.MetadataServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		query := c.Query("query")
		page := c.Query("page")
		language := c.Query("language")

		if page == "" || language == "" {
			page = c.DefaultQuery("page", "1")
			language = c.DefaultQuery("language", "ru-RU")
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := metadataServiceClient.SearchMovies(ctx, &pb.SearchRequest{
			Query:    query,
			Page:     int32(pageInt),
			Language: language,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search movies"})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func movieByIDHandler(client pb.MetadataServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		language := c.DefaultQuery("language", "ru-RU")

		idInt, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie ID parameter"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := client.GetMovieByID(ctx, &pb.GetMovieByIDRequest{
			MovieId:  idInt,
			Language: language,
		})
		if err != nil {
			log.Printf("ошибка при вызове GetMovieByID: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get movie by ID"})
			return
		}
		c.JSON(http.StatusOK, response)
	}
}

func searchTVShowsHandler(metadataServiceClient pb.MetadataServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		query := c.Query("query")
		page := c.Query("page")
		language := c.Query("language")

		if page == "" || language == "" {
			page = c.DefaultQuery("page", "1")
			language = c.DefaultQuery("language", "ru-RU")
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := metadataServiceClient.SearchTVShows(ctx, &pb.SearchRequest{
			Query:    query,
			Page:     int32(pageInt),
			Language: language,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search TV shows"})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
