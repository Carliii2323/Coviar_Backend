// RUTA: coviar-backend/internal/domain/bodega.go
package domain

// Bodega representa una bodega en el sistema COVIAR
type Bodega struct {
	IdBodega       int     `json:"idBodega"`
	Cuit           int64   `json:"cuit"`
	Inv            int     `json:"inv"`
	ViñedosInv     int     `json:"viñedos_inv"`
	Nombre         string  `json:"nombre"`
	Ubicacion      string  `json:"ubicacion"`
	ContactoEmail  string  `json:"contacto_email"`
	RazonSocial    string  `json:"razon_social"`
	NombreFantasia string  `json:"nombre_fantasia"`
	CreatedAt      *string `json:"created_at"`
}
