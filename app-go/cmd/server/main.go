package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Colores para personalizar la terminal
const (
	azul     = "\033[36m"
	verde    = "\033[32m"
	amarillo = "\033[33m"
	rojo     = "\033[31m"
	blanco   = "\033[0m"
)

// Funcion para mostrar la bienvenida
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

	fecha := time.Now().Format("02/01/2006 - 15:04:05")
	fmt.Println("[INFO] Inicio:", fecha)

	fmt.Println(azul + "---------------------------------------" + blanco)
	fmt.Println(amarillo + "[ESPERANDO] Conexiones de usuarios..." + blanco)
}

// Funcion que responde a las solicitudes del navegador
func verificarServidor(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println(verde + "[INFO] Se recibio una solicitud" + blanco)

	w.Header().Set("Content-Type", "application/json")

	respuesta := map[string]string{
		"estado":  "activo",
		"sistema": "CONNECTA Chat Empresarial",
		"mensaje": "Servidor funcionando correctamente",
	}

	err := json.NewEncoder(w).Encode(respuesta)

	if err != nil {
		log.Println(rojo + "[ERROR] No se pudo enviar la respuesta" + blanco)
	}
}

func main() {

	// Configuracion del servidor
	ip := flag.String("ip", "127.0.0.1", "IP del servidor")
	puerto := flag.String("puerto", "8080", "Puerto del servidor")
	flag.Parse()

	// Rutas del servidor
	mux := http.NewServeMux()

	mux.HandleFunc("/salud", verificarServidor)

	// Direccion del servidor
	direccion := *ip + ":" + *puerto

	// Iniciar servidor
	mostrarBienvenida(*ip, *puerto)

	err := http.ListenAndServe(direccion, mux)

	if err != nil {
		log.Println(rojo+"[ERROR] No se pudo iniciar el servidor:"+blanco, err)
	}
}
