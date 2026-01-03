package optimization

import "time"

// QueryMetrics tracks query performance and optimization opportunities
type QueryMetrics struct {
	QueryName      string
	ExecutionTime  float64 // milliseconds
	RowsReturned   int64
	RowsAffected   int64
	CachedResult   bool
	IndexUsed      string
	QueryType      string // SELECT, INSERT, UPDATE, DELETE
	OptimizationID string
}

// OptimizationPattern represents a query optimization pattern
type OptimizationPattern struct {
	Name        string
	Description string
	Before      string
	After       string
	Impact      string
	Difficulty  string // Easy, Medium, Hard
}

// OptimizedQueries contains pre-optimized query patterns
var OptimizedQueries = map[string]OptimizationPattern{
	"list_with_pagination": {
		Name:        "Pagination Query",
		Description: "Efficient list retrieval with pagination using offset/limit",
		Before:      "SELECT * FROM items WHERE status = $1",
		After:       "SELECT id, name, price FROM items WHERE status = $1 ORDER BY id DESC LIMIT $2 OFFSET $3",
		Impact:      "Reduces memory usage by fetching only needed rows",
		Difficulty:  "Easy",
	},
	"count_optimization": {
		Name:        "Count Query Optimization",
		Description: "Use COUNT(*) separately from pagination queries",
		Before:      "SELECT COUNT(*) OVER() as total, * FROM items WHERE status = $1 LIMIT 20",
		After:       "SELECT COUNT(*) FROM items WHERE status = $1; SELECT id, name FROM items WHERE status = $1 ORDER BY id LIMIT 20 OFFSET 0",
		Impact:      "Separates count logic, allows caching of count separately",
		Difficulty:  "Easy",
	},
	"index_usage": {
		Name:        "Index Usage",
		Description: "Ensure WHERE clause columns have indexes",
		Before:      "SELECT * FROM sessions WHERE table_id = $1",
		After:       "CREATE INDEX idx_sessions_table_id ON sessions(table_id); SELECT id, table_id, status FROM sessions WHERE table_id = $1",
		Impact:      "Reduces query execution time from O(n) to O(log n)",
		Difficulty:  "Medium",
	},
	"join_optimization": {
		Name:        "Join Optimization",
		Description: "Use INNER JOIN for required data, LEFT JOIN for optional",
		Before:      "SELECT * FROM orders o LEFT JOIN menu_items m ON o.item_id = m.id WHERE o.session_id = $1",
		After:       "SELECT o.id, o.session_id, m.name FROM orders o INNER JOIN menu_items m ON o.item_id = m.id WHERE o.session_id = $1",
		Impact:      "Clearer intent, better query planner optimization",
		Difficulty:  "Medium",
	},
	"avoid_select_star": {
		Name:        "Avoid SELECT *",
		Description: "Explicitly select only needed columns",
		Before:      "SELECT * FROM menu_items",
		After:       "SELECT id, name, price, category_id FROM menu_items WHERE is_available = true",
		Impact:      "Reduces data transfer and memory usage",
		Difficulty:  "Easy",
	},
	"use_exist_not_in": {
		Name:        "EXISTS instead of IN",
		Description: "Use EXISTS for subqueries instead of IN",
		Before:      "SELECT * FROM sessions WHERE id IN (SELECT session_id FROM orders WHERE status = 'completed')",
		After:       "SELECT * FROM sessions s WHERE EXISTS (SELECT 1 FROM orders o WHERE o.session_id = s.id AND o.status = 'completed')",
		Impact:      "Better performance for large datasets",
		Difficulty:  "Medium",
	},
	"batch_operations": {
		Name:        "Batch Operations",
		Description: "Insert/update multiple rows in single query",
		Before:      "INSERT INTO orders (session_id, item_id) VALUES ($1, $2); INSERT INTO orders (session_id, item_id) VALUES ($3, $4);",
		After:       "INSERT INTO orders (session_id, item_id) VALUES ($1, $2), ($3, $4) RETURNING id",
		Impact:      "Reduces roundtrips and improves throughput",
		Difficulty:  "Medium",
	},
	"connection_pooling": {
		Name:        "Connection Pooling",
		Description: "Reuse database connections instead of creating new ones",
		Before:      "db.SetMaxOpenConns(0) // unlimited",
		After:       "db.SetMaxOpenConns(25); db.SetMaxIdleConns(5)",
		Impact:      "Reduces connection overhead by 80-90%",
		Difficulty:  "Easy",
	},
}

