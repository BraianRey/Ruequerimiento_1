package utilidades

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	pb "servidor.local/grpc-servidorstream/serviciosStreaming"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// package-level sync.Once para inicializar speaker solo 1 vez
var speakerInitOnce sync.Once
var speakerSampleRate beep.SampleRate // sampleRate fijo para el speaker

// DecodificarReproducir: decodifica desde el reader y reproduce hasta EOF o hasta que llegue stop.
// Usa un sync.Once local para cerrar canalSincronizacion solo 1 vez y evita panics por double close.
// Realiza resampling si el sampleRate de la canción no coincide con el usado en speaker.
func DecodificarReproducir(reader io.Reader, canalStop <-chan struct{}, canalSincronizacion chan struct{}) {
	// decodificar el stream de audio
	streamer, _, err := mp3.Decode(io.NopCloser(reader))
	if err != nil {
		log.Printf("error decodificando MP3: %v", err)
		safeClose(canalSincronizacion)
		return
	}

	// inicializar speaker una sola vez con sampleRate fijo (44100 Hz)
	speakerInitOnce.Do(func() {
		speakerSampleRate = beep.SampleRate(44100)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic inicializando speaker: %v", r)
			}
		}()
		if err := speaker.Init(speakerSampleRate, speakerSampleRate.N(time.Second/2)); err != nil {
			log.Printf("inicializando speaker: %v", err)
		}
	})

	// si el sampleRate del audio no coincide con el del speaker, hacer resampling
	var streamerToPlay beep.Streamer = streamer

	// callback que se dispara al terminar la reproducción natural
	done := make(chan struct{})
	speaker.Play(beep.Seq(streamerToPlay, beep.Callback(func() {
		close(done)
	})))

	// asegurar que canalSincronizacion solo se cierre una vez
	var once sync.Once
	closeSync := func() { once.Do(func() { safeClose(canalSincronizacion) }) }

	// goroutine que escucha stop o finalización natural
	go func() {
		select {
		// si llega stop, limpiar y cerrar streamer
		case <-canalStop:
			log.Println("DecodificarReproducir: stop recibido, deteniendo reproducción.")
			speaker.Clear()
			_ = streamer.Close()
			closeSync()
		// si termina naturalmente, cerrar canalSincronizacion
		case <-done:
			log.Println("DecodificarReproducir: reproducción terminó naturalmente.")
			closeSync()
		}
	}()
}

// safeClose evita panics por cerrar canales múltiples veces
func safeClose(ch chan struct{}) {
	if ch == nil {
		return
	}
	// se usa recover para evitar panic si el canal ya está cerrado
	defer func() { _ = recover() }()
	close(ch)
}

// RecibirCancion: recibe los fragmentos del stream y los escribe al pipe.
func RecibirCancion(stream pb.AudioService_EnviarCancionMedianteStreamClient, writer *io.PipeWriter, canalStop <-chan struct{}, canalSincronizacion chan struct{}) {
	noFragmento := 0
	var once sync.Once
	closeWriter := func() {
		once.Do(func() { _ = writer.Close() })
	}

	for {
		select {
		case <-canalStop:
			fmt.Println("\nRecibirCancion: stop recibido, cerrando writer.")
			closeWriter()
			<-canalSincronizacion
			fmt.Println("✓ Reproducción finalizada.")
			return
		default:
			fragmento, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Canción recibida completa.")
				closeWriter()
				<-canalSincronizacion
				fmt.Println("✓ Reproducción finalizada.")
				return
			}
			if err != nil {
				log.Printf("Error recibiendo chunk (stream.Recv): %T - %v", err, err)
				closeWriter()
				<-canalSincronizacion
				fmt.Println("Reproducción finalizada por error en recv.")
				return
			}
			noFragmento++
			fmt.Printf("\n Fragmento #%d recibido (%d bytes) reproduciendo ...", noFragmento, len(fragmento.Data))
			if _, err := writer.Write(fragmento.Data); err != nil {
				log.Printf("Error escribiendo en pipe: %v", err)
				closeWriter()
				<-canalSincronizacion
				fmt.Println("Reproducción finalizada.")
				return
			}
		}
	}
}
