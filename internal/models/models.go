package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system (customers, vendors, admins)
type User struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Email       string     `json:"email" db:"email" validate:"required,email"`
	Password    string     `json:"-" db:"password_hash"`
	FirstName   string     `json:"first_name" db:"first_name" validate:"required"`
	LastName    string     `json:"last_name" db:"last_name" validate:"required"`
	Phone       *string    `json:"phone,omitempty" db:"phone"`
	Role        UserRole   `json:"role" db:"role"`
	Status      UserStatus `json:"status" db:"status"`
	Avatar      *string    `json:"avatar,omitempty" db:"avatar"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleVendor   UserRole = "vendor"
	RoleAdmin    UserRole = "admin"
)

type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusSuspended UserStatus = "suspended"
)

// Vendor represents a business/SME on the platform
type Vendor struct {
	ID              uuid.UUID    `json:"id" db:"id"`
	UserID          uuid.UUID    `json:"user_id" db:"user_id"`
	BusinessName    string       `json:"business_name" db:"business_name" validate:"required"`
	BusinessType    string       `json:"business_type" db:"business_type"`
	Description     *string      `json:"description,omitempty" db:"description"`
	Logo            *string      `json:"logo,omitempty" db:"logo"`
	Website         *string      `json:"website,omitempty" db:"website"`
	Address         Address      `json:"address" db:"address"`
	Status          VendorStatus `json:"status" db:"status"`
	VerifiedAt      *time.Time   `json:"verified_at,omitempty" db:"verified_at"`
	CreatedAt       time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at" db:"updated_at"`
}

type VendorStatus string

const (
	VendorStatusPending  VendorStatus = "pending"
	VendorStatusApproved VendorStatus = "approved"
	VendorStatusRejected VendorStatus = "rejected"
	VendorStatusSuspended VendorStatus = "suspended"
)

// Product represents a product or service
type Product struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	NumericID   int           `json:"numeric_id" db:"numeric_id"`
	VendorID    uuid.UUID     `json:"vendor_id" db:"vendor_id"`
	Name        string        `json:"name" db:"name" validate:"required"`
	Description string        `json:"description" db:"description"`
	Price       float64       `json:"price" db:"price" validate:"required,gt=0"`
	ComparePrice *float64     `json:"compare_price,omitempty" db:"compare_price"`
	SKU         *string       `json:"sku,omitempty" db:"sku"`
	Category    string        `json:"category" db:"category" validate:"required"`
	Tags        []string      `json:"tags" db:"tags"`
	Images      []string      `json:"images" db:"images"`
	Stock       int           `json:"stock" db:"stock"`
	Status      ProductStatus `json:"status" db:"status"`
	Featured    bool          `json:"featured" db:"featured"`
	Weight      *float64      `json:"weight,omitempty" db:"weight"`
	Dimensions  *Dimensions   `json:"dimensions,omitempty" db:"dimensions"`
	SEO         SEOData       `json:"seo" db:"seo"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
}

type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusActive    ProductStatus = "active"
	ProductStatusInactive  ProductStatus = "inactive"
	ProductStatusArchived  ProductStatus = "archived"
)

type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type SEOData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    []string `json:"keywords"`
}

// Order represents a customer order
type Order struct {
	ID              uuid.UUID   `json:"id" db:"id"`
	CustomerID      uuid.UUID   `json:"customer_id" db:"customer_id"`
	OrderNumber     string      `json:"order_number" db:"order_number"`
	Status          OrderStatus `json:"status" db:"status"`
	PaymentStatus   PaymentStatus `json:"payment_status" db:"payment_status"`
	PaymentMethod   string      `json:"payment_method" db:"payment_method"`
	StripePaymentID *string     `json:"stripe_payment_id,omitempty" db:"stripe_payment_id"`
	Subtotal        float64     `json:"subtotal" db:"subtotal"`
	Tax             float64     `json:"tax" db:"tax"`
	Shipping        float64     `json:"shipping" db:"shipping"`
	Discount        float64     `json:"discount" db:"discount"`
	Total           float64     `json:"total" db:"total"`
	Currency        string      `json:"currency" db:"currency"`
	ShippingAddress Address     `json:"shipping_address" db:"shipping_address"`
	BillingAddress  Address     `json:"billing_address" db:"billing_address"`
	Notes           *string     `json:"notes,omitempty" db:"notes"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at" db:"updated_at"`
	Items           []OrderItem `json:"items,omitempty"`
}

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	OrderID   uuid.UUID `json:"order_id" db:"order_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	VendorID  uuid.UUID `json:"vendor_id" db:"vendor_id"`
	Name      string    `json:"name" db:"name"`
	Price     float64   `json:"price" db:"price"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Total     float64   `json:"total" db:"total"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Address represents a shipping/billing address
type Address struct {
	Street     string  `json:"street" validate:"required"`
	City       string  `json:"city" validate:"required"`
	State      string  `json:"state" validate:"required"`
	PostalCode string  `json:"postal_code" validate:"required"`
	Country    string  `json:"country" validate:"required"`
	Phone      *string `json:"phone,omitempty"`
}

// Category represents product categories
type Category struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required"`
	Slug        string    `json:"slug" db:"slug"`
	Description *string   `json:"description,omitempty" db:"description"`
	Image       *string   `json:"image,omitempty" db:"image"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Cart represents a shopping cart
type Cart struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	CustomerID *uuid.UUID `json:"customer_id,omitempty" db:"customer_id"`
	SessionID  *string    `json:"session_id,omitempty" db:"session_id"`
	Items      []CartItem `json:"items,omitempty"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CartID    uuid.UUID `json:"cart_id" db:"cart_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity" validate:"required,gt=0"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Product   *Product  `json:"product,omitempty"`
}

// Review represents a product review
type Review struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ProductID  uuid.UUID `json:"product_id" db:"product_id"`
	CustomerID uuid.UUID `json:"customer_id" db:"customer_id"`
	Rating     int       `json:"rating" db:"rating" validate:"required,min=1,max=5"`
	Title      *string   `json:"title,omitempty" db:"title"`
	Comment    *string   `json:"comment,omitempty" db:"comment"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Database scanning methods for JSONB fields

// Scan implements sql.Scanner for Dimensions
func (d *Dimensions) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into Dimensions", value)
	}

	return json.Unmarshal(bytes, d)
}

// Value implements driver.Valuer for Dimensions
func (d Dimensions) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements sql.Scanner for SEOData
func (s *SEOData) Scan(value interface{}) error {
	if value == nil {
		*s = SEOData{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into SEOData", value)
	}

	return json.Unmarshal(bytes, s)
}

// Value implements driver.Valuer for SEOData
func (s SEOData) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implements sql.Scanner for Address
func (a *Address) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into Address", value)
	}

	return json.Unmarshal(bytes, a)
}

// Value implements driver.Valuer for Address
func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}