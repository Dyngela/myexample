package entity

import "gorm.io/gorm"

type Vehicle struct {
	gorm.Model
	Chassis string `json:"chassis"`
	Immat   string `json:"immat"`
}
