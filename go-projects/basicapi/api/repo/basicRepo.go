package repo

import (
	"errors"
	"gorm.io/gorm"
)

type SQLRequestConfig struct {
	Associations []string
	Where        string
	Id           any
}

type basicRepository struct {
	DB *gorm.DB
}

func NewBasicRepository(conn *gorm.DB) IBasicRepository {
	return &basicRepository{
		DB: conn,
	}
}

type IBasicRepository interface {
	FindAll(ref any, index, size int, conf SQLRequestConfig) (int64, error)
	FindByID(ref any, conf SQLRequestConfig) error
	Create(ref any, tx *gorm.DB) error
	Update(ref any, tx *gorm.DB) error
	DeleteByID(ref any, conf SQLRequestConfig, tx *gorm.DB) error
}

func (v basicRepository) FindAll(ref any, index, size int, conf SQLRequestConfig) (int64, error) {
	var totalRecords int64
	query := v.DB.Model(ref)
	for _, association := range conf.Associations {
		query = query.Preload(association)
	}
	if conf.Where != "" {
		query = query.Where(conf.Where)
	}
	err := query.Count(&totalRecords).Offset((index) * size).Limit(size).Find(ref).Error

	return totalRecords, err
}

func (v basicRepository) FindByID(ref any, conf SQLRequestConfig) error {
	query := v.DB.Model(ref)
	for _, association := range conf.Associations {
		query = query.Preload(association)
	}
	if conf.Where != "" {
		query = query.Where(conf.Where, conf.Id)
	}
	return query.First(ref).Error
}

func (v basicRepository) Create(ref any, tx *gorm.DB) error {
	if tx == nil {
		return v.DB.Create(ref).Error
	}
	return tx.Create(ref).Error
}

func (v basicRepository) Update(ref any, tx *gorm.DB) error {
	if tx == nil {
		return v.DB.Updates(ref).Error
	}
	return tx.Updates(ref).Error
}

func (v basicRepository) DeleteByID(ref any, conf SQLRequestConfig, tx *gorm.DB) error {
	if conf.Where == "" {
		return errors.New("no where condition passed in")
	}
	if conf.Id == nil {
		return errors.New("no key specified")
	}
	if tx == nil {
		v.DB.Delete(ref, conf.Where, conf.Id)
	}
	return tx.Delete(ref, conf.Where, conf.Id).Error
}
