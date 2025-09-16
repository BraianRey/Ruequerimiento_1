package vistas

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	util "cliente.local/grpc-cliente/utilidades"
	pb "servidor.local/grpc-servidorstream/serviciosCancion"
)

func MostrarMenuPrincipal(client pb.AudioServiceClient, ctx context.Context) {

	fmt.Print("\n Menu de audio mediante streaming gRPC \n")

	readerInput := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el título de la canción: ")
	titulo, _ := readerInput.ReadString('\n')
	titulo = strings.TrimSpace(titulo)

	//Invocación del procedimiento remoto
	stream, err := client.EnviarCancionMedianteStream(ctx, &pb.PeticionDTO{Titulo: titulo})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Recibiendo y reproduciendo canción en vivo...")
	reader, writer := io.Pipe()
	canalSincronizacion := make(chan struct{})

	// Arranca la goroutine de decodificación y reproducción de los fragmentos
	go util.DecodificarReproducir(reader, canalSincronizacion)

	// Arranca la recepción de los fragmentos desde el servidor
	util.RecibirCancion(stream, writer, canalSincronizacion)
}
