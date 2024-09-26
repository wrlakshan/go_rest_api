package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"go_rest_api/internal/models"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	db *sql.DB
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	rows, err := h.db.Query("SELECT id, name, price, description FROM products")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch products")
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Description); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to scan product")
		}
		products = append(products, product)
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	product := new(models.Product)
	if err := c.Bind(product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	result, err := h.db.Exec("INSERT INTO products (name, price, description ) VALUES (?, ?, ?)",
		product.Name, product.Price, product.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create product")
	}

	id, _ := result.LastInsertId()
	product.ID = int(id)

	return c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	product := new(models.Product)

	err := h.db.QueryRow("SELECT id, name, price, description FROM products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.Price, &product.Description)
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch product")
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	product := new(models.Product)
	if err := c.Bind(product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	_, err := h.db.Exec("UPDATE products SET name = ?, price = ?, description = ? WHERE id = ?", product.Name, product.Price,product.Description, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update product")
	}

	product.ID = id
	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	result, err := h.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete product")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	return c.NoContent(http.StatusNoContent)
}