package response

import "basicapi/api/entity"

type Vehicle struct {
	Immat   string `json:"immat"`
	Chassis string `json:"chassis"`
}

func (v *Vehicle) MapToDto(vehicle entity.Vehicle) {
	v.Immat = vehicle.Immat
	v.Chassis = vehicle.Chassis
}
