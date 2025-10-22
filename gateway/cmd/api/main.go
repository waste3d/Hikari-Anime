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
			c.DefaultQuery("page", "1")
			c.DefaultQuery("language", "en-EN")
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
