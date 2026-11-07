package models

import (
	"errors"
	"strings"
	"time"
)

// Mensaje representa un mensaje enviado entre usuarios.
type Mensaje struct {
	id             int
	remitenteID    int
	destinatarioID int
	contenido      string
	fecha          time.Time
}

// NewMensaje crea un mensaje y valida sus datos.
func NewMensaje(id int, remitenteID int, destinatarioID int, contenido string) (*Mensaje, error) {
	contenido = strings.TrimSpace(contenido)

	if id <= 0 {
		return nil, errors.New("el ID del mensaje debe ser mayor que cero")
	}

	if remitenteID <= 0 || destinatarioID <= 0 {
		return nil, errors.New("los usuarios deben tener un ID valido")
	}

	if remitenteID == destinatarioID {
		return nil, errors.New("no se permite enviar mensajes a uno mismo")
	}

	if contenido == "" {
		return nil, errors.New("el mensaje no puede estar vacio")
	}

	mensaje := &Mensaje{
		id:             id,
		remitenteID:    remitenteID,
		destinatarioID: destinatarioID,
		contenido:      contenido,
		fecha:          time.Now(),
	}

	return mensaje, nil
}

// ObtenerID devuelve el identificador del mensaje.
func (m *Mensaje) ObtenerID() int {
	return m.id
}

// ObtenerRemitente devuelve el ID del usuario que envia.
func (m *Mensaje) ObtenerRemitente() int {
	return m.remitenteID
}

// ObtenerDestinatario devuelve el ID del usuario que recibe.
func (m *Mensaje) ObtenerDestinatario() int {
	return m.destinatarioID
}

// ObtenerContenido devuelve el texto del mensaje.
func (m *Mensaje) ObtenerContenido() string {
	return m.contenido
}

// ObtenerFecha devuelve la fecha de envio.
func (m *Mensaje) ObtenerFecha() time.Time {
	return m.fecha
}
