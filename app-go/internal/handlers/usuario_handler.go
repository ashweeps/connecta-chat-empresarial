package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/services"
)

// UsuarioHandler recibe las solicitudes HTTP de usuarios.
type UsuarioHandler struct {
	servicio *services.UsuarioService
}

// NewUsuarioHandler crea el controlador de usuarios.
func NewUsuarioHandler(servicio *services.UsuarioService) *UsuarioHandler {
	return &UsuarioHandler{servicio: servicio}
}

// ResponderJSON envia una respuesta en formato JSON.
func ResponderJSON(w http.ResponseWriter, codigo int, datos interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(codigo)
	_ = json.NewEncoder(w).Encode(datos)
}

// autenticar comprueba que la solicitud tenga una sesion valida.
func (h *UsuarioHandler) autenticar(
	w http.ResponseWriter,
	r *http.Request,
) (int, bool) {

	autorizacion := strings.Fields(r.Header.Get("Authorization"))

	if len(autorizacion) != 2 ||
		!strings.EqualFold(autorizacion[0], "Bearer") {

		ResponderJSON(w, http.StatusUnauthorized,
			map[string]string{"error": "Sesion requerida"})
		return 0, false
	}

	usuarioID, err := h.servicio.ValidarToken(autorizacion[1])
	if err != nil {
		ResponderJSON(w, http.StatusUnauthorized,
			map[string]string{"error": "Sesion invalida"})
		return 0, false
	}

	return usuarioID, true
}

// Registrar permite crear usuarios mediante POST.
func (h *UsuarioHandler) Registrar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ResponderJSON(w, http.StatusMethodNotAllowed,
			map[string]string{"error": "Metodo no permitido"})
		return
	}

	var solicitud struct {
		Nombre   string `json:"nombre"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&solicitud); err != nil {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": "JSON incorrecto"})
		return
	}

	if solicitud.Nombre == "" || solicitud.Password == "" {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": "Nombre y contraseña son obligatorios"})
		return
	}

	usuario, err := h.servicio.RegistrarUsuario(
		solicitud.Nombre,
		solicitud.Password,
	)
	if err != nil {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": err.Error()})
		return
	}

	ResponderJSON(w, http.StatusCreated, map[string]interface{}{
		"id":     usuario.ObtenerID(),
		"nombre": usuario.ObtenerNombre(),
		"activo": usuario.EstaActivo(),
	})
}

// Listar permite consultar usuarios solo con sesion valida.
func (h *UsuarioHandler) Listar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ResponderJSON(w, http.StatusMethodNotAllowed,
			map[string]string{"error": "Metodo no permitido"})
		return
	}

	if _, valido := h.autenticar(w, r); !valido {
		return
	}

	usuarios := h.servicio.ListarUsuarios()
	lista := make([]map[string]interface{}, 0, len(usuarios))

	for _, usuario := range usuarios {
		lista = append(lista, map[string]interface{}{
			"id":     usuario.ObtenerID(),
			"nombre": usuario.ObtenerNombre(),
			"activo": usuario.EstaActivo(),
		})
	}

	ResponderJSON(w, http.StatusOK, lista)
}

// Buscar permite consultar un usuario por ID con sesion valida.
func (h *UsuarioHandler) Buscar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ResponderJSON(w, http.StatusMethodNotAllowed,
			map[string]string{"error": "Metodo no permitido"})
		return
	}

	if _, valido := h.autenticar(w, r); !valido {
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": "ID invalido"})
		return
	}

	usuario, err := h.servicio.BuscarUsuario(id)
	if err != nil {
		ResponderJSON(w, http.StatusNotFound,
			map[string]string{"error": "Usuario no encontrado"})
		return
	}

	ResponderJSON(w, http.StatusOK, map[string]interface{}{
		"id":     usuario.ObtenerID(),
		"nombre": usuario.ObtenerNombre(),
		"activo": usuario.EstaActivo(),
	})
}

// IniciarSesion autentica mediante el servicio de usuarios,
// que verifica las credenciales con FastAPI.
func (h *UsuarioHandler) IniciarSesion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ResponderJSON(w, http.StatusMethodNotAllowed,
			map[string]string{"error": "Metodo no permitido"})
		return
	}

	var solicitud struct {
		ID       int    `json:"id"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&solicitud); err != nil {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": "JSON incorrecto"})
		return
	}

	if solicitud.ID <= 0 || solicitud.Password == "" {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": "ID y contraseña son obligatorios"})
		return
	}

	usuario, err := h.servicio.IniciarSesion(
		solicitud.ID,
		solicitud.Password,
	)
	if err != nil {
		ResponderJSON(w, http.StatusUnauthorized,
			map[string]string{
				"error": "No se pudo iniciar sesion: credenciales incorrectas o sesion activa",
			})
		return
	}

	// Crear token para la sesion autenticada.
	token, err := h.servicio.CrearToken(usuario.ObtenerID())
	if err != nil {
		_, _ = h.servicio.CerrarSesion(usuario.ObtenerID())

		ResponderJSON(w, http.StatusInternalServerError,
			map[string]string{"error": "No se pudo crear la sesion segura"})
		return
	}

	ResponderJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Credenciales verificadas correctamente",
		"id":      usuario.ObtenerID(),
		"nombre":  usuario.ObtenerNombre(),
		"activo":  usuario.EstaActivo(),
		"token":   token,
	})
}

// CerrarSesion finaliza la sesion autenticada.
func (h *UsuarioHandler) CerrarSesion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ResponderJSON(w, http.StatusMethodNotAllowed,
			map[string]string{"error": "Metodo no permitido"})
		return
	}

	usuarioID, valido := h.autenticar(w, r)
	if !valido {
		return
	}

	var solicitud struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&solicitud); err != nil {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": "JSON incorrecto"})
		return
	}

	if solicitud.ID <= 0 {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": "ID invalido"})
		return
	}

	if solicitud.ID != usuarioID {
		ResponderJSON(w, http.StatusForbidden,
			map[string]string{"error": "No tienes permiso para cerrar esta sesion"})
		return
	}

	usuario, err := h.servicio.CerrarSesion(usuarioID)
	if err != nil {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": err.Error()})
		return
	}

	ResponderJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Sesion cerrada correctamente",
		"id":      usuario.ObtenerID(),
		"nombre":  usuario.ObtenerNombre(),
		"activo":  usuario.EstaActivo(),
	})
}
