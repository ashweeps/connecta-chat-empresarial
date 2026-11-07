package services

import "testing"

// Comprobar que se puede registrar un usuario.
func TestRegistrarUsuario(t *testing.T) {
	servicio := NewUsuarioService()

	usuario, err := servicio.RegistrarUsuario("Camila")

	if err != nil {
		t.Fatal("Error al registrar usuario:", err)
	}

	if usuario.ObtenerID() != 1 {
		t.Error("El ID deberia ser 1")
	}

	if usuario.ObtenerNombre() != "Camila" {
		t.Error("El nombre es incorrecto")
	}
}

// Comprobar que se puede buscar un usuario.
func TestBuscarUsuario(t *testing.T) {
	servicio := NewUsuarioService()

	servicio.RegistrarUsuario("Ashley")

	usuario, err := servicio.BuscarUsuario(1)

	if err != nil {
		t.Fatal("No se encontro el usuario:", err)
	}

	if usuario.ObtenerNombre() != "Ashley" {
		t.Error("Se encontro un usuario incorrecto")
	}
}

// Comprobar que se pueden listar usuarios.
func TestListarUsuarios(t *testing.T) {
	servicio := NewUsuarioService()

	servicio.RegistrarUsuario("Camila")
	servicio.RegistrarUsuario("Ashley")
	servicio.RegistrarUsuario("Paula")

	usuarios := servicio.ListarUsuarios()

	if len(usuarios) != 3 {
		t.Errorf("Se esperaban 3 usuarios, se encontraron %d", len(usuarios))
	}
}

// Comprobar que no se acepten nombres vacios.
func TestRegistrarUsuarioVacio(t *testing.T) {
	servicio := NewUsuarioService()

	_, err := servicio.RegistrarUsuario("   ")

	if err == nil {
		t.Error("Se permitio registrar un usuario sin nombre")
	}
}

// Comprobar que se detecten usuarios inexistentes.
func TestBuscarUsuarioInexistente(t *testing.T) {
	servicio := NewUsuarioService()

	_, err := servicio.BuscarUsuario(99)

	if err == nil {
		t.Error("Se encontro un usuario que no existe")
	}
}
