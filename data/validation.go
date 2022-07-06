package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidationError wraps the validator's FieldError to avoid exposing this to out code
type ValidationError struct {
	validator.FieldError
}

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"[KEY]: '%s' ERROR: Field Validation for '%s' failed on tag '%s'",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

// Errors converst the slice into a string slice
func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

// Validation contains
type Validation struct {
	validate *validator.Validate
}

// NewValidation creates a new Validation instance
func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("sku", SKUValidation)
	return &Validation{validate}
}

// Validate the item
// for more details the returned error can be cast into a  validator.ValidationErrors collection
//
// if ve, ok := err.(validator.ValidationErrors); ok {
// 						fmt.Println(ve.Namespace())
// 						fmt.Println(ve.Field())
// 						fmt.Println(ve.StructNamespace())
// 						fmt.Println(ve.StructField())
// 						fmt.Println(ve.Tag())
// 						fmt.Println(ve.ActualTag())
// 						fmt.Println(ve.Kind())
// 						fmt.Println(ve.Type())
// 						fmt.Println(ve.Value())
// 						fmt.Println(ve.Param())
// 						fmt.Println()
// }
func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i).(validator.ValidationErrors)

	if len(errs) == 0 {
		return nil
	}
	var returnErrs []ValidationError
	for _, err := range errs {
		// cast the FieldError into our ValidationError and append to the slice
		ve := ValidationError{err.(validator.FieldError)}
		returnErrs = append(returnErrs, ve)
	}
	return returnErrs
}

// SKUValidation
// Custom validation to the product structure
func SKUValidation(fl validator.FieldLevel) bool {
	// SKU format is abc-def-ghi
	// This is a very basic format for validating and the regex below will work only on this format

	regex := regexp.MustCompile(`[a-zA-Z]+-[a-zA-Z]+-[a-zA-Z]+`)
	matches := regex.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}
	return true
}
