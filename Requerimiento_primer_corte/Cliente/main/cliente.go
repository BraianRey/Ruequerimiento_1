package main

import (
	"context"
	"log"
	"time"

	menu "cliente.local/grpc-cliente/vistas"
	"google.golang.org/grpc"
	pbCancion "servidor.local/grpc-servidorcanciones/serviciosCanciones"
	pbStream "servidor.local/grpc-servidorstream/serviciosStreaming"
)

func main() {
	// Conexión al servidor de streaming
	connStream, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer connStream.Close()
	clientStream := pbStream.NewAudioServiceClient(connStream)

	// Conexión al servidor de canciones
	connCancion, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer connCancion.Close()
	clientCancion := pbCancion.NewCancionesServiceClient(connCancion)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	menu.MostrarMenuPrincipal(clientStream, clientCancion, ctx)
}
