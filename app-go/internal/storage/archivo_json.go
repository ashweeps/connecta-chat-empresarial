package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ArchivoJSON struct {
	ruta string
}

func NuevoArchivoJSON(ruta string) *ArchivoJSON {
	return &ArchivoJSON{
		ruta: ruta,
	}
}

func (a *ArchivoJSON) Guardar(datos any) error {

	contenido, err := json.MarshalIndent(datos, "", "  ")
	if err != nil {
		return fmt.Errorf("error al convertir datos a JSON: %w", err)
	}

	carpeta := filepath.Dir(a.ruta)

	if err := os.MkdirAll(carpeta, 0755); err != nil {
		return fmt.Errorf("error al crear carpeta: %w", err)
	}

	archivoTemporal, err := os.CreateTemp(carpeta, ".connecta-*.tmp")
	if err != nil {
		return fmt.Errorf("error al crear archivo temporal: %w", err)
	}

	nombreTemporal := archivoTemporal.Name()
	defer os.Remove(nombreTemporal)

	if _, err := archivoTemporal.Write(contenido); err != nil {
		archivoTemporal.Close()
		return fmt.Errorf("error al escribir datos: %w", err)
	}

	if err := archivoTemporal.Close(); err != nil {
		return fmt.Errorf("error al cerrar archivo: %w", err)
	}

	if err := os.Rename(nombreTemporal, a.ruta); err != nil {
		return fmt.Errorf("error al guardar archivo: %w", err)
	}

	return nil
}

// Cargar recupera la informacion guardada.
func (a *ArchivoJSON) Cargar(destino any) error {

	contenido, err := os.ReadFile(a.ruta)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error al leer archivo: %w", err)
	}

	if err := json.Unmarshal(contenido, destino); err != nil {
		return fmt.Errorf("error al interpretar JSON: %w", err)
	}

	return nil
}
