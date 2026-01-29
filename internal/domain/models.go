package domain

import "time"

// ============================================
// MODELOS DE BODEGA
// ============================================

type Bodega struct {
	ID                 int       `json:"id_bodega,omitempty"`
	RazonSocial        string    `json:"razon_social"`
	NombreFantasia     string    `json:"nombre_fantasia"`
	CUIT               string    `json:"cuit"`              // char(11), check: ^[0-9]{11}$
	InvBod             *string   `json:"inv_bod,omitempty"` // char(6)
	InvVin             *string   `json:"inv_vin,omitempty"` // char(6)
	Calle              string    `json:"calle"`
	Numeracion         string    `json:"numeracion"` // default 'S/N'
	IDLocalidad        int       `json:"id_localidad"`
	Telefono           string    `json:"telefono"`            // check: ^[0-9]+$
	EmailInstitucional string    `json:"email_institucional"` // check: like '%@%'
	FechaRegistro      time.Time `json:"fecha_registro,omitempty"`
}

type BodegaRequest struct {
	RazonSocial        string  `json:"razon_social"`
	NombreFantasia     string  `json:"nombre_fantasia"`
	CUIT               string  `json:"cuit"`
	InvBod             *string `json:"inv_bod,omitempty"`
	InvVin             *string `json:"inv_vin,omitempty"`
	Calle              string  `json:"calle"`
	Numeracion         string  `json:"numeracion"`
	IDLocalidad        int     `json:"id_localidad"`
	Telefono           string  `json:"telefono"`
	EmailInstitucional string  `json:"email_institucional"`
}

type BodegaUpdateDTO struct {
	Telefono           string `json:"telefono"`
	EmailInstitucional string `json:"email_institucional"`
	NombreFantasia     string `json:"nombre_fantasia"`
}

// ============================================
// MODELOS DE CUENTA
// ============================================

type TipoCuenta string

const (
	TipoCuentaBodega           TipoCuenta = "BODEGA"
	TipoCuentaAdministradorApp TipoCuenta = "ADMINISTRADOR_APP"
)

type Cuenta struct {
	ID            int        `json:"id_cuenta,omitempty"`
	Tipo          TipoCuenta `json:"tipo"`                // ENUM: BODEGA, ADMINISTRADOR_APP
	IDBodega      *int       `json:"id_bodega,omitempty"` // nullable, depende de tipo
	EmailLogin    string     `json:"email_login"`
	PasswordHash  string     `json:"-"`
	FechaRegistro time.Time  `json:"fecha_registro,omitempty"`
}

type CuentaRequest struct {
	EmailLogin string `json:"email_login"`
	Password   string `json:"password"`
}

// ============================================
// MODELOS DE RESPONSABLE
// ============================================

type Responsable struct {
	ID            int        `json:"id_responsable,omitempty"`
	IDCuenta      int        `json:"id_cuenta"`
	Nombre        string     `json:"nombre"`
	Apellido      string     `json:"apellido"`
	Cargo         string     `json:"cargo"`
	DNI           string     `json:"dni"` // varchar(8), check: ^[0-9]{7,8}$
	Activo        bool       `json:"activo"`
	FechaRegistro time.Time  `json:"fecha_registro,omitempty"`
	FechaBaja     *time.Time `json:"fecha_baja,omitempty"`
}

type ResponsableRequest struct {
	Nombre   string  `json:"nombre"`
	Apellido string  `json:"apellido"`
	Cargo    string  `json:"cargo"`
	DNI      *string `json:"dni,omitempty"`
}

// ============================================
// MODELOS DE USUARIO
// ============================================

type Usuario struct {
	IdUsuario     int        `json:"idUsuario,omitempty" db:"id_usuario"`
	Email         string     `json:"email" db:"email"`
	PasswordHash  string     `json:"-" db:"password_hash"`
	Nombre        string     `json:"nombre" db:"nombre"`
	Apellido      string     `json:"apellido" db:"apellido"`
	Rol           string     `json:"rol" db:"rol"`
	Activo        bool       `json:"activo" db:"activo"`
	FechaRegistro time.Time  `json:"fecha_registro" db:"fecha_registro"`
	UltimoAcceso  *time.Time `json:"ultimo_acceso" db:"ultimo_acceso"`
}

type UsuarioDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Rol      string `json:"rol"`
}

type UsuarioLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ToPublic limpia campos sensibles antes de enviar al cliente
func (u *Usuario) ToPublic() *Usuario {
	publicUser := *u
	publicUser.PasswordHash = ""
	return &publicUser
}

// ============================================
// MODELOS DE UBICACIÓN
// ============================================

type Provincia struct {
	ID     int    `json:"id_provincia,omitempty"`
	Nombre string `json:"nombre"`
}

type Departamento struct {
	ID          int    `json:"id_departamento,omitempty"`
	IDProvincia int    `json:"id_provincia"`
	Nombre      string `json:"nombre"`
}

type Localidad struct {
	ID             int    `json:"id_localidad,omitempty"`
	IDDepartamento int    `json:"id_departamento"`
	Nombre         string `json:"nombre"`
}

// ============================================
// DTOs DE REGISTRO
// ============================================

type RegistroRequest struct {
	Bodega      BodegaRequest      `json:"bodega"`
	Cuenta      CuentaRequest      `json:"cuenta"`
	Responsable ResponsableRequest `json:"responsable"`
}

type RegistroResponse struct {
	IDBodega      int    `json:"id_bodega"`
	IDCuenta      int    `json:"id_cuenta"`
	IDResponsable int    `json:"id_responsable"`
	Mensaje       string `json:"mensaje"`
}

// ============================================
// MODELOS DE AUTOEVALUACIÓN
// ============================================

type EstadoAutoevaluacion string

const (
	EstadoPendiente  EstadoAutoevaluacion = "PENDIENTE"
	EstadoCompletada EstadoAutoevaluacion = "COMPLETADA"
	EstadoCancelada  EstadoAutoevaluacion = "CANCELADA"
)

type Segmento struct {
	ID          int    `json:"id_segmento"`
	Nombre      string `json:"nombre"`
	MinTuristas int    `json:"min_turistas"`
	MaxTuristas *int   `json:"max_turistas,omitempty"`
}

type NivelSostenibilidad struct {
	ID         int    `json:"id_nivel_sostenibilidad"`
	IDSegmento int    `json:"id_segmento"`
	Nombre     string `json:"nombre"`
	MinPuntaje int    `json:"min_puntaje"`
	MaxPuntaje int    `json:"max_puntaje"`
}

type Capitulo struct {
	ID          int    `json:"id_capitulo"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Orden       int    `json:"orden"`
}

type Indicador struct {
	ID          int    `json:"id_indicador"`
	IDCapitulo  int    `json:"id_capitulo"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Orden       int    `json:"orden"`
}

type IndicadorConHabilitacion struct {
	Indicador        *Indicador        `json:"indicador"`
	NivelesRespuesta []*NivelRespuesta `json:"niveles_respuesta"`
	Habilitado       bool              `json:"habilitado"`
}

type NivelRespuesta struct {
	ID          int    `json:"id_nivel_respuesta"`
	IDIndicador int    `json:"id_indicador"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Puntos      int    `json:"puntos"`
}

type Autoevaluacion struct {
	ID          int                  `json:"id_autoevaluacion"`
	FechaInicio time.Time            `json:"fecha_inicio"`
	FechaFin    *time.Time           `json:"fecha_fin,omitempty"`
	Estado      EstadoAutoevaluacion `json:"estado"`
	IDBodega    int                  `json:"id_bodega"`
	IDSegmento  *int                 `json:"id_segmento,omitempty"`
}

type Respuesta struct {
	ID               int `json:"id_respuesta"`
	IDNivelRespuesta int `json:"id_nivel_respuesta"`
	IDIndicador      int `json:"id_indicador"`
	IDAutoevaluacion int `json:"id_autoevaluacion"`
}

type EstructuraAutoevaluacion struct {
	Capitulos []*CapituloEstructura `json:"capitulos"`
}

type CapituloEstructura struct {
	Capitulo    *Capitulo                   `json:"capitulo"`
	Indicadores []*IndicadorConHabilitacion `json:"indicadores"`
}

// DTOs
type CreateAutoevaluacionRequest struct {
	IDBodega int `json:"id_bodega"`
}

type SeleccionarSegmentoRequest struct {
	IDSegmento int `json:"id_segmento"`
}

type GuardarRespuestaRequest struct {
	IDIndicador      int `json:"id_indicador"`
	IDNivelRespuesta int `json:"id_nivel_respuesta"`
}

type GuardarRespuestasRequest struct {
	Respuestas []GuardarRespuestaRequest `json:"respuestas"`
}
