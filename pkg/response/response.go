package response

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	validatorpkg "github.com/ranggakrisnaa/sharing-vision-backend/pkg/validator"
)

type Meta struct {
	Limit   int   `json:"limit"`
	Page    int   `json:"page"`
	Total   int64 `json:"total"`
	HasNext bool  `json:"has_next"`
}

type Response struct {
	Success bool          `json:"success"`
	Data    interface{}   `json:"data,omitempty"`
	Message string        `json:"message,omitempty"`
	Error   string        `json:"error,omitempty"`
	Errors  []interface{} `json:"errors,omitempty"`
}

func Success(ctx *fiber.Ctx, status int, data interface{}, message string) error {
	return ctx.Status(status).JSON(Response{Success: true, Data: data, Message: message})
}

func Fail(ctx *fiber.Ctx, status int, msg interface{}) error {
	var errors []interface{}
	var errorMsg string

	switch v := msg.(type) {
	case []interface{}:
		errors = v
	case []validatorpkg.FieldError:
		for _, fe := range v {
			errors = append(errors, fe)
		}
	case string:
		errorMsg = v
	default:
		errorMsg = fmt.Sprintf("%v", v)
	}

	return ctx.Status(status).JSON(Response{
		Success: false,
		Error:   errorMsg,
		Errors:  errors,
	})
}

func PageMeta(limit, offset int, total int64) *Meta {
	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}
	return &Meta{
		Limit:   limit,
		Page:    page,
		Total:   total,
		HasNext: int64(offset+limit) < total,
	}
}
