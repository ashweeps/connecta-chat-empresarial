
import csv
import os
import time
import urllib.request
import urllib.error
from datetime import datetime, timezone
from pathlib import Path

# ==========================================
# CONNECTA - PROBABILIDAD Y ESTADISTICA
# INSTRUMENTO DE RECOLECCION DE DATOS
# ==========================================

# Direcciones de los servidores
GO_URL = os.getenv(
    "CONNECTA_GO_URL",
    "http://127.0.0.1:8080"
)

FASTAPI_URL = os.getenv(
    "CONNECTA_AUTH_URL",
    "http://127.0.0.1:8000"
)

# Configuracion de las mediciones
REPETICIONES = 25
PAUSA = 0.1
TIMEOUT = 10

# Token de una sesion valida de CONNECTA
TOKEN = os.getenv("CONNECTA_TOKEN", "").strip()

# Carpeta donde se guardaran los resultados
CARPETA_DATOS = Path(__file__).resolve().parents[1] / "data"
ARCHIVO_CSV = CARPETA_DATOS / "mediciones.csv"

# Operaciones reales de CONNECTA
PRUEBAS = [
    {
        "servicio": "Go",
        "operacion": "Estado del servidor",
        "ruta": "/salud",
        "requiere_token": False
    },
    {
        "servicio": "Go",
        "operacion": "Listado de usuarios",
        "ruta": "/usuarios",
        "requiere_token": True
    },
    {
        "servicio": "Go",
        "operacion": "Busqueda de usuario",
        "ruta": "/usuarios/buscar?id=1",
        "requiere_token": True
    },
    {
        "servicio": "Go",
        "operacion": "Historial de mensajes",
        "ruta": "/mensajes/historial?usuario1=1&usuario2=2",
        "requiere_token": True
    },
    {
        "servicio": "FastAPI",
        "operacion": "Estado de autenticacion",
        "ruta": "/salud",
        "requiere_token": False
    }
]


def realizar_solicitud(prueba):

    if prueba["servicio"] == "Go":
        base = GO_URL
    else:
        base = FASTAPI_URL

    url = base.rstrip("/") + prueba["ruta"]

    encabezados = {
        "Accept": "application/json"
    }

    if prueba["requiere_token"]:
        encabezados["Authorization"] = "Bearer " + TOKEN

    solicitud = urllib.request.Request(
        url,
        headers=encabezados,
        method="GET"
    )

    codigo_http = 0
    error = ""

    inicio = time.perf_counter_ns()

    try:
        with urllib.request.urlopen(
            solicitud,
            timeout=TIMEOUT
        ) as respuesta:

            codigo_http = respuesta.status
            respuesta.read()

    except urllib.error.HTTPError as excepcion:
        codigo_http = excepcion.code
        error = "Error HTTP " + str(codigo_http)

    except Exception as excepcion:
        error = type(excepcion).__name__

    fin = time.perf_counter_ns()

    tiempo_ms = (fin - inicio) / 1_000_000

    exito = 1 if 200 <= codigo_http < 300 else 0

    return {
        "fecha_utc": datetime.now(
            timezone.utc
        ).isoformat(),
        "servicio": prueba["servicio"],
        "operacion": prueba["operacion"],
        "ruta": prueba["ruta"],
        "tiempo_ms": round(tiempo_ms, 4),
        "codigo_http": codigo_http,
        "exito": exito,
        "error": error
    }


def main():

    print("=" * 50)
    print("CONNECTA - RECOLECCION ESTADISTICA")
    print("=" * 50)

    if not TOKEN:
        print("ERROR: No existe un token de sesion.")
        print("Primero debemos preparar una sesion valida.")
        return

    CARPETA_DATOS.mkdir(
        parents=True,
        exist_ok=True
    )

    columnas = [
        "fecha_utc",
        "servicio",
        "operacion",
        "ruta",
        "tiempo_ms",
        "codigo_http",
        "exito",
        "error"
    ]

    registros = []

    total = len(PRUEBAS) * REPETICIONES

    print("Operaciones:", len(PRUEBAS))
    print("Repeticiones:", REPETICIONES)
    print("Total de solicitudes:", total)

    for prueba in PRUEBAS:

        print("\nOperacion:", prueba["operacion"])

        for numero in range(1, REPETICIONES + 1):

            resultado = realizar_solicitud(prueba)

            registros.append(resultado)

            print(
                numero,
                "/",
                REPETICIONES,
                "- Tiempo:",
                resultado["tiempo_ms"],
                "ms - HTTP:",
                resultado["codigo_http"]
            )

            time.sleep(PAUSA)

    existe = (
        ARCHIVO_CSV.exists()
        and ARCHIVO_CSV.stat().st_size > 0
    )

    with open(
        ARCHIVO_CSV,
        "a",
        newline="",
        encoding="utf-8"
    ) as archivo:

        escritor = csv.DictWriter(
            archivo,
            fieldnames=columnas
        )

        if not existe:
            escritor.writeheader()

        escritor.writerows(registros)

    exitos = sum(
        registro["exito"]
        for registro in registros
    )

    fallidos = len(registros) - exitos

    print("\n" + "=" * 50)
    print("RESULTADOS DE LA RECOLECCION")
    print("=" * 50)

    print("Total de solicitudes:", len(registros))
    print("Solicitudes exitosas:", exitos)
    print("Solicitudes fallidas:", fallidos)
    print("Archivo generado:", ARCHIVO_CSV)

    print("\nDatos guardados correctamente.")


if __name__ == "__main__":
    main()
