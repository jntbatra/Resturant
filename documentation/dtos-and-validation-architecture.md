# DTOs and Validation Architecture

## Overview

This document describes the three-layer validation and DTO pattern implemented in the restaurant service:

1. **Shape Validation Layer** (Transport)
2. **Business Logic Layer** (Service)
3. **Data Access Layer** (Repository)

## Architecture Pattern

### Layer 1: Shape Validation (Handler)

**Responsibility**: Ensure the shape, format, and ranges of input data.

**Where it lives**: `handler/` package

**How it works**:
1. Handler receives HTTP request
2. Parse JSON into DTO using `c.ShouldBindJSON(&req)` or `c.ShouldBindQuery(&req)`
3. Call validation function: `validation.ValidateCreateSession(req)`
4. ValidateStruct uses struct tags with `go-playground/validator` to check:
   - Required fields: `validate:"required"`
   - Number ranges: `validate:"gt=0"` (greater than), `validate:"min=0"`, `validate:"max=100"`
   - String lengths: `validate:"min=1,max=255"`
   - Enums: `validate:"oneof=active completed pending cancelled"`
   - UUID not nil

**Example Handler**:
```go
// CreateSession handles POST /sessions
func (h *Handler) CreateSession(c *gin.Context) {
    var req validation.CreateSessionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Validation checks format, ranges, required fields
    if err := validation.ValidateCreateSession(req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Once validated, pass to service
    session, err := h.svc.CreateSession(req.TableID)
    // ...
}
```

### Layer 2: Business Logic (Service)

**Responsibility**: Enforce rules that depend on state or business meaning.

**Where it lives**: `service/` package

**What stays in service**:
- State transitions (e.g., "Can't change completed session to active")
- Resource existence checks (e.g., "Category must exist")
- Cross-resource validation (e.g., "Menu item must be available")
- Orchestration logic

**What does NOT stay in service**:
- Field required checks (already done by DTO validation)
- Range validation (already done by struct tags)
- String length validation (already done by struct tags)
- Enum validation (already done by struct tags)

**Example Service**:
```go
// UpdateSession updates the status of a session
func (s *sessionService) UpdateSession(id uuid.UUID, status models.SessionStatus) (*models.Session, error) {
    // Get current session to validate state transition (BUSINESS LOGIC)
    currentSession, err := s.repo.GetSession(id)
    if err != nil {
        return nil, err
    }

    // Validate state transitions (BUSINESS LOGIC - cannot change completed/cancelled sessions)
    validTransitions := map[models.SessionStatus][]models.SessionStatus{
        "active":    {"pending", "cancelled"},
        "pending":   {"completed", "cancelled"},
        "completed": {},
        "cancelled": {},
    }

    allowedStatuses, exists := validTransitions[currentSession.Status]
    if !exists {
        return nil, errors.New("invalid current status")
    }

    // ... validation continues ...
    
    // Shape validation (format, ranges) already done by handler using ValidateStruct
    err = s.repo.UpdateSession(id, status)
    if err != nil {
        return nil, err
    }
    return s.repo.GetSession(id)
}
```

### DTOs (Data Transfer Objects)

DTOs are defined in the `validation/` package and used as transport models between frontend and backend.

**Benefits**:
- Decouples frontend from database models
- Centralizes shape validation rules
- Works with any frontend (HTTP, gRPC, CLI)

**DTO Naming**:
- `CreateXxxRequest` - Create operations
- `UpdateXxxRequest` - Update operations
- `ListXxxRequest` - List with pagination
- `ChangeXxxRequest` - Other operations

**Example DTO with Validation Tags**:
```go
// Session DTOs
type CreateSessionRequest struct {
    TableID int `json:"table_id" validate:"required,gt=0"`
}

type UpdateSessionRequest struct {
    Status models.SessionStatus `json:"status" validate:"required,oneof=active completed pending cancelled"`
}

type ListSessionsRequest struct {
    Offset int `json:"offset" validate:"min=0"`
    Limit  int `json:"limit" validate:"required,min=1,max=100"`
}

// Menu DTOs
type CreateMenuItemRequest struct {
    Name        string  `json:"name" validate:"required,min=1,max=255"`
    Description string  `json:"description" validate:"max=1000"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Category    string  `json:"category" validate:"required,min=1,max=100"`
}

