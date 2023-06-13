package storage

type ConfigStorage struct {
}

func NewConfigStorage() *ConfigStorage {
	return &ConfigStorage{}
}

func (s *ConfigStorage) Valid() bool {
	return true
}
