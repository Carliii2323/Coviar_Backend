// RUTA: coviar-backend/internal/usuario/service.go
package usuario

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/carli/coviar-backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// Service contiene la lógica de negocio de Usuario
type Service struct {
	repo *Repository
}

// NewService crea una nueva instancia del servicio
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create crea un nuevo usuario con validaciones
func (s *Service) Create(dto *domain.UsuarioDTO) (*domain.Usuario, error) {
	// Validar email
	if !isValidEmail(dto.Email) {
		return nil, fmt.Errorf("email inválido")
	}

	// Validar que el email no exista
	existing, _ := s.repo.FindByEmail(dto.Email)
	if existing != nil {
		return nil, fmt.Errorf("el email ya está registrado")
	}

	// Validar password
	if len(dto.Password) < 6 {
		return nil, fmt.Errorf("la contraseña debe tener al menos 6 caracteres")
	}

	// Validar nombre y apellido
	if strings.TrimSpace(dto.Nombre) == "" {
		return nil, fmt.Errorf("el nombre es requerido")
	}
	if strings.TrimSpace(dto.Apellido) == "" {
		return nil, fmt.Errorf("el apellido es requerido")
	}

	// Validar rol
	validRoles := map[string]bool{"admin": true, "bodega": true, "auditor": true}
	if !validRoles[dto.Rol] {
		dto.Rol = "bodega" // Rol por defecto
	}

	// Hash de la contraseña
	hashedPassword, err := hashPassword(dto.Password)
	if err != nil {
		return nil, fmt.Errorf("error al procesar contraseña")
	}

	// Crear usuario
	usuario := &domain.Usuario{
		Email:        strings.ToLower(strings.TrimSpace(dto.Email)),
		PasswordHash: hashedPassword,
		Nombre:       strings.TrimSpace(dto.Nombre),
		Apellido:     strings.TrimSpace(dto.Apellido),
		Telefono:     dto.Telefono,
		Rol:          dto.Rol,
		Activo:       true,
	}

	if err := s.repo.Create(usuario); err != nil {
		return nil, err
	}

	return usuario, nil
}

// Verify verifica las credenciales de un usuario
func (s *Service) Verify(login *domain.UsuarioLogin) (*domain.Usuario, error) {
	if login.Email == "" || login.Password == "" {
		return nil, fmt.Errorf("email y contraseña son requeridos")
	}

	usuario, err := s.repo.FindByEmail(strings.ToLower(login.Email))
	if err != nil {
		return nil, fmt.Errorf("credenciales inválidas")
	}

	// Verificar que esté activo
	if !usuario.Activo {
		return nil, fmt.Errorf("usuario desactivado")
	}

	// Verificar contraseña
	if !checkPasswordHash(login.Password, usuario.PasswordHash) {
		return nil, fmt.Errorf("credenciales inválidas")
	}

	// Actualizar último acceso
	s.repo.UpdateLastAccess(usuario.IdUsuario)

	return usuario, nil
}

// GetByID obtiene un usuario por ID
func (s *Service) GetByID(id int) (*domain.Usuario, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID inválido")
	}

	return s.repo.FindByID(id)
}

// Deactivate da de baja a un usuario
func (s *Service) Deactivate(id int) error {
	if id <= 0 {
		return fmt.Errorf("ID inválido")
	}

	// Verificar que existe
	usuario, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("usuario no encontrado")
	}

	if !usuario.Activo {
		return fmt.Errorf("el usuario ya está desactivado")
	}

	return s.repo.Deactivate(id)
}

// GetAll obtiene todos los usuarios activos
func (s *Service) GetAll() ([]domain.Usuario, error) {
	return s.repo.FindAll()
}

// Utilidades

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
