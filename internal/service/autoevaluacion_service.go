package service

import (
	"context"
	"fmt"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
)

type AutoevaluacionService struct {
	autoevaluacionRepo repository.AutoevaluacionRepository
	segmentoRepo       repository.SegmentoRepository
	capituloRepo       repository.CapituloRepository
	indicadorRepo      repository.IndicadorRepository
	nivelRespuestaRepo repository.NivelRespuestaRepository
	respuestaRepo      repository.RespuestaRepository
}

func NewAutoevaluacionService(
	autoevaluacionRepo repository.AutoevaluacionRepository,
	segmentoRepo repository.SegmentoRepository,
	capituloRepo repository.CapituloRepository,
	indicadorRepo repository.IndicadorRepository,
	nivelRespuestaRepo repository.NivelRespuestaRepository,
	respuestaRepo repository.RespuestaRepository,
) *AutoevaluacionService {
	return &AutoevaluacionService{
		autoevaluacionRepo: autoevaluacionRepo,
		segmentoRepo:       segmentoRepo,
		capituloRepo:       capituloRepo,
		indicadorRepo:      indicadorRepo,
		nivelRespuestaRepo: nivelRespuestaRepo,
		respuestaRepo:      respuestaRepo,
	}
}

// CreateAutoevaluacion crea una nueva autoevaluación para una bodega
func (s *AutoevaluacionService) CreateAutoevaluacion(ctx context.Context, idBodega int) (*domain.Autoevaluacion, error) {
	auto := &domain.Autoevaluacion{
		IDBodega: idBodega,
	}

	id, err := s.autoevaluacionRepo.Create(ctx, nil, auto)
	if err != nil {
		return nil, fmt.Errorf("error creating autoevaluacion: %w", err)
	}

	auto.ID = id
	return auto, nil
}

// GetSegmentos obtiene todos los segmentos disponibles
func (s *AutoevaluacionService) GetSegmentos(ctx context.Context) ([]*domain.Segmento, error) {
	segmentos, err := s.segmentoRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting segmentos: %w", err)
	}
	return segmentos, nil
}

// SeleccionarSegmento selecciona un segmento para la autoevaluación
func (s *AutoevaluacionService) SeleccionarSegmento(ctx context.Context, idAutoevaluacion int, idSegmento int) error {
	// Verificar que la autoevaluación existe
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	if auto == nil {
		return domain.ErrNotFound
	}

	// Verificar que el segmento existe
	seg, err := s.segmentoRepo.FindByID(ctx, idSegmento)
	if err != nil {
		return fmt.Errorf("error finding segmento: %w", err)
	}

	if seg == nil {
		return domain.ErrNotFound
	}

	// Actualizar la autoevaluación con el segmento
	err = s.autoevaluacionRepo.UpdateSegmento(ctx, idAutoevaluacion, idSegmento)
	if err != nil {
		return fmt.Errorf("error selecting segmento: %w", err)
	}

	return nil
}

// GetEstructura obtiene la estructura del cuestionario con indicadores habilitados según el segmento
func (s *AutoevaluacionService) GetEstructura(ctx context.Context, idAutoevaluacion int) (*domain.EstructuraAutoevaluacion, error) {
	// Obtener autoevaluación
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return nil, fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	if auto == nil || auto.IDSegmento == nil {
		return nil, fmt.Errorf("autoevaluacion not found or segmento not selected")
	}

	// Obtener indicadores habilitados para este segmento
	habilitadosIds, err := s.indicadorRepo.FindBySegmento(ctx, *auto.IDSegmento)
	if err != nil {
		return nil, fmt.Errorf("error getting enabled indicators: %w", err)
	}

	// Convertir a mapa para búsqueda rápida
	habilitadosMap := make(map[int]bool)
	for _, id := range habilitadosIds {
		habilitadosMap[id] = true
	}

	// Obtener todos los capítulos
	capitulos, err := s.capituloRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting capitulos: %w", err)
	}

	estructura := &domain.EstructuraAutoevaluacion{
		Capitulos: make([]*domain.CapituloEstructura, 0),
	}

	for _, cap := range capitulos {
		// Obtener indicadores del capítulo
		indicadores, err := s.indicadorRepo.FindByCapitulo(ctx, cap.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting indicadores: %w", err)
		}

		capEstructura := &domain.CapituloEstructura{
			Capitulo:    cap,
			Indicadores: make([]*domain.IndicadorConHabilitacion, 0),
		}

		for _, ind := range indicadores {
			// Obtener niveles de respuesta
			niveles, err := s.nivelRespuestaRepo.FindByIndicador(ctx, ind.ID)
			if err != nil {
				return nil, fmt.Errorf("error getting niveles_respuesta: %w", err)
			}

			habilitado := habilitadosMap[ind.ID]

			indConHab := &domain.IndicadorConHabilitacion{
				Indicador:        ind,
				NivelesRespuesta: niveles,
				Habilitado:       habilitado,
			}

			capEstructura.Indicadores = append(capEstructura.Indicadores, indConHab)
		}

		estructura.Capitulos = append(estructura.Capitulos, capEstructura)
	}

	return estructura, nil
}

// GuardarRespuestas guarda las respuestas de la autoevaluación
func (s *AutoevaluacionService) GuardarRespuestas(ctx context.Context, idAutoevaluacion int, respuestas []domain.GuardarRespuestaRequest) error {
	// Verificar que la autoevaluación existe
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	if auto == nil {
		return domain.ErrNotFound
	}

	// Limpiar respuestas anteriores
	err = s.respuestaRepo.DeleteByAutoevaluacion(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error clearing previous respuestas: %w", err)
	}

	// Guardar nuevas respuestas
	for _, respReq := range respuestas {
		respuesta := &domain.Respuesta{
			IDNivelRespuesta: respReq.IDNivelRespuesta,
			IDIndicador:      respReq.IDIndicador,
			IDAutoevaluacion: idAutoevaluacion,
		}

		_, err := s.respuestaRepo.Create(ctx, nil, respuesta)
		if err != nil {
			return fmt.Errorf("error creating respuesta: %w", err)
		}
	}

	return nil
}

// CompletarAutoevaluacion marca la autoevaluación como completada
func (s *AutoevaluacionService) CompletarAutoevaluacion(ctx context.Context, idAutoevaluacion int) error {
	// Verificar que la autoevaluación existe
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	if auto == nil {
		return domain.ErrNotFound
	}

	// Obtener respuestas para validar que todas las preguntas fueron respondidas
	respuestas, err := s.respuestaRepo.FindByAutoevaluacion(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error getting respuestas: %w", err)
	}

	// Validación básica: debe haber al menos una respuesta
	if len(respuestas) == 0 {
		return fmt.Errorf("autoevaluacion must have at least one respuesta")
	}

	// Marcar como completada
	err = s.autoevaluacionRepo.Complete(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error completing autoevaluacion: %w", err)
	}

	return nil
}
