# CONNECTA | Chat Empresarial

## 1. Descripción del proyecto

CONNECTA es un prototipo de chat empresarial desarrollado como parte del Proyecto Integrador de tercer semestre de Ingeniería en Sistemas.

El sistema permite la comunicación mediante mensajes de texto entre usuarios registrados, utilizando una arquitectura cliente-servidor desarrollada en Go.

La aplicación busca facilitar el intercambio de mensajes entre usuarios, permitiendo registrar personas, iniciar sesiones, enviar mensajes y consultar conversaciones desde una interfaz de terminal.

También incorpora almacenamiento mediante archivos JSON, lo que permite conservar la información de los usuarios y mensajes después de apagar y volver a iniciar el servidor.

Durante el desarrollo del proyecto aplicamos conceptos de programación en Go, estructuras de datos, constructores, métodos, manejo de errores, comunicación HTTP, almacenamiento de archivos y control de versiones con GitHub.

## 2. Objetivo del proyecto

Desarrollar un prototipo funcional de chat empresarial mediante una arquitectura cliente-servidor, utilizando Go para implementar las principales operaciones de comunicación entre usuarios.

El sistema permite registrar usuarios, iniciar y cerrar sesiones, enviar mensajes y consultar el historial de conversaciones.

Además, se implementó un mecanismo de almacenamiento mediante archivos JSON, con el propósito de conservar los datos y recuperarlos al reiniciar el servidor.

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
- Validación de datos ingresados.
- Manejo básico de errores.
- Almacenamiento de usuarios en JSON.
- Almacenamiento de mensajes en JSON.
- Recuperación de usuarios después de reiniciar el servidor.
- Recuperación de mensajes después de reiniciar el servidor.
- Conservación de identificadores y fechas originales.
- Comunicación entre cliente y servidor mediante HTTP.

Estas funcionalidades fueron desarrolladas y verificadas durante la segunda fase del proyecto.

## 4. Tecnologías y herramientas utilizadas

Para desarrollar CONNECTA utilizamos las siguientes herramientas:

| Tecnología | Uso dentro del proyecto |
|---|---|
| Go | Desarrollo del cliente y servidor. |
| HTTP | Comunicación entre cliente y servidor. |
| JSON | Intercambio y almacenamiento de información. |
| Git | Control de versiones del código. |
| GitHub | Organización del repositorio y registro de contribuciones. |
| GitHub Codespaces | Entorno de desarrollo y pruebas. |
| Visual Studio Code | Edición y organización del código. |
| Terminal Linux | Ejecución de comandos y pruebas. |

El módulo de Go declara la versión 1.27.1 y las pruebas documentadas se realizaron utilizando Go 1.27.3.

Para el almacenamiento implementado no se necesita instalar una base de datos externa, ya que se utilizan archivos JSON.

## 5. Requisitos del sistema

Para ejecutar el prototipo se necesita:

- Go 1.27.1 o una versión compatible posterior.
- Visual Studio Code, GitHub Codespaces u otro entorno compatible.
- Una terminal de comandos.
- Acceso al código fuente del repositorio.
- Un puerto disponible para ejecutar el servidor.
- Conexión entre el cliente y el servidor.

El sistema fue probado en GitHub Codespaces, utilizando un entorno Linux.

## 6. Instalación del proyecto

Para obtener el proyecto desde GitHub, abrir una terminal y ejecutar:

```bash
git clone https://github.com/ashweeps/connecta-chat-empresarial.git
```

Ingresar a la carpeta del repositorio:

```bash
cd connecta-chat-empresarial
```

La implementación de almacenamiento de la Persona 3 se encuentra actualmente en la rama:

```bash
git switch feature/almacenamiento-documentacion
```

Ingresar a la aplicación Go:

```bash
cd app-go
```

Verificar la versión instalada:

```bash
go version
```

Comprobar que el proyecto compile:

```bash
go build ./...
```

Si el comando termina sin mostrar errores, significa que la compilación se realizó correctamente.

Cuando la rama de almacenamiento se integre a `main`, ya no será necesario cambiar a la rama de desarrollo para ejecutar esa funcionalidad.

## 7. Ejecución del servidor

Para iniciar el servidor CONNECTA, abrir una terminal y ubicarse dentro de la carpeta `app-go`.

Ejecutar:

```bash
go run ./cmd/server
```

El servidor utiliza la siguiente configuración predeterminada:

| Configuración | Valor |
|---|---|
| Dirección IP | 127.0.0.1 |
| Puerto | 8080 |
| Protocolo | HTTP |
| Estado | Activo después de iniciarse correctamente |

También se permite configurar la dirección IP y el puerto mediante parámetros.

Ejemplo:

```bash
go run ./cmd/server -ip 127.0.0.1 -puerto 8080
```

