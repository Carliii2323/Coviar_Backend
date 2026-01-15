// RUTA: coviar-backend/internal/platform/database/postgres.go
package database

import (
	supa "github.com/supabase-community/supabase-go"
)

// Connect establece conexión con Supabase
func Connect(url, key string) (*supa.Client, error) {
	client, err := supa.NewClient(url, key, &supa.ClientOptions{})
	if err != nil {
		return nil, err
	}

	return client, nil
}
