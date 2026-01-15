// RUTA: coviar-backend/internal/bodega/handler.go
// REEMPLAZA EL ARCHIVO ANTERIOR
package bodega

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/carli/coviar-backend/internal/domain"
)

// Handler maneja las peticiones HTTP para Bodega
type Handler struct {
	service *Service
}

// NewHandler crea una nueva instancia del handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListBodegas maneja GET /api/bodegas
func (h *Handler) ListBodegas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		sendError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	bodegas, err := h.service.GetAll()
	if err != nil {
		log.Printf("Error al obtener bodegas: %v", err)
		sendError(w, "Error al obtener bodegas", http.StatusInternalServerError)
		return
	}

	sendSuccess(w, bodegas)
}

// GetBodega maneja GET /api/bodegas/{id}
func (h *Handler) GetBodega(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		sendError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer ID de la URL
	path := strings.TrimPrefix(r.URL.Path, "/api/bodegas/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	bodega, err := h.service.GetByID(id)
	if err != nil {
		log.Printf("Error al obtener bodega: %v", err)
		if err.Error() == "bodega no encontrada" {
			sendError(w, "Bodega no encontrada", http.StatusNotFound)
		} else {
			sendError(w, "Error al obtener bodega", http.StatusInternalServerError)
		}
		return
	}

	sendSuccess(w, bodega)
}

// CreateBodega maneja POST /api/bodegas - Guardar los datos de la bodega
func (h *Handler) CreateBodega(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var bodega domain.Bodega
	if err := json.NewDecoder(r.Body).Decode(&bodega); err != nil {
		sendError(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(&bodega); err != nil {
		log.Printf("Error al crear bodega: %v", err)
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendSuccess(w, bodega)
}

// Utilidades para respuestas JSON

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type successResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{
		Error:   "error",
		Message: message,
	})
}

func sendSuccess(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(successResponse{
		Success: true,
		Data:    data,
	})
}
