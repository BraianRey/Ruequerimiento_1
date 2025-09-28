package utilidades

import (
	"fmt"
	"io"
	"log"
	"time"

	pb "servidor.local/grpc-servidorstream/serviciosStreaming"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// DecodificarReproducir: decodifica desde el reader y reproduce hasta EOF o hasta que llegue stop.
func DecodificarReproducir(reader io.Reader, canalStop <-chan struct{}, canalSincronizacion chan struct{}) {
	// decodificar el stream de audio
	streamer, format, err := mp3.Decode(io.NopCloser(reader))
	if err != nil {
		log.Fatalf("error decodificando MP3: %v", err)
	}

	// inicializar el altavoz
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/2))

	// callback que se dispara al terminar la reproducción natural
	done := make(chan struct{})
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))

	// goroutine que escucha stop o finalización natural
	go func() {
		select {
		case <-canalStop:
			// detener inmediatamente
			speaker.Clear()
			streamer.Close()
			close(canalSincronizacion)
		case <-done:
			// terminó la canción sola
			close(canalSincronizacion)
		}
	}()
}

// RecibirCancion: recibe los fragmentos del stream y los escribe al pipe.
// Se interrumpe si canalStop se cierra.
func RecibirCancion(stream pb.AudioService_EnviarCancionMedianteStreamClient, writer *io.PipeWriter, canalStop <-chan struct{}, canalSincronizacion chan struct{}) {
	noFragmento := 0
	defer writer.Close()

	for {
		select {
		case <-canalStop:
			fmt.Println("\nRecibirCancion: stop recibido, cerrando writer.")
			return
		default:
			fragmento, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Canción recibida completa.")
				return
			}
			if err != nil {
				log.Printf("Error recibiendo chunk: %v", err)
				return
			}
			noFragmento++
			fmt.Printf("\n Fragmento #%d recibido (%d bytes) reproduciendo ...", noFragmento, len(fragmento.Data))
			if _, err := writer.Write(fragmento.Data); err != nil {
				log.Printf("Error escribiendo en pipe: %v", err)
				return
			}
		}
	}

	// al terminar, esperar fin de reproducción
	<-canalSincronizacion
	fmt.Println("✓ Reproducción finalizada.")
}
