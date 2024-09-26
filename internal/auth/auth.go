package auth

import (
	"database/sql"
	"net/http"
	"time"

	"go_rest_api/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("dsjfksdafkjsadkjfksdafkjsad")

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginReq LoginRequest

		if err := c.Bind(&loginReq); err!= nil {
            return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
        }
		
		username := loginReq.Username
		password := loginReq.Password

		user, err := GetUserByUsername(db, username)
		if err != nil {
			return echo.ErrUnauthorized
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return echo.ErrUnauthorized
		}

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = user.ID
		claims["username"] = user.Username
		claims["role"] = user.Role
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString(jwtSecret)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
}


func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow("SELECT id, username, password, role FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}

	return user, nil
}
