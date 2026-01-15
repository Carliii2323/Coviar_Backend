// RUTA: coviar-backend/internal/domain/indicador.go
package domain

// Indicador representa un indicador de sostenibilidad
type Indicador struct {
	IdIndicador int     `json:"idIndicador"`
	Codigo      string  `json:"codigo"`
	Nombre      string  `json:"nombre"`
	Descripcion *string `json:"descripcion"`
	Vigente     bool    `json:"vigente"`
}
