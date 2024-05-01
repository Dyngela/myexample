package service

import (
	"basicapi/api/entity"
	"basicapi/api/handler/request"
	"basicapi/api/handler/response"
	"basicapi/api/repo"
	"basicapi/config"
	"basicapi/types"
)

type vehicleService struct {
	basicRepo repo.IBasicRepository
}

func NewVehicleService() IVehicleService {
	return &vehicleService{
		basicRepo: repo.NewBasicRepository(config.DB),
	}
}

type IVehicleService interface {
	FindAll(index, size int) (types.Page[response.Vehicle], error)
	FindByID(id uint) (response.Vehicle, error)
	Create(request request.CreateVehicle) error
	Update(request request.UpdateVehicle) error
	DeleteByID(id uint) error
}

func (v vehicleService) FindAll(index, size int) (types.Page[response.Vehicle], error) {
	var vehicles []entity.Vehicle
	totalRecords, err := v.basicRepo.FindAll(&vehicles, index, size, repo.SQLRequestConfig{})
	if err != nil {
		return types.Page[response.Vehicle]{}, err
	}

	var resp []response.Vehicle
	for _, item := range vehicles {
		var vehicule response.Vehicle
		vehicule.MapToDto(item)
		resp = append(resp, vehicule)
	}

	return types.Page[response.Vehicle]{
		TotalRecords: totalRecords,
		TotalPages:   types.CalculateTotalPages(totalRecords, size),
		CurrentPage:  index,
		Data:         resp,
	}, nil

}

func (v vehicleService) FindByID(id uint) (response.Vehicle, error) {
	var e entity.Vehicle
	var resp response.Vehicle

	err := v.basicRepo.FindByID(&e, repo.SQLRequestConfig{
		Associations: nil,
		Where:        "id = ?",
		Id:           id,
	})
	if err != nil {
		return response.Vehicle{}, err
	}
	resp.MapToDto(e)
	return resp, nil
}

func (v vehicleService) Create(request request.CreateVehicle) error {

	return v.basicRepo.Create(request.MapToEntity(), nil)
}

func (v vehicleService) Update(request request.UpdateVehicle) error {
	//We could, in case we want update to be transactional invoke a tx := config.DB and start a transaction from here.
	// Useful in case we want several update/create to be unified. If one fail everyone fail.
	tx := config.DB
	tx.Begin()
	err := v.basicRepo.Update(request.MapToEntity(), tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (v vehicleService) DeleteByID(id uint) error {
	return v.basicRepo.DeleteByID(entity.Vehicle{}, repo.SQLRequestConfig{
		Associations: nil,
		Where:        "id = ?",
		Id:           id,
	}, nil)
}
