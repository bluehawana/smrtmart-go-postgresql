package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"smrtmart-go-postgresql/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id uuid.UUID) (*models.Product, error)
	GetByNumericID(numericID int) (*models.Product, error)
	GetAll(filters ProductFilters) ([]*models.Product, int, error)
	Update(product *models.Product) error
	Delete(id uuid.UUID) error
	GetByVendor(vendorID uuid.UUID, filters ProductFilters) ([]*models.Product, int, error)
	Search(query string, filters ProductFilters) ([]*models.Product, int, error)
	GetFeatured(limit int) ([]*models.Product, error)
	UpdateStock(id uuid.UUID, stock int) error
}

type ProductFilters struct {
	Category string
	Status   string
	Featured *bool
	MinPrice *float64
	MaxPrice *float64
	Page     int
	Limit    int
	SortBy   string
	SortDir  string
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	query := `
		INSERT INTO products (id, vendor_id, name, description, price, compare_price, sku, 
			category, tags, images, stock, status, featured, weight, dimensions, seo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING numeric_id, created_at, updated_at`

	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}

	err := r.db.QueryRow(query,
		product.ID, product.VendorID, product.Name, product.Description,
		product.Price, product.ComparePrice, product.SKU, product.Category,
		pq.Array(product.Tags), pq.Array(product.Images), product.Stock,
		product.Status, product.Featured, product.Weight, product.Dimensions,
		product.SEO,
	).Scan(&product.NumericID, &product.CreatedAt, &product.UpdatedAt)

	return err
}

func (r *productRepository) GetByID(id uuid.UUID) (*models.Product, error) {
	query := `
		SELECT id, numeric_id, vendor_id, name, description, price, compare_price, sku,
			category, tags, images, stock, status, featured, weight, dimensions,
			seo, created_at, updated_at
		FROM products WHERE id = $1`

	product := &models.Product{}
	var dimensionsJSON, seoJSON []byte
	
	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.NumericID, &product.VendorID, &product.Name, &product.Description,
		&product.Price, &product.ComparePrice, &product.SKU, &product.Category,
		pq.Array(&product.Tags), pq.Array(&product.Images), &product.Stock,
		&product.Status, &product.Featured, &product.Weight, &dimensionsJSON,
		&seoJSON, &product.CreatedAt, &product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	if dimensionsJSON != nil {
		json.Unmarshal(dimensionsJSON, &product.Dimensions)
	}
	if seoJSON != nil {
		json.Unmarshal(seoJSON, &product.SEO)
	}

	return product, nil
}