Una vez iniciado, el servidor permanece funcionando y esperando solicitudes.

Para detenerlo, se utiliza la combinación de teclas `Ctrl + C`.

## 8. Ejecución del cliente

Para ejecutar el cliente, primero se debe mantener el servidor encendido.

Después, abrir una segunda terminal y ubicarse dentro de la carpeta `app-go`.

Ejecutar:

```bash
go run ./cmd/client
```

El cliente muestra un menú con las siguientes opciones:

```text
================================
       CONNECTA - CLIENTE GO
================================

1. Registrar usuario
2. Iniciar sesion
3. Ver usuarios
4. Enviar mensaje
5. Consultar mensajes
6. Ver historial
7. Cerrar sesion
0. Salir

Seleccione una opcion:
```

El usuario debe ingresar el número correspondiente a la operación que desea realizar.

Por ejemplo, para iniciar sesión se selecciona la opción 2 y se ingresa el identificador del usuario registrado.

Para enviar un mensaje se selecciona la opción 4, se escribe el identificador del destinatario y posteriormente el contenido del mensaje.

## 9. Configuración de IP y puerto

El cliente Go utiliza por defecto la dirección:

```text
http://localhost:8080
```

Esta dirección corresponde al servidor ejecutado en el mismo entorno.

También se permite configurar la dirección mediante la variable de entorno `CONNECTA_SERVER`.

En Linux o GitHub Codespaces:

```bash
export CONNECTA_SERVER=http://127.0.0.1:8080
```

Después se ejecuta el cliente:

```bash
go run ./cmd/client
```

Para que la comunicación funcione, el cliente debe apuntar a la dirección y puerto donde se encuentra ejecutándose el servidor.

En las pruebas realizadas se utilizó la dirección local `127.0.0.1` y el puerto `8080`.

## 10. Arquitectura del sistema

CONNECTA utiliza una arquitectura cliente-servidor organizada en componentes.

El cliente Go permite interactuar con las funcionalidades mediante un menú.

El servidor Go recibe las solicitudes HTTP, valida la información y ejecuta las operaciones correspondientes.

Los servicios administran los usuarios y mensajes, mientras que el módulo de almacenamiento permite guardar y recuperar los registros mediante archivos JSON.

La comunicación general se representa de la siguiente manera:

```text
                USUARIO
                   |
                   v
              CLIENTE GO
                   |
                   | HTTP / JSON
                   v
              SERVIDOR GO
                   |
                   v
           CONTROLADORES HTTP
                   |
                   v
                SERVICIOS
                   |
          ---------------------
          |                   |
          v                   v
    SERVICIO DE         SERVICIO DE
      USUARIOS            MENSAJES
          |                   |
          ---------------------
                   |
                   v
           ALMACENAMIENTO JSON
                   |
          ---------------------
          |                   |
          v                   v
     usuarios.json       mensajes.json
```

Esta organización permite mantener separadas las responsabilidades de cada componente.

## 11. Almacenamiento de datos

Como parte de la segunda fase del proyecto se implementó almacenamiento persistente mediante archivos JSON.

Anteriormente, la información se almacenaba únicamente en memoria, por lo que los usuarios y mensajes desaparecían después de detener el servidor.

Para solucionar esta limitación se desarrolló un módulo de almacenamiento encargado de guardar, cargar, convertir y reconstruir la información.

Los archivos utilizados son:

```text
app-go/data/
    usuarios.json
    mensajes.json
```

### 11.1. Almacenamiento de usuarios

El archivo `usuarios.json` almacena el identificador y nombre de cada usuario.

Ejemplo:

```json
[
  {
    "id": 1,
    "nombre": "Ashley"
  },
  {
    "id": 2,
    "nombre": "Paula"
  }
]
```

Cuando se registra un nuevo usuario, el servidor actualiza el archivo JSON.

Al reiniciar CONNECTA, el servidor recupera los usuarios registrados y conserva sus identificadores.

Las sesiones no se restauran como activas, por lo que cada usuario debe iniciar sesión nuevamente.

### 11.2. Almacenamiento de mensajes

El archivo `mensajes.json` conserva los siguientes datos:

- Identificador del mensaje.
- Identificador del remitente.
- Identificador del destinatario.
- Contenido.
- Fecha de envío.

Ejemplo:

```json
[
  {
    "id": 1,
    "remitente_id": 1,
    "destinatario_id": 2,
    "contenido": "Hola Paula",
    "fecha": "2026-11-10T02:47:09Z"
  }
]
```

El módulo permite recuperar los mensajes al iniciar el servidor, conservando su contenido y fecha original.

También verifica que los identificadores de los usuarios relacionados con cada mensaje existan.

### 11.3. Manejo de errores

El almacenamiento implementa controles para detectar problemas durante la lectura y escritura de archivos.

