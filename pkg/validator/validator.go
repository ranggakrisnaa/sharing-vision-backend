package validatorpkg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ranggakrisnaa/sharing-vision-backend/pkg/logger"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	v *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{v: validator.New()}
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Param   string `json:"param,omitempty"`
}

func (v *Validator) ValidateStructDetailed(s interface{}) ([]FieldError, error) {
	err := v.v.Struct(s)
	if err == nil {
		return nil, nil
	}
	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil, err
	}
	logger.Log.WithError(err).Error("validator error")

	// Build formatted error list
	formatted := make([]FieldError, 0, len(verrs))
	for _, fe := range verrs {
		fieldName := jsonFieldName(s, fe.StructField())
		msg := humanMessage(fieldName, fe)
		formatted = append(formatted, FieldError{
			Field:   fieldName,
			Message: msg,
			Tag:     fe.Tag(),
			Param:   fe.Param(),
		})
	}

	return formatted, nil
}

// json field return prettier error
func jsonFieldName(s interface{}, structField string) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return structField
	}
	f, ok := t.FieldByName(structField)
	if !ok {
		return structField
	}
	tag := f.Tag.Get("json")
	if tag == "-" || tag == "" {
		return strings.ToLower(structField)
	}

	// Get field name before comma (optional tag options)
	if idx := strings.Index(tag, ","); idx >= 0 {
		tag = tag[:idx]
	}
	return tag
}

// generate readable error message
func humanMessage(field string, fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi", field)
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", field, fe.Param())
	case "oneof":
		// Join options with commas for readability
		opts := strings.Join(strings.Fields(fe.Param()), ", ")
		return fmt.Sprintf("%s harus salah satu dari: %s", field, opts)
	default:
		return fmt.Sprintf("%s tidak valid (%s)", field, fe.Tag())
	}
}
