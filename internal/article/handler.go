package article

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/ranggakrisnaa/sharing-vision-backend/pkg/response"
	validatorpkg "github.com/ranggakrisnaa/sharing-vision-backend/pkg/validator"
)

type Handler struct {
	svc       *Service
	validator *validatorpkg.Validator
}

func NewHandler(svc *Service, validator *validatorpkg.Validator) *Handler {
	return &Handler{svc: svc, validator: validator}
}

func (h *Handler) Register(r fiber.Router) {
	r.Post("/", h.create)
	r.Get("/", h.list)
	r.Get("/:id", h.getByID)
	r.Put("/:id", h.update)
	r.Delete("/:id", h.delete)
}

func (h *Handler) create(c *fiber.Ctx) error {
	var req CreateArticleRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid JSON body")
	}
	errors, _ := h.validator.ValidateStructDetailed(req)
	if len(errors) > 0 {
		return response.Fail(c, fiber.StatusUnprocessableEntity, errors)
	}
	art, err := h.svc.Create(c.Context(), req)
	if err != nil {
		return response.Fail(c, fiber.StatusBadRequest, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, art, "article created successfully")
}

func (h *Handler) list(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// Filters
	var filter ListFilter
	if st := strings.ToLower(c.Query("status")); st != "" {
		switch st {
		case "publish", "draft", "thrash":
			filter.Status = st
		default:
			return response.Fail(c, fiber.StatusBadRequest, "status filter invalid: pilih publish | draft | thrash")
		}
	}
	if cat := strings.TrimSpace(c.Query("category")); cat != "" {
		filter.Category = cat
	}
	if title := strings.ToLower(c.Query("title")); title != "" {
		filter.Title = title
	}

	items, meta, err := h.svc.List(c.Context(), limit, page, filter)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, map[string]interface{}{
		"items": items,
		"meta":  meta,
	}, "articles retrieved successfully")
}

func (h *Handler) getByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Fail(c, fiber.StatusUnprocessableEntity, "id harus integer")
	}

	art, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.Fail(c, fiber.StatusNotFound, "article tidak ditemukan")
		}
		return response.Fail(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, art, "article retrieved successfully")
}

func (h *Handler) update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Fail(c, fiber.StatusUnprocessableEntity, "id harus integer")
	}

	var req UpdateArticleRequest
	if errBody := c.BodyParser(&req); errBody != nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid JSON body")
	}

	errors, _ := h.validator.ValidateStructDetailed(req)
	if len(errors) > 0 {
		return response.Fail(c, fiber.StatusUnprocessableEntity, errors)
	}

	art, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.Fail(c, fiber.StatusNotFound, "article tidak ditemukan")
		}
		return response.Fail(c, fiber.StatusBadRequest, err.Error())
	}
	return response.Success(c, fiber.StatusOK, art, "article updated successfully")
}

func (h *Handler) delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Fail(c, fiber.StatusUnprocessableEntity, "id harus integer")
	}

	err = h.svc.Delete(c.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.Fail(c, fiber.StatusNotFound, "article tidak ditemukan")
		}
		return response.Fail(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, nil, "article deleted successfully")
}
