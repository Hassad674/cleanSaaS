package storage

import (
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
)

type Service struct {
	storage service.StorageService
}

func NewService(storage service.StorageService) *Service {
	return &Service{storage: storage}
}

// Upload, Delete, GetURL will be implemented
