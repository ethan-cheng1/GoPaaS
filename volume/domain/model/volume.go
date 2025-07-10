package model

type Volume struct{
	ID int64 `gorm:"primary_key;not_null;auto_increment"`
	// Volume name
	VolumeName string `json:"volume_name"`
	// Volume namespace
	VolumeNamespace string `json:"volume_namespace"`
	// Volume access mode: RWO, ROX, RWX
	VolumeAccessMode string `json:"volume_access_mode"`
	// Storage class name
	VolumeStorageClassName string `json:"volume_storage_class_name"`
	// Requested resource size
	VolumeRequest float32 `json:"volume_request"`
	// Storage type: Block, filesystem
	VolumePersistentVolumeMode string `json:"volume_persistent_volume_mode"`
}

