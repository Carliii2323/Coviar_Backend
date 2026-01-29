package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
)

type AutoevaluacionRepository struct {
	db *sql.DB
}

func NewAutoevaluacionRepository(db *sql.DB) repository.AutoevaluacionRepository {
	return &AutoevaluacionRepository{db: db}
}

func (r *AutoevaluacionRepository) Create(ctx context.Context, tx repository.Transaction, auto *domain.Autoevaluacion) (int, error) {
	query := `
		INSERT INTO autoevaluaciones (fecha_inicio, estado, id_bodega)
		VALUES (NOW(), $1, $2)
		RETURNING id_autoevaluacion
	`

	var id int
	err := r.db.QueryRowContext(ctx, query, domain.EstadoPendiente, auto.IDBodega).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating autoevaluacion: %w", err)
	}

	auto.ID = id
	auto.FechaInicio = time.Now()
	auto.Estado = domain.EstadoPendiente
	return id, nil
}

func (r *AutoevaluacionRepository) FindByID(ctx context.Context, id int) (*domain.Autoevaluacion, error) {
	query := `
		SELECT id_autoevaluacion, fecha_inicio, fecha_fin, estado, id_bodega, id_segmento
		FROM autoevaluaciones WHERE id_autoevaluacion = $1
	`

	auto := &domain.Autoevaluacion{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&auto.ID, &auto.FechaInicio, &auto.FechaFin, &auto.Estado, &auto.IDBodega, &auto.IDSegmento,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	return auto, nil
}

func (r *AutoevaluacionRepository) UpdateSegmento(ctx context.Context, id int, idSegmento int) error {
	query := `UPDATE autoevaluaciones SET id_segmento = $1 WHERE id_autoevaluacion = $2`

	_, err := r.db.ExecContext(ctx, query, idSegmento, id)
	if err != nil {
		return fmt.Errorf("error updating segmento: %w", err)
	}

	return nil
}

func (r *AutoevaluacionRepository) Complete(ctx context.Context, id int) error {
	query := `UPDATE autoevaluaciones SET estado = $1, fecha_fin = NOW() WHERE id_autoevaluacion = $2`

	_, err := r.db.ExecContext(ctx, query, domain.EstadoCompletada, id)
	if err != nil {
		return fmt.Errorf("error completing autoevaluacion: %w", err)
	}

	return nil
}
