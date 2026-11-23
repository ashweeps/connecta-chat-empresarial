package services

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/models"
	"github.com/ashweeps/connecta-chat-empresarial/app-go/internal/storage"
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

// Cliente HTTP para comunicarse con FastAPI.
var clienteAuth = &http.Client{
	Timeout: 5 * time.Second,
}

// URL del servicio FastAPI.
func obtenerURLAuth() string {
	url := os.Getenv("CONNECTA_AUTH_URL")

	if url == "" {
		return "http://127.0.0.1:8000"
	}

	return strings.TrimRight(url, "/")
}

// Solicita a FastAPI un hash Argon2id.
func generarHashFastAPI(password string) (string, error) {
	datos, err := json.Marshal(map[string]string{
		"password": password,
	})
	if err != nil {
		return "", errors.New("error al preparar credenciales")
	}

	respuesta, err := clienteAuth.Post(
		obtenerURLAuth()+"/auth/hash",
		"application/json",
		bytes.NewReader(datos),
	)
	if err != nil {
		return "", errors.New("servicio de autenticacion no disponible")
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		return "", errors.New("no se pudo proteger la contraseña")
	}

	var resultado struct {
		PasswordHash string `json:"password_hash"`
	}

	if err := json.NewDecoder(respuesta.Body).Decode(&resultado); err != nil {
		return "", errors.New("respuesta de autenticacion invalida")
	}

	if !strings.HasPrefix(resultado.PasswordHash, "$argon2id$") {
		return "", errors.New("hash de contraseña invalido")
	}

	return resultado.PasswordHash, nil
}

// Solicita a FastAPI verificar una contraseña.
func verificarPasswordFastAPI(password, hash string) error {
	datos, err := json.Marshal(map[string]string{
		"password":      password,
		"password_hash": hash,
	})
	if err != nil {
		return errors.New("error al preparar credenciales")
	}

	respuesta, err := clienteAuth.Post(
		obtenerURLAuth()+"/auth/verificar",
		"application/json",
		bytes.NewReader(datos),
	)
	if err != nil {
		return errors.New("servicio de autenticacion no disponible")
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode == http.StatusUnauthorized {
		return errors.New("credenciales incorrectas")
	}

	if respuesta.StatusCode != http.StatusOK {
		return errors.New("error del servicio de autenticacion")
	}

	var resultado struct {
		Valido bool `json:"valido"`
	}

	if err := json.NewDecoder(respuesta.Body).Decode(&resultado); err != nil {
		return errors.New("respuesta de autenticacion invalida")
	}

	if !resultado.Valido {
		return errors.New("credenciales incorrectas")
	}

	return nil
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

	// FastAPI genera el hash Argon2id.
	hash, err := generarHashFastAPI(password)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	usuario, err := models.NewUsuario(s.siguiente, nombre)
	if err != nil {
		return nil, err
	}

	usuario.EstablecerPasswordHash(hash)

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

// IniciarSesion verifica las credenciales mediante FastAPI.
func (s *UsuarioService) IniciarSesion(
	id int,
	password string,
) (*models.Usuario, error) {

	s.mu.Lock()

	usuario, existe := s.usuarios[id]
	if !existe {
		s.mu.Unlock()
		return nil, errors.New("credenciales incorrectas")
	}

	hash := usuario.ObtenerPasswordHash()

	s.mu.Unlock()

	if hash == "" {
		return nil, errors.New(
			"la cuenta necesita configurar una contraseña",
		)
	}

	// Verificacion de contraseña con FastAPI.
	if err := verificarPasswordFastAPI(password, hash); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	usuario, existe = s.usuarios[id]
	if !existe {
		return nil, errors.New("credenciales incorrectas")
	}

	if usuario.ObtenerPasswordHash() != hash {
		return nil, errors.New("credenciales modificadas")
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

	datos := make([]byte, 32)

	if _, err := rand.Read(datos); err != nil {
		return "", errors.New("no se pudo generar el token")
	}

	token := hex.EncodeToString(datos)

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

	for token, usuarioID := range s.sesiones {
		if usuarioID == id {
			delete(s.sesiones, token)
		}
	}

	return usuario, nil
}
