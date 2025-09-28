package vistas

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	util "cliente.local/grpc-cliente/utilidades"
	"google.golang.org/protobuf/types/known/emptypb"
	pbCancion "servidor.local/grpc-servidorcanciones/serviciosCanciones"
	pbStream "servidor.local/grpc-servidorstream/serviciosStreaming"
)

func MostrarMenuPrincipal(clientStream pbStream.AudioServiceClient, clientCancion pbCancion.CancionesServiceClient, ctx context.Context) {
	// lector de entrada estándar
	readerInput := bufio.NewReader(os.Stdin)
	// menú principal
	for {
		fmt.Print("\n1) Ver géneros\n0) Salir\nSeleccione una opción: ")
		opcion, _ := readerInput.ReadString('\n')
		opcion = strings.TrimSpace(opcion)
		if opcion == "0" {
			fmt.Println("Saliendo...")
			return
		}
		if opcion != "1" {
			fmt.Println("Opción inválida.")
			continue
		}

		// Paso 1: Listar géneros
		respGeneros, err := clientCancion.ListarGeneros(ctx, &emptypb.Empty{})
		if err != nil {
			fmt.Println("Error al obtener géneros:", err)
			continue
		}
		if len(respGeneros.Generos) == 0 {
			fmt.Println("No hay géneros disponibles.")
			continue
		}
		fmt.Println("\nGéneros disponibles:")
		for i, g := range respGeneros.Generos {
			fmt.Printf("%d) %s\n", i+1, g.Nombre)
		}
		fmt.Print("Seleccione el género por número (o 0 para volver): ")
		genStr, _ := readerInput.ReadString('\n')
		genStr = strings.TrimSpace(genStr)
		if genStr == "0" {
			continue
		}
		genIdx, err := strconv.Atoi(genStr)
		if err != nil || genIdx < 1 || genIdx > len(respGeneros.Generos) {
			fmt.Println("Opción inválida.")
			continue
		}
		idGenero := respGeneros.Generos[genIdx-1].Id

		// Paso 2: Listar canciones del género seleccionado
		respCanciones, err := clientCancion.ListarCancionesPorGenero(ctx, &pbCancion.GeneroId{Id: idGenero})
		if err != nil {
			fmt.Println("Error al obtener canciones:", err)
			continue
		}
		if len(respCanciones.Canciones) == 0 {
			fmt.Println("No hay canciones en este género.")
			continue
		}
		fmt.Println("\nCanciones disponibles:")
		for i, c := range respCanciones.Canciones {
			fmt.Printf("%d) %s - %s\n", i+1, c.Titulo, c.Artista)
		}
		fmt.Print("Seleccione la canción por número (o 0 para volver): ")
		cancStr, _ := readerInput.ReadString('\n')
		cancStr = strings.TrimSpace(cancStr)
		if cancStr == "0" {
			continue
		}
		cancIdx, err := strconv.Atoi(cancStr)
		if err != nil || cancIdx < 1 || cancIdx > len(respCanciones.Canciones) {
			fmt.Println("Opción inválida.")
			continue
		}
		idCancion := respCanciones.Canciones[cancIdx-1].Id

		// Paso 3: Mostrar información de la canción seleccionada
		respDetalle, err := clientCancion.ObtenerDetallesCancion(ctx, &pbCancion.CancionId{Id: idCancion})
		if err != nil || respDetalle == nil {
			fmt.Println("Error al obtener detalles o canción no encontrada.")
			continue
		}
		fmt.Printf("\nDetalles de la canción:\nTítulo: %s\nArtista: %s\nAño: %d\nDuración: %s\nGénero: %s\n",
			respDetalle.Titulo, respDetalle.Artista, respDetalle.Anio, respDetalle.Duracion, respDetalle.Genero.Nombre)

		// Paso 4: Opción de reproducir o volver
		for {
			fmt.Print("\n1) Reproducir\n0) Volver\nSeleccione una opción: ")
			opc, _ := readerInput.ReadString('\n')
			opc = strings.TrimSpace(opc)
			if opc == "0" {
				break
			}
			if opc == "1" {
				stream, err := clientStream.EnviarCancionMedianteStream(ctx, &pbStream.PeticionDTO{Titulo: respDetalle.Titulo})
				if err != nil {
					log.Fatal(err)
				}
				readerPipe, writerPipe := io.Pipe()
				canalStop := make(chan struct{})
				canalSincronizacion := make(chan struct{})

				// Goroutine para reproducir
				go util.DecodificarReproducir(readerPipe, canalStop, canalSincronizacion)

				// Goroutine para recibir datos de la canción
				go util.RecibirCancion(stream, writerPipe, canalStop, canalSincronizacion)

				// Menú de reproducción
				for {
					fmt.Print("\nReproduciendo...\n0) Salir\nSeleccione una opción: ")
					subOpc, _ := readerInput.ReadString('\n')
					subOpc = strings.TrimSpace(subOpc)
					if subOpc == "0" {
						// detener reproducción
						close(canalStop)
						<-canalSincronizacion // esperar a que todo termine
						fmt.Println("Reproducción detenida.")
						break
					}
				}
				break
			}
			fmt.Println("Opción inválida.")
		}
	}
}
