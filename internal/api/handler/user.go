package handler

import (
	"RIP_lab1/internal/models"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SignUp godoc
// @Summary      Sign up a new user
// @Description Register a new user account. The user will receive a JWT token upon successful registration.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.UserSignUp true "User sign up data"
// @Success 201 {object} map[string]interface{} "User successfully registered"
// @Failure 400 {object} map[string]interface{} "Bad request error"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /signup [post]
func (h *Handler) SignUp(c *gin.Context) {
	var newClient models.UserSignUp
	var err error

	if err = c.BindJSON(&newClient); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных о новом пользователе"})
		return
	}

	if newClient.Password, err = h.hasher.Hash(newClient.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Неверный формат пароля"})
		return
	}

	userId, isAdmin, err := h.repo.SignUp(c.Request.Context(), models.User{
		Login:    newClient.Login,
		Password: newClient.Password,
		Email:    newClient.Email,
	})

	h.logger.Println("Created user: ", userId, isAdmin, newClient.Login)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Нельзя создать пользователя с таким логином"})

		return
	}

	token, err := h.tokenManager.NewJWT(userId, isAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка при формировании токена"})
		return
	}

	c.SetCookie("AccessToken", "Bearer "+token, 0, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно создан"})
}

// CheckAuth godoc
// @Summary      Check user authentication
// @Description  Retrieves user information based on the provided user context
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      500  {object}  string
// @Router       /check-auth [get]
func (h *Handler) CheckAuth(c *gin.Context) {
	var userInfo = models.User{UserId: c.GetInt(userCtx)}

	userInfo, err := h.repo.GetUserInfo(c.Request.Context(), userInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Неверный формат данных")
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

// SignIn godoc
// @Summary      User sign-in
// @Description  Authenticates a user and generates an access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body  models.UserLogin  true  "User information"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Failure      401  {object}  error
// @Failure      500  {object}  error
// @Router       /sign_in [post]
func (h *Handler) SignIn(c *gin.Context) {
	var clientInfo models.UserLogin
	var err error

	if err = c.BindJSON(&clientInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Неверный формат данных")
		return
	}

	if clientInfo.Password, err = h.hasher.Hash(clientInfo.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "неверный формат пароля"})
		return
	}

	user, err := h.repo.GetByCredentials(c.Request.Context(), models.User{Password: clientInfo.Password, Login: clientInfo.Login})
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка авторизации"})
		return
	}

	token, err := h.tokenManager.NewJWT(user.UserId, user.IsAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка при формировании токена"})
		return
	}

	c.SetCookie("AccessToken", "Bearer "+token, 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "клиент успешно авторизован", "is_admin": user.IsAdmin, "login": user.Login, "userId": user.UserId})
}

// Logout godoc
// @Summary      Logout
// @Description  Logs out the user by blacklisting the access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Router       /logout [post]
func (h *Handler) Logout(c *gin.Context) {
	jwtStr, err := c.Cookie("AccessToken")
	if !strings.HasPrefix(jwtStr, jwtPrefix) || err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	jwtStr = jwtStr[len(jwtPrefix):]

	_, _, err = h.tokenManager.Parse(jwtStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	err = h.redis.WriteJWTToBlacklist(c.Request.Context(), jwtStr, time.Hour)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// ChangeProfile godoc
// @Summary Update Profile
// @Description Update the profile of the currently logged-in user. Allows changing login, password, and email.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.UserSignUp true "Updated user profile data"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request error"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /profile [put]
func (h *Handler) ChangeProfile(c *gin.Context) {
	var changedClient models.UserSignUp
	var err error

	if err = c.BindJSON(&changedClient); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Неверный формат обновлённых данных пользователе"})
		return
	}

	if changedClient.Password != "" {
		if changedClient.Password, err = h.hasher.Hash(changedClient.Password); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Неверный формат пароля"})
			return
		}
	}

	err = h.repo.ChangeProfile(c.Request.Context(), models.User{
		UserId:   c.GetInt(userCtx),
		Login:    changedClient.Login,
		Password: changedClient.Password,
		Email:    changedClient.Email,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Нельзя изменить пользователя на такой логин"})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Данные пользователя успешно изменены"})
}
