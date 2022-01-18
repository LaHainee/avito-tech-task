package utils

import (
	"avito-tech-task/internal/pkg/constants"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Validation struct {
	validate *validator.Validate
}

func NewValidator() *Validation {
	validate := validator.New()

	if err := validate.RegisterValidation("operation_type", validateOperationTypeUpdateBalance); err != nil {
		logrus.Fatalf("Could not register validation func by \"operation_type\" tag: %s", err)
	}

	return &Validation{validate: validate}
}

func (v *Validation) Validate(i interface{}) validator.ValidationErrors {
	if err := v.validate.Struct(i); err != nil {
		//nolint:errorlint
		return err.(validator.ValidationErrors)
	}

	return nil
}

func validateOperationTypeUpdateBalance(fl validator.FieldLevel) bool {
	return fl.Field().Int() <= constants.REDUCE && fl.Field().Int() >= constants.ADD
}
