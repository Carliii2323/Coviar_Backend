package service

import (
	"context"
	"fmt"
	"strings"

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
/*func (s *AutoevaluacionService) CreateAutoevaluacion(ctx context.Context, idBodega int) (*domain.Autoevaluacion, error) {
	auto := &domain.Autoevaluacion{
		IDBodega: idBodega,
	}

	id, err := s.autoevaluacionRepo.Create(ctx, nil, auto)
	if err != nil {
		return nil, fmt.Errorf("error creating autoevaluacion: %w", err)
	}

	auto.ID = id
	return auto, nil
}*/

// CreateAutoevaluacion crea una nueva autoevaluación para una bodega o retorna la pendiente
func (s *AutoevaluacionService) CreateAutoevaluacion(ctx context.Context, idBodega int) (*domain.AutoevaluacionPendienteResponse, error) {
	// Verificar si ya existe una autoevaluación pendiente para esta bodega
	autoPendiente, err := s.autoevaluacionRepo.FindPendienteByBodega(ctx, idBodega)
	if err != nil {
		return nil, fmt.Errorf("error checking for pending autoevaluacion: %w", err)
	}

	// Si existe una autoevaluación pendiente, retornarla con sus respuestas
	if autoPendiente != nil {
		// Obtener las respuestas existentes
		respuestas, err := s.respuestaRepo.FindByAutoevaluacion(ctx, autoPendiente.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting respuestas: %w", err)
		}

		// Convertir respuestas al formato DTO
		respuestasDTO := make([]domain.GuardarRespuestaRequest, len(respuestas))
		for i, resp := range respuestas {
			respuestasDTO[i] = domain.GuardarRespuestaRequest{
				IDIndicador:      resp.IDIndicador,
				IDNivelRespuesta: resp.IDNivelRespuesta,
			}
		}

		return &domain.AutoevaluacionPendienteResponse{
			AutoevaluacionPendiente: autoPendiente,
			Respuestas:              respuestasDTO,
			Mensaje:                 "Ya existe una autoevaluación pendiente. Para crear una nueva, primero debe cancelar la actual.",
		}, nil
	}

	// No existe autoevaluación pendiente, crear una nueva
	auto := &domain.Autoevaluacion{
		IDBodega: idBodega,
	}

	id, err := s.autoevaluacionRepo.Create(ctx, nil, auto)
	if err != nil {
		return nil, fmt.Errorf("error creating autoevaluacion: %w", err)
	}

	auto.ID = id
	return &domain.AutoevaluacionPendienteResponse{
		AutoevaluacionPendiente: auto,
		Respuestas:              []domain.GuardarRespuestaRequest{},
		Mensaje:                 "Autoevaluación creada correctamente",
	}, nil
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

	// Guardar nuevas respuestas o actualizar existentes
	for _, respReq := range respuestas {
		respuesta := &domain.Respuesta{
			IDNivelRespuesta: respReq.IDNivelRespuesta,
			IDIndicador:      respReq.IDIndicador,
			IDAutoevaluacion: idAutoevaluacion,
		}

		_, err := s.respuestaRepo.Upsert(ctx, nil, respuesta)
		if err != nil {
			return fmt.Errorf("error guarding respuesta: %w", err)
		}
	}

	return nil
}

// CompletarAutoevaluacion marca la autoevaluación como completada
/*func (s *AutoevaluacionService) CompletarAutoevaluacion(ctx context.Context, idAutoevaluacion int) error {
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
}*/

// CompletarAutoevaluacion marca la autoevaluación como completada
func (s *AutoevaluacionService) CompletarAutoevaluacion(ctx context.Context, idAutoevaluacion int) error {
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	if auto == nil {
		return domain.ErrNotFound
	}

	// Verificar que tenga segmento seleccionado
	if auto.IDSegmento == nil {
		return fmt.Errorf("autoevaluacion must have segmento selected")
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

	// === VALIDACIÓN ESTRICTA DE COMPLETITUD ===
	// Obtener los indicadores requeridos para este segmento
	requiredIndicators, err := s.indicadorRepo.FindBySegmento(ctx, *auto.IDSegmento)
	if err != nil {
		return fmt.Errorf("error getting required indicators: %w", err)
	}

	// Verificar que coincida la cantidad
	if len(respuestas) != len(requiredIndicators) {
		return fmt.Errorf("autoevaluacion incomplete: expected %d answers, got %d", len(requiredIndicators), len(respuestas))
	}
	// ==========================================

	// Calcular puntaje total
	puntajeTotal, err := s.respuestaRepo.CalculateTotalScore(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error calculating total score: %w", err)
	}

	// Obtener niveles de sostenibilidad para el segmento
	niveles, err := s.segmentoRepo.FindNivelesSostenibilidadBySegmento(ctx, *auto.IDSegmento)
	if err != nil {
		return fmt.Errorf("error getting niveles sostenibilidad: %w", err)
	}

	// Si no hay niveles configurados, completar sin asignar nivel
	if len(niveles) == 0 {
		// Marcar como completada solo con puntaje (sin nivel de sostenibilidad)
		err = s.autoevaluacionRepo.CompleteWithScore(ctx, idAutoevaluacion, puntajeTotal, 0)
		if err != nil {
			return fmt.Errorf("error completing autoevaluacion: %w", err)
		}
		return nil
	}

	// Determinar el nivel de sostenibilidad según el puntaje
	var nivelAsignado *domain.NivelSostenibilidad
	for _, nivel := range niveles {
		if puntajeTotal >= nivel.MinPuntaje && puntajeTotal <= nivel.MaxPuntaje {
			nivelAsignado = nivel
			break
		}
	}

	// Si no encontró nivel exacto, asignar el más cercano
	if nivelAsignado == nil {
		// Buscar el nivel más alto si el puntaje supera todos los rangos
		for _, nivel := range niveles {
			if puntajeTotal >= nivel.MinPuntaje {
				nivelAsignado = nivel
			}
		}
		// Si aún no hay nivel, usar el primero (puntaje menor al mínimo)
		if nivelAsignado == nil && len(niveles) > 0 {
			nivelAsignado = niveles[0]
		}
	}

	var idNivelSostenibilidad int
	if nivelAsignado != nil {
		idNivelSostenibilidad = nivelAsignado.ID
	}

	// Marcar como completada con puntaje y nivel de sostenibilidad
	err = s.autoevaluacionRepo.CompleteWithScore(ctx, idAutoevaluacion, puntajeTotal, idNivelSostenibilidad)
	if err != nil {
		return fmt.Errorf("error completing autoevaluacion: %w", err)
	}

	return nil
}

// CancelarAutoevaluacion marca la autoevaluación como cancelada
func (s *AutoevaluacionService) CancelarAutoevaluacion(ctx context.Context, idAutoevaluacion int) error {
	// Verificar que la autoevaluación existe
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	if auto == nil {
		return domain.ErrNotFound
	}

	// Verificar que esté en estado PENDIENTE
	if auto.Estado != domain.EstadoPendiente {
		return fmt.Errorf("solo se pueden cancelar autoevaluaciones en estado PENDIENTE")
	}

	// Marcar como cancelada
	err = s.autoevaluacionRepo.Cancel(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error canceling autoevaluacion: %w", err)
	}

	return nil
}

// GetResultadosUltimaAutoevaluacion obtiene los resultados de la última autoevaluación completada
func (s *AutoevaluacionService) GetResultadosUltimaAutoevaluacion(ctx context.Context, idBodega int, bodegaRepo repository.BodegaRepository) (*domain.ResultadoAutoevaluacionResponse, error) {
	// Obtener la bodega
	bodega, err := bodegaRepo.FindByID(ctx, idBodega)
	if err != nil {
		return nil, fmt.Errorf("error finding bodega: %w", err)
	}
	if bodega == nil {
		return nil, domain.ErrNotFound
	}

	// Obtener la última autoevaluación completada
	auto, err := s.autoevaluacionRepo.FindUltimaCompletadaByBodega(ctx, idBodega)
	if err != nil {
		return nil, fmt.Errorf("error finding completed autoevaluacion: %w", err)
	}
	if auto == nil {
		return nil, domain.ErrNotFound
	}

	// Verificar que tenga segmento
	if auto.IDSegmento == nil {
		return nil, fmt.Errorf("autoevaluacion does not have segmento")
	}

	// Verificar que tenga puntaje final
	if auto.PuntajeFinal == nil {
		return nil, fmt.Errorf("autoevaluacion does not have puntaje_final")
	}

	// Verificar que tenga nivel de sostenibilidad
	if auto.IDNivelSostenibilidad == nil {
		return nil, fmt.Errorf("autoevaluacion does not have nivel_sostenibilidad")
	}

	// Verificar que tenga fecha_fin
	if auto.FechaFin == nil {
		return nil, fmt.Errorf("autoevaluacion does not have fecha_fin")
	}

	// Obtener el segmento
	segmento, err := s.segmentoRepo.FindByID(ctx, *auto.IDSegmento)
	if err != nil {
		return nil, fmt.Errorf("error finding segmento: %w", err)
	}

	// Obtener niveles de sostenibilidad para encontrar el nombre
	niveles, err := s.segmentoRepo.FindNivelesSostenibilidadBySegmento(ctx, *auto.IDSegmento)
	if err != nil {
		return nil, fmt.Errorf("error finding niveles sostenibilidad: %w", err)
	}

	var nombreNivelSustentabilidad string
	for _, nivel := range niveles {
		if nivel.ID == *auto.IDNivelSostenibilidad {
			nombreNivelSustentabilidad = nivel.Nombre
			break
		}
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

	// Obtener las respuestas de esta autoevaluación
	respuestas, err := s.respuestaRepo.FindByAutoevaluacion(ctx, auto.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting respuestas: %w", err)
	}

	// Crear mapa de respuestas por indicador
	respuestasPorIndicador := make(map[int]*domain.Respuesta)
	for _, resp := range respuestas {
		respuestasPorIndicador[resp.IDIndicador] = resp
	}

	// Obtener todos los capítulos
	capitulos, err := s.capituloRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting capitulos: %w", err)
	}

	// Construir la estructura de resultados
	resultadoCapitulos := make([]domain.ResultadoCapitulo, 0)

	for _, cap := range capitulos {
		// Obtener indicadores del capítulo
		indicadores, err := s.indicadorRepo.FindByCapitulo(ctx, cap.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting indicadores: %w", err)
		}

		resultadoIndicadores := make([]domain.ResultadoIndicador, 0)

		for _, ind := range indicadores {
			// Solo incluir indicadores habilitados para el segmento
			if !habilitadosMap[ind.ID] {
				continue
			}

			// Solo incluir indicadores que tienen respuesta
			respuesta, tieneRespuesta := respuestasPorIndicador[ind.ID]
			if !tieneRespuesta {
				continue
			}

			// Obtener todos los niveles de respuesta del indicador
			nivelesRespuesta, err := s.nivelRespuestaRepo.FindByIndicador(ctx, ind.ID)
			if err != nil {
				return nil, fmt.Errorf("error getting niveles_respuesta: %w", err)
			}

			// Filtrar solo el nivel que fue seleccionado
			resultadoNiveles := make([]domain.ResultadoNivelRespuesta, 0)
			for _, nivel := range nivelesRespuesta {
				if nivel.ID == respuesta.IDNivelRespuesta {
					resultadoNiveles = append(resultadoNiveles, domain.ResultadoNivelRespuesta{
						Nombre:      nivel.Nombre,
						Descripcion: nivel.Descripcion,
						Puntos:      nivel.Puntos,
					})
					break
				}
			}

			resultadoIndicadores = append(resultadoIndicadores, domain.ResultadoIndicador{
				Nombre:           ind.Nombre,
				Descripcion:      ind.Descripcion,
				Orden:            ind.Orden,
				NivelesRespuesta: resultadoNiveles,
			})
		}

		// Solo incluir capítulos que tengan al menos un indicador con respuesta
		if len(resultadoIndicadores) > 0 {
			resultadoCapitulos = append(resultadoCapitulos, domain.ResultadoCapitulo{
				Nombre:      cap.Nombre,
				Orden:       cap.Orden,
				Indicadores: resultadoIndicadores,
			})
		}
	}

	// Construir la respuesta final
	response := &domain.ResultadoAutoevaluacionResponse{
		Bodega: domain.ResultadoBodega{
			NombreFantasia: bodega.NombreFantasia,
		},
		Autoevaluacion: domain.ResultadoAutoevaluacion{
			FechaFin:     *auto.FechaFin,
			PuntajeFinal: *auto.PuntajeFinal,
		},
		Segmento: domain.ResultadoSegmento{
			Nombre: segmento.Nombre,
		},
		NivelSustentabilidad: domain.ResultadoNivelSustentabilidad{
			Nombre: nombreNivelSustentabilidad,
		},
		Capitulos: resultadoCapitulos,
	}

	return response, nil
}

// GetHistorialAutoevaluaciones obtiene la lista resumida de todas las autoevaluaciones completadas de una bodega
func (s *AutoevaluacionService) GetHistorialAutoevaluaciones(ctx context.Context, idBodega int) ([]domain.HistorialItemResponse, error) {
	// Obtener todas las autoevaluaciones completadas
	autoevaluaciones, err := s.autoevaluacionRepo.FindCompletadasByBodega(ctx, idBodega)
	if err != nil {
		return nil, fmt.Errorf("error finding completed autoevaluaciones: %w", err)
	}

	if len(autoevaluaciones) == 0 {
		return []domain.HistorialItemResponse{}, nil
	}

	// Cache de max puntos por segmento para evitar queries repetidas
	maxPuntosCache := make(map[int]map[int]int)

	historial := make([]domain.HistorialItemResponse, 0, len(autoevaluaciones))

	for _, auto := range autoevaluaciones {
		item := domain.HistorialItemResponse{
			IDAutoevaluacion:      auto.ID,
			FechaInicio:           auto.FechaInicio,
			FechaFinalizacion:     auto.FechaFin,
			Estado:                strings.ToLower(string(auto.Estado)),
			IDBodega:              auto.IDBodega,
			IDSegmento:            auto.IDSegmento,
			PuntajeFinal:          auto.PuntajeFinal,
			IDNivelSostenibilidad: auto.IDNivelSostenibilidad,
		}

		// Resolver segmento y calcular puntaje máximo
		if auto.IDSegmento != nil {
			segmento, err := s.segmentoRepo.FindByID(ctx, *auto.IDSegmento)
			if err == nil && segmento != nil {
				item.NombreSegmento = segmento.Nombre
			}

			// Obtener max puntos (con cache)
			maxPuntos, exists := maxPuntosCache[*auto.IDSegmento]
			if !exists {
				maxPuntos, err = s.nivelRespuestaRepo.FindMaxPuntosBySegmento(ctx, *auto.IDSegmento)
				if err != nil {
					return nil, fmt.Errorf("error getting max puntos for segmento %d: %w", *auto.IDSegmento, err)
				}
				maxPuntosCache[*auto.IDSegmento] = maxPuntos
			}

			// Calcular puntaje máximo total
			totalMax := 0
			for _, mp := range maxPuntos {
				totalMax += mp
			}
			if totalMax > 0 {
				item.PuntajeMaximo = &totalMax
				if auto.PuntajeFinal != nil {
					porcentaje := (*auto.PuntajeFinal * 100) / totalMax
					item.Porcentaje = &porcentaje
				}
			}
		}

		// Resolver nivel de sostenibilidad
		if auto.IDNivelSostenibilidad != nil && auto.IDSegmento != nil {
			niveles, err := s.segmentoRepo.FindNivelesSostenibilidadBySegmento(ctx, *auto.IDSegmento)
			if err == nil {
				for _, nivel := range niveles {
					if nivel.ID == *auto.IDNivelSostenibilidad {
						item.NivelSostenibilidad = &domain.NivelSostenibilidadInfo{
							ID:     nivel.ID,
							Nombre: nivel.Nombre,
						}
						break
					}
				}
			}
		}

		historial = append(historial, item)
	}

	return historial, nil
}

// GetResultadosByID obtiene los resultados detallados de una autoevaluación por su ID
func (s *AutoevaluacionService) GetResultadosByID(ctx context.Context, idAutoevaluacion int, bodegaRepo repository.BodegaRepository) (*domain.ResultadoDetalladoResponse, error) {
	// Obtener la autoevaluación
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return nil, fmt.Errorf("error finding autoevaluacion: %w", err)
	}
	if auto == nil {
		return nil, domain.ErrNotFound
	}

	// Verificar que esté completada
	if auto.Estado != domain.EstadoCompletada {
		return nil, fmt.Errorf("la autoevaluación no está completada")
	}

	if auto.IDSegmento == nil {
		return nil, fmt.Errorf("autoevaluacion does not have segmento")
	}

	// Obtener segmento
	segmento, err := s.segmentoRepo.FindByID(ctx, *auto.IDSegmento)
	if err != nil {
		return nil, fmt.Errorf("error finding segmento: %w", err)
	}

	// Obtener max puntos por indicador para el segmento
	maxPuntosPorIndicador, err := s.nivelRespuestaRepo.FindMaxPuntosBySegmento(ctx, *auto.IDSegmento)
	if err != nil {
		return nil, fmt.Errorf("error getting max puntos: %w", err)
	}

	// Obtener indicadores habilitados para el segmento
	habilitadosIds, err := s.indicadorRepo.FindBySegmento(ctx, *auto.IDSegmento)
	if err != nil {
		return nil, fmt.Errorf("error getting enabled indicators: %w", err)
	}
	habilitadosMap := make(map[int]bool)
	for _, id := range habilitadosIds {
		habilitadosMap[id] = true
	}

	// Obtener respuestas de la autoevaluación
	respuestas, err := s.respuestaRepo.FindByAutoevaluacion(ctx, auto.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting respuestas: %w", err)
	}

	// Crear mapa de respuestas por indicador
	respuestasPorIndicador := make(map[int]*domain.Respuesta)
	for _, resp := range respuestas {
		respuestasPorIndicador[resp.IDIndicador] = resp
	}

	// Obtener todos los capítulos
	capitulos, err := s.capituloRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting capitulos: %w", err)
	}

	// Construir resultados por capítulo
	resultadoCapitulos := make([]domain.ResultadoCapituloDetallado, 0)

	totalMaxGlobal := 0
	for _, cap := range capitulos {
		indicadores, err := s.indicadorRepo.FindByCapitulo(ctx, cap.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting indicadores: %w", err)
		}

		puntajeObtenido := 0
		puntajeMaximo := 0
		indicadoresTotal := 0
		indicadoresCompletados := 0
		indicadoresDetalle := make([]domain.ResultadoIndicadorDetalle, 0)

		for _, ind := range indicadores {
			if !habilitadosMap[ind.ID] {
				continue
			}

			indicadoresTotal++
			indMaxPuntos := maxPuntosPorIndicador[ind.ID]
			puntajeMaximo += indMaxPuntos

			indDetalle := domain.ResultadoIndicadorDetalle{
				IDIndicador:   ind.ID,
				Nombre:        ind.Nombre,
				Descripcion:   ind.Descripcion,
				Orden:         ind.Orden,
				PuntajeMaximo: indMaxPuntos,
			}

			if resp, tiene := respuestasPorIndicador[ind.ID]; tiene {
				indicadoresCompletados++
				// Obtener datos de la respuesta seleccionada
				nivelesResp, err := s.nivelRespuestaRepo.FindByIndicador(ctx, ind.ID)
				if err == nil {
					for _, nivel := range nivelesResp {
						if nivel.ID == resp.IDNivelRespuesta {
							puntajeObtenido += nivel.Puntos
							indDetalle.RespuestaNombre = nivel.Nombre
							indDetalle.RespuestaDescripcion = nivel.Descripcion
							indDetalle.RespuestaPuntos = nivel.Puntos
							break
						}
					}
				}
			}

			indicadoresDetalle = append(indicadoresDetalle, indDetalle)
		}

		totalMaxGlobal += puntajeMaximo

		if indicadoresTotal > 0 {
			porcentaje := 0
			if puntajeMaximo > 0 {
				porcentaje = (puntajeObtenido * 100) / puntajeMaximo
			}

			resultadoCapitulos = append(resultadoCapitulos, domain.ResultadoCapituloDetallado{
				IDCapitulo:             cap.ID,
				Nombre:                 cap.Nombre,
				PuntajeObtenido:        puntajeObtenido,
				PuntajeMaximo:          puntajeMaximo,
				Porcentaje:             porcentaje,
				IndicadoresCompletados: indicadoresCompletados,
				IndicadoresTotal:       indicadoresTotal,
				Indicadores:            indicadoresDetalle,
			})
		}
	}

	// Construir item de autoevaluación
	autoItem := domain.HistorialItemResponse{
		IDAutoevaluacion:      auto.ID,
		FechaInicio:           auto.FechaInicio,
		FechaFinalizacion:     auto.FechaFin,
		Estado:                strings.ToLower(string(auto.Estado)),
		IDBodega:              auto.IDBodega,
		IDSegmento:            auto.IDSegmento,
		PuntajeFinal:          auto.PuntajeFinal,
		IDNivelSostenibilidad: auto.IDNivelSostenibilidad,
	}

	if segmento != nil {
		autoItem.NombreSegmento = segmento.Nombre
	}

	if totalMaxGlobal > 0 {
		autoItem.PuntajeMaximo = &totalMaxGlobal
		if auto.PuntajeFinal != nil {
			porcentaje := (*auto.PuntajeFinal * 100) / totalMaxGlobal
			autoItem.Porcentaje = &porcentaje
		}
	}

	// Resolver nivel de sostenibilidad
	if auto.IDNivelSostenibilidad != nil && auto.IDSegmento != nil {
		niveles, err := s.segmentoRepo.FindNivelesSostenibilidadBySegmento(ctx, *auto.IDSegmento)
		if err == nil {
			for _, nivel := range niveles {
				if nivel.ID == *auto.IDNivelSostenibilidad {
					autoItem.NivelSostenibilidad = &domain.NivelSostenibilidadInfo{
						ID:     nivel.ID,
						Nombre: nivel.Nombre,
					}
					break
				}
			}
		}
	}

	return &domain.ResultadoDetalladoResponse{
		Autoevaluacion: autoItem,
		Capitulos:      resultadoCapitulos,
	}, nil
}
