package services

import (
	"errors"
	"sync"

	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/models"
	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/storage"
)

// ChatService administra los mensajes de CONNECTA.
type ChatService struct {
	usuarios  *UsuarioService
	mensajes  []*models.Mensaje
	siguiente int
	mu        sync.Mutex
	archivo   *storage.ArchivoJSON
}

// NewChatService crea el servicio de mensajeria.
func NewChatService(usuarios *UsuarioService) *ChatService {
	return &ChatService{
		usuarios:  usuarios,
		mensajes:  make([]*models.Mensaje, 0),
		siguiente: 1,
	}
}

// EnviarMensaje valida y guarda un mensaje.
func (s *ChatService) EnviarMensaje(
	remitenteID int,
	destinatarioID int,
	contenido string,
) (*models.Mensaje, error) {

	// Comprobar que el remitente exista.
	_, err := s.usuarios.BuscarUsuario(remitenteID)

	if err != nil {
		return nil, errors.New("el remitente no existe")
	}

	// Comprobar que el destinatario exista.
	_, err = s.usuarios.BuscarUsuario(destinatarioID)

	if err != nil {
		return nil, errors.New("el destinatario no existe")
	}

	// Bloquear el servicio durante el registro.
	s.mu.Lock()
	defer s.mu.Unlock()

	// Crear y validar el mensaje.
	mensaje, err := models.NewMensaje(
		s.siguiente,
		remitenteID,
		destinatarioID,
		contenido,
	)

	if err != nil {
		return nil, err
	}

	// Guardar el mensaje en memoria.
	s.mensajes = append(s.mensajes, mensaje)
	s.siguiente++

	// Guardar el mensaje en JSON cuando
	// el almacenamiento persistente este configurado.
	if err := s.guardarMensajesBloqueados(); err != nil {

		// Revertir el registro en memoria si falla
		// el almacenamiento del archivo.
		s.mensajes = s.mensajes[:len(s.mensajes)-1]
		s.siguiente--

		return nil, err
	}

	return mensaje, nil
}

// RecibirMensajes devuelve los mensajes de un destinatario.
func (s *ChatService) RecibirMensajes(
	usuarioID int,
) ([]*models.Mensaje, error) {

	// Comprobar que el usuario exista.
	_, err := s.usuarios.BuscarUsuario(usuarioID)

	if err != nil {
		return nil, errors.New("el usuario no existe")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	resultado := make([]*models.Mensaje, 0)

	// Buscar los mensajes recibidos por el usuario.
	for _, mensaje := range s.mensajes {

		if mensaje.ObtenerDestinatario() == usuarioID {
			resultado = append(resultado, mensaje)
		}
	}

	return resultado, nil
}

// ObtenerHistorial devuelve los mensajes entre dos usuarios.
func (s *ChatService) ObtenerHistorial(
	usuario1 int,
	usuario2 int,
) ([]*models.Mensaje, error) {

	// Comprobar que exista el primer usuario.
	_, err := s.usuarios.BuscarUsuario(usuario1)

	if err != nil {
		return nil, errors.New("el primer usuario no existe")
	}

	// Comprobar que exista el segundo usuario.
	_, err = s.usuarios.BuscarUsuario(usuario2)

	if err != nil {
		return nil, errors.New("el segundo usuario no existe")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	historial := make([]*models.Mensaje, 0)

	// Recorrer los mensajes registrados.
	for _, mensaje := range s.mensajes {

		ida := mensaje.ObtenerRemitente()
		idb := mensaje.ObtenerDestinatario()

		// Recuperar mensajes enviados en ambas direcciones.
		if (ida == usuario1 && idb == usuario2) ||
			(ida == usuario2 && idb == usuario1) {

			historial = append(historial, mensaje)
		}
	}

	return historial, nil
}