func (r *productRepository) GetByNumericID(numericID int) (*models.Product, error) {
	query := `
		SELECT id, numeric_id, vendor_id, name, description, price, compare_price, sku,
			category, tags, images, stock, status, featured, weight, dimensions,
			seo, created_at, updated_at
		FROM products WHERE numeric_id = $1`

	product := &models.Product{}
	var dimensionsJSON, seoJSON []byte
	
	err := r.db.QueryRow(query, numericID).Scan(
		&product.ID, &product.NumericID, &product.VendorID, &product.Name, &product.Description,
		&product.Price, &product.ComparePrice, &product.SKU, &product.Category,
		pq.Array(&product.Tags), pq.Array(&product.Images), &product.Stock,
		&product.Status, &product.Featured, &product.Weight, &dimensionsJSON,
		&seoJSON, &product.CreatedAt, &product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	if dimensionsJSON != nil {
		json.Unmarshal(dimensionsJSON, &product.Dimensions)
	}
	if seoJSON != nil {
		json.Unmarshal(seoJSON, &product.SEO)
	}

	return product, nil
}

func (r *productRepository) GetAll(filters ProductFilters) ([]*models.Product, int, error) {
	whereClause, args := r.buildWhereClause(filters)
	
	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Main query
	orderClause := r.buildOrderClause(filters)
	limitClause := r.buildLimitClause(filters)
	
	query := fmt.Sprintf(`
		SELECT id, numeric_id, vendor_id, name, description, price, compare_price, sku,
			category, tags, images, stock, status, featured, weight, dimensions,
			seo, created_at, updated_at
		FROM products %s %s %s`, whereClause, orderClause, limitClause)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		var dimensionsJSON, seoJSON []byte
		
		err := rows.Scan(
			&product.ID, &product.NumericID, &product.VendorID, &product.Name, &product.Description,
			&product.Price, &product.ComparePrice, &product.SKU, &product.Category,
			pq.Array(&product.Tags), pq.Array(&product.Images), &product.Stock,
			&product.Status, &product.Featured, &product.Weight, &dimensionsJSON,
			&seoJSON, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Parse JSON fields
		if dimensionsJSON != nil {
			json.Unmarshal(dimensionsJSON, &product.Dimensions)
		}
		if seoJSON != nil {
			json.Unmarshal(seoJSON, &product.SEO)
		}

		products = append(products, product)
	}

	return products, total, nil
}

func (r *productRepository) Update(product *models.Product) error {
	query := `
		UPDATE products SET
			name = $2, description = $3, price = $4, compare_price = $5,
			sku = $6, category = $7, tags = $8, images = $9, stock = $10,
			status = $11, featured = $12, weight = $13, dimensions = $14, seo = $15
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(query,
		product.ID, product.Name, product.Description, product.Price,
		product.ComparePrice, product.SKU, product.Category,
		pq.Array(product.Tags), pq.Array(product.Images), product.Stock,
		product.Status, product.Featured, product.Weight, product.Dimensions,
		product.SEO,
	).Scan(&product.UpdatedAt)

	return err
}

func (r *productRepository) Delete(id uuid.UUID) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *productRepository) GetByVendor(vendorID uuid.UUID, filters ProductFilters) ([]*models.Product, int, error) {
	filters.Category = "" // Reset category filter for vendor-specific queries
	whereClause, args := r.buildWhereClause(filters)
	
	// Add vendor filter
	if whereClause == "" {
		whereClause = "WHERE vendor_id = $1"
		args = []interface{}{vendorID}
	} else {
		whereClause += " AND vendor_id = $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, vendorID)
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Main query
	orderClause := r.buildOrderClause(filters)
	limitClause := r.buildLimitClause(filters)
	
	query := fmt.Sprintf(`
		SELECT id, numeric_id, vendor_id, name, description, price, compare_price, sku,
			category, tags, images, stock, status, featured, weight, dimensions,
			seo, created_at, updated_at
		FROM products %s %s %s`, whereClause, orderClause, limitClause)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		var dimensionsJSON, seoJSON []byte
		
		err := rows.Scan(
			&product.ID, &product.NumericID, &product.VendorID, &product.Name, &product.Description,
			&product.Price, &product.ComparePrice, &product.SKU, &product.Category,
			pq.Array(&product.Tags), pq.Array(&product.Images), &product.Stock,
			&product.Status, &product.Featured, &product.Weight, &dimensionsJSON,
			&seoJSON, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Parse JSON fields
		if dimensionsJSON != nil {
			json.Unmarshal(dimensionsJSON, &product.Dimensions)
		}
		if seoJSON != nil {
			json.Unmarshal(seoJSON, &product.SEO)
		}

		products = append(products, product)
	}

	return products, total, nil
}

func (r *productRepository) Search(query string, filters ProductFilters) ([]*models.Product, int, error) {
	whereClause, args := r.buildWhereClause(filters)
	
	// Add search condition
	searchCondition := `(
		to_tsvector('english', name) @@ plainto_tsquery('english', $%d) OR
		to_tsvector('english', description) @@ plainto_tsquery('english', $%d) OR
		name ILIKE $%d OR
		description ILIKE $%d
	)`
	
	searchArg := len(args) + 1
	searchCondition = fmt.Sprintf(searchCondition, searchArg, searchArg, searchArg+1, searchArg+1)
	args = append(args, query, "%"+query+"%")
	
	if whereClause == "" {
		whereClause = "WHERE " + searchCondition
	} else {
		whereClause += " AND " + searchCondition
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Main query with relevance scoring
	orderClause := r.buildOrderClause(filters)
	if orderClause == "" {
		orderClause = `ORDER BY (
			ts_rank(to_tsvector('english', name), plainto_tsquery('english', $` + fmt.Sprintf("%d", searchArg) + `)) +
			ts_rank(to_tsvector('english', description), plainto_tsquery('english', $` + fmt.Sprintf("%d", searchArg) + `))
		) DESC`
	}
	
	limitClause := r.buildLimitClause(filters)
	
	query = fmt.Sprintf(`
		SELECT id, vendor_id, name, description, price, compare_price, sku,
			category, tags, images, stock, status, featured, weight, dimensions,
			seo, created_at, updated_at
		FROM products %s %s %s`, whereClause, orderClause, limitClause)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		var dimensionsJSON, seoJSON []byte
		
		err := rows.Scan(
			&product.ID, &product.NumericID, &product.VendorID, &product.Name, &product.Description,
			&product.Price, &product.ComparePrice, &product.SKU, &product.Category,
			pq.Array(&product.Tags), pq.Array(&product.Images), &product.Stock,
			&product.Status, &product.Featured, &product.Weight, &dimensionsJSON,
			&seoJSON, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Parse JSON fields
		if dimensionsJSON != nil {
			json.Unmarshal(dimensionsJSON, &product.Dimensions)
		}
		if seoJSON != nil {
			json.Unmarshal(seoJSON, &product.SEO)
		}

		products = append(products, product)
	}

	return products, total, nil
}

func (r *productRepository) GetFeatured(limit int) ([]*models.Product, error) {
	query := `
		SELECT id, numeric_id, vendor_id, name, description, price, compare_price, sku,
			category, tags, images, stock, status, featured, weight, dimensions,
			seo, created_at, updated_at
		FROM products 
		WHERE featured = true AND status = 'active' AND stock > 0
		ORDER BY created_at DESC
		LIMIT $1`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		var dimensionsJSON, seoJSON []byte
		
		err := rows.Scan(
			&product.ID, &product.NumericID, &product.VendorID, &product.Name, &product.Description,
			&product.Price, &product.ComparePrice, &product.SKU, &product.Category,
			pq.Array(&product.Tags), pq.Array(&product.Images), &product.Stock,
			&product.Status, &product.Featured, &product.Weight, &dimensionsJSON,
			&seoJSON, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON fields
		if dimensionsJSON != nil {
			json.Unmarshal(dimensionsJSON, &product.Dimensions)
		}
		if seoJSON != nil {
			json.Unmarshal(seoJSON, &product.SEO)
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *productRepository) UpdateStock(id uuid.UUID, stock int) error {
	query := "UPDATE products SET stock = $2 WHERE id = $1"
	result, err := r.db.Exec(query, id, stock)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *productRepository) buildWhereClause(filters ProductFilters) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, filters.Category)
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if filters.Featured != nil {
		conditions = append(conditions, fmt.Sprintf("featured = $%d", argIndex))
		args = append(args, *filters.Featured)
		argIndex++
	}

	if filters.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argIndex))
		args = append(args, *filters.MinPrice)
		argIndex++
	}

	if filters.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argIndex))
		args = append(args, *filters.MaxPrice)
		argIndex++
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

func (r *productRepository) buildOrderClause(filters ProductFilters) string {
	if filters.SortBy == "" {
		return "ORDER BY created_at DESC"
	}

	sortDir := "ASC"
	if filters.SortDir == "desc" {
		sortDir = "DESC"
	}

	switch filters.SortBy {
	case "name":
		return fmt.Sprintf("ORDER BY name %s", sortDir)
	case "price":
		return fmt.Sprintf("ORDER BY price %s", sortDir)
	case "created_at":
		return fmt.Sprintf("ORDER BY created_at %s", sortDir)
	case "updated_at":
		return fmt.Sprintf("ORDER BY updated_at %s", sortDir)
	default:
		return "ORDER BY created_at DESC"
	}
}

func (r *productRepository) buildLimitClause(filters ProductFilters) string {
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}

	offset := (filters.Page - 1) * filters.Limit
	return fmt.Sprintf("LIMIT %d OFFSET %d", filters.Limit, offset)
}