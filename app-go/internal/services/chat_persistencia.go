package services

import (
	"fmt"

	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/storage"
)

// NewChatServiceConArchivo crea el servicio de mensajes
// y recupera los mensajes guardados anteriormente.
func NewChatServiceConArchivo(
	usuarios *UsuarioService,
	ruta string,
) (*ChatService, error) {

	// Crear el servicio de mensajeria existente.
	servicio := NewChatService(usuarios)

	// Preparar el archivo JSON.
	archivo := storage.NuevoArchivoJSON(ruta)

	// Crear una lista para recuperar los registros.
	var registros []storage.MensajeGuardado

	// Leer los mensajes almacenados.
	err := archivo.Cargar(&registros)

	if err != nil {
		return nil, fmt.Errorf(
			"error al cargar mensajes: %w",
			err,
		)
	}

	// Reconstruir los mensajes originales.
	mensajes, err := storage.ReconstruirMensajes(registros)

	if err != nil {
		return nil, fmt.Errorf(
			"error al reconstruir mensajes: %w",
			err,
		)
	}

	// Recuperar los mensajes en memoria.
	for _, mensaje := range mensajes {

		// Verificar que el remitente exista.
		if _, err := usuarios.BuscarUsuario(
			mensaje.ObtenerRemitente(),
		); err != nil {
			return nil, fmt.Errorf(
				"mensaje %d: remitente no existe",
				mensaje.ObtenerID(),
			)
		}

		// Verificar que el destinatario exista.
		if _, err := usuarios.BuscarUsuario(
			mensaje.ObtenerDestinatario(),
		); err != nil {
			return nil, fmt.Errorf(
				"mensaje %d: destinatario no existe",
				mensaje.ObtenerID(),
			)
		}

		servicio.mensajes = append(
			servicio.mensajes,
			mensaje,
		)

		// Conservar la numeracion de los mensajes.
		if mensaje.ObtenerID() >= servicio.siguiente {
			servicio.siguiente = mensaje.ObtenerID() + 1
		}
	}

	// Asociar el archivo JSON al servicio.
	servicio.archivo = archivo

	return servicio, nil
}

// guardarMensajesBloqueados guarda los mensajes
// mientras el servicio mantiene su bloqueo.
func (s *ChatService) guardarMensajesBloqueados() error {

	if s.archivo == nil {
		return nil
	}

	// Convertir los mensajes a registros JSON.
	registros := storage.ConvertirMensajes(s.mensajes)

	// Guardar los registros.
	if err := s.archivo.Guardar(registros); err != nil {
		return fmt.Errorf(
			"error al guardar mensajes: %w",
			err,
		)
	}

	return nil
}
