package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Cliente struct {
	servidor  string
	http      *http.Client
	usuarioID int
	token     string
}

func NuevoCliente(servidor string) (*Cliente, error) {
	servidor = strings.TrimRight(strings.TrimSpace(servidor), "/")

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

// crearSolicitud prepara una solicitud HTTP.
// Cuando se requiere autenticacion, incluye el token Bearer.
func (c *Cliente) crearSolicitud(
	metodo string,
	ruta string,
	datos interface{},
	autenticada bool,
) (*http.Request, error) {

	var cuerpo *bytes.Reader

	if datos != nil {
		contenido, err := json.Marshal(datos)
		if err != nil {
			return nil, fmt.Errorf("error al preparar datos: %w", err)
		}
		cuerpo = bytes.NewReader(contenido)
	} else {
		cuerpo = bytes.NewReader(nil)
	}

	solicitud, err := http.NewRequest(
		metodo,
		c.servidor+ruta,
		cuerpo,
	)
	if err != nil {
		return nil, err
	}

	if datos != nil {
		solicitud.Header.Set("Content-Type", "application/json")
	}

	if autenticada {
		if c.usuarioID <= 0 || c.token == "" {
			return nil, errors.New("debe iniciar sesion primero")
		}

		solicitud.Header.Set("Authorization", "Bearer "+c.token)
	}

	return solicitud, nil
}

// ejecutarSolicitud envia una solicitud al servidor.
func (c *Cliente) ejecutarSolicitud(
	metodo string,
	ruta string,
	datos interface{},
	autenticada bool,
) (*http.Response, error) {

	solicitud, err := c.crearSolicitud(
		metodo,
		ruta,
		datos,
		autenticada,
	)
	if err != nil {
		return nil, err
	}

	respuesta, err := c.http.Do(solicitud)
	if err != nil {
		return nil, fmt.Errorf("error de conexion: %w", err)
	}

	return respuesta, nil
}

// ListarUsuarios consulta los usuarios con una sesion valida.
func (c *Cliente) ListarUsuarios() error {
	respuesta, err := c.ejecutarSolicitud(
		http.MethodGet,
		"/usuarios",
		nil,
		true,
	)
	if err != nil {
		return err
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		return fmt.Errorf("error del servidor: %s", respuesta.Status)
	}

	var usuarios []Usuario

	if err := json.NewDecoder(respuesta.Body).Decode(&usuarios); err != nil {
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

// RegistrarUsuario registra un usuario con contraseña.
func (c *Cliente) RegistrarUsuario(nombre string, password string) error {
	nombre = strings.TrimSpace(nombre)

	if nombre == "" {
		return errors.New("el nombre no puede estar vacio")
	}

	if len(password) < 12 {
		return errors.New("la contraseña debe tener al menos 12 caracteres")
	}

	respuesta, err := c.ejecutarSolicitud(
		http.MethodPost,
		"/usuarios/registrar",
		map[string]string{
			"nombre":   nombre,
			"password": password,
		},
		false,
	)
	if err != nil {
		return err
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

	fmt.Println("\nUsuario registrado correctamente!")
	fmt.Printf("ID: %d | Nombre: %s\n", usuario.ID, usuario.Nombre)

	return nil
}

// IniciarSesion valida las credenciales y guarda el token.
func (c *Cliente) IniciarSesion(id int, password string) error {
	if id <= 0 {
		return errors.New("el ID del usuario debe ser mayor que cero")
	}

	respuesta, err := c.ejecutarSolicitud(
		http.MethodPost,
		"/usuarios/iniciar-sesion",
		map[string]interface{}{
			"id":       id,
			"password": password,
		},
		false,
	)
	if err != nil {
		return err
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		return fmt.Errorf("no se pudo iniciar sesion: %s", respuesta.Status)
	}

	var usuario struct {
		ID      int    `json:"id"`
		Nombre  string `json:"nombre"`
		Mensaje string `json:"mensaje"`
		Token   string `json:"token"`
	}

	if err := json.NewDecoder(respuesta.Body).Decode(&usuario); err != nil {
		return fmt.Errorf("error al leer la respuesta: %w", err)
	}

	if usuario.ID != id || usuario.Token == "" {
		return errors.New("el servidor devolvio una sesion invalida")
	}

	c.usuarioID = usuario.ID
	c.token = usuario.Token

	fmt.Printf("\n%s\n", usuario.Mensaje)
	fmt.Printf("Bienvenida, %s! (ID: %d)\n", usuario.Nombre, usuario.ID)

	return nil
}

// EnviarMensaje envia un mensaje con autenticacion.
func (c *Cliente) EnviarMensaje(destinatarioID int, contenido string) error {
	if c.usuarioID <= 0 || c.token == "" {
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

	respuesta, err := c.ejecutarSolicitud(
		http.MethodPost,
		"/mensajes/enviar",
		map[string]interface{}{
			"remitente_id":    c.usuarioID,
			"destinatario_id": destinatarioID,
			"contenido":       contenido,
		},
		true,
	)
	if err != nil {
		return err
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusCreated {
		return fmt.Errorf("no se pudo enviar el mensaje: %s", respuesta.Status)
	}

	fmt.Println("Mensaje enviado correctamente!")
	return nil
}

type Mensaje struct {
	ID             int    `json:"id"`
	RemitenteID    int    `json:"remitente_id"`
	DestinatarioID int    `json:"destinatario_id"`
	Contenido      string `json:"contenido"`
	Fecha          string `json:"fecha"`
}

// ConsultarMensajes consulta los mensajes recibidos.
func (c *Cliente) ConsultarMensajes() error {
	if c.usuarioID <= 0 || c.token == "" {
		return errors.New("debe iniciar sesion para consultar mensajes")
	}

	ruta := "/mensajes/recibir?usuario_id=" +
		strconv.Itoa(c.usuarioID)

	respuesta, err := c.ejecutarSolicitud(
		http.MethodGet,
		ruta,
		nil,
		true,
	)
	if err != nil {
		return err
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		return fmt.Errorf("error al consultar mensajes: %s",
			respuesta.Status)
	}

	var mensajes []Mensaje

	if err := json.NewDecoder(respuesta.Body).Decode(&mensajes); err != nil {
		return fmt.Errorf("error al leer mensajes: %w", err)
	}

	if len(mensajes) == 0 {
		fmt.Println("No tienes mensajes recibidos.")
		return nil
	}

	fmt.Println("\n--- MENSAJES RECIBIDOS ---")

	for _, mensaje := range mensajes {
		fmt.Printf("\nDe usuario: %d\n", mensaje.RemitenteID)
		fmt.Printf("Mensaje: %s\n", mensaje.Contenido)
		fmt.Printf("Fecha: %s\n", mensaje.Fecha)
	}

	return nil
}

// ConsultarHistorial consulta una conversacion autorizada.
func (c *Cliente) ConsultarHistorial(otroUsuarioID int) error {
	if c.usuarioID <= 0 || c.token == "" {
		return errors.New("debe iniciar sesion para consultar el historial")
	}

	if otroUsuarioID <= 0 {
		return errors.New("el ID del usuario no es valido")
	}

	if otroUsuarioID == c.usuarioID {
		return errors.New("debe seleccionar otro usuario")
	}

	ruta := "/mensajes/historial?usuario1=" +
		strconv.Itoa(c.usuarioID) +
		"&usuario2=" + strconv.Itoa(otroUsuarioID)

	respuesta, err := c.ejecutarSolicitud(
		http.MethodGet,
		ruta,
		nil,
		true,
	)
	if err != nil {
		return err
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		return fmt.Errorf("error al consultar historial: %s",
			respuesta.Status)
	}

	var mensajes []Mensaje

	if err := json.NewDecoder(respuesta.Body).Decode(&mensajes); err != nil {
		return fmt.Errorf("error al leer historial: %w", err)
	}

	if len(mensajes) == 0 {
		fmt.Println("No hay mensajes en esta conversacion.")
		return nil
	}

	fmt.Println("\n--- HISTORIAL DE CONVERSACION ---")

	for _, mensaje := range mensajes {
		fmt.Printf("\nDe: %d | Para: %d\n",
			mensaje.RemitenteID, mensaje.DestinatarioID)
		fmt.Printf("Mensaje: %s\n", mensaje.Contenido)
		fmt.Printf("Fecha: %s\n", mensaje.Fecha)
	}

	return nil
}

// CerrarSesion invalida la sesion y limpia el token local.
func (c *Cliente) CerrarSesion() error {
	if c.usuarioID <= 0 || c.token == "" {
		return errors.New("no hay una sesion activa")
	}

	respuesta, err := c.ejecutarSolicitud(
		http.MethodPost,
		"/usuarios/cerrar-sesion",
		map[string]int{
			"id": c.usuarioID,
		},
		true,
	)
	if err != nil {
		return err
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		return fmt.Errorf("no se pudo cerrar sesion: %s", respuesta.Status)
	}

	c.usuarioID = 0
	c.token = ""

	fmt.Println("Sesion cerrada correctamente.")
	return nil
}
