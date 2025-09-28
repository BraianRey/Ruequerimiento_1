package controladores

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"servidor.local/grpc-servidorcanciones/fachada"
	pb "servidor.local/grpc-servidorcanciones/serviciosCanciones"
)

// controlador de canciones
// contiene la fachada para acceder a la lógica de negocio
type CancionesController struct {
	pb.UnimplementedCancionesServiceServer
	fachada *fachada.FachadaCanciones
}

// constructor del controlador de canciones
func NuevoCancionesController(f *fachada.FachadaCanciones) *CancionesController {
	return &CancionesController{fachada: f}
}

// métodos gRPC que implementan la interfaz del servicio de canciones
// ListarGeneros devuelve la lista de géneros disponibles
func (c *CancionesController) ListarGeneros(ctx context.Context, in *emptypb.Empty) (*pb.ListaGeneros, error) {
	var generos []*pb.Genero
	for _, g := range c.fachada.ListarGeneros() {
		generos = append(generos, &pb.Genero{
			Id:     int32(g.ID),
			Nombre: g.Nombre,
		})
	}
	return &pb.ListaGeneros{Generos: generos}, nil
}

// ListarCancionesPorGenero devuelve las canciones de un género específico
func (c *CancionesController) ListarCancionesPorGenero(ctx context.Context, in *pb.GeneroId) (*pb.ListaCanciones, error) {
	var canciones []*pb.Cancion
	for _, cgo := range c.fachada.ListarCancionesPorGenero(in.Id) {
		canciones = append(canciones, &pb.Cancion{
			Id:       int32(cgo.ID),
			Titulo:   cgo.Titulo,
			Artista:  cgo.Artista,
			Album:    cgo.Album,
			Anio:     int32(cgo.Anio),
			Duracion: cgo.Duracion,
			Genero: &pb.Genero{
				Id:     int32(cgo.Genero.ID),
				Nombre: cgo.Genero.Nombre,
			},
		})
	}
	return &pb.ListaCanciones{Canciones: canciones}, nil
}

// ObtenerDetallesCancion devuelve los detalles de una canción específica
func (c *CancionesController) ObtenerDetallesCancion(ctx context.Context, in *pb.CancionId) (*pb.Cancion, error) {
	cgo := c.fachada.ObtenerDetallesCancion(in.Id)
	if cgo == nil {
		return nil, nil
	}
	return &pb.Cancion{
		Id:       int32(cgo.ID),
		Titulo:   cgo.Titulo,
		Artista:  cgo.Artista,
		Album:    cgo.Album,
		Anio:     int32(cgo.Anio),
		Duracion: cgo.Duracion,
		Genero: &pb.Genero{
			Id:     int32(cgo.Genero.ID),
			Nombre: cgo.Genero.Nombre,
		},
	}, nil
}
