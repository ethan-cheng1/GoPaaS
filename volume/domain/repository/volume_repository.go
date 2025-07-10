package repository
import (
	"github.com/jinzhu/gorm"
	"git.imooc.com/coding-535/volume/domain/model"
)
// Interface that needs to be implemented
type IVolumeRepository interface{
    // Initialize table
    InitTable() error
    // Find data by ID
    FindVolumeByID(int64) (*model.Volume, error)
    // Create a volume record
	CreateVolume(*model.Volume) (int64, error)
    // Delete a volume record by ID
	DeleteVolumeByID(int64) error
    // Update data
	UpdateVolume(*model.Volume) error
    // Find all volume data
	FindAll()([]model.Volume,error)

}
// Create volumeRepository
func NewVolumeRepository(db *gorm.DB) IVolumeRepository  {
	return &VolumeRepository{mysqlDb:db}
}

type VolumeRepository struct {
	mysqlDb *gorm.DB
}

// Initialize table
func (u *VolumeRepository)InitTable() error  {
	return u.mysqlDb.CreateTable(&model.Volume{}).Error
}

// Find Volume information by ID
func (u *VolumeRepository)FindVolumeByID(volumeID int64) (volume *model.Volume,err error) {
	volume = &model.Volume{}
	return volume, u.mysqlDb.First(volume,volumeID).Error
}

// Create Volume information
func (u *VolumeRepository) CreateVolume(volume *model.Volume) (int64, error) {
	return volume.ID, u.mysqlDb.Create(volume).Error
}

// Delete Volume information by ID
func (u *VolumeRepository) DeleteVolumeByID(volumeID int64) error {
	return u.mysqlDb.Where("id = ?",volumeID).Delete(&model.Volume{}).Error
}

// Update Volume information
func (u *VolumeRepository) UpdateVolume(volume *model.Volume) error {
	return u.mysqlDb.Model(volume).Update(volume).Error
}

// Get result set
func (u *VolumeRepository) FindAll()(volumeAll []model.Volume,err error) {
	return volumeAll, u.mysqlDb.Find(&volumeAll).Error
}

