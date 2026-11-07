package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/handlers"
	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/services"
)

// Colores para personalizar la terminal
const (
	azul     = "\033[36m"
	verde    = "\033[32m"
	amarillo = "\033[33m"
	rojo     = "\033[31m"
	blanco   = "\033[0m"
)

// EstadoServidor guarda la informacion del servidor.
type EstadoServidor struct {
	Estado    string `json:"estado"`
	Sistema   string `json:"sistema"`
	Mensaje   string `json:"mensaje"`
	Protocolo string `json:"protocolo"`
	IP        string `json:"ip"`
	Puerto    string `json:"puerto"`
	Inicio    string `json:"inicio"`
}

// Mostrar la bienvenida del servidor.
func mostrarBienvenida(ip string, puerto string) {
	fmt.Println(azul)
	fmt.Println("=======================================")
	fmt.Println("           C O N N E C T A")
	fmt.Println("           CHAT EMPRESARIAL")
	fmt.Println("=======================================")
	fmt.Print(blanco)

	fmt.Println(verde + "[OK] Servidor iniciado correctamente" + blanco)
	fmt.Println("[INFO] Direccion IP:", ip)
	fmt.Println("[INFO] Puerto:", puerto)
	fmt.Println("[INFO] Protocolo: HTTP")
	fmt.Println("[INFO] Estado: ACTIVO")
	fmt.Println("[INFO] Inicio:", time.Now().Format("02/01/2006 - 15:04:05"))
	fmt.Println(azul + "---------------------------------------" + blanco)
	fmt.Println(amarillo + "[ESPERANDO] Conexiones de usuarios..." + blanco)
}

// Verificar que el servidor esta funcionando.
func verificarServidor(estado EstadoServidor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
			return
		}

		fmt.Println(verde + "[INFO] Estado del servidor consultado" + blanco)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if err := json.NewEncoder(w).Encode(estado); err != nil {
			log.Println(rojo+"[ERROR] No se pudo enviar la respuesta:"+blanco, err)
		}
	}
}

// Mostrar el panel visual de CONNECTA.
func mostrarPanel(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println(azul + "[INFO] Panel visual consultado" + blanco)
	http.ServeFile(w, r, "web/panel.html")
}

func main() {
	// Configurar IP y puerto.
	ip := flag.String("ip", "127.0.0.1", "IP del servidor")
	puerto := flag.String("puerto", "8080", "Puerto del servidor")
	flag.Parse()

	// Datos del estado del servidor.
	estado := EstadoServidor{
		Estado:    "activo",
		Sistema:   "CONNECTA Chat Empresarial",
		Mensaje:   "Servidor funcionando correctamente",
		Protocolo: "HTTP",
		IP:        *ip,
		Puerto:    *puerto,
		Inicio:    time.Now().Format("02/01/2006 - 15:04:05"),
	}

	// Servicio de usuarios.
	usuarioService := services.NewUsuarioService()
	usuarioHandler := handlers.NewUsuarioHandler(usuarioService)

	// Servicio de mensajeria.
	// Utiliza el mismo servicio de usuarios para validar IDs.
	chatService := services.NewChatService(usuarioService)
	chatHandler := handlers.NewChatHandler(chatService)

	// Configurar las rutas HTTP.
	mux := http.NewServeMux()

	// Rutas del servidor y panel.
	mux.HandleFunc("/salud", verificarServidor(estado))
	mux.HandleFunc("/", mostrarPanel)

	// Rutas de usuarios.
	mux.HandleFunc("/usuarios/registrar", usuarioHandler.Registrar)
	mux.HandleFunc("/usuarios", usuarioHandler.Listar)
	mux.HandleFunc("/usuarios/buscar", usuarioHandler.Buscar)

	// Rutas de sesiones.
	mux.HandleFunc("/usuarios/iniciar-sesion", usuarioHandler.IniciarSesion)
	mux.HandleFunc("/usuarios/cerrar-sesion", usuarioHandler.CerrarSesion)

	// Rutas de mensajeria.
	mux.HandleFunc("/mensajes/enviar", chatHandler.Enviar)
	mux.HandleFunc("/mensajes/recibir", chatHandler.Recibir)
	mux.HandleFunc("/mensajes/historial", chatHandler.Historial)

	// Iniciar servidor.
	direccion := *ip + ":" + *puerto
	mostrarBienvenida(*ip, *puerto)

	servidor := &http.Server{
		Addr:              direccion,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := servidor.ListenAndServe(); err != nil {
		log.Println(rojo+"[ERROR] No se pudo iniciar el servidor:"+blanco, err)
	}
}
