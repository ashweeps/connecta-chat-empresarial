package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Colores para la terminal
const (
	azul     = "\033[36m"
	verde    = "\033[32m"
	amarillo = "\033[33m"
	rojo     = "\033[31m"
	blanco   = "\033[0m"
)

// Informacion del servidor para el panel
type EstadoServidor struct {
	Estado    string `json:"estado"`
	Sistema   string `json:"sistema"`
	Mensaje   string `json:"mensaje"`
	Protocolo string `json:"protocolo"`
	IP        string `json:"ip"`
	Puerto    string `json:"puerto"`
	Inicio    string `json:"inicio"`
}

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

// Comprobar que el servidor funciona
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

// Mostrar el panel de control
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
	// Configuracion de IP y puerto
	ip := flag.String("ip", "127.0.0.1", "IP del servidor")
	puerto := flag.String("puerto", "8080", "Puerto del servidor")
	flag.Parse()

	// Datos que se mostraran en el panel
	estado := EstadoServidor{
		Estado:    "activo",
		Sistema:   "CONNECTA Chat Empresarial",
		Mensaje:   "Servidor funcionando correctamente",
		Protocolo: "HTTP",
		IP:        *ip,
		Puerto:    *puerto,
		Inicio:    time.Now().Format("02/01/2006 - 15:04:05"),
	}

	// Rutas del servidor
	mux := http.NewServeMux()
	mux.HandleFunc("/salud", verificarServidor(estado))
	mux.HandleFunc("/", mostrarPanel)

	// Iniciar el servidor
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
