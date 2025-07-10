package service

import (
	"git.imooc.com/coding-535/base/domain/model"
	"git.imooc.com/coding-535/base/domain/repository"
)

// Interface type definition
type IBaseDataService interface {
	AddBase(*model.Base) (int64 , error)
	DeleteBase(int64) error
	UpdateBase(*model.Base) error
	FindBaseByID(int64) (*model.Base, error)
	FindAllBase() ([]model.Base, error)
}


// Create service instance
// Note: Returns IBaseDataService interface type
func NewBaseDataService(baseRepository repository.IBaseRepository) IBaseDataService{
	return &BaseDataService{ baseRepository }
}

type BaseDataService struct {
    // Note: This is IBaseRepository type
	BaseRepository repository.IBaseRepository
}


// Insert
func (u *BaseDataService) AddBase(base *model.Base) (int64 ,error) {
	 return u.BaseRepository.CreateBase(base)
}

// Delete
func (u *BaseDataService) DeleteBase(baseID int64) error {
	return u.BaseRepository.DeleteBaseByID(baseID)
}

// Update
func (u *BaseDataService) UpdateBase(base *model.Base) error {
	return u.BaseRepository.UpdateBase(base)
}

// Find by ID
func (u *BaseDataService) FindBaseByID(baseID int64) (*model.Base, error) {
	return u.BaseRepository.FindBaseByID(baseID)
}

// Find all
func (u *BaseDataService) FindAllBase() ([]model.Base, error) {
	return u.BaseRepository.FindAll()
}

