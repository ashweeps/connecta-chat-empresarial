package models

import "testing"

// Comprobar que se puede crear un mensaje.
func TestCrearMensaje(t *testing.T) {
	mensaje, err := NewMensaje(1, 1, 2, "Hola Ashley")

	if err != nil {
		t.Fatal("Error al crear mensaje:", err)
	}

	if mensaje.ObtenerID() != 1 {
		t.Error("El ID del mensaje es incorrecto")
	}

	if mensaje.ObtenerRemitente() != 1 {
		t.Error("El remitente es incorrecto")
	}

	if mensaje.ObtenerDestinatario() != 2 {
		t.Error("El destinatario es incorrecto")
	}

	if mensaje.ObtenerContenido() != "Hola Ashley" {
		t.Error("El contenido es incorrecto")
	}

	if mensaje.ObtenerFecha().IsZero() {
		t.Error("El mensaje no tiene fecha")
	}
}

// Comprobar que no se permiten mensajes vacios.
func TestMensajeVacio(t *testing.T) {
	_, err := NewMensaje(2, 1, 2, "   ")

	if err == nil {
		t.Error("Se permitio crear un mensaje vacio")
	}
}

// Comprobar que los IDs sean validos.
func TestUsuariosInvalidos(t *testing.T) {
	_, err := NewMensaje(3, 0, 2, "Hola")

	if err == nil {
		t.Error("Se permitio un remitente invalido")
	}
}

// Comprobar que no se permita enviar mensajes a uno mismo.
func TestMensajeMismoUsuario(t *testing.T) {
	_, err := NewMensaje(4, 1, 1, "Hola")

	if err == nil {
		t.Error("Se permitio enviar un mensaje al mismo usuario")
	}
}
