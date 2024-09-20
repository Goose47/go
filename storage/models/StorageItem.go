package models

import (
	"Goose47/storage/config"
	"path"
)

type StorageItem struct {
	Key          string
	Ttl          int
	Path         string
	OriginalName string
}

func (i *StorageItem) GetFullPath() string {
	return path.Join(config.FSConfig.Base, i.Path)
}
