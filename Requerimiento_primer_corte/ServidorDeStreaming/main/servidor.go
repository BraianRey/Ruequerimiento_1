package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	capacontroladores "servidor.local/grpc-servidorstream/capaControladores"
	pb "servidor.local/grpc-servidorstream/serviciosStreaming"
)

func main() {

	// Escuchar en el puerto 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error escuchando en el puerto: %v", err)
	}

	// Crear un servidor gRPC
	grpcServer := grpc.NewServer()

	// Se registra el controlador que ofrece el procedimiento remoto
	pb.RegisterAudioServiceServer(grpcServer, &capacontroladores.ControladorServidor{})

	fmt.Println("Servidor gRPC escuchando en :50051...")

	// Iniciar el servidor
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error al iniciar servidor gRPC: %v", err)
	}
}
