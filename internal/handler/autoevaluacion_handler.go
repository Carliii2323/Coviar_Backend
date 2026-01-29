package handler

import (
	"net/http"
	"strconv"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/service"
	"coviar_backend/pkg/httputil"
	"coviar_backend/pkg/router"
)

type AutoevaluacionHandler struct {
	service *service.AutoevaluacionService
}

func NewAutoevaluacionHandler(service *service.AutoevaluacionService) *AutoevaluacionHandler {
	return &AutoevaluacionHandler{service: service}
}

// CreateAutoevaluacion POST /api/autoevaluaciones
func (h *AutoevaluacionHandler) CreateAutoevaluacion(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateAutoevaluacionRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	auto, err := h.service.CreateAutoevaluacion(r.Context(), req.IDBodega)
	if err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusCreated, auto)
}

// GetSegmentos GET /api/autoevaluaciones/{id_autoevaluacion}/segmentos
func (h *AutoevaluacionHandler) GetSegmentos(w http.ResponseWriter, r *http.Request) {
	segmentos, err := h.service.GetSegmentos(r.Context())
	if err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, segmentos)
}

// SeleccionarSegmento PUT /api/autoevaluaciones/{id_autoevaluacion}/segmento
func (h *AutoevaluacionHandler) SeleccionarSegmento(w http.ResponseWriter, r *http.Request) {
	idStr := router.GetParam(r, "id_autoevaluacion")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req domain.SeleccionarSegmentoRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.SeleccionarSegmento(r.Context(), id, req.IDSegmento); err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "Segmento seleccionado correctamente"})
}

// GetEstructura GET /api/autoevaluaciones/{id_autoevaluacion}/estructura
func (h *AutoevaluacionHandler) GetEstructura(w http.ResponseWriter, r *http.Request) {
	idStr := router.GetParam(r, "id_autoevaluacion")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	estructura, err := h.service.GetEstructura(r.Context(), id)
	if err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, estructura)
}

// GuardarRespuestas POST /api/autoevaluaciones/{id_autoevaluacion}/respuestas
func (h *AutoevaluacionHandler) GuardarRespuestas(w http.ResponseWriter, r *http.Request) {
	idStr := router.GetParam(r, "id_autoevaluacion")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req domain.GuardarRespuestasRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.GuardarRespuestas(r.Context(), id, req.Respuestas); err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "Respuestas guardadas correctamente"})
}

// CompletarAutoevaluacion POST /api/autoevaluaciones/{id_autoevaluacion}/completar
func (h *AutoevaluacionHandler) CompletarAutoevaluacion(w http.ResponseWriter, r *http.Request) {
	idStr := router.GetParam(r, "id_autoevaluacion")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.service.CompletarAutoevaluacion(r.Context(), id); err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "Autoevaluación completada correctamente"})
}
