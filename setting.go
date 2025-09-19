package go_setting

import (
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/parsidev/go_setting/models"
)

var (
	instance *settingData
)

type settingData struct {
	data map[string]any
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

	instance = &settingData{data: make(map[string]any), db: db}

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

func Set(data map[string]any) (err error) {
	for key, value := range data {
		var s models.Setting
		err = instance.db.Where(fmt.Sprintf("%v.key = ?", models.Setting{}.TableName()), key).First(&s).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = instance.db.Create(&models.Setting{
				Key:       key,
				PlainValue: value,
			}).Error; err != nil {
				return fmt.Errorf("failed to create setting '%s': %w", key, err)
			}
		} else if err == nil {
			s.PlainValue = value
			if err = instance.db.Save(&s).Error; err != nil {
				return fmt.Errorf("failed to update setting '%s': %w", key, err)
			}
		} else {
			return fmt.Errorf("failed to fetch setting '%s': %w", key, err)
		}

		instance.data[key] = value
	}

	return nil
}

func Get(key string, def any) (val any) {
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

	return true, nil
}

func GetAll() map[string]any {
	return instance.data
}
