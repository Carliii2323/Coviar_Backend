// RUTA: coviar-backend/internal/bodega/repository.go
package bodega

import (
	"encoding/json"
	"fmt"

	"github.com/carli/coviar-backend/internal/domain"
	supa "github.com/supabase-community/supabase-go"
)

// Repository maneja el acceso a datos de Bodega
type Repository struct {
	db *supa.Client
}

// NewRepository crea una nueva instancia del repositorio
func NewRepository(db *supa.Client) *Repository {
	return &Repository{db: db}
}

// FindAll obtiene todas las bodegas
func (r *Repository) FindAll() ([]domain.Bodega, error) {
	data, _, err := r.db.From("bodega").
		Select("*", "", false).
		Order("nombre", nil).
		Execute()

	if err != nil {
		return nil, err
	}

	var bodegas []domain.Bodega
	if err := json.Unmarshal(data, &bodegas); err != nil {
		return nil, err
	}

	return bodegas, nil
}

// FindByID obtiene una bodega por ID
func (r *Repository) FindByID(id int) (*domain.Bodega, error) {
	data, _, err := r.db.From("bodega").
		Select("*", "", false).
		Eq("idBodega", fmt.Sprintf("%d", id)).
		Execute()

	if err != nil {
		return nil, err
	}

	var bodegas []domain.Bodega
	if err := json.Unmarshal(data, &bodegas); err != nil {
		return nil, err
	}

	if len(bodegas) == 0 {
		return nil, fmt.Errorf("bodega no encontrada")
	}

	return &bodegas[0], nil
}

// Create crea una nueva bodega
func (r *Repository) Create(bodega *domain.Bodega) error {
	bodegaMap := map[string]interface{}{
		"cuit":            bodega.Cuit,
		"inv":             bodega.Inv,
		"viñedos_inv":     bodega.ViñedosInv,
		"nombre":          bodega.Nombre,
		"ubicacion":       bodega.Ubicacion,
		"contacto_email":  bodega.ContactoEmail,
		"razon_social":    bodega.RazonSocial,
		"nombre_fantasia": bodega.NombreFantasia,
	}

	data, _, err := r.db.From("bodega").
		Insert(bodegaMap, false, "", "", "").
		Execute()

	if err != nil {
		return err
	}

	var result []domain.Bodega
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	if len(result) > 0 {
		*bodega = result[0]
	}

	return nil
}
