package models

import (
	"errors"
	"strings"
)

// Usuario representa a una persona que utiliza CONNECTA.
type Usuario struct {
	id     int
	nombre string
	activo bool
}

// NewUsuario crea un usuario nuevo.
func NewUsuario(id int, nombre string) (*Usuario, error) {
	nombre = strings.TrimSpace(nombre)

	// Validar que el ID sea correcto.
	if id <= 0 {
		return nil, errors.New("el ID debe ser mayor que cero")
	}

	// Validar que el nombre no este vacio.
	if nombre == "" {
		return nil, errors.New("el nombre del usuario no puede estar vacio")
	}

	usuario := &Usuario{
		id:     id,
		nombre: nombre,
		activo: false,
	}

	return usuario, nil
}

// ObtenerID devuelve el identificador del usuario.
func (u *Usuario) ObtenerID() int {
	return u.id
}

// ObtenerNombre devuelve el nombre del usuario.
func (u *Usuario) ObtenerNombre() string {
	return u.nombre
}

// EstaActivo indica si el usuario tiene una sesion activa.
func (u *Usuario) EstaActivo() bool {
	return u.activo
}

// IniciarSesion activa la sesion del usuario.
func (u *Usuario) IniciarSesion() {
	u.activo = true
}

// CerrarSesion finaliza la sesion del usuario.
func (u *Usuario) CerrarSesion() {
	u.activo = false
}
