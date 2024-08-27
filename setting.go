package go_setting

import (
	"database/sql"
	"github.com/parsidev/go_setting/models"
	"gorm.io/gorm"
)

var (
	instance *settingData
)

type settingData struct {
	data map[string]interface{}
	db   *gorm.DB
}

func Init(db *gorm.DB) (err error) {
	var (
		d *sql.DB
	)

	if d, err = db.DB(); d == nil || err != nil {
		return err
	}

	if err = d.Ping(); err != nil {
		return err
	}

	instance = &settingData{data: make(map[string]interface{}), db: db}

	if !db.Migrator().HasTable(&models.Setting{}) {
		if err = db.AutoMigrate(&models.Setting{}); err != nil {
			return err
		}
	}

	settings := make([]*models.Setting, 0)

	if err = db.Model(&models.Setting{}).Find(&settings).Error; err != nil {
		return err
	}

	for _, s := range settings {
		instance.data[s.Key] = s.PlainValue
	}

	return nil
}

func Set(data map[string]interface{}) {
	for key, value := range data {
		instance.data[key] = value
		s := new(models.Setting)
		_ = instance.db.First(&s, "key = ?", key).Error

		if !s.IsValid() {
			_ = instance.db.Create(&models.Setting{Key: key, PlainValue: value}).Error
		} else {
			s.PlainValue = value
			_ = instance.db.Save(s).Error
		}
	}
}

func Get(key string, def interface{}) (val interface{}) {
	var (
		ok bool
	)

	if val, ok = instance.data[key]; !ok || val == nil {
		val = def
	}

	return val
}

func Has(key string) (bool, error) {
	if _, ok := instance.data[key]; !ok {
		return false, ErrKeyNotFound
	}
}

func GetAll() map[string]interface{} {
	return instance.data
}
