package apimethod

type Interactor interface {
	GetWhere(apiMethod *ApiMethod, where map[string]interface{}) error
}
