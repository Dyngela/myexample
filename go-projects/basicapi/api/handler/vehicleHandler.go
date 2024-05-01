package handler

import (
	"basicapi/api/handler/request"
	"basicapi/api/handler/validators"
	"basicapi/types"
	"errors"
	"github.com/gin-contrib/graceful"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

func VehicleHandler(router *graceful.Graceful) {
	v1 := router.Group("demo")
	{
		v1.GET("/vehicules/:index/:size", findAllVehicle)
		// Kinda annoying but if we do not put an id tag in the url gin is actually not happy cuz it says path is ambiguous with the findAll one.
		v1.GET("/vehicules/id/:id", findVehicleById)
		v1.POST("/vehicules", createVehicle)
		v1.PUT("/vehicules", updateVehicle)
		v1.DELETE("/vehicules/:id", deleteVehicleById)
	}
}

// @Summary Find all vehicles
// @Description Get all vehicles in a paginated way
// @Tags vehicules
// @Produce json
// @Param index path int true "Index of the page. Start at 0"
// @Param size path int true "Size of the page"
// @Success 200 {array} []response.Vehicle
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /demo/vehicules/{index}/{size} [get]
func findAllVehicle(ctx *gin.Context) {
	index, err := strconv.Atoi(ctx.Param("index"))
	if err != nil || index < 0 {
		ctx.JSON(http.StatusBadRequest, types.APIError{Message: "Bad args"})
		return
	}
	size, err := strconv.Atoi(ctx.Param("size"))
	if err != nil || size <= 0 {
		ctx.JSON(http.StatusBadRequest, types.APIError{Message: "Bad args"})
		return
	}
	vehicles, err := vehiculeService.FindAll(index, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.APIError{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, vehicles)
}

// @Summary Find a vehicle by its id
// @Description Get a vehicle by its id
// @Tags vehicules
// @Produce json
// @Param id path int true "Id of the vehicle"
// @Success 200 {object} response.Vehicle
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /demo/vehicules/id/{id} [get]
func findVehicleById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, types.APIError{Message: "Bad args"})
		return
	}
	vehicle, err := vehiculeService.FindByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.APIError{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, vehicle)
}

// @Summary Create a vehicle
// @Description Create a vehicle
// @Tags vehicules
// @Accept json
// @Produce json
// @Param vehicle body request.CreateVehicle true "Vehicle to create"
// @Success 201 {object} types.GenericResponse
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /demo/vehicules [post]
func createVehicle(ctx *gin.Context) {
	var req request.CreateVehicle
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			errMap := validators.HandleValidationError(errs)
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": errMap})
			return
		}
		ctx.JSON(http.StatusBadRequest, types.APIError{Message: err.Error()})
		return
	}
	err := vehiculeService.Create(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.APIError{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, types.GenericResponse{Message: "Vehicle created"})
}

// @Summary Update a vehicle
// @Description Update a vehicle
// @Tags vehicules
// @Accept json
// @Produce json
// @Param vehicle body request.UpdateVehicle true "Vehicle to update"
// @Success 200 {object} types.GenericResponse
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /demo/vehicules [put]
func updateVehicle(ctx *gin.Context) {
	var req request.UpdateVehicle
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			errMap := validators.HandleValidationError(errs)
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": errMap})
			return
		}
		ctx.JSON(http.StatusBadRequest, types.APIError{Message: err.Error()})
		return
	}
	err := vehiculeService.Update(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.APIError{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, types.GenericResponse{Message: "Vehicle updated"})
}

// @Summary Delete a vehicle
// @Description Delete a vehicle
// @Tags vehicules
// @Produce json
// @Param id path int true "Id of the vehicle"
// @Success 200 {object} types.GenericResponse
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /demo/vehicules/{id} [delete]
func deleteVehicleById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, types.APIError{Message: "Bad args"})
		return
	}
	err = vehiculeService.DeleteByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.APIError{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, types.GenericResponse{Message: "Vehicle deleted"})
}
