package model

import "github.com/jinzhu/gorm"

// @dao
type StringModel struct {
  gorm.Model

  Key string
  Value string
}
