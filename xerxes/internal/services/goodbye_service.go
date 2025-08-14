package services

type GoodbyeService struct{}

func NewGoodbyeService() *GoodbyeService {
    return &GoodbyeService{}
}

func (s *GoodbyeService) GoodbyeWorld() string {
    return "Goodbye, World!"
}
