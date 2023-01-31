package service

import (
	"L0task/pkg/model"
	"L0task/pkg/repository"
)

type ItemService struct {
	repo repository.Item
}

func NewItemService(repo repository.Item) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) SetItem(item model.Item) (int, error) {
	return s.repo.SetItem(item)
}

func (s *ItemService) GetItemById(itemId int) (model.Item, error) {
	return s.repo.GetItemById(itemId)
}
func (s *ItemService) GetAllItems() ([]model.Item, error) {
	return s.repo.GetAllItems()
}

func (s *ItemService) DeleteItem(itemId int) error {
	return s.repo.DeleteItem(itemId)
}
