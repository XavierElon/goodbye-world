package services

type HealthService struct{}

func NewHealthService() *HealthService {
    return &HealthService{}
}

func (s *HealthService) Health() string {
    return "OK"
} 