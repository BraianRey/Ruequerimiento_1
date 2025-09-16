package main

import (
	"context"
	"log"
	"time"

	menu "cliente.local/grpc-cliente/vistas"
	"google.golang.org/grpc"
	pb "servidor.local/grpc-servidor/serviciosCancion"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewAudioServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	menu.MostrarMenuPrincipal(client, ctx)
}
