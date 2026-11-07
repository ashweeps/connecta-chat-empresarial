package services

import "testing"

// Comprobar que se puede enviar un mensaje.
func TestEnviarMensaje(t *testing.T) {
	usuarios := NewUsuarioService()
	camila, _ := usuarios.RegistrarUsuario("Camila")
	ashley, _ := usuarios.RegistrarUsuario("Ashley")

	chat := NewChatService(usuarios)

	mensaje, err := chat.EnviarMensaje(
		camila.ObtenerID(),
		ashley.ObtenerID(),
		"Hola Ashley",
	)

	if err != nil {
		t.Fatal("Error al enviar mensaje:", err)
	}

	if mensaje.ObtenerContenido() != "Hola Ashley" {
		t.Error("El contenido del mensaje es incorrecto")
	}
}

// Comprobar que no se permiten mensajes vacios.
func TestEnviarMensajeVacio(t *testing.T) {
	usuarios := NewUsuarioService()
	usuarios.RegistrarUsuario("Camila")
	usuarios.RegistrarUsuario("Ashley")

	chat := NewChatService(usuarios)

	_, err := chat.EnviarMensaje(1, 2, "   ")

	if err == nil {
		t.Error("Se permitio enviar un mensaje vacio")
	}
}

// Comprobar que no se puede enviar a un usuario inexistente.
func TestDestinatarioInexistente(t *testing.T) {
	usuarios := NewUsuarioService()
	usuarios.RegistrarUsuario("Camila")

	chat := NewChatService(usuarios)

	_, err := chat.EnviarMensaje(1, 99, "Hola")

	if err == nil {
		t.Error("Se permitio enviar a un usuario inexistente")
	}
}

// Comprobar que se pueden recibir mensajes.
func TestRecibirMensajes(t *testing.T) {
	usuarios := NewUsuarioService()
	usuarios.RegistrarUsuario("Camila")
	usuarios.RegistrarUsuario("Ashley")

	chat := NewChatService(usuarios)
	chat.EnviarMensaje(1, 2, "Hola Ashley")
	chat.EnviarMensaje(1, 2, "Como estas?")

	mensajes, err := chat.RecibirMensajes(2)

	if err != nil {
		t.Fatal("Error al recibir mensajes:", err)
	}

	if len(mensajes) != 2 {
		t.Errorf("Se esperaban 2 mensajes, llegaron %d", len(mensajes))
	}
}

// Comprobar el historial de una conversacion.
func TestHistorialConversacion(t *testing.T) {
	usuarios := NewUsuarioService()
	usuarios.RegistrarUsuario("Camila")
	usuarios.RegistrarUsuario("Ashley")

	chat := NewChatService(usuarios)
	chat.EnviarMensaje(1, 2, "Hola Ashley")
	chat.EnviarMensaje(2, 1, "Hola Camila")

	historial, err := chat.ObtenerHistorial(1, 2)

	if err != nil {
		t.Fatal("Error al consultar historial:", err)
	}

	if len(historial) != 2 {
		t.Errorf("Se esperaban 2 mensajes, se encontraron %d", len(historial))
	}
}
