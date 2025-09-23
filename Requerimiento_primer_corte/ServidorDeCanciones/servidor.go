package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"servidor.local/grpc-servidorcanciones/modelos"
	"servidor.local/grpc-servidorcanciones/servicios"
)

const addr = "localhost:9000"

func handleConnection(conn net.Conn, vectorCanciones []modelos.Cancion) {

	remote := conn.RemoteAddr().String()

	fmt.Printf("Conexión entrante desde %s\n", remote)

	//Leer del canal el titulo enviado por el cliente
	buffer := make([]byte, 1024) // tamaño máximo esperado
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error leyendo desde %s: %v\n", remote, err)
		return
	}
	titulo := string(buffer[:n])
	titulo = strings.TrimSpace(titulo)

	fmt.Printf("Título recibido de %s: %q\n", remote, titulo)

	objRespuesta := servicios.BuscarCancion(titulo, vectorCanciones)
	jsonResp, _ := json.Marshal(objRespuesta)

	_, err = conn.Write(append(jsonResp, '\n'))
	if err != nil {
		fmt.Printf("Error enviado a %s: %v\n", remote, err)
		return
	}

	//Escribir en el canal la respuesta
	time.Sleep(10 * time.Second)

	fmt.Printf("Atendido %s, título: %q\n", remote, titulo)

	//Cerrar el canal virtual
	conn.Close()
}

func main() {
	//Al iniciar el servidor se cargan las canciones
	vectorCanciones := make([]modelos.Cancion, 3)
	servicios.CargarCanciones(vectorCanciones)

	//Colocar el servidor en estado de escucha
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Sprintf("No se pudo escuchar en %s: %v", addr, err))
	}

	defer ln.Close()

	fmt.Printf("Servidor escuchando en %s\n", addr)
	for {
		//Aceptar conexión a un canal virtual
		//La referencia al canal virtual se almacena en la varible objReferenciaCanal
		objReferenciaCanal, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error al aceptar conexión: %v\n", err)
			continue
		}
		go handleConnection(objReferenciaCanal, vectorCanciones)

		// Enviamos la referencia al canal virtual y el vector de canciones
		// a una go rutine para que lo atienda
		// y el servidor pueda seguir escuchando otras peticiones
		// de otros clientes

	}
}
