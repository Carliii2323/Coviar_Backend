// RUTA: coviar-backend/cmd/api/main.go
// REEMPLAZA EL ARCHIVO ANTERIOR - VERSION SIN LOCALIDAD
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/carli/coviar-backend/internal/bodega"
	"github.com/carli/coviar-backend/internal/config"
	"github.com/carli/coviar-backend/internal/platform/database"
	"github.com/carli/coviar-backend/internal/usuario"
)

func main() {
	// 1. Cargar configuración
	cfg := config.Load()
	fmt.Println("🔧 Configuración cargada")

	// 2. Conectar a base de datos
	db, err := database.Connect(cfg.SupabaseURL, cfg.SupabaseKey)
	if err != nil {
		log.Fatal("❌ Error al conectar con Supabase:", err)
	}
	fmt.Println("✅ Conectado a Supabase")

	// 3. Inicializar módulos

	// Módulo Bodega
	bodegaRepo := bodega.NewRepository(db)
	bodegaService := bodega.NewService(bodegaRepo)
	bodegaHandler := bodega.NewHandler(bodegaService)

	// Módulo Usuario
	usuarioRepo := usuario.NewRepository(db)
	usuarioService := usuario.NewService(usuarioRepo)
	usuarioHandler := usuario.NewHandler(usuarioService)

	// 4. Configurar rutas

	// Rutas generales
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)

	// Rutas de Bodega
	http.HandleFunc("/api/bodegas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			bodegaHandler.ListBodegas(w, r)
		} else if r.Method == http.MethodPost {
			bodegaHandler.CreateBodega(w, r)
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/bodegas/", bodegaHandler.GetBodega)

	// Rutas de Usuario
	http.HandleFunc("/api/usuarios", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			usuarioHandler.ListAll(w, r)
		} else if r.Method == http.MethodPost {
			usuarioHandler.Create(w, r)
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/usuarios/verificar", usuarioHandler.Verify)
	http.HandleFunc("/api/usuarios/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			usuarioHandler.GetByID(w, r)
		} else if r.Method == http.MethodDelete {
			usuarioHandler.Deactivate(w, r)
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// 5. Iniciar servidor
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	fmt.Printf("\n🚀 Servidor corriendo en http://localhost:%s\n", port)
	fmt.Println("📖 Endpoints disponibles:")
	fmt.Println()
	fmt.Println("   GENERAL:")
	fmt.Println("   GET    /health                      - Estado del servidor")
	fmt.Println()
	fmt.Println("   BODEGAS:")
	fmt.Println("   GET    /api/bodegas                 - Listar bodegas")
	fmt.Println("   GET    /api/bodegas/{id}            - Obtener bodega por ID")
	fmt.Println("   POST   /api/bodegas                 - Crear bodega")
	fmt.Println()
	fmt.Println("   USUARIOS:")
	fmt.Println("   POST   /api/usuarios                - Crear usuario")
	fmt.Println("   POST   /api/usuarios/verificar      - Verificar credenciales")
	fmt.Println("   GET    /api/usuarios                - Listar usuarios")
	fmt.Println("   GET    /api/usuarios/{id}           - Obtener usuario por ID")
	fmt.Println("   DELETE /api/usuarios/{id}           - Dar de baja usuario")
	fmt.Println()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message":"API COVIAR - Clean Architecture","version":"2.0.0","endpoints":{
		"usuarios":"/api/usuarios",
		"bodegas":"/api/bodegas"
	}}`)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","database":"connected","version":"2.0.0"}`)
}
