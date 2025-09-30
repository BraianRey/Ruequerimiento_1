package capafachadaservices

import (
	"fmt"
	"io"
	"log"
	"os"
)

// StreamAudioFile lee un archivo de audio y envía sus fragmentos usando la función proporcionada
func StreamAudioFile(tituloCancion string, funcionParaEnviarFragmento func([]byte) error) error {
	tituloCancion = tituloCancion + ".mp3"
	log.Printf("Canción solicitada: %s", tituloCancion)
	// abrir el archivo de audio
	file, err := os.Open("canciones/" + tituloCancion)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	// asegurar el cierre del archivo al finalizar
	defer file.Close()

	// leer y enviar en fragmentos
	buffer := make([]byte, 64*1024) // 64 KB se envian por fragmento
	fragmento := 0

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			log.Println("Canción enviada completamente desde la fachada.")
			break
		}
		if err != nil {
			return fmt.Errorf("error leyendo el archivo: %w", err)
		}
		fragmento++
		log.Printf("Fragmento #%d leido (%d bytes) y enviando", fragmento, n)

		// ejecutamos la función para enviar el fragmento cliente
		err = funcionParaEnviarFragmento(buffer[:n])
		if err != nil {
			return fmt.Errorf("error enviando fragmento #%d: %w", fragmento, err)
		}
	}

	return nil
}
