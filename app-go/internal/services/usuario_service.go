package services

import (
	"errors"
	"sync"

	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/models"
	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/storage"
)

// UsuarioService administra los usuarios de CONNECTA.
type UsuarioService struct {
	usuarios  map[int]*models.Usuario
	siguiente int
	mu        sync.Mutex
	archivo   *storage.ArchivoJSON
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

	// Crear el usuario utilizando el modelo existente.
	usuario, err := models.NewUsuario(s.siguiente, nombre)

	if err != nil {
		return nil, err
	}

	// Guardar el usuario en memoria.
	s.usuarios[s.siguiente] = usuario
	s.siguiente++

	// Guardar los usuarios en JSON cuando el almacenamiento
	// persistente se encuentre configurado.
	if err := s.guardarUsuariosBloqueados(); err != nil {

		// Revertir el registro si ocurre un error.
		s.siguiente--

		delete(s.usuarios, usuario.ObtenerID())

		return nil, err
	}

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

// IniciarSesion activa la sesion de un usuario registrado.
func (s *UsuarioService) IniciarSesion(id int) (*models.Usuario, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	usuario, existe := s.usuarios[id]
	if !existe {
		return nil, errors.New("el usuario no existe")
	}

	if usuario.EstaActivo() {
		return nil, errors.New("el usuario ya tiene una sesion activa")
	}

	usuario.IniciarSesion()
	return usuario, nil
}

// CerrarSesion finaliza la sesion de un usuario registrado.
func (s *UsuarioService) CerrarSesion(id int) (*models.Usuario, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	usuario, existe := s.usuarios[id]
	if !existe {
		return nil, errors.New("el usuario no existe")
	}

	if !usuario.EstaActivo() {
		return nil, errors.New("el usuario no tiene una sesion activa")
	}

	usuario.CerrarSesion()
	return usuario, nil
}
