# CONNECTA | Chat Empresarial

## 1. Descripción del proyecto

CONNECTA es un prototipo de chat empresarial desarrollado como parte del Proyecto Integrador de tercer semestre de Ingeniería en Sistemas.

El sistema permite la comunicación mediante mensajes de texto entre usuarios registrados, utilizando una arquitectura cliente-servidor desarrollada en Go.

El proyecto busca aplicar conocimientos de programación, comunicación mediante HTTP, estructuras de datos, manejo de archivos y organización de código.

Además, se implementó almacenamiento mediante archivos JSON para conservar los usuarios y mensajes después de reiniciar el servidor.

## 2. Objetivo

Desarrollar un prototipo funcional de chat empresarial que permita registrar usuarios, iniciar y cerrar sesiones, enviar mensajes y consultar conversaciones mediante una arquitectura cliente-servidor.

El sistema también busca mantener la información organizada y conservar los registros de usuarios y mensajes mediante archivos JSON.

## 3. Funcionalidades implementadas

CONNECTA cuenta con las siguientes funcionalidades:

- Registro de usuarios.
- Inicio de sesión mediante identificador.
- Consulta de usuarios registrados.
- Búsqueda de usuarios por identificador.
- Envío de mensajes entre usuarios.
- Consulta de mensajes recibidos.
- Consulta del historial de conversaciones.
- Cierre de sesión.
- Almacenamiento de usuarios en archivos JSON.
- Almacenamiento de mensajes en archivos JSON.
- Recuperación de información después de reiniciar el servidor.
- Validación de información y manejo básico de errores.

## 4. Tecnologías utilizadas

El prototipo utiliza las siguientes tecnologías:

- Go, para el desarrollo del cliente y servidor.
- HTTP, para la comunicación cliente-servidor.
- JSON, para el intercambio y almacenamiento de información.
- Git, para el control de versiones.
- GitHub, para organizar el repositorio y registrar contribuciones.
- Visual Studio Code y GitHub Codespaces, como herramientas de desarrollo.

El módulo Go declara la versión 1.27.1 y las pruebas documentadas se realizaron con Go 1.27.3.

## 5. Requisitos de ejecución

Para ejecutar CONNECTA se necesita:

- Go 1.27.1 o una versión compatible posterior.
- Una terminal de comandos.
- Acceso al código fuente del proyecto.
- Un entorno compatible con Go, como Linux o Windows.
- Un puerto disponible para ejecutar el servidor.

No se necesita instalar una base de datos externa para el almacenamiento JSON implementado.

## 6. Instalación

Primero, descargar o clonar el repositorio:

```bash
git clone https://github.com/ashweeps/connecta-chat-empresarial.git
```

Ingresar a la carpeta de la aplicación:

```bash
cd connecta-chat-empresarial/app-go
```

Comprobar la versión instalada de Go:

```bash
go version
```

Compilar el proyecto:

```bash
go build ./...
```

## 7. Ejecución del servidor

Desde la carpeta `app-go`, ejecutar:

```bash
go run ./cmd/server
```

Por defecto, el servidor utiliza:

- IP: 127.0.0.1
- Puerto: 8080
- Protocolo: HTTP

También se permite configurar la IP y el puerto:

```bash
go run ./cmd/server -ip 127.0.0.1 -puerto 8080
```

Cuando el servidor se inicia correctamente, muestra su estado en la terminal y queda preparado para recibir solicitudes.

## 8. Ejecución del cliente

Abrir una segunda terminal y ubicarse dentro de `app-go`.

Ejecutar:

```bash
go run ./cmd/client
```

El cliente muestra un menú con las operaciones disponibles.

Para utilizar las funcionalidades, se debe mantener el servidor funcionando.

## 9. Configuración de la conexión

El cliente utiliza por defecto la dirección:

```text
http://localhost:8080
```

También permite configurar la dirección del servidor mediante la variable de entorno `CONNECTA_SERVER`.

En Linux o GitHub Codespaces:

```bash
export CONNECTA_SERVER=http://127.0.0.1:8080
go run ./cmd/client
```

El servidor y el cliente deben utilizar la misma dirección IP y el mismo puerto para comunicarse correctamente.

## 10. Almacenamiento de información

CONNECTA utiliza archivos JSON para conservar la información de usuarios y mensajes.

Al ejecutar el servidor desde `app-go`, los archivos se almacenan en:

```text
app-go/data/
├── usuarios.json
└── mensajes.json
```

El archivo `usuarios.json` conserva los identificadores y nombres de los usuarios registrados.

El archivo `mensajes.json` conserva los identificadores, remitentes, destinatarios, contenidos y fechas de los mensajes.

Los registros se recuperan al iniciar nuevamente el servidor.

Los archivos de datos locales se excluyen de Git mediante `.gitignore`.

## 11. Estructura del proyecto

```text
connecta-chat-empresarial/
│
├── app-go/
│   ├── cmd/
│   │   ├── client/
│   │   └── server/
│   │
│   ├── internal/
│   │   ├── handlers/
│   │   ├── models/
│   │   ├── services/
│   │   └── storage/
│   │
│   ├── data/
│   ├── go.mod
│   └── README.md
│
├── auth-service/
├── database/
├── docs/
├── tests/
├── .gitignore
└── README.md
```

