package validator

type CategoryValidator interface {
	ValidateCategoryName(name string) error
}
