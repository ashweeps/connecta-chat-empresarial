package main

import "fmt"

func main() {

	cliente, err := NuevoCliente("http://localhost:8080")

	if err != nil {
		fmt.Println("Error al configurar el cliente:", err)
		return
	}

	fmt.Println("Servidor configurado:", cliente.servidor)

	var opcion int

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

		_, err := fmt.Scanln(&opcion)

		if err != nil {
			fmt.Println("Error: ingrese un numero valido.")
			var entrada string
			fmt.Scanln(&entrada)
			continue
		}

		switch opcion {
		case 1:
			fmt.Println("Registrar usuario")
		case 2:
			fmt.Println("Iniciar sesion")
		case 3:
			fmt.Println("Ver usuarios")
		case 4:
			fmt.Println("Enviar mensaje")
		case 5:
			fmt.Println("Consultar mensajes")
		case 6:
			fmt.Println("Ver historial")
		case 7:
			fmt.Println("Cerrar sesion")
		case 0:
			fmt.Println("Gracias por usar CONNECTA.")
			return
		default:
			fmt.Println("Opcion no valida.")
		}
	}
}
