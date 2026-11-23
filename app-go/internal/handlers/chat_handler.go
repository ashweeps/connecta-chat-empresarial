package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/models"
	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/services"
)

// ChatHandler recibe las solicitudes HTTP de mensajeria.
type ChatHandler struct {
	servicio *services.ChatService
	usuarios *services.UsuarioService
}

// NewChatHandler crea el controlador de mensajeria.
func NewChatHandler(
	servicio *services.ChatService,
	usuarios *services.UsuarioService,
) *ChatHandler {
	return &ChatHandler{
		servicio: servicio,
		usuarios: usuarios,
	}
}

// autenticar identifica al usuario mediante su token.
func (h *ChatHandler) autenticar(
	w http.ResponseWriter,
	r *http.Request,
) (int, bool) {
	partes := strings.Fields(r.Header.Get("Authorization"))

	if len(partes) != 2 ||
		!strings.EqualFold(partes[0], "Bearer") {
		ResponderJSON(w, http.StatusUnauthorized,
			map[string]string{"error": "Sesion requerida"})
		return 0, false
	}

	id, err := h.usuarios.ValidarToken(partes[1])
	if err != nil {
		ResponderJSON(w, http.StatusUnauthorized,
			map[string]string{"error": "Sesion invalida"})
		return 0, false
	}

	return id, true
}

// MensajeJSON prepara un mensaje para enviarlo como JSON.
func MensajeJSON(m *models.Mensaje) map[string]interface{} {
	return map[string]interface{}{
		"id":              m.ObtenerID(),
		"remitente_id":    m.ObtenerRemitente(),
		"destinatario_id": m.ObtenerDestinatario(),
		"contenido":       m.ObtenerContenido(),
		"fecha":           m.ObtenerFecha().Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Enviar recibe un mensaje mediante POST.
func (h *ChatHandler) Enviar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ResponderJSON(w, http.StatusMethodNotAllowed,
			map[string]string{"error": "Metodo no permitido"})
		return
	}

	usuarioAutenticado, valido := h.autenticar(w, r)
	if !valido {
		return
	}

	var solicitud struct {
		RemitenteID    int    `json:"remitente_id"`
		DestinatarioID int    `json:"destinatario_id"`
		Contenido      string `json:"contenido"`
	}

	if err := json.NewDecoder(r.Body).Decode(&solicitud); err != nil {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": "JSON incorrecto"})
		return
	}

	if solicitud.RemitenteID != usuarioAutenticado {
		ResponderJSON(w, http.StatusForbidden,
			map[string]string{"error": "No puedes enviar mensajes como otro usuario"})
		return
	}

	mensaje, err := h.servicio.EnviarMensaje(
		usuarioAutenticado,
		solicitud.DestinatarioID,
		solicitud.Contenido,
	)
	if err != nil {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": err.Error()})
		return
	}

	ResponderJSON(w, http.StatusCreated, MensajeJSON(mensaje))
}

// Recibir consulta los mensajes recibidos mediante GET.
func (h *ChatHandler) Recibir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ResponderJSON(w, http.StatusMethodNotAllowed,
			map[string]string{"error": "Metodo no permitido"})
		return
	}

	usuarioAutenticado, valido := h.autenticar(w, r)
	if !valido {
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("usuario_id"))
	if err != nil || id <= 0 {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": "ID invalido"})
		return
	}

	if id != usuarioAutenticado {
		ResponderJSON(w, http.StatusForbidden,
			map[string]string{"error": "No puedes consultar mensajes de otro usuario"})
		return
	}

	mensajes, err := h.servicio.RecibirMensajes(id)
	if err != nil {
		ResponderJSON(w, http.StatusNotFound,
			map[string]string{"error": err.Error()})
		return
	}

	lista := make([]map[string]interface{}, 0, len(mensajes))
	for _, mensaje := range mensajes {
		lista = append(lista, MensajeJSON(mensaje))
	}

	ResponderJSON(w, http.StatusOK, lista)
}

// Historial consulta los mensajes entre dos usuarios mediante GET.
func (h *ChatHandler) Historial(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ResponderJSON(w, http.StatusMethodNotAllowed,
			map[string]string{"error": "Metodo no permitido"})
		return
	}

	usuarioAutenticado, valido := h.autenticar(w, r)
	if !valido {
		return
	}

	usuario1, err1 := strconv.Atoi(r.URL.Query().Get("usuario1"))
	usuario2, err2 := strconv.Atoi(r.URL.Query().Get("usuario2"))

	if err1 != nil || err2 != nil || usuario1 <= 0 || usuario2 <= 0 {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": "IDs invalidos"})
		return
	}

	// Solo los participantes pueden consultar la conversacion.
	if usuarioAutenticado != usuario1 &&
		usuarioAutenticado != usuario2 {
		ResponderJSON(w, http.StatusForbidden,
			map[string]string{"error": "No tienes permiso para consultar esta conversacion"})
		return
	}

	mensajes, err := h.servicio.ObtenerHistorial(usuario1, usuario2)
	if err != nil {
		ResponderJSON(w, http.StatusNotFound,
			map[string]string{"error": err.Error()})
		return
	}

	lista := make([]map[string]interface{}, 0, len(mensajes))
	for _, mensaje := range mensajes {
		lista = append(lista, MensajeJSON(mensaje))
	}

	ResponderJSON(w, http.StatusOK, lista)
}
