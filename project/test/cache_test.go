package test

import (
	"github.com/google/uuid"
	"project/api"
	"project/cmd/my_service/repository"
	"reflect"
	"testing"
)
import "project/cmd/my_service/cache"

func TestAddItemToCacheName(t *testing.T) {
	// arrange
	newUuid, _ := uuid.NewUUID()
	addedItem := api.Item{}
	item := repository.ItemRow{
		Uuid:    newUuid,
		Payload: addedItem,
	}

	var items []repository.ItemRow
	mp := cache.NewMapCache(items)

	// act
	mp.AddItem(item.Uuid, addedItem)

	// assert
	expectedItem, err := mp.GetItem(newUuid)
	if !reflect.DeepEqual(expectedItem, addedItem) {
		t.Errorf("expectedItem != addedItem: err %s", err)
		t.Fail()
	}
}
