package modelos

type RespuestaDTO[T any] struct {
	Data    T
	Codigo  int
	Mensaje string
}
