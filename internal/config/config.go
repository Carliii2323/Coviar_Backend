// RUTA: coviar-backend/internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	SupabaseURL string
	SupabaseKey string
	Port        string
}

// Load carga las variables de entorno desde .env
func Load() *Config {
	// Cargar .env (opcional)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No se encontró .env, usando variables del sistema")
	}

	cfg := &Config{
		SupabaseURL: os.Getenv("SUPABASE_URL"),
		SupabaseKey: os.Getenv("SUPABASE_KEY"),
		Port:        os.Getenv("APP_PORT"),
	}

	// Validar variables críticas
	if cfg.SupabaseURL == "" || cfg.SupabaseKey == "" {
		log.Fatal("❌ ERROR: SUPABASE_URL y SUPABASE_KEY son requeridas")
	}

	return cfg
}
