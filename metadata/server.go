package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/waste3d/Hikari-Anime/metadata/proto"
	"github.com/waste3d/Hikari-Anime/metadata/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TMDbPopularResponse struct {
	Page       int         `json:"page"`
	Results    []TMDbMovie `json:"results"`
	TotalPages int         `json:"total_pages"`
}

type TMDbMovie struct {
	ID            int64   `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	PosterPath    string  `json:"poster_path"`
	ReleaseDate   string  `json:"release_date"`
	VoteAverage   float64 `json:"vote_average"`
}

type TMDbTVShowSearchResponse struct {
	Page       int          `json:"page"`
	Results    []TMDbTVShow `json:"results"`
	TotalPages int          `json:"total_pages"`
}

type TMDbTVShow struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	OriginalName string  `json:"original_name"`
	Overview     string  `json:"overview"`
	PosterPath   string  `json:"poster_path"`
	FirstAirDate string  `json:"first_air_date"`
	VoteAverage  float64 `json:"vote_average"`
}

const (
	tmdbAPIKey  = "59a65ec73d0fbe5eca1f931db3031d3f"
	tmdbBaseURL = "https://api.themoviedb.org/3"
)

type Server struct {
	pb.UnimplementedMetadataServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetPopularMovies(ctx context.Context, req *pb.GetPopularMoviesRequest) (*pb.GetPopularMoviesResponse, error) {
	url := fmt.Sprintf("%s/movie/popular?api_key=%s&language=%s&page=%d",
		tmdbBaseURL, tmdbAPIKey, req.GetLanguage(), req.GetPage())
	log.Printf("Выполняю запрос к TMDb по URL: %s", url)

	resp, err := utils.GetRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tmdbResponse TMDbPopularResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResponse); err != nil {
		log.Printf("ОШИБКА при декодировании JSON: %v", err)
		return nil, fmt.Errorf("ошибка при декодировании JSON: %w", err)
	}

	log.Printf("Получено %d фильмов от TMDb", len(tmdbResponse.Results))

	var movies []*pb.Movie
	for _, movie := range tmdbResponse.Results {
		movies = append(movies, &pb.Movie{
			Id:            movie.ID,
			Title:         movie.Title,
			OriginalTitle: movie.OriginalTitle,
			PosterPath:    "https://image.tmdb.org/t/p/w500" + movie.PosterPath,
			Overview:      movie.Overview,
			ReleaseDate:   movie.ReleaseDate,
			VoteAverage:   movie.VoteAverage,
		})
	}

	response := &pb.GetPopularMoviesResponse{
		Results:    movies,
		Page:       int32(tmdbResponse.Page),
		TotalPages: int32(tmdbResponse.TotalPages),
	}

	return response, nil
}

func (s *Server) SearchMovies(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	query := req.GetQuery()
	if query == "" {
		return nil, status.Errorf(codes.InvalidArgument, "поисковый запрос (query) не может быть пустым")
	}

	url := fmt.Sprintf("%s/search/movie?api_key=%s&language=%s&query=%s&page=%d",
		tmdbBaseURL, tmdbAPIKey, req.GetLanguage(), query, req.GetPage())
	log.Printf("Выполняю запрос к TMDb по URL: %s", url)

	resp, err := utils.GetRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tmdbResponse TMDbPopularResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResponse); err != nil {
		log.Printf("ОШИБКА при декодировании JSON: %v", err)
		return nil, fmt.Errorf("ошибка при декодировании JSON: %w", err)
	}

	log.Printf("Найдено %d фильмов от TMDb", len(tmdbResponse.Results))

	var movies []*pb.Movie
	for _, movie := range tmdbResponse.Results {
		movies = append(movies, &pb.Movie{
			Id:            movie.ID,
			Title:         movie.Title,
			OriginalTitle: movie.OriginalTitle,
			PosterPath:    "https://image.tmdb.org/t/p/w500" + movie.PosterPath,
			Overview:      movie.Overview,
			ReleaseDate:   movie.ReleaseDate,
			VoteAverage:   movie.VoteAverage,
		})
	}

	response := &pb.SearchResponse{
		Results:    movies,
		Page:       int32(tmdbResponse.Page),
		TotalPages: int32(tmdbResponse.TotalPages),
	}

	return response, nil
}

func (s *Server) GetMovieByID(ctx context.Context, req *pb.GetMovieByIDRequest) (*pb.Movie, error) {
	movieID := req.GetMovieId()
	if movieID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "ID фильма (movie_id) не может быть равен 0")
	}

	url := fmt.Sprintf("%s/movie/%d?api_key=%s&language=%s",
		tmdbBaseURL, movieID, tmdbAPIKey, req.GetLanguage())
	log.Printf("Выполняю запрос к TMDb по URL: %s", url)

	resp, err := utils.GetRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tmdbResponse TMDbMovie
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResponse); err != nil {
		log.Printf("ОШИБКА при декодировании JSON: %v", err)
		return nil, fmt.Errorf("ошибка при декодировании ответа от TMDb: %w", err)
	}

	log.Printf("Получен фильм от TMDb: %s", tmdbResponse.Title)
	return &pb.Movie{
		Id:            tmdbResponse.ID,
		Title:         tmdbResponse.Title,
		OriginalTitle: tmdbResponse.OriginalTitle,
		PosterPath:    "https://image.tmdb.org/t/p/w500" + tmdbResponse.PosterPath,
		Overview:      tmdbResponse.Overview,
		VoteAverage:   tmdbResponse.VoteAverage,
		ReleaseDate:   tmdbResponse.ReleaseDate,
	}, nil
}

func (s *Server) SearchTVShows(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	query := req.GetQuery()
	page := req.GetPage()
	language := req.GetLanguage()
	if query == "" || page == 0 || language == "" {
		return nil, status.Errorf(codes.InvalidArgument, "поисковый запрос (query) не может быть пустым")
	}

	url := fmt.Sprintf("%s/search/tv?api_key=%s&language=%s&query=%s&page=%d",
		tmdbBaseURL, tmdbAPIKey, req.GetLanguage(), query, req.GetPage())
	log.Printf("Выполняю запрос к TMDb по URL: %s", url)

	resp, err := utils.GetRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tmdbResponse TMDbTVShowSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResponse); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании JSON: %w", err)
	}

	log.Printf("Найдено %d сериалов от TMDb", len(tmdbResponse.Results))

	var tvShows []*pb.Movie
	for _, tvShow := range tmdbResponse.Results {
		tvShows = append(tvShows, &pb.Movie{
			Id:            tvShow.ID,
			Title:         tvShow.Name,
			OriginalTitle: tvShow.OriginalName,
			PosterPath:    "https://image.tmdb.org/t/p/w500" + tvShow.PosterPath,
			Overview:      tvShow.Overview,
			ReleaseDate:   tvShow.FirstAirDate,
			VoteAverage:   tvShow.VoteAverage,
		})
	}

	response := &pb.SearchResponse{
		Results:    tvShows,
		Page:       int32(tmdbResponse.Page),
		TotalPages: int32(tmdbResponse.TotalPages),
	}

	return response, nil
}
