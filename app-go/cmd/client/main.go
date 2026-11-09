package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func leerTexto(lector *bufio.Reader) string {
	texto, err := lector.ReadString('\n')

	if err != nil && len(texto) == 0 {
		return ""
	}

	return strings.TrimSpace(texto)
}

func leerNumero(lector *bufio.Reader) (int, error) {
	texto := leerTexto(lector)
	return strconv.Atoi(texto)
}

func main() {

	servidor := strings.TrimSpace(os.Getenv("CONNECTA_SERVER"))

	if servidor == "" {
		servidor = "http://localhost:8080"
	}

	cliente, err := NuevoCliente(servidor)

	if err != nil {
		fmt.Println("Error al configurar el cliente:", err)
		return
	}

	fmt.Println("Servidor configurado:", cliente.servidor)

	lector := bufio.NewReader(os.Stdin)

	for {

		fmt.Println("\n================================")
		fmt.Println("       CONNECTA - CLIENTE GO")
		fmt.Println("================================")
		fmt.Println("1. Registrar usuario")
		fmt.Println("2. Iniciar sesion")
		fmt.Println("3. Ver usuarios")
		fmt.Println("4. Enviar mensaje")
		fmt.Println("5. Consultar mensajes")
		fmt.Println("6. Ver historial")
		fmt.Println("7. Cerrar sesion")
		fmt.Println("0. Salir")
		fmt.Print("Seleccione una opcion: ")

		opcion, err := leerNumero(lector)

		if err != nil {
			fmt.Println("Error: ingrese un numero valido.")
			continue
		}

		switch opcion {

		case 1:
			fmt.Print("Ingrese el nombre del usuario: ")

			nombre := leerTexto(lector)

			if nombre == "" {
				fmt.Println("Error: el nombre no puede estar vacio.")
				continue
			}

			if err := cliente.RegistrarUsuario(nombre); err != nil {
				fmt.Println("Error:", err)
			}

		case 2:
			fmt.Print("Ingrese su ID de usuario: ")

			id, err := leerNumero(lector)

			if err != nil || id <= 0 {
				fmt.Println("Error: ingrese un ID numerico valido.")
				continue
			}

			if err := cliente.IniciarSesion(id); err != nil {
				fmt.Println("Error:", err)
			}

		case 3:
			if err := cliente.ListarUsuarios(); err != nil {
				fmt.Println("Error:", err)
			}

		case 4:
			if cliente.usuarioID <= 0 {
				fmt.Println("Error: debe iniciar sesion antes de enviar mensajes.")
				continue
			}

			fmt.Print("Ingrese el ID del destinatario: ")

			destinatarioID, err := leerNumero(lector)

			if err != nil || destinatarioID <= 0 {
				fmt.Println("Error: ingrese un ID numerico valido.")
				continue
			}

			if destinatarioID == cliente.usuarioID {
				fmt.Println("Error: no puede enviarse mensajes a si mismo.")
				continue
			}

			fmt.Print("Escriba su mensaje: ")

			contenido := leerTexto(lector)

			if contenido == "" {
				fmt.Println("Error: el mensaje no puede estar vacio.")
				continue
			}

			if err := cliente.EnviarMensaje(destinatarioID, contenido); err != nil {
				fmt.Println("Error:", err)
			}

		case 5:
			if cliente.usuarioID <= 0 {
				fmt.Println("Error: debe iniciar sesion para consultar mensajes.")
				continue
			}

			if err := cliente.ConsultarMensajes(); err != nil {
				fmt.Println("Error:", err)
			}

		case 6:
			if cliente.usuarioID <= 0 {
				fmt.Println("Error: debe iniciar sesion para consultar el historial.")
				continue
			}

			fmt.Print("Ingrese el ID del usuario para consultar el historial: ")

			otroUsuarioID, err := leerNumero(lector)

			if err != nil || otroUsuarioID <= 0 {
				fmt.Println("Error: ingrese un ID numerico valido.")
				continue
			}

			if otroUsuarioID == cliente.usuarioID {
				fmt.Println("Error: debe seleccionar otro usuario.")
				continue
			}

			if err := cliente.ConsultarHistorial(otroUsuarioID); err != nil {
				fmt.Println("Error:", err)
			}

		case 7:
			if cliente.usuarioID <= 0 {
				fmt.Println("No hay una sesion activa.")
				continue
			}

			if err := cliente.CerrarSesion(); err != nil {
				fmt.Println("Error:", err)
			}

		case 0:
			if cliente.usuarioID > 0 {
				fmt.Println("Cerrando sesion activa...")

				if err := cliente.CerrarSesion(); err != nil {
					fmt.Println("Error al cerrar sesion:", err)
					fmt.Println("No se pudo confirmar el cierre en el servidor.")
					continue
				}
			}

			fmt.Println("Gracias por usar CONNECTA.")
			return

		default:
			fmt.Println("Opcion no valida.")
		}
	}
}
