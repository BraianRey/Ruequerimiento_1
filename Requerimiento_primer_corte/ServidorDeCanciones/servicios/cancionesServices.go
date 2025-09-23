package servicios

import (
	. "servidor.local/grpc-servidorcanciones/modelos"
)

func CargarCanciones(vec []Cancion) {
	var objCancion1, objCancion2, objCancion3 Cancion

	objCancion1.Titulo = "Canción 1"
	objCancion1.Tamanio = 10
	objCancion1.Url = "Ruta canción 1"
	objCancion1.EsActivada = true

	objCancion2.Titulo = "Canción 2"
	objCancion2.Tamanio = 20
	objCancion2.Url = "Ruta canción 2"
	objCancion2.EsActivada = false

	objCancion3.Titulo = "Canción 3"
	objCancion3.Tamanio = 30
	objCancion3.Url = "Ruta cancion 3"
	objCancion3.EsActivada = true

	vec[0] = objCancion1
	vec[1] = objCancion2
	vec[2] = objCancion3
}

func BuscarCancion(titulo string, vectorCanciones []Cancion) RespuestaDTO[Cancion] {
	for i := 0; i < len(vectorCanciones); i++ {
		if vectorCanciones[i].Titulo == titulo {
			var resp RespuestaDTO[Cancion]
			resp.Data = vectorCanciones[i]
			resp.Codigo = 200
			resp.Mensaje = "Canción encontrada"
			return resp
		}
	}
	var resp RespuestaDTO[Cancion]
	resp.Codigo = 404
	resp.Mensaje = "La canción no se encontró"
	return resp
}
