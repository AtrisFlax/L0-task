package cache

import (
	"fmt"
	"github.com/google/uuid"
	"math"
	"project/api"
	"project/cmd/my_service/repository"
	"sync"
)

type MapCache struct {
	cache map[uuid.UUID]api.Item
	sync.RWMutex
}

func NewMapCache(items []repository.ItemRow) *MapCache {
	mp := &MapCache{}
	mp.cache = make(map[uuid.UUID]api.Item, cacheSize(len(items)))

	for _, item := range items {
		mp.cache[item.Uuid] = item.Payload
	}

	return mp
}

func (mp *MapCache) GetItem(uuid uuid.UUID) (api.Item, error) {
	mp.RLock()
	defer mp.RUnlock()

	item, ok := mp.cache[uuid]
	if !ok {
		return api.Item{}, fmt.Errorf("item with id=%d not found in ICache", uuid)
	}

	return item, nil
}

func (mp *MapCache) AddItem(uuid uuid.UUID, item api.Item) {
	mp.Lock()
	defer mp.Unlock()

	mp.cache[uuid] = item
}

func cacheSize(repoSize int) int {
	near2powN := math.Log(float64(repoSize)) / math.Log(float64(2))
	size := int((math.Ceil(near2powN) + 1) * 2)
	return size
}