Entre las situaciones consideradas se encuentran los archivos inexistentes, contenido JSON incorrecto, identificadores duplicados y registros inválidos.

También se controla el fallo de escritura para evitar confirmar operaciones que no se hayan guardado correctamente.

## 12. Estructura del repositorio

El proyecto está organizado de la siguiente manera:

```text
connecta-chat-empresarial/
│
├── app-go/
│   │
│   ├── cmd/
│   │   ├── client/
│   │   │   ├── cliente.go
│   │   │   └── main.go
│   │   │
│   │   └── server/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── handlers/
│   │   │   ├── chat_handler.go
│   │   │   └── usuario_handler.go
│   │   │
│   │   ├── models/
│   │   │   ├── mensaje.go
│   │   │   └── usuario.go
│   │   │
│   │   ├── services/
│   │   │   ├── chat_service.go
│   │   │   ├── chat_persistencia.go
│   │   │   ├── usuario_service.go
│   │   │   └── usuario_persistencia.go
│   │   │
│   │   └── storage/
│   │       ├── archivo_json.go
│   │       ├── archivo_json_test.go
│   │       ├── registros.go
│   │       └── registros_test.go
│   │
│   ├── data/
│   │   ├── usuarios.json
│   │   └── mensajes.json
│   │
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

Los archivos JSON de la carpeta `data` se generan durante la ejecución del servidor y no se incluyen en GitHub.

Los archivos de pruebas adicionales de los modelos y servicios se encuentran junto a sus respectivos archivos Go.

Las carpetas `auth-service` y `database` forman parte de la organización prevista del repositorio, pero todavía no contienen servicios externos funcionales integrados al prototipo.

## 13. Comunicación cliente-servidor

CONNECTA utiliza HTTP para intercambiar información entre el cliente y el servidor.

El cliente envía solicitudes y el servidor procesa las operaciones mediante controladores y servicios.

Las principales rutas implementadas son:

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/salud` | Consultar el estado del servidor. |
| POST | `/usuarios/registrar` | Registrar usuarios. |
| GET | `/usuarios` | Listar usuarios. |
| GET | `/usuarios/buscar` | Buscar usuarios por ID. |
| POST | `/usuarios/iniciar-sesion` | Iniciar sesión. |
| POST | `/usuarios/cerrar-sesion` | Cerrar sesión. |
| POST | `/mensajes/enviar` | Enviar mensajes. |
| GET | `/mensajes/recibir` | Consultar mensajes recibidos. |
| GET | `/mensajes/historial` | Consultar conversaciones. |

Las solicitudes y respuestas utilizan JSON cuando corresponde.

## 14. Datos de prueba

Durante las verificaciones se registraron los siguientes usuarios:

| ID | Nombre |
|---|---|
| 1 | Ashley |
| 2 | Paula |
| 3 | Camila |

También se realizaron pruebas de envío de mensajes desde Ashley hacia Paula.

Uno de los mensajes utilizados fue:

```text
Hola Paula, estamos probando el almacenamiento de CONNECTA
```

Posteriormente se envió otro mensaje desde el cliente Go:

```text
Hola Paula, este mensaje lo envio desde el cliente Go de CONNECTA.
```

Ambos mensajes aparecieron en el historial de conversación.

Estas pruebas permitieron verificar el funcionamiento del cliente, servidor y almacenamiento JSON.

Los datos corresponden al entorno de pruebas y no forman parte de la instalación inicial del proyecto.

## 15. Pruebas realizadas

Durante el desarrollo se ejecutaron pruebas de compilación y funcionamiento.

Para compilar el proyecto:

```bash
go build ./...
```

Para ejecutar las pruebas automáticas:

```bash
go test ./... -count=1
```

Para comprobar específicamente el módulo de almacenamiento:

```bash
go test ./internal/storage -v -count=1
```

Se verificaron las siguientes operaciones:

- Registro y consulta de usuarios.
- Envío de mensajes.
- Consulta del historial.
- Lectura y escritura de archivos JSON.
- Manejo de archivos inexistentes.
- Detección de JSON incorrecto.
- Validación de identificadores duplicados.
- Reconstrucción de usuarios.
- Reconstrucción de mensajes.
- Conservación de fechas originales.
- Persistencia de información después de reiniciar el servidor.
- Comunicación entre el cliente y servidor.
- Inicio y cierre de sesión.

Las pruebas ejecutadas durante la integración finalizaron correctamente.

## 16. Limitaciones actuales

CONNECTA es un prototipo académico y todavía presenta algunas limitaciones.

Actualmente:

