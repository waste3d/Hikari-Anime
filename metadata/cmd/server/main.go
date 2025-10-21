// file: metadata/cmd/server/main.go
package main

import (
	"log"
	"net"

	"github.com/waste3d/Hikari-Anime/metadata"
	pb "github.com/waste3d/Hikari-Anime/metadata/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = ":50051"

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("не удалось начать прослушивание порта: %v", err)
	}

	grpcServer := grpc.NewServer()

	metadataServer := metadata.NewServer()

	pb.RegisterMetadataServiceServer(grpcServer, metadataServer)
	reflection.Register(grpcServer)

	log.Printf("Сервер запущен и слушает порт %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}
}
