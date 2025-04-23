package common

type AppResponse[T any, P any, F any] struct {
	Data   T `json:"data"`
	Paging P `json:"paging,omitempty"`
	Filter F `json:"filter,omitempty"`
}

func NewSuccessResponse[T, P, F any](data T, paging P, filter F) *AppResponse[T, P, F] {
	return &AppResponse[T, P, F]{Data: data, Paging: paging, Filter: filter}

}

func SimpleSuccessResponse[T any](data T) *AppResponse[T, any, any] {
	return NewSuccessResponse[T, any, any](data, nil, nil)
}
