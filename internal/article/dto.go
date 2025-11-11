package article

type CreateArticleRequest struct {
	Title    string `json:"title" validate:"required,min=20"`
	Content  string `json:"content" validate:"required,min=200"`
	Category string `json:"category" validate:"required,min=3"`
	Status   string `json:"status" validate:"required,oneof=publish draft thrash"`
}

type UpdateArticleRequest struct {
	Title    string `json:"title" validate:"omitempty,min=20"`
	Content  string `json:"content" validate:"omitempty,min=200"`
	Category string `json:"category" validate:"omitempty,min=3"`
	Status   string `json:"status" validate:"omitempty,oneof=publish draft thrash"`
}

type ListResponse struct {
	Items []Article      `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	Total   int64 `json:"total"`
	HasNext bool  `json:"has_next"`
}

// title (substring match), category (exact), dan status (exact).
type ListFilter struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Status   string `json:"status"`
}
