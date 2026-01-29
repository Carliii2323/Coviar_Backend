package main

import (
	"log"
	"net/http"

	"coviar_backend/internal/handler"
	"coviar_backend/internal/middleware"
	"coviar_backend/internal/repository/postgres"
	"coviar_backend/internal/service"
	"coviar_backend/pkg/config"
	"coviar_backend/pkg/database"
	"coviar_backend/pkg/httputil"
	"coviar_backend/pkg/router"
)

func main() {
	// 1. Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Error cargando configuración: %v", err)
	}
	log.Println("✓ Configuración cargada")

	// 2. Conectar a Supabase
	db, err := database.ConnectSupabase(cfg.Supabase.URL, cfg.Supabase.Key, cfg.Supabase.DBPassword)
	if err != nil {
		log.Fatalf("❌ Error conectando a Supabase: %v", err)
	}
	defer db.Close()
	log.Println("✓ Conexión a Supabase establecida")

	// 3. Inicializar repositorios
	bodegaRepo := postgres.NewBodegaRepository(db.DB)
	cuentaRepo := postgres.NewCuentaRepository(db.DB)
	responsableRepo := postgres.NewResponsableRepository(db.DB)
	ubicacionRepo := postgres.NewUbicacionRepository(db.DB)
	autoevaluacionRepo := postgres.NewAutoevaluacionRepository(db.DB)
	segmentoRepo := postgres.NewSegmentoRepository(db.DB)
	capituloRepo := postgres.NewCapituloRepository(db.DB)
	indicadorRepo := postgres.NewIndicadorRepository(db.DB)
	nivelRespuestaRepo := postgres.NewNivelRespuestaRepository(db.DB)
	respuestaRepo := postgres.NewRespuestaRepository(db.DB)
	txManager := postgres.NewTransactionManager(db.DB)

	log.Println("✓ Repositorios inicializados")

	// 4. Inicializar servicios
	registroService := service.NewRegistroService(bodegaRepo, cuentaRepo, responsableRepo, txManager)
	ubicacionService := service.NewUbicacionService(ubicacionRepo)
	cuentaService := service.NewCuentaService(cuentaRepo, bodegaRepo)
	bodegaService := service.NewBodegaService(bodegaRepo)
	autoevaluacionService := service.NewAutoevaluacionService(autoevaluacionRepo, segmentoRepo, capituloRepo, indicadorRepo, nivelRespuestaRepo, respuestaRepo)

	log.Println("✓ Servicios inicializados")

	// 5. Inicializar handlers (con JWT secret para autenticación)
	registroHandler := handler.NewRegistroHandler(registroService)
	ubicacionHandler := handler.NewUbicacionHandler(ubicacionService)
	cuentaHandler := handler.NewCuentaHandler(cuentaService, cfg.JWT.Secret)
	bodegaHandler := handler.NewBodegaHandler(bodegaService)
	autoevaluacionHandler := handler.NewAutoevaluacionHandler(autoevaluacionService)

	log.Println("✓ Handlers inicializados")

	// 6. Configurar router
	r := router.New()

	// Middlewares globales
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS)

	// ===== RUTAS PÚBLICAS =====

	// Registro y autenticación (no requieren autenticación)
	r.POST("/api/registro", registroHandler.RegistrarBodega)
	r.POST("/api/login", cuentaHandler.Login)

	// Logout (elimina cookies)
	r.POST("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🔓 Logout request recibido")

		// Eliminar cookie auth_token
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1, // Eliminar cookie
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		// Eliminar cookie refresh_token
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1, // Eliminar cookie
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		log.Printf("✅ Cookies eliminadas")
		httputil.RespondJSON(w, http.StatusOK, map[string]string{
			"mensaje": "Logout exitoso",
		})
	})

	// Ubicaciones (públicas - necesarias para registro)
	r.GET("/api/provincias", ubicacionHandler.GetProvincias)
	r.GET("/api/departamentos", ubicacionHandler.GetDepartamentos)
	r.GET("/api/localidades", ubicacionHandler.GetLocalidades)

	// Health check
	r.GET("/health", func(w http.ResponseWriter, r *http.Request) {
		httputil.RespondJSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"version": "2.0.0",
			"message": "Coviar Backend - Integrado y Funcional con JWT",
		})
	})

	// ===== RUTAS PROTEGIDAS (requieren autenticación) =====

	authMiddleware := middleware.AuthMiddleware(cfg.JWT.Secret)

	// Helper para convertir http.Handler a http.HandlerFunc
	protect := func(handler http.HandlerFunc) http.HandlerFunc {
		return authMiddleware(handler).ServeHTTP
	}

	// Cuentas (protegidas)
	r.GET("/api/cuentas/{id}", protect(cuentaHandler.GetByID))
	r.PUT("/api/cuentas/{id}", protect(cuentaHandler.UpdatePassword))

	// Bodegas (protegidas)
	r.GET("/api/bodegas/{id}", protect(bodegaHandler.GetByID))
	r.PUT("/api/bodegas/{id}", protect(bodegaHandler.Update))

	// Autoevaluaciones (protegidas)
	r.POST("/api/autoevaluaciones", protect(autoevaluacionHandler.CreateAutoevaluacion))
	r.GET("/api/autoevaluaciones/{id_autoevaluacion}/segmentos", protect(autoevaluacionHandler.GetSegmentos))
	r.PUT("/api/autoevaluaciones/{id_autoevaluacion}/segmento", protect(autoevaluacionHandler.SeleccionarSegmento))
	r.GET("/api/autoevaluaciones/{id_autoevaluacion}/estructura", protect(autoevaluacionHandler.GetEstructura))
	r.POST("/api/autoevaluaciones/{id_autoevaluacion}/respuestas", protect(autoevaluacionHandler.GuardarRespuestas))
	r.POST("/api/autoevaluaciones/{id_autoevaluacion}/completar", protect(autoevaluacionHandler.CompletarAutoevaluacion))

	// 7. Iniciar servidor
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("🚀 Servidor iniciando en http://%s", addr)
	log.Printf("📍 Entorno: %s", cfg.App.Environment)
	log.Printf("🔗 Supabase URL: %s", cfg.Supabase.URL)
	log.Printf("🔐 JWT Secret configurado: %s", maskSecret(cfg.JWT.Secret))
	log.Printf("🍪 Autenticación basada en cookies HttpOnly habilitada")

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}

// maskSecret enmascara el secret para no mostrarlo completo en logs
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "***"
	}
	return secret[:4] + "..." + secret[len(secret)-4:]
}
