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

// DecodificarReproducir: decodifica desde el reader y reproduce hasta EOF o hasta que llegue stop.
// Ahora usa un sync.Once local para cerrar canalSincronizacion solo 1 vez y evita panics por double close.
func DecodificarReproducir(reader io.Reader, canalStop <-chan struct{}, canalSincronizacion chan struct{}) {
	// decodificar el stream de audio
	streamer, format, err := mp3.Decode(io.NopCloser(reader))
	if err != nil {
		// No terminamos la aplicación con Fatalf: devolvemos el error vía log y cerramos la sincronización.
		log.Printf("error decodificando MP3: %v", err)

		// asegurar que canalSincronizacion se cierra para no bloquear al llamador
		select {
		case <-canalStop:
		default:
		}
		// cerrar canalSincronizacion de forma segura
		safeClose(canalSincronizacion)
		return
	}

	// inicializar el altavoz solo la primera vez con el sample rate actual
	speakerInitOnce.Do(func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic initializing speaker: %v", r)
			}
		}()
		// buffer de segundo/2 para latencia similar al código original
		if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/2)); err != nil {
			log.Printf("inicializando speaker: %v", err)
		}
	})

	// callback que se dispara al terminar la reproducción natural
	done := make(chan struct{})
	// reproducir secuencia y avisar con callback
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))

	// asegurar que canalSincronizacion solo se cierre una vez
	var once sync.Once
	closeSync := func() { once.Do(func() { safeClose(canalSincronizacion) }) }

	// goroutine que escucha stop o finalización natural
	go func() {
		select {
		case <-canalStop:
			// detener inmediatamente
			log.Println("DecodificarReproducir: stop recibido, deteniendo reproducción.")
			// Clear detiene lo que está sonando
			speaker.Clear()
			// intentar cerrar streamer (puede devolver error)
			_ = streamer.Close()
			closeSync()
		case <-done:
			// terminó la canción sola
			log.Println("DecodificarReproducir: reproducción terminó naturalmente.")
			closeSync()
		}
	}()
}

// safeClose cierra un canal de señal (chan struct{}) sin provocar panic por doble close.
// Si canal es nil o ya cerrado, no hace nada.
func safeClose(ch chan struct{}) {
	if ch == nil {
		return
	}
	// cerrarlo de forma no bloqueante: intentar enviar y recuperar panic
	defer func() {
		_ = recover()
	}()
	close(ch)
}

// RecibirCancion: recibe los fragmentos del stream y los escribe al pipe.
// Se interrumpe si canalStop se cierra.
func RecibirCancion(stream pb.AudioService_EnviarCancionMedianteStreamClient, writer *io.PipeWriter, canalStop <-chan struct{}, canalSincronizacion chan struct{}) {
	noFragmento := 0

	// función auxiliar para cerrar writer solo una vez
	var once sync.Once
	closeWriter := func() {
		once.Do(func() {
			_ = writer.Close()
		})
	}

	for {
		select {
		case <-canalStop:
			fmt.Println("\nRecibirCancion: stop recibido, cerrando writer.")
			// cerrar writer para que el reader/decoder reciba EOF
			closeWriter()
			// esperar a que la reproducción confirme finalización
			<-canalSincronizacion
			fmt.Println("✓ Reproducción finalizada.")
			return
		default:
			fragmento, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Canción recibida completa.")
				// cerrar writer para que el decoder lea EOF y termine
				closeWriter()
				// esperar a que el decodificador confirme finalización
				<-canalSincronizacion
				fmt.Println("✓ Reproducción finalizada.")
				return
			}
			if err != nil {
				// Loguear error con detalle para diagnosticar si es de transporte/ctx/cancel
				log.Printf("Error recibiendo chunk (stream.Recv): %T - %v", err, err)
				// cerrar writer para liberar al reader
				closeWriter()
				// Esperar al decodificador para que libere recursos (si aplica)
				<-canalSincronizacion
				fmt.Println("Reproducción finalizada por error en recv.")
				return
			}
			noFragmento++
			fmt.Printf("\n Fragmento #%d recibido (%d bytes) reproduciendo ...", noFragmento, len(fragmento.Data))
			// Intentar escribir todo el fragmento (handle de error)
			if _, err := writer.Write(fragmento.Data); err != nil {
				log.Printf("Error escribiendo en pipe: %v", err)
				// cerrar writer y esperar sincronización del decodificador
				closeWriter()
				<-canalSincronizacion
				fmt.Println("Reproducción finalizada.")
				return
			}
		}
	}
}
