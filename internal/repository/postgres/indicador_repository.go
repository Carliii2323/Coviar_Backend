package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
)

type IndicadorRepository struct {
	db *sql.DB
}

func NewIndicadorRepository(db *sql.DB) repository.IndicadorRepository {
	return &IndicadorRepository{db: db}
}

func (r *IndicadorRepository) FindByCapitulo(ctx context.Context, idCapitulo int) ([]*domain.Indicador, error) {
	query := `
		SELECT id_indicador, id_capitulo, nombre, descripcion, orden
		FROM indicadores
		WHERE id_capitulo = $1
		ORDER BY orden
	`

	rows, err := r.db.QueryContext(ctx, query, idCapitulo)
	if err != nil {
		return nil, fmt.Errorf("error querying indicadores: %w", err)
	}
	defer rows.Close()

	var indicadores []*domain.Indicador
	for rows.Next() {
		ind := &domain.Indicador{}
		if err := rows.Scan(&ind.ID, &ind.IDCapitulo, &ind.Nombre, &ind.Descripcion, &ind.Orden); err != nil {
			return nil, fmt.Errorf("error scanning indicador: %w", err)
		}
		indicadores = append(indicadores, ind)
	}

	return indicadores, rows.Err()
}

func (r *IndicadorRepository) FindBySegmento(ctx context.Context, idSegmento int) ([]int, error) {
	query := `SELECT id_indicador FROM segmento_indicador WHERE id_segmento = $1`

	rows, err := r.db.QueryContext(ctx, query, idSegmento)
	if err != nil {
		return nil, fmt.Errorf("error querying segmento_indicador: %w", err)
	}
	defer rows.Close()

	var indicadorIds []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error scanning indicador id: %w", err)
		}
		indicadorIds = append(indicadorIds, id)
	}

	return indicadorIds, rows.Err()
}
