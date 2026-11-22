package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"sync"

	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/models"
	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

// UsuarioService administra los usuarios de CONNECTA.
type UsuarioService struct {
	usuarios  map[int]*models.Usuario
	siguiente int
	mu        sync.Mutex
	archivo   *storage.ArchivoJSON

	// Tokens de las sesiones autenticadas.
	sesiones map[string]int
}

// NewUsuarioService crea el servicio de usuarios.
func NewUsuarioService() *UsuarioService {
	return &UsuarioService{
		usuarios:  make(map[int]*models.Usuario),
		siguiente: 1,
		sesiones:  make(map[string]int),
	}
}

// RegistrarUsuario crea y guarda un usuario.
func (s *UsuarioService) RegistrarUsuario(
	nombre string,
	password string,
) (*models.Usuario, error) {

	nombre = strings.TrimSpace(nombre)

	if nombre == "" {
		return nil, errors.New("el nombre no puede estar vacio")
	}

	if len(password) < 12 {
		return nil, errors.New(
			"la contraseña debe tener al menos 12 caracteres",
		)
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, errors.New(
			"no se pudo proteger la contraseña",
		)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	usuario, err := models.NewUsuario(s.siguiente, nombre)
	if err != nil {
		return nil, err
	}

	usuario.EstablecerPasswordHash(string(hash))

	s.usuarios[s.siguiente] = usuario
	s.siguiente++

	if err := s.guardarUsuariosBloqueados(); err != nil {
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
func (s *UsuarioService) IniciarSesion(
	id int,
	password string,
) (*models.Usuario, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	usuario, existe := s.usuarios[id]

	if !existe {
		return nil, errors.New("credenciales incorrectas")
	}

	hash := usuario.ObtenerPasswordHash()

	if hash == "" {
		return nil, errors.New(
			"la cuenta necesita configurar una contraseña",
		)
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	if err != nil {
		return nil, errors.New("credenciales incorrectas")
	}

	if usuario.EstaActivo() {
		return nil, errors.New(
			"el usuario ya tiene una sesion activa",
		)
	}

	usuario.IniciarSesion()

	return usuario, nil
}

// CrearToken genera un token seguro para un usuario autenticado.
func (s *UsuarioService) CrearToken(id int) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	usuario, existe := s.usuarios[id]
	if !existe || !usuario.EstaActivo() {
		return "", errors.New("el usuario no tiene una sesion activa")
	}

	// Generar 32 bytes aleatorios.
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", errors.New("no se pudo generar el token")
	}

	token := hex.EncodeToString(bytes)

	if s.sesiones == nil {
		s.sesiones = make(map[string]int)
	}

	s.sesiones[token] = id

	return token, nil
}

// ValidarToken identifica al usuario de una sesion.
func (s *UsuarioService) ValidarToken(token string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if token == "" {
		return 0, errors.New("token no proporcionado")
	}

	id, existe := s.sesiones[token]
	if !existe {
		return 0, errors.New("token invalido")
	}

	usuario, existe := s.usuarios[id]
	if !existe || !usuario.EstaActivo() {
		return 0, errors.New("sesion no activa")
	}

	return id, nil
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

	// Invalidar los tokens asociados al usuario.
	for token, usuarioID := range s.sesiones {
		if usuarioID == id {
			delete(s.sesiones, token)
		}
	}

	return usuario, nil
}
