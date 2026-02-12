package service

import (
	"context"

	"coviar_backend/internal/repository/postgres"
)

type AdminService struct {
	adminRepo *postgres.AdminRepository
}

func NewAdminService(adminRepo *postgres.AdminRepository) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
	}
}

// GetStats obtiene las estadísticas generales del sistema
func (s *AdminService) GetStats(ctx context.Context) (*postgres.AdminStats, error) {
	return s.adminRepo.GetStats(ctx)
}

// GetAllEvaluaciones obtiene todas las evaluaciones con filtros opcionales
func (s *AdminService) GetAllEvaluaciones(ctx context.Context, estado string, idBodega *int) ([]postgres.EvaluacionListItem, error) {
	return s.adminRepo.GetAllEvaluaciones(ctx, estado, idBodega)
}
