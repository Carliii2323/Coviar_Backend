package repository

import (
	"context"

	"coviar_backend/internal/domain"
)

// Transaction representa una transacción de base de datos
type Transaction interface {
	Commit() error
	Rollback() error
}

// Repositorios para Bodega
type BodegaRepository interface {
	Create(ctx context.Context, tx Transaction, bodega *domain.Bodega) (int, error)
	FindByID(ctx context.Context, id int) (*domain.Bodega, error)
	FindByCUIT(ctx context.Context, cuit string) (*domain.Bodega, error)
	Update(ctx context.Context, tx Transaction, bodega *domain.Bodega) error
	Delete(ctx context.Context, tx Transaction, id int) error
	GetAll(ctx context.Context) ([]*domain.Bodega, error)
}

// Repositorios para Cuenta
type CuentaRepository interface {
	Create(ctx context.Context, tx Transaction, cuenta *domain.Cuenta) (int, error)
	FindByID(ctx context.Context, id int) (*domain.Cuenta, error)
	FindByEmail(ctx context.Context, email string) (*domain.Cuenta, error)
	Update(ctx context.Context, tx Transaction, cuenta *domain.Cuenta) error
	Delete(ctx context.Context, tx Transaction, id int) error
}

// Repositorios para Responsable
type ResponsableRepository interface {
	Create(ctx context.Context, tx Transaction, responsable *domain.Responsable) (int, error)
	FindByID(ctx context.Context, id int) (*domain.Responsable, error)
	FindByCuentaID(ctx context.Context, cuentaID int) ([]*domain.Responsable, error)
	Update(ctx context.Context, tx Transaction, responsable *domain.Responsable) error
	Delete(ctx context.Context, tx Transaction, id int) error
}

// Repositorios para Usuario
type UsuarioRepository interface {
	Create(ctx context.Context, usuario *domain.Usuario) (int, error)
	FindByID(ctx context.Context, id int) (*domain.Usuario, error)
	FindByEmail(ctx context.Context, email string) (*domain.Usuario, error)
	GetAll(ctx context.Context) ([]*domain.Usuario, error)
	Update(ctx context.Context, usuario *domain.Usuario) error
	Delete(ctx context.Context, id int) error
}

// Repositorios para Ubicación
type UbicacionRepository interface {
	// Provincias
	GetProvincias(ctx context.Context) ([]*domain.Provincia, error)
	GetProvinciaByID(ctx context.Context, id int) (*domain.Provincia, error)

	// Departamentos
	GetDepartamentos(ctx context.Context) ([]*domain.Departamento, error)
	GetDepartamentoByID(ctx context.Context, id int) (*domain.Departamento, error)
	GetDepartamentosByProvinciaID(ctx context.Context, provinciaID int) ([]*domain.Departamento, error)

	// Localidades
	GetLocalidades(ctx context.Context) ([]*domain.Localidad, error)
	GetLocalidadByID(ctx context.Context, id int) (*domain.Localidad, error)
	GetLocalidadesByDepartamentoID(ctx context.Context, departamentoID int) ([]*domain.Localidad, error)
}

// TransactionManager maneja transacciones
type TransactionManager interface {
	BeginTx(ctx context.Context) (Transaction, error)
}

// Repositorios para Autoevaluación
type SegmentoRepository interface {
	FindAll(ctx context.Context) ([]*domain.Segmento, error)
	FindByID(ctx context.Context, id int) (*domain.Segmento, error)
}

type AutoevaluacionRepository interface {
	Create(ctx context.Context, tx Transaction, auto *domain.Autoevaluacion) (int, error)
	FindByID(ctx context.Context, id int) (*domain.Autoevaluacion, error)
	UpdateSegmento(ctx context.Context, id int, idSegmento int) error
	Complete(ctx context.Context, id int) error
}

type CapituloRepository interface {
	FindAll(ctx context.Context) ([]*domain.Capitulo, error)
}

type IndicadorRepository interface {
	FindByCapitulo(ctx context.Context, idCapitulo int) ([]*domain.Indicador, error)
	FindBySegmento(ctx context.Context, idSegmento int) ([]int, error)
}

type NivelRespuestaRepository interface {
	FindByIndicador(ctx context.Context, idIndicador int) ([]*domain.NivelRespuesta, error)
}

type RespuestaRepository interface {
	Create(ctx context.Context, tx Transaction, respuesta *domain.Respuesta) (int, error)
	FindByAutoevaluacion(ctx context.Context, idAutoevaluacion int) ([]*domain.Respuesta, error)
	DeleteByAutoevaluacion(ctx context.Context, idAutoevaluacion int) error
}