La carpeta `data` se genera durante la ejecución local cuando se guardan registros.

El servidor y el cliente se encuentran separados, mientras que la lógica de usuarios, mensajes y almacenamiento se organiza mediante paquetes.

Las carpetas `auth-service` y `database` corresponden a componentes previstos en el proyecto, pero no representan servicios externos funcionales en este prototipo.

## 12. Comunicación cliente-servidor

CONNECTA utiliza solicitudes HTTP para intercambiar información entre el cliente Go y el servidor.

El cliente envía una solicitud, el servidor procesa la operación y devuelve una respuesta.

Las principales rutas son:

| Método | Ruta | Función |
|---|---|---|
| GET | `/salud` | Consultar el estado del servidor. |
| POST | `/usuarios/registrar` | Registrar usuarios. |
| GET | `/usuarios` | Consultar usuarios. |
| GET | `/usuarios/buscar` | Buscar un usuario. |
| POST | `/usuarios/iniciar-sesion` | Iniciar sesión. |
| POST | `/usuarios/cerrar-sesion` | Cerrar sesión. |
| POST | `/mensajes/enviar` | Enviar mensajes. |
| GET | `/mensajes/recibir` | Consultar mensajes recibidos. |
| GET | `/mensajes/historial` | Consultar conversaciones. |

Las solicitudes y respuestas utilizan JSON cuando corresponde.

## 13. Datos de prueba

Durante las pruebas del prototipo se registraron tres usuarios:

| ID | Usuario |
|---|---|
| 1 | Ashley |
| 2 | Paula |
| 3 | Camila |

También se realizaron pruebas de envío de mensajes entre Ashley y Paula, consulta del historial y recuperación de información después de reiniciar el servidor.

Estos datos corresponden al entorno de pruebas y no se incluyen como registros obligatorios para nuevas instalaciones.

## 14. Pruebas realizadas

Para ejecutar las pruebas automáticas:

```bash
go test ./... -count=1
```

Para comprobar la compilación:

```bash
go build ./...
```

Durante el desarrollo se verificaron:

- Validación de usuarios y mensajes.
- Escritura y lectura de archivos JSON.
- Manejo de archivos inexistentes.
- Detección de información JSON incorrecta.
- Reconstrucción de usuarios y mensajes.
- Detección de identificadores duplicados.
- Conservación de fechas originales.
- Recuperación de información después de reiniciar el servidor.
- Comunicación entre el cliente y el servidor.

Las pruebas documentadas finalizaron correctamente.

## 15. Limitaciones actuales

CONNECTA es un prototipo académico.

Actualmente presenta las siguientes limitaciones:

- El inicio de sesión se realiza mediante identificador, sin autenticación mediante contraseña.
- No se ha integrado el servicio FastAPI de autenticación.
- La comunicación utiliza HTTP sin cifrado TLS.
- No existe una interfaz gráfica independiente.
- La comunicación se realiza mediante un cliente de terminal.
- No se ha implementado mensajería grupal.
- El almacenamiento se realiza mediante archivos JSON y no mediante una base de datos relacional.
- No se ha implementado un sistema avanzado de permisos.
- No está preparado para su implementación en un entorno empresarial de producción.

Estas funcionalidades representan posibles mejoras para futuras etapas.

## 16. Autoras

El proyecto CONNECTA fue desarrollado por:

- Ashley Michelle Espinoza Rosero.
- Nadia Camila Lascano Alemán.
- Paula Gabriela Simbaña Pérez.

## 17. Contribuciones

El equipo distribuyó las responsabilidades del proyecto de la siguiente manera:

**Persona 1. Servidor**

Desarrollo del servidor Go, configuración de rutas HTTP, controladores y procesamiento de solicitudes.

**Persona 2. Cliente**

Desarrollo del cliente Go, menú interactivo, comunicación con el servidor y operaciones disponibles para los usuarios.

**Persona 3. Almacenamiento, integración y documentación**

Ashley Michelle Espinoza Rosero.

Implementación del almacenamiento JSON, funciones de conversión y reconstrucción de registros, integración con los servicios existentes, pruebas de persistencia y organización de las evidencias técnicas.

La documentación y demostración del proyecto corresponden al trabajo conjunto del equipo.

## 18. Créditos

El desarrollo se realizó utilizando Go y sus bibliotecas estándar.

Se utilizaron GitHub y GitHub Codespaces para organizar el código, comprobar su funcionamiento y mantener un registro de contribuciones.

## 19. Licencia

El repositorio todavía no contiene un archivo de licencia definido.

Hasta que el equipo seleccione una licencia, el proyecto conserva los derechos correspondientes a sus autoras.

## 20. Declaración del uso de inteligencia artificial

Durante el desarrollo de CONNECTA se utilizó ChatGPT como herramienta de apoyo para comprender conceptos de programación, revisar fragmentos de código, orientar la implementación del almacenamiento JSON y organizar la documentación técnica.

Las funciones implementadas se verificaron mediante pruebas ejecutadas en GitHub Codespaces.

El uso de inteligencia artificial se declara como parte del proceso de apoyo académico y no sustituye la revisión y participación de las integrantes.

## 21. Repositorio

Código fuente del proyecto:

https://github.com/ashweeps/connecta-chat-empresarial

## 22. Video demostrativo
