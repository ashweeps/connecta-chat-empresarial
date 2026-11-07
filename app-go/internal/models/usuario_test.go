package models

import "testing"

// Comprobar que se puede crear un usuario
func TestCrearUsuario(t *testing.T) {
	usuario, err := NewUsuario(1, "Camila")

	if err != nil {
		t.Fatal("Error al crear usuario:", err)
	}

	if usuario.ObtenerID() != 1 {
		t.Error("El ID del usuario es incorrecto")
	}

	if usuario.ObtenerNombre() != "Camila" {
		t.Error("El nombre del usuario es incorrecto")
	}
}

// Comprobar que no se acepten nombres vacios
func TestUsuarioSinNombre(t *testing.T) {
	_, err := NewUsuario(2, "   ")

	if err == nil {
		t.Error("Se permitio crear un usuario sin nombre")
	}
}

// Comprobar el inicio y cierre de sesion
func TestSesionUsuario(t *testing.T) {
	usuario, err := NewUsuario(3, "Ashley")

	if err != nil {
		t.Fatal("Error al crear usuario:", err)
	}

	usuario.IniciarSesion()

	if !usuario.EstaActivo() {
		t.Error("El usuario deberia estar activo")
	}

	usuario.CerrarSesion()

	if usuario.EstaActivo() {
		t.Error("El usuario deberia estar inactivo")
	}
}
