package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
)

type ResponsableRepository struct {
	db *sql.DB
}

func NewResponsableRepository(db *sql.DB) repository.ResponsableRepository {
	return &ResponsableRepository{db: db}
}

func (r *ResponsableRepository) Create(ctx context.Context, tx repository.Transaction, responsable *domain.Responsable) (int, error) {
	query := `
	       INSERT INTO responsables (id_cuenta, nombre, apellido, cargo, dni, activo, fecha_registro, fecha_baja)
	       VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7)
	       RETURNING id_responsable
       `

	var id int
	err := r.db.QueryRowContext(ctx, query,
		responsable.IDCuenta, responsable.Nombre, responsable.Apellido, responsable.Cargo, responsable.DNI, responsable.Activo, responsable.FechaBaja,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("error creating responsable: %w", err)
	}

	responsable.ID = id
	responsable.FechaRegistro = time.Now()
	return id, nil
}

func (r *ResponsableRepository) FindByID(ctx context.Context, id int) (*domain.Responsable, error) {
	query := `
	       SELECT id_responsable, id_cuenta, nombre, apellido, cargo, dni, activo, fecha_registro, fecha_baja
	       FROM responsables WHERE id_responsable = $1
       `

	responsable := &domain.Responsable{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&responsable.ID, &responsable.IDCuenta, &responsable.Nombre, &responsable.Apellido,
		&responsable.Cargo, &responsable.DNI, &responsable.Activo, &responsable.FechaRegistro, &responsable.FechaBaja,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("error finding responsable: %w", err)
	}

	return responsable, nil
}

// NOTA: Cambiar FindByBodegaID a FindByCuentaID si corresponde a la nueva lógica de negocio
func (r *ResponsableRepository) FindByCuentaID(ctx context.Context, cuentaID int) ([]*domain.Responsable, error) {
	query := `
	       SELECT id_responsable, id_cuenta, nombre, apellido, cargo, dni, activo, fecha_registro, fecha_baja
	       FROM responsables WHERE id_cuenta = $1 ORDER BY id_responsable
       `

	rows, err := r.db.QueryContext(ctx, query, cuentaID)
	if err != nil {
		return nil, fmt.Errorf("error finding responsables: %w", err)
	}
	defer rows.Close()

	var responsables []*domain.Responsable
	for rows.Next() {
		responsable := &domain.Responsable{}
		err := rows.Scan(
			&responsable.ID, &responsable.IDCuenta, &responsable.Nombre, &responsable.Apellido,
			&responsable.Cargo, &responsable.DNI, &responsable.Activo, &responsable.FechaRegistro, &responsable.FechaBaja,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning responsable: %w", err)
		}
		responsables = append(responsables, responsable)
	}

	return responsables, rows.Err()
}

func (r *ResponsableRepository) Update(ctx context.Context, tx repository.Transaction, responsable *domain.Responsable) error {
	query := `
	       UPDATE responsables
	       SET nombre = $1, apellido = $2, cargo = $3, dni = $4, activo = $5, fecha_baja = $6
	       WHERE id_responsable = $7
       `

	_, err := r.db.ExecContext(ctx, query,
		responsable.Nombre, responsable.Apellido, responsable.Cargo, responsable.DNI, responsable.Activo, responsable.FechaBaja, responsable.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating responsable: %w", err)
	}

	return nil
}

func (r *ResponsableRepository) Delete(ctx context.Context, tx repository.Transaction, id int) error {
	query := `DELETE FROM responsables WHERE id_responsable = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting responsable: %w", err)
	}

	return nil
}
