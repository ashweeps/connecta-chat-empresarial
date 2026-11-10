package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// UsuarioPrueba representa un usuario utilizado para comprobar el almacenamiento JSON.
type UsuarioPrueba struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

// Comprobar que los datos se guarden y recuperen.
func TestGuardarYCargar(t *testing.T) {

	ruta := filepath.Join(t.TempDir(), "usuarios.json")

	archivo := NuevoArchivoJSON(ruta)

	usuarios := []UsuarioPrueba{
		{ID: 1, Nombre: "Ana"},
		{ID: 2, Nombre: "Maria"},
	}

	// Guardar usuarios.
	err := archivo.Guardar(usuarios)

	if err != nil {
		t.Fatalf("Error al guardar usuarios: %v", err)
	}

	// Crear otro administrador para leer el archivo.
	otroArchivo := NuevoArchivoJSON(ruta)

	var recuperados []UsuarioPrueba

	err = otroArchivo.Cargar(&recuperados)

	if err != nil {
		t.Fatalf("Error al cargar usuarios: %v", err)
	}

	// Comparar los datos originales con los recuperados.
	if !reflect.DeepEqual(usuarios, recuperados) {
		t.Errorf("Los datos recuperados son diferentes")
	}
}

// Comprobar que el programa controle un archivo inexistente.
func TestArchivoInexistente(t *testing.T) {

	ruta := filepath.Join(t.TempDir(), "no_existe.json")

	archivo := NuevoArchivoJSON(ruta)

	var usuarios []UsuarioPrueba

	err := archivo.Cargar(&usuarios)

	if err != nil {
		t.Errorf("No deberia producir error: %v", err)
	}
}

// Comprobar que se detecte un JSON incorrecto.
func TestJSONInvalido(t *testing.T) {

	ruta := filepath.Join(t.TempDir(), "datos.json")

	err := os.WriteFile(ruta, []byte("{invalido"), 0600)

	if err != nil {
		t.Fatalf("Error al preparar prueba: %v", err)
	}

	archivo := NuevoArchivoJSON(ruta)

	var usuarios []UsuarioPrueba

	err = archivo.Cargar(&usuarios)

	if err == nil {
		t.Error("Se esperaba un error por JSON incorrecto")
	}
}
