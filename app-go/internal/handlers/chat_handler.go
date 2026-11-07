package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/models"
	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/services"
)

// ChatHandler recibe las solicitudes HTTP de mensajeria.
type ChatHandler struct {
	servicio *services.ChatService
}

// NewChatHandler crea el controlador de mensajeria.
func NewChatHandler(servicio *services.ChatService) *ChatHandler {
	return &ChatHandler{servicio: servicio}
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

	mensaje, err := h.servicio.EnviarMensaje(
		solicitud.RemitenteID,
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

	id, err := strconv.Atoi(r.URL.Query().Get("usuario_id"))
	if err != nil || id <= 0 {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": "ID invalido"})
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

	usuario1, err1 := strconv.Atoi(r.URL.Query().Get("usuario1"))
	usuario2, err2 := strconv.Atoi(r.URL.Query().Get("usuario2"))

	if err1 != nil || err2 != nil || usuario1 <= 0 || usuario2 <= 0 {
		ResponderJSON(w, http.StatusBadRequest,
			map[string]string{"error": "IDs invalidos"})
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
