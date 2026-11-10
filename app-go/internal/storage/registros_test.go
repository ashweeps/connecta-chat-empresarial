package storage

import (
	"testing"
	"time"
)

func TestReconstruirUsuarios(t *testing.T) {

	registros := []UsuarioGuardado{
		{ID: 1, Nombre: "Ashley"},
		{ID: 2, Nombre: "Pau"},
	}

	usuarios, err := ReconstruirUsuarios(registros)

	if err != nil {
		t.Fatalf("Error al recuperar usuarios: %v", err)
	}

	if len(usuarios) != 2 {
		t.Fatal("La cantidad de usuarios es incorrecta")
	}

	if usuarios[0].ObtenerNombre() != "Ashley" {
		t.Error("El primer usuario no coincide")
	}

	if usuarios[1].ObtenerID() != 2 {
		t.Error("El ID del segundo usuario es incorrecto")
	}

	if usuarios[0].EstaActivo() {
		t.Error("El usuario recuperado no debe tener sesion activa")
	}
}

func TestUsuariosDuplicados(t *testing.T) {

	registros := []UsuarioGuardado{
		{ID: 1, Nombre: "Ashley"},
		{ID: 1, Nombre: "Pau"},
	}

	_, err := ReconstruirUsuarios(registros)

	if err == nil {
		t.Error("Se esperaba un error por ID duplicado")
	}
}

func TestUsuarioInvalido(t *testing.T) {

	registros := []UsuarioGuardado{
		{ID: 1, Nombre: ""},
	}

	_, err := ReconstruirUsuarios(registros)

	if err == nil {
		t.Error("Se esperaba un error por nombre vacio")
	}
}

// Comprueba que los mensajes recuperados
// conserven sus datos y fecha original.
func TestReconstruirMensajes(t *testing.T) {

	fechaOriginal := time.Date(
		2026, 11, 9, 10, 30, 0, 0, time.UTC,
	)

	registros := []MensajeGuardado{
		{
			ID:             1,
			RemitenteID:    1,
			DestinatarioID: 2,
			Contenido:      "Hola Ashley",
			Fecha:          fechaOriginal,
		},
	}

	mensajes, err := ReconstruirMensajes(registros)

	if err != nil {
		t.Fatalf("Error al recuperar mensajes: %v", err)
	}

	if len(mensajes) != 1 {
		t.Fatal("La cantidad de mensajes es incorrecta")
	}

	if mensajes[0].ObtenerContenido() != "Hola Ashley" {
		t.Error("El contenido recuperado es incorrecto")
	}

	if mensajes[0].ObtenerRemitente() != 1 {
		t.Error("El remitente recuperado es incorrecto")
	}

	if mensajes[0].ObtenerDestinatario() != 2 {
		t.Error("El destinatario recuperado es incorrecto")
	}

	if !mensajes[0].ObtenerFecha().Equal(fechaOriginal) {
		t.Error("La fecha original no se conservo")
	}
}

// Comprueba que se detecten mensajes
// con identificadores repetidos.
func TestMensajesDuplicados(t *testing.T) {

	fecha := time.Now()

	registros := []MensajeGuardado{
		{
			ID:             1,
			RemitenteID:    1,
			DestinatarioID: 2,
			Contenido:      "Primer mensaje",
			Fecha:          fecha,
		},
		{
			ID:             1,
			RemitenteID:    2,
			DestinatarioID: 1,
			Contenido:      "Segundo mensaje",
			Fecha:          fecha,
		},
	}

	_, err := ReconstruirMensajes(registros)

	if err == nil {
		t.Error("Se esperaba un error por ID duplicado")
	}
}

// Comprueba que se rechacen mensajes sin una fecha valida.
func TestMensajeSinFecha(t *testing.T) {

	registros := []MensajeGuardado{
		{
			ID:             1,
			RemitenteID:    1,
			DestinatarioID: 2,
			Contenido:      "Mensaje de prueba",
		},
	}

	_, err := ReconstruirMensajes(registros)

	if err == nil {
		t.Error("Se esperaba un error por fecha vacia")
	}
}
