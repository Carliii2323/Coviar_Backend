package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
)

type SegmentoRepository struct {
	db *sql.DB
}

func NewSegmentoRepository(db *sql.DB) repository.SegmentoRepository {
	return &SegmentoRepository{db: db}
}

func (r *SegmentoRepository) FindAll(ctx context.Context) ([]*domain.Segmento, error) {
	query := `SELECT id_segmento, nombre, min_turistas, max_turistas FROM segmentos ORDER BY id_segmento`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying segmentos: %w", err)
	}
	defer rows.Close()

	var segmentos []*domain.Segmento
	for rows.Next() {
		seg := &domain.Segmento{}
		if err := rows.Scan(&seg.ID, &seg.Nombre, &seg.MinTuristas, &seg.MaxTuristas); err != nil {
			return nil, fmt.Errorf("error scanning segmento: %w", err)
		}
		segmentos = append(segmentos, seg)
	}

	return segmentos, rows.Err()
}

func (r *SegmentoRepository) FindByID(ctx context.Context, id int) (*domain.Segmento, error) {
	query := `SELECT id_segmento, nombre, min_turistas, max_turistas FROM segmentos WHERE id_segmento = $1`

	seg := &domain.Segmento{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&seg.ID, &seg.Nombre, &seg.MinTuristas, &seg.MaxTuristas)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("error finding segmento: %w", err)
	}

	return seg, nil
}
