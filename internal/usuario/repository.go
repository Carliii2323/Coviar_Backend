// RUTA: coviar-backend/internal/usuario/repository.go
// REEMPLAZA EL ARCHIVO ANTERIOR
package usuario

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/carli/coviar-backend/internal/domain"
	supa "github.com/supabase-community/supabase-go"
)

// Repository maneja el acceso a datos de Usuario
type Repository struct {
	db *supa.Client
}

// NewRepository crea una nueva instancia del repositorio
func NewRepository(db *supa.Client) *Repository {
	return &Repository{db: db}
}

// Create crea un nuevo usuario
func (r *Repository) Create(usuario *domain.Usuario) error {
	usuarioMap := map[string]interface{}{
		"email":          usuario.Email,
		"password_hash":  usuario.PasswordHash,
		"nombre":         usuario.Nombre,
		"apellido":       usuario.Apellido,
		"telefono":       usuario.Telefono,
		"rol":            usuario.Rol,
		"activo":         true,
		"fecha_registro": time.Now(),
	}

	data, _, err := r.db.From("usuario").
		Insert(usuarioMap, false, "", "", "").
		Execute()

	if err != nil {
		return err
	}

	var result []domain.Usuario
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	if len(result) > 0 {
		*usuario = result[0]
	}

	return nil
}

// FindByEmail busca un usuario por email
func (r *Repository) FindByEmail(email string) (*domain.Usuario, error) {
	data, _, err := r.db.From("usuario").
		Select("*", "", false).
		Eq("email", email).
		Execute()

	if err != nil {
		return nil, err
	}

	var usuarios []domain.Usuario
	if err := json.Unmarshal(data, &usuarios); err != nil {
		return nil, err
	}

	if len(usuarios) == 0 {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	return &usuarios[0], nil
}

// FindByID busca un usuario por ID
func (r *Repository) FindByID(id int) (*domain.Usuario, error) {
	data, _, err := r.db.From("usuario").
		Select("*", "", false).
		Eq("idUsuario", fmt.Sprintf("%d", id)).
		Execute()

	if err != nil {
		return nil, err
	}

	var usuarios []domain.Usuario
	if err := json.Unmarshal(data, &usuarios); err != nil {
		return nil, err
	}

	if len(usuarios) == 0 {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	return &usuarios[0], nil
}

// Deactivate da de baja (desactiva) a un usuario
func (r *Repository) Deactivate(id int) error {
	updateMap := map[string]interface{}{
		"activo": false,
	}

	_, _, err := r.db.From("usuario").
		Update(updateMap, "", "").
		Eq("idUsuario", fmt.Sprintf("%d", id)).
		Execute()

	return err
}

// UpdateLastAccess actualiza la fecha de último acceso
func (r *Repository) UpdateLastAccess(id int) error {
	updateMap := map[string]interface{}{
		"ultimo_acceso": time.Now(),
	}

	_, _, err := r.db.From("usuario").
		Update(updateMap, "", "").
		Eq("idUsuario", fmt.Sprintf("%d", id)).
		Execute()

	return err
}

// FindAll obtiene todos los usuarios activos
func (r *Repository) FindAll() ([]domain.Usuario, error) {
	data, _, err := r.db.From("usuario").
		Select("*", "", false).
		Eq("activo", "true").
		Order("fecha_registro", nil). // Sintaxis correcta para v0.0.4
		Execute()

	if err != nil {
		return nil, err
	}

	var usuarios []domain.Usuario
	if err := json.Unmarshal(data, &usuarios); err != nil {
		return nil, err
	}

	return usuarios, nil
}
