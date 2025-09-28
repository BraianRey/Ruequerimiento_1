package modelos

// Cancion representa una canción con sus detalles
type Cancion struct {
	ID       int
	Titulo   string
	Artista  string
	Album    string
	Anio     int
	Duracion string
	Genero   Genero
}
