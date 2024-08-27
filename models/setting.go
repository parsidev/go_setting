package models

import (
	serilizer "github.com/parsidev/go_serializer"
	"gorm.io/gorm"
)

type Setting struct {
	Base
	Key        string      `json:"key"`
	Value      string      `json:"value"`
	PlainValue interface{} `json:"plain_value" gorm:"-:all"`
}

func (s *Setting) AfterFind(tx *gorm.DB) error {
	s.PlainValue = serilizer.UnSerialize(s.Value)
	return nil
}

func (s *Setting) BeforeSave(tx *gorm.DB) error {
	s.Value = serilizer.Serialize(s.PlainValue)
	return nil
}
