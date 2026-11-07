package services

import (
	"errors"
	"sync"

	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/models"
)

// UsuarioService administra los usuarios de CONNECTA.
type UsuarioService struct {
	usuarios  map[int]*models.Usuario
	siguiente int
	mu        sync.Mutex
}

// NewUsuarioService crea el servicio de usuarios.
func NewUsuarioService() *UsuarioService {
	return &UsuarioService{
		usuarios:  make(map[int]*models.Usuario),
		siguiente: 1,
	}
}

// RegistrarUsuario crea y guarda un usuario.
func (s *UsuarioService) RegistrarUsuario(nombre string) (*models.Usuario, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	usuario, err := models.NewUsuario(s.siguiente, nombre)

	if err != nil {
		return nil, err
	}

	s.usuarios[s.siguiente] = usuario
	s.siguiente++

	return usuario, nil
}

// BuscarUsuario busca un usuario por su ID.
func (s *UsuarioService) BuscarUsuario(id int) (*models.Usuario, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	usuario, existe := s.usuarios[id]

	if !existe {
		return nil, errors.New("el usuario no existe")
	}

	return usuario, nil
}

// ListarUsuarios devuelve todos los usuarios registrados.
func (s *UsuarioService) ListarUsuarios() []*models.Usuario {
	s.mu.Lock()
	defer s.mu.Unlock()

	lista := make([]*models.Usuario, 0, len(s.usuarios))

	for _, usuario := range s.usuarios {
		lista = append(lista, usuario)
	}

	return lista
}