// QueryOptimizationGuide provides comprehensive optimization guidelines
const QueryOptimizationGuide = `
# Query Optimization Implementation Guide

## 1. Indexing Strategy
- Create indexes on columns used in WHERE clauses
- Create indexes on foreign key columns
- Create composite indexes for common WHERE+ORDER BY patterns
- Monitor slow queries and create indexes on them

Example:
CREATE INDEX idx_sessions_table_id ON sessions(table_id);
CREATE INDEX idx_orders_session_id ON orders(session_id);
CREATE INDEX idx_orders_status ON orders(status);

## 2. Query Patterns
- Use SELECT with specific columns, never SELECT *
- Use LIMIT/OFFSET for pagination, not all data
- Use COUNT(*) separately for total count
- Use INNER JOIN for required relationships
- Use LEFT JOIN for optional data with NULL handling

## 3. Connection Pooling
- MaxOpenConns: 25 (default), 100 (high traffic), 10 (low traffic)
- MaxIdleConns: 5 (default), 20 (high traffic), 2 (low traffic)
- ConnMaxLifetime: 30 minutes to prevent stale connections
- ConnMaxIdleTime: 5 minutes to free up idle connections

## 4. Caching Strategies
- Cache list queries with pagination: 5 minute TTL
- Cache individual lookups: 15 minute TTL
- Invalidate cache on write operations (INSERT/UPDATE/DELETE)
- Use cache for menu items and categories (rarely change)

## 5. Transaction Best Practices
- Use transactions for multi-step operations
- Minimize transaction scope (keep them short)
- Use appropriate isolation levels
- Rollback on any error

## 6. Prepared Statements
- Always use parameterized queries ($1, $2, etc.)
- Prevents SQL injection attacks
- Improves query plan caching

## 7. Monitoring and Profiling
- Track query execution time
- Monitor connection pool stats
- Log slow queries (> 100ms)
- Analyze query plans with EXPLAIN

## 8. Common Pitfalls to Avoid
- N+1 query problem: Use JOIN instead of separate queries
- Missing indexes: Monitor slow query log
- Unbounded queries: Always use LIMIT
- Connection exhaustion: Set proper pool limits
- Large transactions: Break into smaller operations
`

// GetOptimizationTip returns optimization tips based on query characteristics
func GetOptimizationTip(queryType string, rowCount int64) string {
	switch queryType {
	case "list":
		if rowCount > 1000 {
			return "Consider pagination or filtering to reduce result set"
		}
		return "Ensure pagination is implemented (LIMIT/OFFSET)"
	case "count":
		return "Run COUNT separately from main query for better performance"
	case "join":
		return "Verify all JOIN columns have indexes"
	case "aggregate":
		return "Use aggregate functions in database, not application code"
	default:
		return "Review query plan with EXPLAIN to identify bottlenecks"
	}
}

// PerformanceThreshold defines performance monitoring thresholds
type PerformanceThreshold struct {
	SlowQueryThreshold   time.Duration // Queries slower than this are logged
	CacheHitTarget       float64       // Target cache hit percentage
	PoolUtilizationLimit float64       // Alert if pool usage exceeds this
}

// DefaultThresholds provides default performance thresholds
func DefaultThresholds() PerformanceThreshold {
	return PerformanceThreshold{
		SlowQueryThreshold:   100 * time.Millisecond,
		CacheHitTarget:       0.70, // 70%
		PoolUtilizationLimit: 0.80, // 80%
	}
}
