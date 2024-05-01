package validators

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
)

func HandleValidationError(err error) gin.H {
	var ve validator.ValidationErrors
	errs := make(gin.H)
	if errors.As(err, &ve) {
		for _, fe := range ve {
			fieldNameToLower := strings.ToLower(fe.Field())
			field := fe.Field()
			switch fe.Tag() {
			case "required":
				errs[field] = "The " + fieldNameToLower + " field is required"
			case "min":
				errs[field] = "The " + fieldNameToLower + " must be at least " + fe.Param() + " characters long"
			case "max":
				errs[field] = "The " + fieldNameToLower + " can have a maximum of " + fe.Param() + " characters"
			case "opt_str_len":
				params := strings.Split(fe.Param(), ",")
				if len(params) == 2 {
					errs[field] = "The " + fieldNameToLower + " must be between " + params[0] + " and " + params[1] + " characters long"
				} else {
					errs[field] = "Invalid parameter format for " + fieldNameToLower
				}
			default:
				errs[field] = "Validation failed on the " + fieldNameToLower + " with error: " + fe.Tag()
			}
		}
	}
	return errs
}

func stringLengthValidator(field *string, min, max int) bool {
	if field == nil {
		return true // if nil, it's valid as it's not required
	}
	length := len(strings.TrimSpace(*field))
	return length >= min && length <= max
}

func RegisterCustomValidations(v *validator.Validate) {
	err := v.RegisterValidation("opt_str_len", func(fl validator.FieldLevel) bool {
		params := strings.Split(fl.Param(), ",")
		minLen, err := strconv.Atoi(params[0])
		if err != nil {
			panic(errors.New("min length must be an integer"))
		}
		maxLen, err := strconv.Atoi(params[1])
		if err != nil {
			panic(errors.New("max length must be an integer"))
		}
		field, _ := fl.Field().Interface().(*string)
		return stringLengthValidator(field, minLen, maxLen)
	})
	if err != nil {
		panic(err)
	}
}
