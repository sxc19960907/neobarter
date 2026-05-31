# Database Guidelines

> Database patterns and conventions for NeoBarter (PostgreSQL + GORM).

---

## ORM

- **Library**: GORM v1.25+ (`gorm.io/gorm`)
- **Driver**: `gorm.io/driver/postgres`
- **Migrations**: `db.AutoMigrate()` in `cmd/migrate/main.go` (development); raw SQL in `docs/schema.sql` (production reference)

## Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Table | snake_case, plural | `users`, `trade_requests`, `wallet_transactions` |
| Column | snake_case | `user_id`, `created_at`, `credit_score` |
| Index | `idx_<table>_<column(s)>` | `idx_items_user`, `idx_trade_status` |
| Constraint | `chk_<table>_<description>` | `chk_balance_non_negative` |
| Foreign key | `<referenced_table_singular>_id` | `user_id`, `wallet_id` |

## Model Conventions

Every model must define:

```go
type Item struct {
    ID        int64     `json:"id" gorm:"primaryKey"`
    // ... fields
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (Item) TableName() string {
    return "items"
}
```

- Always explicit `TableName()` — never rely on GORM's auto-pluralization
- Use `BIGSERIAL` (int64) for primary keys
- Timestamps: `created_at` on all tables, `updated_at` on mutable tables
- Soft delete via status field (`status = 'deleted'`), not GORM's `DeletedAt`

## Query Patterns

### Repository constructor

```go
type ItemRepository struct {
    db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
    return &ItemRepository{db: db}
}
```

### Pagination

```go
func (r *Repo) List(page, pageSize int) ([]Model, int64, error) {
    var items []Model
    var total int64
    query := r.db.Model(&Model{})
    query.Count(&total)
    err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
    return items, total, err
}
```

### Row locking for financial operations

```go
func (r *WalletRepository) GetByUserIDForUpdate(tx *gorm.DB, userID int64) (*model.Wallet, error) {
    var wallet model.Wallet
    err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        Where("user_id = ?", userID).First(&wallet).Error
    return &wallet, err
}
```

### Transactions

Transactions are managed at the service layer using `db.Transaction()`:

```go
db.Transaction(func(tx *gorm.DB) error {
    // pass tx to repository methods that need atomicity
    return nil
})
```

## Decimal Handling

- Use `github.com/shopspring/decimal` for all monetary values (巴特币)
- GORM column type: `type:decimal(12,2)`
- Never use `float64` for money

## Forbidden Patterns

- ❌ `float64` for monetary amounts
- ❌ Transactions in repository layer (service owns transaction boundaries)
- ❌ `db.Exec()` with string interpolation (use parameterized queries)
- ❌ Missing `WHERE` clause on `UPDATE` / `DELETE`
- ❌ N+1 queries — use `Preload()` for associations
