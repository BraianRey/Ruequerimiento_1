package capacontroladores

import (
	capafachadaservices "servidor.local/grpc-servidor/capaFachadaServices"
	pb "servidor.local/grpc-servidor/serviciosCancion"
)

type ControladorServidor struct {
	pb.UnimplementedAudioServiceServer
}

// Implementación del procedimiento remoto
func (s *ControladorServidor) EnviarCancionMedianteStream(req *pb.PeticionDTO, stream pb.AudioService_EnviarCancionMedianteStreamServer) error {
	return capafachadaservices.StreamAudioFile(
		req.Titulo,
		func(data []byte) error {
			return stream.Send(&pb.FragmentoCancion{Data: data})
		})
}
