package web

type CategoryCreateRequest struct {
	Name string `validate:"required,max=255,min=1" json:"name"`
}

type CategoryUpdateRequest struct {
	Id   int    `validate:"required" json:"id"`
	Name string `validate:"required,max=255,min=1" json:"name"`
}
