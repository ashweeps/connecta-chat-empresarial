package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Cliente struct {
	servidor  string
	http      *http.Client
	usuarioID int
}

func NuevoCliente(servidor string) (*Cliente, error) {
	servidor = strings.TrimSpace(servidor)

	if servidor == "" {
		return nil, errors.New("la direccion del servidor es obligatoria")
	}

	cliente := &Cliente{
		servidor: servidor,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	return cliente, nil

}

type Usuario struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Activo bool   `json:"activo"`
}

func (c *Cliente) ListarUsuarios() error {
	respuesta, err := c.http.Get(c.servidor + "/usuarios")
	if err != nil {
		return fmt.Errorf("no se pudo conectar al servidor: %w", err)
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		return fmt.Errorf("error del servidor: %s", respuesta.Status)
	}

	var usuarios []Usuario

	err = json.NewDecoder(respuesta.Body).Decode(&usuarios)
	if err != nil {
		return fmt.Errorf("error al leer los usuarios: %w", err)
	}

	if len(usuarios) == 0 {
		fmt.Println("No hay usuarios registrados.")
		return nil
	}

	fmt.Println("\n--- USUARIOS REGISTRADOS ---")

	for _, usuario := range usuarios {
		fmt.Printf("ID: %d | Nombre: %s\n",
			usuario.ID, usuario.Nombre)
	}

	return nil
}

func (c *Cliente) RegistrarUsuario(nombre string) error {
	nombre = strings.TrimSpace(nombre)

	if nombre == "" {
		return errors.New("el nombre no puede estar vacio")
	}

	datos, err := json.Marshal(map[string]string{
		"nombre": nombre,
	})
	if err != nil {
		return fmt.Errorf("error al preparar los datos: %w", err)
	}

	respuesta, err := c.http.Post(
		c.servidor+"/usuarios/registrar",
		"application/json",
		bytes.NewReader(datos),
	)
	if err != nil {
		return fmt.Errorf("no se pudo conectar al servidor: %w", err)
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusCreated {
		return fmt.Errorf("el servidor rechazo el registro: %s",
			respuesta.Status)
	}

	var usuario Usuario

	if err := json.NewDecoder(respuesta.Body).Decode(&usuario); err != nil {
		return fmt.Errorf("error al leer la respuesta: %w", err)
	}

	fmt.Printf("\nUsuario registrado correctamente!\n")
	fmt.Printf("ID: %d | Nombre: %s\n", usuario.ID, usuario.Nombre)

	return nil
}

func (c *Cliente) IniciarSesion(id int) error {
	if id <= 0 {
		return errors.New("el ID del usuario debe ser mayor que cero")
	}

	datos, err := json.Marshal(map[string]int{
		"id": id,
	})
	if err != nil {
		return fmt.Errorf("error al preparar la solicitud: %w", err)
	}

	respuesta, err := c.http.Post(
		c.servidor+"/usuarios/iniciar-sesion",
		"application/json",
		bytes.NewReader(datos),
	)
	if err != nil {
		return fmt.Errorf("no se pudo conectar al servidor: %w", err)
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		return fmt.Errorf("no se pudo iniciar sesion: %s", respuesta.Status)
	}

	var usuario struct {
		ID      int    `json:"id"`
		Nombre  string `json:"nombre"`
		Mensaje string `json:"mensaje"`
	}

	if err := json.NewDecoder(respuesta.Body).Decode(&usuario); err != nil {
		return fmt.Errorf("error al leer la respuesta: %w", err)
	}

	c.usuarioID = usuario.ID

	fmt.Printf("\n%s\n", usuario.Mensaje)
	fmt.Printf("Bienvenida, %s! (ID: %d)\n", usuario.Nombre, usuario.ID)

	return nil
}

func (c *Cliente) EnviarMensaje(destinatarioID int, contenido string) error {
	if c.usuarioID <= 0 {
		return errors.New("debe iniciar sesion antes de enviar mensajes")
	}

	if destinatarioID <= 0 {
		return errors.New("el ID del destinatario no es valido")
	}

	if destinatarioID == c.usuarioID {
		return errors.New("no puede enviarse mensajes a si mismo")
	}

	contenido = strings.TrimSpace(contenido)

	if contenido == "" {
		return errors.New("el mensaje no puede estar vacio")
	}

	datos, err := json.Marshal(map[string]interface{}{
		"remitente_id":    c.usuarioID,
		"destinatario_id": destinatarioID,
		"contenido":       contenido,
	})
	if err != nil {
		return fmt.Errorf("error al preparar el mensaje: %w", err)
	}

	respuesta, err := c.http.Post(
		c.servidor+"/mensajes/enviar",
		"application/json",
		bytes.NewReader(datos),
	)
	if err != nil {
		return fmt.Errorf("error de conexion: %w", err)
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusCreated {
		return fmt.Errorf("no se pudo enviar el mensaje: %s", respuesta.Status)
	}

	fmt.Println("Mensaje enviado correctamente!")
	return nil
}
