// RUTA: coviar-backend/internal/usuario/handler.go
package usuario

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/carli/coviar-backend/internal/domain"
)

// Handler maneja las peticiones HTTP para Usuario
type Handler struct {
	service *Service
}

// NewHandler crea una nueva instancia del handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create maneja POST /api/usuarios - Guardar datos del usuario
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var dto domain.UsuarioDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		sendError(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	usuario, err := h.service.Create(&dto)
	if err != nil {
		log.Printf("Error al crear usuario: %v", err)
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendSuccess(w, usuario)
}

// Verify maneja POST /api/usuarios/verificar - Verificar datos del usuario
func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var login domain.UsuarioLogin
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		sendError(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	usuario, err := h.service.Verify(&login)
	if err != nil {
		log.Printf("Error en verificación: %v", err)
		sendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sendSuccess(w, usuario)
}

// Deactivate maneja DELETE /api/usuarios/{id} - Dar de baja al usuario
func (h *Handler) Deactivate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		sendError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer ID de la URL
	path := strings.TrimPrefix(r.URL.Path, "/api/usuarios/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.Deactivate(id); err != nil {
		log.Printf("Error al desactivar usuario: %v", err)
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendSuccess(w, map[string]string{"message": "Usuario desactivado correctamente"})
}

// GetByID maneja GET /api/usuarios/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		sendError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer ID de la URL
	path := strings.TrimPrefix(r.URL.Path, "/api/usuarios/")
	
	// Si la ruta es "verificar", no es un ID
	if path == "verificar" {
		sendError(w, "Ruta no válida", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	usuario, err := h.service.GetByID(id)
	if err != nil {
		log.Printf("Error al obtener usuario: %v", err)
		sendError(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	sendSuccess(w, usuario)
}

// ListAll maneja GET /api/usuarios
func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		sendError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	usuarios, err := h.service.GetAll()
	if err != nil {
		log.Printf("Error al obtener usuarios: %v", err)
		sendError(w, "Error al obtener usuarios", http.StatusInternalServerError)
		return
	}

	sendSuccess(w, usuarios)
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
