// RUTA: coviar-backend/internal/bodega/service.go
package bodega

import (
	"fmt"

	"github.com/carli/coviar-backend/internal/domain"
)

// Service contiene la lógica de negocio de Bodega
type Service struct {
	repo *Repository
}

// NewService crea una nueva instancia del servicio
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetAll obtiene todas las bodegas
func (s *Service) GetAll() ([]domain.Bodega, error) {
	return s.repo.FindAll()
}

// GetByID obtiene una bodega por ID
func (s *Service) GetByID(id int) (*domain.Bodega, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID inválido")
	}

	return s.repo.FindByID(id)
}

// Create crea una nueva bodega con validaciones
func (s *Service) Create(bodega *domain.Bodega) error {
	// Validaciones de negocio
	if bodega.Nombre == "" {
		return fmt.Errorf("el nombre es requerido")
	}

	if bodega.Cuit <= 0 {
		return fmt.Errorf("CUIT inválido")
	}

	// Aquí podrías agregar más validaciones:
	// - Verificar que el CUIT no exista
	// - Validar formato de email
	// - etc.

	return s.repo.Create(bodega)
}
