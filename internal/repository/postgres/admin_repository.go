package postgres

import (
	"context"
	"database/sql"
	"strconv"
	"time"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// AdminStats representa las estadísticas del sistema
type AdminStats struct {
	TotalBodegas             int     `json:"totalBodegas"`
	EvaluacionesCompletadas  int     `json:"evaluacionesCompletadas"`
	PromedioSostenibilidad   float64 `json:"promedioSostenibilidad"`
}

// GetStats obtiene las estadísticas generales del sistema
func (r *AdminRepository) GetStats(ctx context.Context) (*AdminStats, error) {
	stats := &AdminStats{}

	// 1. Total de bodegas registradas
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM bodegas
	`).Scan(&stats.TotalBodegas)
	if err != nil {
		return nil, err
	}

	// 2. Total de evaluaciones completadas
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM autoevaluaciones
		WHERE estado = 'COMPLETADA'
	`).Scan(&stats.EvaluacionesCompletadas)
	if err != nil {
		return nil, err
	}

	// 3. Promedio de sostenibilidad (calculado del puntaje final promedio)
	// Por ahora devolvemos 0 - se puede mejorar calculando el promedio real desde respuestas
	stats.PromedioSostenibilidad = 0

	return stats, nil
}

// EvaluacionListItem representa un item en la lista de evaluaciones
type EvaluacionListItem struct {
	IDAutoevaluacion int       `json:"id_autoevaluacion"`
	IDBodega         int       `json:"id_bodega"`
	NombreBodega     string    `json:"nombre_bodega"`
	RazonSocial      string    `json:"razon_social"`
	Estado           string    `json:"estado"`
	Porcentaje       *float64  `json:"porcentaje"`
	FechaInicio      time.Time `json:"fecha_inicio"`
	FechaFin         *time.Time `json:"fecha_fin"`
	Responsable      string    `json:"responsable"`
}

// GetAllEvaluaciones obtiene todas las evaluaciones del sistema con filtros opcionales
func (r *AdminRepository) GetAllEvaluaciones(ctx context.Context, estado string, idBodega *int) ([]EvaluacionListItem, error) {
	query := `
		SELECT
			a.id_autoevaluacion,
			a.id_bodega,
			b.nombre_fantasia,
			b.razon_social,
			a.estado,
			a.fecha_inicio,
			a.fecha_fin,
			COALESCE(resp.nombre || ' ' || resp.apellido, 'Sin responsable') as responsable
		FROM autoevaluaciones a
		INNER JOIN bodegas b ON a.id_bodega = b.id_bodega
		LEFT JOIN cuentas c ON b.id_bodega = c.id_bodega
		LEFT JOIN responsables resp ON c.id_cuenta = resp.id_cuenta AND resp.activo = true
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	// Filtro por estado
	if estado != "" && estado != "TODOS" {
		query += ` AND a.estado = $` + strconv.Itoa(argIndex)
		args = append(args, estado)
		argIndex++
	}

	// Filtro por bodega
	if idBodega != nil {
		query += ` AND a.id_bodega = $` + strconv.Itoa(argIndex)
		args = append(args, *idBodega)
		argIndex++
	}

	query += ` ORDER BY a.fecha_inicio DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evaluaciones []EvaluacionListItem
	for rows.Next() {
		var eval EvaluacionListItem
		err := rows.Scan(
			&eval.IDAutoevaluacion,
			&eval.IDBodega,
			&eval.NombreBodega,
			&eval.RazonSocial,
			&eval.Estado,
			&eval.FechaInicio,
			&eval.FechaFin,
			&eval.Responsable,
		)
		if err != nil {
			return nil, err
		}
		// Porcentaje se deja como nil por ahora - se puede calcular después
		eval.Porcentaje = nil
		evaluaciones = append(evaluaciones, eval)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return evaluaciones, nil
}
