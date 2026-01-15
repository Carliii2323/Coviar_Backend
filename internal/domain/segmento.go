// RUTA: coviar-backend/internal/domain/segmento.go
package domain

// Segmento representa la categoría de bodega según turistas anuales
type Segmento struct {
	IdSegmento  int     `json:"idSegmento"`
	Nombre      string  `json:"nombre"`
	MinTuristas *int    `json:"min_turistas"`
	MaxTuristas *int    `json:"max_turistas"`
	Descripcion *string `json:"descripcion"`
}
