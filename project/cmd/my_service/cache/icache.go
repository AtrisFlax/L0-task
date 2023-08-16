package cache

import (
	"github.com/google/uuid"
	"project/api"
)

type ICache interface {
	GetItem(uuid uuid.UUID) (api.Item, error)
	AddItem(uuid uuid.UUID, item api.Item)
}
