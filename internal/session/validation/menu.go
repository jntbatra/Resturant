package validation

// CreateMenuItemRequest represents the request to create a menu item
type CreateMenuItemRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description" validate:"max=1000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Category    string  `json:"category" validate:"required,min=1,max=100"`
	Status      string  `json:"status" validate:"required,oneof=available unavailable discontinued"`
}

// UpdateMenuItemRequest represents the request to update a menu item
type UpdateMenuItemRequest struct {
	Name        string  `json:"name" validate:"omitempty,min=1,max=255"`
	Description string  `json:"description" validate:"omitempty,max=1000"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
	Category    string  `json:"category" validate:"omitempty,min=1,max=100"`
}

// ListMenuItemsRequest represents the request to list menu items with pagination
type ListMenuItemsRequest struct {
	Offset   int    `json:"offset" validate:"min=0"`
	Limit    int    `json:"limit" validate:"required,min=1,max=100"`
}

// CreateCategoryRequest represents the request to create a category
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// CategoryIDByNameRequest represents the request to get category ID by name
type CategoryIDByNameRequest struct {
	Name string `query:"name" validate:"required,min=1,max=100"`
}





// ValidateCreateMenuItem validates the create menu item request
func ValidateCreateMenuItem(req CreateMenuItemRequest) error {
	return ValidateStruct(req)
}

// ValidateUpdateMenuItem validates the update menu item request
func ValidateUpdateMenuItem(req UpdateMenuItemRequest) error {
	return ValidateStruct(req)
}

// ValidateListMenuItems validates the list menu items request
func ValidateListMenuItems(req ListMenuItemsRequest) error {
	return ValidateStruct(req)
}

// ValidateCreateCategory validates the create category request
func ValidateCreateCategory(req CreateCategoryRequest) error {
	return ValidateStruct(req)
}

// ValidateUpdateCategory validates the update category request
func ValidateUpdateCategory(req UpdateCategoryRequest) error {
	return ValidateStruct(req)
}

// ValidateCategoryIDByName validates the category ID by name request
func ValidateCategoryIDByName(req CategoryIDByNameRequest) error {
	return ValidateStruct(req)
}
