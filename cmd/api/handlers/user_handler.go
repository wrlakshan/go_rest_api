package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"go_rest_api/internal/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) SignUp(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	result, err := h.db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", 
		user.Username, hashedPassword, user.Role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create signup user")
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user ID")
	}
	user.ID = int(id)
	// user.Password = ""

	return c.JSON(http.StatusCreated, user)
}


func (h *UserHandler) GetUsers(c echo.Context) error {
	rows, err := h.db.Query("SELECT id, username, role FROM users")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch users")
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Role); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to scan user")
		}
		users = append(users, user)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	result, err := h.db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)",
		user.Username, hashedPassword, user.Role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)
	user.Password = "" 

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user := new(models.User)

	err := h.db.QueryRow("SELECT id, username, role FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username, &user.Role)
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user")
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	_, err := h.db.Exec("UPDATE users SET username = ?, role = ? WHERE id = ?", user.Username, user.Role, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}

	user.ID = id
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	result, err := h.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.NoContent(http.StatusNoContent)
}