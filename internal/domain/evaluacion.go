// RUTA: coviar-backend/internal/domain/evaluacion.go
package domain

// Evaluacion representa una autoevaluación de sostenibilidad
type Evaluacion struct {
	IdEvaluacion    int     `json:"idEvaluacion"`
	IdBodega        int     `json:"idBodega"`
	IdSegmento      int     `json:"idSegmento"`
	FechaInicio     string  `json:"fecha_inicio"`
	FechaCompletado *string `json:"fecha_completado"`
	Estado          string  `json:"estado"`
	PuntajeTotal    *int    `json:"puntaje_total"`
	IdNvSos         *int    `json:"idNvSos"`
}
