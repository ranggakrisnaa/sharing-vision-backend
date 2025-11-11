package article

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
)

type Repository interface {
	Insert(ctx context.Context, title, content, category, status string) (Article, error)
	List(ctx context.Context, limit, offset int, filter ListFilter) ([]Article, int64, error)
	FindByID(ctx context.Context, id int64) (Article, error)
	UpdateAll(ctx context.Context, id int64, title, content, category, status string) (Article, error)
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context, filter ListFilter) (int64, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Insert(ctx context.Context, title, content, category, status string) (Article, error) {
	q := `
    INSERT INTO articles (title, content, category, status)
    VALUES (?, ?, ?, ?)
    `
	res, err := r.db.ExecContext(ctx, q, title, content, category, status)
	if err != nil {
		return Article{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Article{}, err
	}
	return r.FindByID(ctx, id)
}

func (r *MySQLRepository) List(ctx context.Context, limit, offset int, filter ListFilter) ([]Article, int64, error) {
	base := `SELECT id, title, content, category, status, created_at, updated_at FROM articles`
	where, args := buildFilterClause(filter)
	q := base + where + " ORDER BY id LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return []Article{}, 0, err
	}
	defer rows.Close()

	res := make([]Article, 0)
	for rows.Next() {
		var a Article
		if errScan := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.Status, &a.CreatedAt, &a.UpdatedAt); errScan != nil {
			return []Article{}, 0, errScan
		}
		res = append(res, a)
	}
	if errRows := rows.Err(); errRows != nil {
		return []Article{}, 0, errRows
	}
	count, err := r.Count(ctx, filter)
	if err != nil {
		return []Article{}, 0, err
	}
	return res, count, nil
}

func (r *MySQLRepository) FindByID(ctx context.Context, id int64) (Article, error) {
	q := `
    SELECT id, title, content, category, status, created_at, updated_at
    FROM articles
    WHERE id = ?
    `
	var a Article
	err := r.db.QueryRowContext(ctx, q, id).Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Article{}, sql.ErrNoRows
		}
		return Article{}, err
	}
	return a, nil
}

func (r *MySQLRepository) UpdateAll(ctx context.Context, id int64, title, content, category, status string) (Article, error) {
	q := `
    UPDATE articles
    SET title = ?, content = ?, category = ?, status = ?, updated_at = CURRENT_TIMESTAMP
    WHERE id = ?
    `
	res, err := r.db.ExecContext(ctx, q, title, content, category, status, id)
	if err != nil {
		return Article{}, err
	}
	n, err := res.RowsAffected()
	if err == nil && n == 0 {
		return Article{}, sql.ErrNoRows
	}
	return r.FindByID(ctx, id)
}

func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	q := `DELETE FROM articles WHERE id = ?`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err == nil && n == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (r *MySQLRepository) Count(ctx context.Context, filter ListFilter) (int64, error) {
	base := `SELECT COUNT(*) FROM articles`
	where, args := buildFilterClause(filter)
	q := base + where
	var total int64
	err := r.db.QueryRowContext(ctx, q, args...).Scan(&total)
	return total, err
}

// buildFilterClause to build WHERE clause and args based on ListFilter
func buildFilterClause(filter ListFilter) (string, []interface{}) {
	conds := make([]string, 0)
	args := make([]interface{}, 0)

	trimmedTitle := "%" + strings.TrimSpace(filter.Title) + "%"
	if trimmedTitle != "" {
		conds = append(conds, "LOWER(title) LIKE ?")
		args = append(args, trimmedTitle)
	}

	trimmedCategory := filter.Category
	if trimmedCategory != "" {
		conds = append(conds, "category = ?")
		args = append(args, trimmedCategory)
	}

	trimmedStatus := strings.TrimSpace(filter.Status)
	if trimmedStatus != "" {
		conds = append(conds, "status = ?")
		args = append(args, trimmedStatus)
	}

	log.Printf("conds: %v, args: %v", conds, args)

	if len(conds) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}