- El inicio de sesión utiliza identificadores y no contraseñas.
- No existe autenticación avanzada.
- El servicio de autenticación FastAPI todavía no se encuentra integrado.
- La comunicación utiliza HTTP sin cifrado TLS.
- La aplicación funciona mediante terminal y no cuenta con una interfaz gráfica independiente.
- No se han implementado conversaciones grupales.
- No se permite enviar archivos multimedia.
- El almacenamiento utiliza archivos JSON y no una base de datos relacional.
- No existe un sistema avanzado de permisos.
- No se ha implementado un mecanismo de recuperación automática de archivos dañados.
- El prototipo no está preparado para utilizarse como sistema empresarial de producción.

Estas limitaciones representan aspectos que se podrán mejorar en futuras etapas del proyecto.

## 17. Integrantes y contribuciones

CONNECTA fue desarrollado por tres estudiantes de Ingeniería en Sistemas, quienes distribuyeron las responsabilidades del proyecto.

### Persona 1. Paula Gabriela Simbaña Pérez

Responsable del desarrollo del servidor en Go, configuración de rutas HTTP, controladores y procesamiento de solicitudes.

Su participación permitió establecer el funcionamiento del servidor y la comunicación con los demás componentes del sistema.

### Persona 2. Nadia Camila Lascano Alemán

Responsable del desarrollo del cliente Go, implementación del menú interactivo y comunicación con el servidor.

Su trabajo permitió que los usuarios pudieran registrarse, iniciar sesión, consultar usuarios, enviar mensajes, revisar conversaciones y cerrar sesión desde la terminal.

### Persona 3. Ashley Michelle Espinoza Rosero

Responsable del almacenamiento, integración y documentación del proyecto.

Su participación incluyó la implementación del módulo de almacenamiento mediante archivos JSON, funciones para guardar y recuperar información, conversión y reconstrucción de usuarios y mensajes, manejo de errores y pruebas automatizadas.

También realizó la integración del almacenamiento con los servicios existentes, verificó la persistencia después de reiniciar el servidor y comprobó la comunicación entre el cliente y servidor.

Además, organizó las evidencias de ejecución y participó en la elaboración de la documentación técnica del proyecto.

La preparación del documento técnico y del video demostrativo corresponde al trabajo conjunto de las tres integrantes.

## 18. Registro de contribuciones en GitHub

Para mantener organizadas las modificaciones, cada integrante trabajó en las partes correspondientes a su responsabilidad.

La Persona 3 utilizó la rama:

```text
feature/almacenamiento-documentacion
```

Los commits relacionados con su contribución fueron:

| Commit | Descripción |
|---|---|
| `d85ebca` | Implementar módulo de almacenamiento JSON y pruebas. |
| `a36c102` | Agregar conversión y reconstrucción de usuarios y mensajes. |
| `2a1cbcf` | Integrar persistencia JSON de usuarios y mensajes al servidor. |

Estos commits registran el desarrollo del almacenamiento y su integración con el proyecto existente.

La incorporación de los cambios a la rama principal se realizará mediante el proceso de revisión e integración del repositorio.

## 19. Créditos

El proyecto fue desarrollado utilizando el lenguaje Go y sus bibliotecas estándar.

Se utilizaron GitHub y GitHub Codespaces como herramientas de desarrollo, organización y control de versiones.

El equipo también utilizó materiales académicos y herramientas de apoyo para comprender los conceptos necesarios durante la implementación.

## 20. Licencia

Actualmente, el repositorio no contiene un archivo de licencia específico.

La selección de una licencia de distribución queda pendiente de acuerdo con las decisiones del equipo.

Mientras no se establezca una licencia, se conservan los derechos correspondientes a las autoras del proyecto.

## 21. Declaración del uso de inteligencia artificial

Durante el desarrollo de CONNECTA se utilizó ChatGPT como herramienta de apoyo académico.

Se empleó para comprender conceptos de programación en Go, revisar fragmentos de código, orientar la implementación del almacenamiento JSON, explicar errores y apoyar la organización de la documentación técnica.

Las funcionalidades implementadas se comprobaron mediante compilaciones y pruebas ejecutadas en GitHub Codespaces.

El uso de inteligencia artificial se declara como parte del proceso de apoyo al desarrollo y aprendizaje, manteniendo la responsabilidad de las integrantes sobre la revisión, ejecución y presentación del proyecto.

## 22. Repositorio de GitHub

El proyecto se encuentra disponible en:

https://github.com/ashweeps/connecta-chat-empresarial

Rama de desarrollo del almacenamiento:

```text
feature/almacenamiento-documentacion
```

## 23. Video demostrativo
https://drive.google.com/file/d/1sKFkosUU5ZP26_I583_l2wZUEM4krBM6/view?usp=drivesdk 

El video demostrativo presentará el funcionamiento del servidor y cliente Go, las operaciones principales de CONNECTA y la persistencia de usuarios y mensajes.

Las tres integrantes participarán explicando sus respectivas contribuciones.

Enlace al video: pendiente de incorporar después de la grabación.
