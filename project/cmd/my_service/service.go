package my_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"project/api"
	"project/cmd/my_service/cache"
	"project/cmd/my_service/repository"
)

type ItemService struct {
	repository repository.IRepository
	cache      cache.ICache
}

func (s *ItemService) GetItem(ctx context.Context, uuid uuid.UUID) (api.Item, error) {
	item, err := s.cache.GetItem(uuid)
	if err != nil {
		item, err = s.repository.GetItem(ctx, uuid)
		if err != nil {
			return api.Item{}, fmt.Errorf("can't get item from item service repository err : %w", err)
		}
	}
	return item, nil
}

func (s *ItemService) AddItem(ctx context.Context, item api.Item) (uuid.UUID, error) {
	marshal, _ := json.Marshal(item)
	log.Println("adding item %", marshal)
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil,
			fmt.Errorf("can't create uuid while AddItem in ICache. uuid.NewUUID() err:%w", err)
	}

	s.cache.AddItem(newUUID, item)

	s.repository.AddItem(ctx, newUUID, item)

	return newUUID, nil
}

func New(itemRepository repository.IRepository, cache cache.ICache) *ItemService {
	newItemStore := &ItemService{}
	newItemStore.cache = cache
	newItemStore.repository = itemRepository
	return newItemStore
}
