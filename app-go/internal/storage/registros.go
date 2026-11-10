package storage

import (
	"fmt"
	"time"

	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/models"
)

type UsuarioGuardado struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type MensajeGuardado struct {
	ID             int       `json:"id"`
	RemitenteID    int       `json:"remitente_id"`
	DestinatarioID int       `json:"destinatario_id"`
	Contenido      string    `json:"contenido"`
	Fecha          time.Time `json:"fecha"`
}

func ConvertirUsuarios(usuarios []*models.Usuario) []UsuarioGuardado {

	resultado := make([]UsuarioGuardado, 0, len(usuarios))

	for _, usuario := range usuarios {

		registro := UsuarioGuardado{
			ID:     usuario.ObtenerID(),
			Nombre: usuario.ObtenerNombre(),
		}

		resultado = append(resultado, registro)
	}

	return resultado
}

func ConvertirMensajes(mensajes []*models.Mensaje) []MensajeGuardado {

	resultado := make([]MensajeGuardado, 0, len(mensajes))

	for _, mensaje := range mensajes {

		registro := MensajeGuardado{
			ID:             mensaje.ObtenerID(),
			RemitenteID:    mensaje.ObtenerRemitente(),
			DestinatarioID: mensaje.ObtenerDestinatario(),
			Contenido:      mensaje.ObtenerContenido(),
			Fecha:          mensaje.ObtenerFecha(),
		}

		resultado = append(resultado, registro)
	}

	return resultado
}

// ReconstruirUsuarios recupera los usuarios utilizando la informacion almacenada en JSON.
func ReconstruirUsuarios(registros []UsuarioGuardado) ([]*models.Usuario, error) {

	usuarios := make([]*models.Usuario, 0, len(registros))
	ids := make(map[int]bool)

	for _, registro := range registros {

		if ids[registro.ID] {
			return nil, fmt.Errorf("ID de usuario duplicado: %d", registro.ID)
		}

		usuario, err := models.NewUsuario(
			registro.ID,
			registro.Nombre,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"error al recuperar usuario %d: %w",
				registro.ID,
				err,
			)
		}

		usuarios = append(usuarios, usuario)
		ids[registro.ID] = true
	}

	return usuarios, nil
}

// ReconstruirMensajes recupera los mensajes utilizando la informacion almacenada en JSON.
func ReconstruirMensajes(registros []MensajeGuardado) ([]*models.Mensaje, error) {

	mensajes := make([]*models.Mensaje, 0, len(registros))
	ids := make(map[int]bool)

	for _, registro := range registros {

		if ids[registro.ID] {
			return nil, fmt.Errorf(
				"ID de mensaje duplicado: %d",
				registro.ID,
			)
		}

		mensaje, err := models.NewMensajeConFecha(
			registro.ID,
			registro.RemitenteID,
			registro.DestinatarioID,
			registro.Contenido,
			registro.Fecha,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"error al recuperar mensaje %d: %w",
				registro.ID,
				err,
			)
		}

		mensajes = append(mensajes, mensaje)
		ids[registro.ID] = true
	}

	return mensajes, nil
}
