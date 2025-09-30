package capacontroladores

import (
	capafachadaservices "servidor.local/grpc-servidorstream/capaFachadaServices"
	pb "servidor.local/grpc-servidorstream/serviciosStreaming"
)

// ControladorServidor implementa el servicio de streaming de audio
type ControladorServidor struct {
	pb.UnimplementedAudioServiceServer
}

// Implementación del procedimiento remoto
func (s *ControladorServidor) EnviarCancionMedianteStream(req *pb.PeticionDTO, stream pb.AudioService_EnviarCancionMedianteStreamServer) error {
	return capafachadaservices.StreamAudioFile(
		req.Titulo,
		// función para enviar fragmento al cliente
		func(data []byte) error {
			return stream.Send(&pb.FragmentoCancion{Data: data})
		})
}