type ListMenuItemsRequest struct {
    Offset   int    `json:"offset" validate:"min=0"`
    Limit    int    `json:"limit" validate:"required,min=1,max=100"`
    Category string `json:"category" validate:"omitempty,min=1,max=100"`
}

// Order DTOs
type CreateOrderRequest struct {
    SessionID uuid.UUID `json:"session_id" validate:"required"`
}

type UpdateOrderRequest struct {
    Status models.OrderStatus `json:"status" validate:"required,oneof=cart pending preparing served cancelled"`
}

type ListOrdersRequest struct {
    Offset    int       `json:"offset" validate:"min=0"`
    Limit     int       `json:"limit" validate:"required,min=1,max=100"`
    SessionID uuid.UUID `json:"session_id" validate:"omitempty"`
}
```

## Validator Implementation

The validator uses a **lazy singleton pattern** with `sync.Once` to ensure thread-safety and single initialization:

```go
package validation

import (
    "sync"
    "github.com/go-playground/validator/v10"
)

var (
    validate *validator.Validate
    once     sync.Once
)

// Init initializes the validator
func Init() {
    once.Do(func() {
        validate = validator.New()
    })
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
    if validate == nil {
        Init()
    }
    return validate
}

// ValidateStruct validates a struct using the validator
func ValidateStruct(s interface{}) error {
    return GetValidator().Struct(s)
}
```

**Why this pattern**:
- Thread-safe: `sync.Once` ensures init happens exactly once
- Lazy: Validator created only when first needed
- Reusable: All handlers use the same validator instance
- Efficient: No repeated allocations

## Request Flow Example: Create Session

```
HTTP Request (POST /sessions)
    ↓
Handler.CreateSession()
    ↓
1. Parse JSON → CreateSessionRequest DTO
    ├─ TableID: string → int
    └─ Unknown fields rejected
    ↓
2. Validate shape using struct tags
    ├─ Check TableID != nil
    ├─ Check TableID > 0
    └─ Return error if invalid
    ↓
3. Call Service.CreateSession(tableID)
    ├─ REMOVED: Check TableID > 0 (already done)
    ├─ Generate UUID for session
    └─ Call Repository.CreateSession()
    ↓
4. Repository persists to database
    ↓
5. Return HTTP 200/500
```

## Workflow for Any Frontend

Whether the request comes from HTTP, gRPC, CLI, or background job:

1. **Frontend**: Construct input → Build DTO
2. **Validation**: Call `validation.ValidateXxx(dto)`
3. **Service**: Call `service.Xxx(dto_fields...)`
4. **Service** enforces business rules only
5. **Repository**: Persists data

This ensures:
- ✅ Shape validation is centralized and reused
- ✅ Business logic is frontend-agnostic
- ✅ No duplicate validation checks
- ✅ Clear separation of concerns

## Validation Rules by Layer

| Validation Type | Layer | Tool | Example |
|---|---|---|---|
| Required fields | DTO | struct tag | `validate:"required"` |
| Ranges | DTO | struct tag | `validate:"gt=0,max=100"` |
| String length | DTO | struct tag | `validate:"min=1,max=255"` |
| Enums | DTO | struct tag | `validate:"oneof=active pending"` |
| State transitions | Service | Code logic | Can't go from completed to active |
| Resource exists | Service | Code logic | Category must exist before use |
| Cross-resource validation | Service | Code logic | Menu item must be in stock |

## Key Takeaways

1. **DTOs in validation package**: Transport models with struct tags
2. **ValidateXxx functions**: Use struct tags via ValidateStruct()
3. **Handlers are thin**: Parse, validate, call service
4. **Services are smart**: Only business logic, no shape checks
5. **Lazy singleton validator**: Thread-safe, efficient, reusable
6. **Same pattern for all frontends**: HTTP, gRPC, CLI, jobs
