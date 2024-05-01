package request

import "basicapi/api/entity"

type CreateVehicle struct {
	Immat   string `json:"immat" binding:"required,min=10,max=10"`
	Chassis string `json:"chassis" binding:"required,min=15,max=15"`
}

func (v *CreateVehicle) MapToEntity() entity.Vehicle {
	return entity.Vehicle{
		Immat:   v.Immat,
		Chassis: v.Chassis,
	}
}

type UpdateVehicle struct {
	Immat   *string `json:"immat" validate:"opt_str_len=10,10"`
	Chassis *string `json:"chassis" validate:"opt_str_len=15,15"`
}

func (v *UpdateVehicle) MapToEntity() entity.Vehicle {
	var vehicle entity.Vehicle
	if v.Immat != nil {
		vehicle.Immat = *v.Immat
	}
	if v.Chassis != nil {
		vehicle.Chassis = *v.Chassis
	}
	return vehicle
}
