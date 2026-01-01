# Implementation Checklist: DTOs & Validation Pattern

This checklist ensures every handler follows the clean architecture pattern.

## ✅ What has been implemented:

### Session Module
- [x] DTOs: CreateSessionRequest, UpdateSessionRequest, ListSessionsRequest, ChangeSessionTableRequest
- [x] Validation functions: ValidateCreateSession, ValidateUpdateSession, ValidateListSessions, ValidateChangeSessionTable
- [x] Handlers: CreateSession, GetSession, UpdateSession, ListSessions, ListActiveSessions, ChangeSessionTable
- [x] Service: Clean layer with only business logic (state transitions)
- [x] Removed duplicate validation from service layer

### Menu Module
- [x] DTOs: CreateMenuItemRequest, UpdateMenuItemRequest, ListMenuItemsRequest, CreateCategoryRequest, UpdateCategoryRequest
- [x] Validation functions: All DTOs use ValidateStruct only
- [x] Handlers: CreateMenuItem, GetMenuItem, ListMenuItems, UpdateMenuItem, DeleteMenuItem, ListCategories, CreateCategory, UpdateCategory, DeleteCategory
- [x] Service: Cleaned up to keep only business logic (category existence check)
- [x] Removed duplicate field validation from service layer

### Order Module
- [x] DTOs: CreateOrderRequest, UpdateOrderRequest, ListOrdersRequest, CreateOrderItemRequest, UpdateOrderItemRequest
- [x] Validation functions: All DTOs use ValidateStruct only
- [x] Handlers: CreateOrder, GetOrder, ListOrders, UpdateOrder, CreateOrderItem, GetOrderItems, GetOrdersBySession, GetOrderItemsBySessionIDs
- [x] Service: Cleaned up to keep only business logic (menu item availability check)
- [x] Removed duplicate validation from service layer

### Validator
- [x] Singleton pattern with sync.Once
- [x] Thread-safe initialization
- [x] Single validator instance for all handlers

## 📋 Handler Template (Copy for new operations)

```go
// OperationName handles HTTP REQUEST /path
func (h *Handler) OperationName(c *gin.Context) {
    // 1. Parse request into DTO
    var req validation.DtoNameRequest
    if err := c.ShouldBindJSON(&req); err != nil {  // or ShouldBindQuery
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 2. Validate using struct tags
    if err := validation.ValidateDtoName(req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 3. Extract path params if needed
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid ID"})
        return
    }

    // 4. Call service (service will do business logic validation)
    result, err := h.svc.ServiceMethod(req.Field1, req.Field2)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // 5. Return response
    c.JSON(200, result)
}
```

## 📋 DTO Template (Copy for new requests)

```go
// DtoNameRequest represents the request for operation
type DtoNameRequest struct {
    Field1 string    `json:"field_1" validate:"required,min=1,max=255"`
    Field2 int       `json:"field_2" validate:"required,gt=0"`
    Field3 uuid.UUID `json:"field_3" validate:"required"`
    Field4 string    `json:"field_4" validate:"oneof=value1 value2 value3"`
    Field5 int       `json:"field_5" validate:"min=0,max=100"`
}

// ValidateDtoName validates the request
func ValidateDtoName(req DtoNameRequest) error {
    return ValidateStruct(req)
}
```

## 📋 Service Template (Clean layer)

```go
// ServiceMethod does the actual business logic
func (s *serviceImpl) ServiceMethod(field1 string, field2 int) (Result, error) {
    // ONLY business logic here:
    
    // 1. Check state/resource existence
    existing, err := s.repo.GetById(field1)
    if err != nil {
        return nil, err
    }
    if existing == nil {
        return nil, errors.New("resource not found")
    }
    
    // 2. Validate state transitions
    if existing.Status != "active" {
        return nil, errors.New("cannot perform action in current state")
    }
    
    // 3. Check cross-resource constraints
    related, err := s.relatedService.GetRelated(existing.RelatedID)
    if err != nil {
        return nil, errors.New("related resource not found")
    }
    
    // 4. Create/update resource
    result := &Result{
        Field1: field1,
        Field2: field2,
    }
    
    // 5. Persist
    err = s.repo.Save(result)
    if err != nil {
        return nil, err
    }
    
    return result, nil
}
```

## 🔍 Validation Struct Tag Reference

| Rule | Syntax | Example |
|------|--------|---------|
| Required | `required` | `validate:"required"` |
| Greater than | `gt=N` | `validate:"gt=0"` |
| Less than | `lt=N` | `validate:"lt=100"` |
| Min value | `min=N` | `validate:"min=1"` |
| Max value | `max=N` | `validate:"max=100"` |
| Min length | `min=N` | `validate:"min=1"` |
| Max length | `max=N` | `validate:"max=255"` |
| Enum | `oneof=v1 v2 v3` | `validate:"oneof=active pending completed"` |
| Optional | `omitempty` | `validate:"omitempty,min=1"` |
| Multiple | Combine with comma | `validate:"required,min=1,max=255"` |

## 🧪 Example: Complete Flow for Creating an Order

### 1. DTO Definition (validation/order.go)
```go
type CreateOrderRequest struct {
    SessionID uuid.UUID `json:"session_id" validate:"required"`
}

func ValidateCreateOrder(req CreateOrderRequest) error {
    return ValidateStruct(req)
}
```

### 2. Handler (handler/order.go)
```go
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    // Parse
    var req validation.CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Validate (checks UUID not nil)
    if err := validation.ValidateCreateOrder(req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Call service (no shape validation needed here)
    err := h.svc.CreateOrder(req.SessionID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, gin.H{"message": "Order created successfully"})
}
```

### 3. Service (service/order.go)
```go
func (s *orderService) CreateOrder(sessionID uuid.UUID) error {
    // ONLY business logic - no UUID nil check needed (already done)
    
    order := &models.Order{
        ID:        uuid.New(),
        SessionID: sessionID,
        Status:    "cart",
        CreatedAt: time.Now(),
    }
    return s.repo.CreateOrder(order)
}
```

### 4. Result
✅ Shape validation: Handled by DTO + ValidateStruct
✅ Business logic: Only in service
✅ Clean separation: Handler → Service → Repository
✅ No duplication: Each check happens exactly once

## 📌 Key Rules

1. **Never validate shape in service** - Let struct tags handle it
2. **Always use DTOs in handlers** - Decouple from response models
3. **Keep validation functions simple** - Most just call ValidateStruct()
4. **Business logic only in service** - State, existence, constraints
5. **One validator instance** - Use lazy singleton pattern
6. **Clear error messages** - Tell client what to fix

## 🚀 Next Steps for Remaining Handlers

All handlers have been updated! The following stub methods remain and need service implementation:
- Menu: UpdateCategory, DeleteCategory
- Order: UpdateOrderItem, DeleteOrderItem, GetOrderItemsByOrderIDs

When implementing these:
1. Create DTO in validation package
2. Add handler with parse → validate → service pattern
3. Add service method with business logic only
4. Update service interface
