// RUTA: coviar-backend/internal/domain/usuario.go
package domain

import "time"

// Usuario representa un usuario del sistema COVIAR
type Usuario struct {
	IdUsuario     int        `json:"idUsuario"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"` // No se expone en JSON
	Nombre        string     `json:"nombre"`
	Apellido      string     `json:"apellido"`
	Telefono      *string    `json:"telefono"`
	Rol           string     `json:"rol"` // "admin", "bodega", "auditor"
	Activo        bool       `json:"activo"`
	FechaRegistro time.Time  `json:"fecha_registro"`
	UltimoAcceso  *time.Time `json:"ultimo_acceso"`
}

// UsuarioDTO para recibir datos sin campos sensibles
type UsuarioDTO struct {
	Email    string  `json:"email"`
	Password string  `json:"password"` // Solo para crear/login
	Nombre   string  `json:"nombre"`
	Apellido string  `json:"apellido"`
	Telefono *string `json:"telefono"`
	Rol      string  `json:"rol"`
}

// UsuarioLogin para autenticación
type UsuarioLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
