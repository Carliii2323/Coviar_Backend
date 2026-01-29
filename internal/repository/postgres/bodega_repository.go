package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
)

type BodegaRepository struct {
	db *sql.DB
}

func NewBodegaRepository(db *sql.DB) repository.BodegaRepository {
	return &BodegaRepository{db: db}
}

func (r *BodegaRepository) Create(ctx context.Context, tx repository.Transaction, bodega *domain.Bodega) (int, error) {
	// Validar restricciones antes de insertar si es necesario
	// - cuit debe ser 11 dígitos
	// - telefono solo números
	// - email_institucional debe contener '@'
	// - numeracion default 'S/N' si está vacío
	if bodega.Numeracion == "" {
		bodega.Numeracion = "S/N"
	}
	query := `
	       INSERT INTO bodegas (razon_social, nombre_fantasia, cuit, inv_bod, inv_vin, calle, numeracion, id_localidad, telefono, email_institucional, fecha_registro)
	       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
	       RETURNING id_bodega
       `

	var id int
	err := r.db.QueryRowContext(ctx, query,
		bodega.RazonSocial, bodega.NombreFantasia, bodega.CUIT, bodega.InvBod, bodega.InvVin,
		bodega.Calle, bodega.Numeracion, bodega.IDLocalidad, bodega.Telefono, bodega.EmailInstitucional,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("error creating bodega: %w", err)
	}

	bodega.ID = id
	bodega.FechaRegistro = time.Now()
	return id, nil
}

func (r *BodegaRepository) FindByID(ctx context.Context, id int) (*domain.Bodega, error) {
	query := `
		SELECT id_bodega, razon_social, nombre_fantasia, cuit, inv_bod, inv_vin, calle, numeracion, id_localidad, telefono, email_institucional, fecha_registro
		FROM bodegas WHERE id_bodega = $1
	`

	bodega := &domain.Bodega{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&bodega.ID, &bodega.RazonSocial, &bodega.NombreFantasia, &bodega.CUIT,
		&bodega.InvBod, &bodega.InvVin, &bodega.Calle, &bodega.Numeracion,
		&bodega.IDLocalidad, &bodega.Telefono, &bodega.EmailInstitucional, &bodega.FechaRegistro,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("error finding bodega: %w", err)
	}

	return bodega, nil
}

func (r *BodegaRepository) FindByCUIT(ctx context.Context, cuit string) (*domain.Bodega, error) {
	query := `
		SELECT id_bodega, razon_social, nombre_fantasia, cuit, inv_bod, inv_vin, calle, numeracion, id_localidad, telefono, email_institucional, fecha_registro
		FROM bodegas WHERE cuit = $1
	`

	bodega := &domain.Bodega{}
	err := r.db.QueryRowContext(ctx, query, cuit).Scan(
		&bodega.ID, &bodega.RazonSocial, &bodega.NombreFantasia, &bodega.CUIT,
		&bodega.InvBod, &bodega.InvVin, &bodega.Calle, &bodega.Numeracion,
		&bodega.IDLocalidad, &bodega.Telefono, &bodega.EmailInstitucional, &bodega.FechaRegistro,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding bodega by CUIT: %w", err)
	}

	return bodega, nil
}

func (r *BodegaRepository) Update(ctx context.Context, tx repository.Transaction, bodega *domain.Bodega) error {
	query := `
		UPDATE bodegas
		SET razon_social = $1, nombre_fantasia = $2, cuit = $3, inv_bod = $4, inv_vin = $5,
		    calle = $6, numeracion = $7, id_localidad = $8, telefono = $9, email_institucional = $10
		WHERE id_bodega = $11
	`

	_, err := r.db.ExecContext(ctx, query,
		bodega.RazonSocial, bodega.NombreFantasia, bodega.CUIT, bodega.InvBod, bodega.InvVin,
		bodega.Calle, bodega.Numeracion, bodega.IDLocalidad, bodega.Telefono, bodega.EmailInstitucional,
		bodega.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating bodega: %w", err)
	}

	return nil
}

func (r *BodegaRepository) Delete(ctx context.Context, tx repository.Transaction, id int) error {
	query := `DELETE FROM bodegas WHERE id_bodega = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting bodega: %w", err)
	}

	return nil
}

func (r *BodegaRepository) GetAll(ctx context.Context) ([]*domain.Bodega, error) {
	query := `
		SELECT id_bodega, razon_social, nombre_fantasia, cuit, inv_bod, inv_vin, calle, numeracion, id_localidad, telefono, email_institucional, fecha_registro
		FROM bodegas ORDER BY id_bodega
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting all bodegas: %w", err)
	}
	defer rows.Close()

	var bodegas []*domain.Bodega
	for rows.Next() {
		bodega := &domain.Bodega{}
		err := rows.Scan(
			&bodega.ID, &bodega.RazonSocial, &bodega.NombreFantasia, &bodega.CUIT,
			&bodega.InvBod, &bodega.InvVin, &bodega.Calle, &bodega.Numeracion,
			&bodega.IDLocalidad, &bodega.Telefono, &bodega.EmailInstitucional, &bodega.FechaRegistro,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning bodega: %w", err)
		}
		bodegas = append(bodegas, bodega)
	}

	return bodegas, rows.Err()
}
