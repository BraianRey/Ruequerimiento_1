package servicios

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"servidor.local/grpc-servidorcanciones/modelos"
	pb "servidor.local/grpc-servidorcanciones/serviciosCanciones"
)

// implementación del servicio de canciones
type CancionesServiceServer struct {
	pb.UnimplementedCancionesServiceServer
	Generos   []modelos.Genero
	Canciones []modelos.Cancion
}

// métodos gRPC que implementan la interfaz del servicio de canciones
// ListarGeneros devuelve la lista de géneros disponibles
func (s *CancionesServiceServer) ListarGeneros(ctx context.Context, in *emptypb.Empty) (*pb.ListaGeneros, error) {
	var generos []*pb.Genero
	for _, g := range s.Generos {
		generos = append(generos, &pb.Genero{
			Id:     int32(g.ID),
			Nombre: g.Nombre,
		})
	}
	return &pb.ListaGeneros{Generos: generos}, nil
}

// ListarCancionesPorGenero devuelve las canciones de un género específico
func (s *CancionesServiceServer) ListarCancionesPorGenero(ctx context.Context, in *pb.GeneroId) (*pb.ListaCanciones, error) {
	var canciones []*pb.Cancion
	for _, c := range s.Canciones {
		if int32(c.Genero.ID) == in.Id {
			canciones = append(canciones, &pb.Cancion{
				Id:       int32(c.ID),
				Titulo:   c.Titulo,
				Artista:  c.Artista,
				Anio:     int32(c.Anio),
				Duracion: c.Duracion,
				Genero: &pb.Genero{
					Id:     int32(c.Genero.ID),
					Nombre: c.Genero.Nombre,
				},
			})
		}
	}
	return &pb.ListaCanciones{Canciones: canciones}, nil
}

// ObtenerDetallesCancion devuelve los detalles de una canción específica
func (s *CancionesServiceServer) ObtenerDetallesCancion(ctx context.Context, in *pb.CancionId) (*pb.Cancion, error) {
	for _, c := range s.Canciones {
		if int32(c.ID) == in.Id {
			return &pb.Cancion{
				Id:       int32(c.ID),
				Titulo:   c.Titulo,
				Artista:  c.Artista,
				Anio:     int32(c.Anio),
				Duracion: c.Duracion,
				Genero: &pb.Genero{
					Id:     int32(c.Genero.ID),
					Nombre: c.Genero.Nombre,
				},
			}, nil
		}
	}
	return nil, nil // O retorna un error si no se encuentra
}
