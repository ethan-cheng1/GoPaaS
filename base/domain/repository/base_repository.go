package repository
import (
	"github.com/jinzhu/gorm"
	"git.imooc.com/coding-535/base/domain/model"
)
// Interface to be implemented
type IBaseRepository interface{
    // Initialize table
    InitTable() error
    // Find data by ID
    FindBaseByID(int64) (*model.Base, error)
    // Create a base record
	CreateBase(*model.Base) (int64, error)
    // Delete a base record by ID
	DeleteBaseByID(int64) error
    // Update data
	UpdateBase(*model.Base) error
    // Find all base records
	FindAll()([]model.Base,error)

}
// Create baseRepository
func NewBaseRepository(db *gorm.DB) IBaseRepository  {
	return &BaseRepository{mysqlDb:db}
}

type BaseRepository struct {
	mysqlDb *gorm.DB
}

// Initialize table
func (u *BaseRepository)InitTable() error  {
	return u.mysqlDb.CreateTable(&model.Base{}).Error
}

// Find Base by ID
func (u *BaseRepository)FindBaseByID(baseID int64) (base *model.Base,err error) {
	base = &model.Base{}
	return base, u.mysqlDb.First(base,baseID).Error
}

// Create Base
func (u *BaseRepository) CreateBase(base *model.Base) (int64, error) {
	return base.ID, u.mysqlDb.Create(base).Error
}

// Delete Base by ID
func (u *BaseRepository) DeleteBaseByID(baseID int64) error {
	return u.mysqlDb.Where("id = ?",baseID).Delete(&model.Base{}).Error
}

// Update Base
func (u *BaseRepository) UpdateBase(base *model.Base) error {
	return u.mysqlDb.Model(base).Update(base).Error
}

// Get all results
func (u *BaseRepository) FindAll()(baseAll []model.Base,err error) {
	return baseAll, u.mysqlDb.Find(&baseAll).Error
}

