package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"servidor.local/grpc-servidorcanciones/controladores"
	"servidor.local/grpc-servidorcanciones/fachada"
	pb "servidor.local/grpc-servidorcanciones/serviciosCanciones"
)

const addr = ":9000"

func main() {
	// inicializar la fachada y el servidor gRPC
	fachadaCanciones := fachada.NuevaFachadaCanciones()

	// iniciar el servidor gRPC
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("No se pudo escuchar en %s: %v", addr, err)
	}
	// crear servidor gRPC y registrar el servicio de canciones
	grpcServer := grpc.NewServer()
	// registrar el servicio de canciones con su controlador
	pb.RegisterCancionesServiceServer(grpcServer, controladores.NuevoCancionesController(fachadaCanciones))

	log.Printf("ServidorDeCanciones escuchando en %s (gRPC)", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
