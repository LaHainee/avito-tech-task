package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"avito-tech-task/internal/app/models"
)

type Validation struct {
	validate *validator.Validate
}

func NewValidator() *Validation {
	validate := validator.New()

	if err := validate.RegisterValidation("operation_type", validateOperationType); err != nil {
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

func validateOperationType(fl validator.FieldLevel) bool {
	return fl.Field().Int() <= models.REDUCE && fl.Field().Int() >= 0
}
