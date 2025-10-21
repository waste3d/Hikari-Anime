package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	pb "github.com/waste3d/Hikari-Anime/metadata/proto"
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

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("ОШИБКА при выполнении HTTP-запроса к TMDb: %v", err)
		return nil, fmt.Errorf("ошибка при выполнении запроса к TMDb: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Статус ответа от TMDb: %s", resp.Status)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Получен не-OK статус от TMDb: %d", resp.StatusCode)
		return nil, fmt.Errorf("TMDb API вернул ошибку: %s", resp.Status)
	}

	// 4. Декодируем JSON-ответ
	var tmdbResponse TMDbPopularResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResponse); err != nil {
		log.Printf("ОШИБКА при декодировании JSON ответа от TMDb: %v", err)
		return nil, fmt.Errorf("ошибка при декодировании ответа от TMDb: %w", err)
	}

	// 5. ЛОГИРУЕМ КОЛИЧЕСТВО ПОЛУЧЕННЫХ ФИЛЬМОВ
	log.Printf("Получено %d фильмов от TMDb", len(tmdbResponse.Results))

	// Конвертируем данные
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

	log.Printf("Отправляю ответ клиенту. Количество фильмов: %d", len(response.Results))
	return response, nil
}

func (s *Server) SearchMovies(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	query := req.GetQuery()
	page := req.GetPage()
	language := req.GetLanguage()

	url := fmt.Sprintf("%s/search/movie?api_key=%s&language=%s&query=%s&page=%d",
		tmdbBaseURL, tmdbAPIKey, language, query, page)
	log.Printf("Выполняю запрос к TMDb по URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("ОШИБКА при выполнении HTTP-запроса к TMDb: %v", err)
		return nil, fmt.Errorf("ошибка при выполнении запроса к TMDb: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Статус ответа от TMDb: %s", resp.Status)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Получен не-OK статус от TMDb: %d", resp.StatusCode)
		return nil, fmt.Errorf("TMDb API вернул ошибку: %s", resp.Status)
	}

	var tmdbResponse TMDbPopularResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResponse); err != nil {
		log.Printf("ОШИБКА при декодировании JSON ответа от TMDb: %v", err)
		return nil, fmt.Errorf("ошибка при декодировании ответа от TMDb: %w", err)
	}

	log.Printf("Получено %d фильмов от TMDb", len(tmdbResponse.Results))

	var movies []*pb.Movie
	for _, movie := range tmdbResponse.Results {
		movies = append(movies, &pb.Movie{
			Id:            movie.ID,
			Title:         movie.Title,
			PosterPath:    "https://image.tmdb.org/t/p/w500" + movie.PosterPath,
			OriginalTitle: movie.OriginalTitle,
			VoteAverage:   movie.VoteAverage,
			Overview:      movie.Overview,
			ReleaseDate:   movie.ReleaseDate,
		})
	}

	response := &pb.SearchResponse{
		Results:    movies,
		Page:       int32(tmdbResponse.Page),
		TotalPages: int32(tmdbResponse.TotalPages),
	}

	log.Printf("Отправляю ответ клиенту. Количество фильмов: %d", len(response.Results))
	return response, nil
}

func (s *Server) GetMovieByID(ctx context.Context, req *pb.GetMovieByIDRequest) (*pb.Movie, error) {
	movieID := req.GetMovieId()

	url := fmt.Sprintf("%s/movie/%d?api_key=%s&language=%s",
		tmdbBaseURL, movieID, tmdbAPIKey, req.GetLanguage())
	log.Printf("Выполняю запрос к TMDb по URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("ОШИБКА при выполнении HTTP-запроса к TMDb: %v", err)
		return nil, fmt.Errorf("ошибка при выполнении запроса к TMDb: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Статус ответа от TMDb: %s", resp.Status)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Получен не-OK статус от TMDb: %d", resp.StatusCode)
		return nil, fmt.Errorf("TMDb API вернул ошибку: %s", resp.Status)
	}

	var tmdbResponse TMDbMovie
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResponse); err != nil {
		log.Printf("ОШИБКА при декодировании JSON ответа от TMDb: %v", err)
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
