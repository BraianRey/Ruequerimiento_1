package modelos

// RespuestaDTO es una estructura genérica para respuestas con datos, código y mensaje
type RespuestaDTO[T any] struct {
	Data    T
	Codigo  int
	Mensaje string
}
