package fachada

import (
	"log"

	"servidor.local/grpc-servidorcanciones/modelos"
)

// FachadaCanciones proporciona una interfaz simplificada para acceder a la lógica de negocio relacionada con las canciones y géneros
type FachadaCanciones struct {
	Generos   []modelos.Genero  // lista de géneros disponibles
	Canciones []modelos.Cancion // lista de canciones disponibles
}

// constructor de la fachada de canciones
func NuevaFachadaCanciones() *FachadaCanciones {
	generos := []modelos.Genero{
		{ID: 1, Nombre: "Rock"},
		{ID: 2, Nombre: "Pop"},
		{ID: 3, Nombre: "Clásica"},
	}
	canciones := []modelos.Cancion{
		{
			ID:       1,
			Titulo:   "Lamento Boliviano",
			Artista:  "Enanitos Verdes",
			Album:    "Contrarreloj",
			Anio:     1994,
			Duracion: "4:20",
			Genero:   generos[0],
		},
		{
			ID:       2,
			Titulo:   "De Música Ligera",
			Artista:  "Soda Stereo",
			Album:    "Canción Animal",
			Anio:     1990,
			Duracion: "3:30",
			Genero:   generos[0],
		},
		{
			ID:       3,
			Titulo:   "La Flaca",
			Artista:  "Jarabe de Palo",
			Album:    "La Flaca",
			Anio:     1996,
			Duracion: "4:00",
			Genero:   generos[0],
		},
		{
			ID:       4,
			Titulo:   "Hey Ya!",
			Artista:  "OutKast",
			Album:    "Speakerboxxx/The Love Below",
			Anio:     2003,
			Duracion: "3:55",
			Genero:   generos[1],
		},
		{
			ID:       5,
			Titulo:   "Umbrella",
			Artista:  "Rihanna",
			Album:    "Good Girl Gone Bad",
			Anio:     2007,
			Duracion: "4:36",
			Genero:   generos[1],
		},
		{
			ID:       6,
			Titulo:   "Sunflower",
			Artista:  "Post Malone & Swae Lee",
			Album:    "Spiderman Into the Spider-Verse",
			Anio:     2018,
			Duracion: "2:38",
			Genero:   generos[1],
		},
		{
			ID:       7,
			Titulo:   "Moonlight Sonata",
			Artista:  "Ludwig van Beethoven",
			Album:    "Piano Sonata No. 14",
			Anio:     1791,
			Duracion: "14:59",
			Genero:   generos[2],
		},
		{
			ID:       8,
			Titulo:   "Fur Elise",
			Artista:  "Ludwig van Beethoven",
			Album:    "Bagatelle No. 25",
			Anio:     1792,
			Duracion: "5:06",
			Genero:   generos[2],
		},
		{
			ID:       9,
			Titulo:   "Hungarian Rhapsody No2",
			Artista:  "Franz Liszt",
			Album:    "Hungarian Rhapsodies",
			Anio:     1847,
			Duracion: "10:31",
			Genero:   generos[2],
		},
	}
	return &FachadaCanciones{
		Generos:   generos,   // inicializar géneros
		Canciones: canciones, // inicializar canciones
	}
}

// ListarGeneros devuelve la lista de géneros disponibles
func (f *FachadaCanciones) ListarGeneros() []modelos.Genero {
	log.Printf("Se a solicitado la lista de géneros, se tienen %d géneros", len(f.Generos))
	return f.Generos
}

// ListarCancionesPorGenero devuelve las canciones de un género específico
func (f *FachadaCanciones) ListarCancionesPorGenero(id int32) []modelos.Cancion {
	var canciones []modelos.Cancion
	for _, c := range f.Canciones {
		if int32(c.Genero.ID) == id {
			canciones = append(canciones, c)
		}
	}
	log.Printf("Se a solicitado la lista de canciones del género %d, se tienen %d canciones", id, len(canciones))
	return canciones
}

// ObtenerDetallesCancion devuelve los detalles de una canción específica
func (f *FachadaCanciones) ObtenerDetallesCancion(id int32) *modelos.Cancion {
	for _, c := range f.Canciones {
		if int32(c.ID) == id {
			log.Printf("Se a solicitado los detalles de la canción '%s'", c.Titulo)
			return &c
		}
	}
	log.Printf("No se encontraron detalles para la canción con ID %d", id)
	return nil
}
