package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
)

type NivelRespuestaRepository struct {
	db *sql.DB
}

func NewNivelRespuestaRepository(db *sql.DB) repository.NivelRespuestaRepository {
	return &NivelRespuestaRepository{db: db}
}

func (r *NivelRespuestaRepository) FindByIndicador(ctx context.Context, idIndicador int) ([]*domain.NivelRespuesta, error) {
	query := `
		SELECT id_nivel_respuesta, id_indicador, nombre, descripcion, puntos
		FROM niveles_respuesta
		WHERE id_indicador = $1
		ORDER BY puntos
	`

	rows, err := r.db.QueryContext(ctx, query, idIndicador)
	if err != nil {
		return nil, fmt.Errorf("error querying niveles_respuesta: %w", err)
	}
	defer rows.Close()

	var niveles []*domain.NivelRespuesta
	for rows.Next() {
		nivel := &domain.NivelRespuesta{}
		if err := rows.Scan(&nivel.ID, &nivel.IDIndicador, &nivel.Nombre, &nivel.Descripcion, &nivel.Puntos); err != nil {
			return nil, fmt.Errorf("error scanning nivel_respuesta: %w", err)
		}
		niveles = append(niveles, nivel)
	}

	return niveles, rows.Err()
}
