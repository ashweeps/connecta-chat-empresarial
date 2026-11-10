package services

import (
	"fmt"
	"sort"

	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/models"
	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/storage"
)

// NewUsuarioServiceConArchivo crea el servicio y recupera los usuarios guardados.
func NewUsuarioServiceConArchivo(ruta string) (*UsuarioService, error) {

	servicio := NewUsuarioService()

	archivo := storage.NuevoArchivoJSON(ruta)

	var registros []storage.UsuarioGuardado

	err := archivo.Cargar(&registros)

	if err != nil {
		return nil, fmt.Errorf(
			"error al cargar usuarios: %w", err,
		)
	}

	usuarios, err := storage.ReconstruirUsuarios(registros)

	if err != nil {
		return nil, fmt.Errorf(
			"error al reconstruir usuarios: %w", err,
		)
	}

	for _, usuario := range usuarios {

		id := usuario.ObtenerID()

		servicio.usuarios[id] = usuario

		if id >= servicio.siguiente {
			servicio.siguiente = id + 1
		}
	}

	servicio.archivo = archivo

	return servicio, nil
}

// guardarUsuariosBloqueados guarda los usuarios mientras el servicio mantiene su bloqueo.
func (s *UsuarioService) guardarUsuariosBloqueados() error {

	if s.archivo == nil {
		return nil
	}

	usuarios := make(
		[]*models.Usuario,
		0,
		len(s.usuarios),
	)

	for _, usuario := range s.usuarios {
		usuarios = append(usuarios, usuario)
	}

	sort.Slice(usuarios, func(i, j int) bool {
		return usuarios[i].ObtenerID() <
			usuarios[j].ObtenerID()
	})

	registros := storage.ConvertirUsuarios(usuarios)

	if err := s.archivo.Guardar(registros); err != nil {
		return fmt.Errorf(
			"error al guardar usuarios: %w", err,
		)
	}

	return nil
}
