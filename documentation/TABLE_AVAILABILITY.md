# Table Availability Management

## Overview

The restaurant management system implements **dynamic table availability tracking** - availability is calculated in real-time from session data rather than stored as a separate field. This ensures perfect accuracy and eliminates sync issues.

## How Availability Works

### Dynamic Calculation
Table availability is determined by checking if there are any active or pending sessions for that table:

```sql
SELECT NOT EXISTS (
  SELECT 1 FROM sessions
  WHERE table_id = ?
    AND status IN ('active', 'pending')
)
```

- **Available**: `true` when no active or pending sessions exist
- **Occupied**: `false` when at least one active or pending session exists

## Business Rules

1. **One Active Session Per Table**: A table cannot have more than one active or pending session simultaneously
2. **Real-time Accuracy**: Availability is always calculated from current session data
3. **No Data Duplication**: No separate `available` field to keep in sync
4. **Automatic Updates**: Changes automatically when session status changes

## API Behavior

### Creating Sessions
```bash
# This will succeed if table 1 has no active/pending sessions
POST /sessions
{
  "table_id": 1
}

# This will fail with 409 if table 1 already has an active/pending session
POST /sessions
{
  "table_id": 1
}
```

### Session Status Changes
```bash
# When session becomes active/pending → table becomes occupied
PUT /sessions/{session-id}
{
  "status": "active"
}

# When session becomes cancelled/completed → table becomes available
PUT /sessions/{session-id}
{
  "status": "completed"
}
```

### Table Responses
All table endpoints return basic table information:

```bash
GET /sessions/tables
# Returns: [{"id": 1}, {"id": 2}, {"id": 3}]

GET /sessions/tables/1
# Returns: {"id": 1}
```

## State Machine

```
Table Available (no active/pending sessions)
    ↓ Session created or status → active/pending
Table Occupied (has active/pending session)
    ↓ Session status → cancelled/completed OR session deleted
Table Available (no active/pending sessions)
```

## Error Handling

- **409 Conflict**: Returned when attempting to create a session on a table with active/pending sessions
- **Message**: "table is not available (has active or pending session)"

## Database Schema

**No changes needed** - availability is calculated dynamically from existing session data.

## Benefits of Dynamic Approach

✅ **Always Accurate**: Calculated from real session data
✅ **No Sync Issues**: No separate field to maintain
✅ **No Extra Updates**: No database writes for availability changes
✅ **Simpler Schema**: No additional columns
✅ **Better Performance**: Single query vs multiple updates
✅ **Data Integrity**: Impossible to have inconsistent states

## Implementation

### Repository Method
```go
func (r *postgresRepository) IsTableAvailable(ctx context.Context, tableID int) (bool, error) {
    query := `SELECT NOT EXISTS (
        SELECT 1 FROM sessions
        WHERE table_id = $1
            AND status IN ('active', 'pending')
    )`
    var available bool
    err := r.db.QueryRowContext(ctx, query, tableID).Scan(&available)
    return available, nil
}
```

### Service Logic
- **CreateSession**: Check `IsTableAvailable()` before creating
- **UpdateSession**: No additional logic needed - availability auto-updates
- **DeleteSession**: No additional logic needed - availability auto-updates

## Testing Scenarios

### Scenario 1: Normal Session Flow
1. Create table 1 → Available (no sessions)
2. Create session on table 1 → Occupied (has active session)
3. Complete session → Available (no active sessions)

### Scenario 2: Multiple Session Attempt
1. Create session on table 1 → Occupied
2. Try creating another session on table 1 → `409 Conflict`
3. Complete first session → Available
4. Now can create new session on table 1

### Scenario 3: Session Deletion
1. Create session on table 1 → Occupied
2. Delete session → Available

## Migration Notes

The previous migration that added the `available` column has been removed. If you have existing data with the column, you can safely drop it:

```sql
ALTER TABLE tables DROP COLUMN available;
```

The system will work correctly with or without the column since availability is now calculated dynamically.