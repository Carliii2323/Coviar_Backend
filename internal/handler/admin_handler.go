package handler

import (
	"net/http"
	"strconv"

	"coviar_backend/internal/service"
	"coviar_backend/pkg/httputil"
)

type AdminHandler struct {
	service *service.AdminService
}

func NewAdminHandler(service *service.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

// GetStats maneja GET /api/admin/stats
func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, stats)
}

// GetAllEvaluaciones maneja GET /api/admin/evaluaciones?estado=...&id_bodega=...
func (h *AdminHandler) GetAllEvaluaciones(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de query
	estado := r.URL.Query().Get("estado")
	if estado == "" {
		estado = "TODOS"
	}

	var idBodega *int
	idBodegaStr := r.URL.Query().Get("id_bodega")
	if idBodegaStr != "" {
		id, err := strconv.Atoi(idBodegaStr)
		if err == nil {
			idBodega = &id
		}
	}

	evaluaciones, err := h.service.GetAllEvaluaciones(r.Context(), estado, idBodega)
	if err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, evaluaciones)
}
