package my_service

import (
	"context"
	"github.com/google/uuid"
	"project/api"
)

//go:generate mockgen -source=.\iservice.go -destination=mocks\service_mock.go

type IService interface {
	AddItem(ctx context.Context, item api.Item) (uuid.UUID, error)
	GetItem(ctx context.Context, uuid uuid.UUID) (api.Item, error)
}
