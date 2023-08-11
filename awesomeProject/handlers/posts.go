package handlers

import (
	"awesomeProject/models"
	"awesomeProject/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// @Security ApiKeyAuth
// @tags Posts
// @Summary Create a new post
// @Description Create a new post with the given data
// @Accept json
// @Produce json
// @Param Authorization header string true "your_token_here"
// @Param post body models.Post true "Post data"
// @Success 200 {object} models.Post
// @Router /api/posts [post]
func CreatePostHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var post models.Post
		err := c.Bind(&post)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		// Get the userID from the context (added by JWT middleware)
		userID := c.Get("userID").(uint)
		post.UserID = int(userID)

		err = models.InsertPost(db, post)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return utils.JsonResponse(c, post)
	}
}

// @tags Posts
// @Summary Get all posts
// @Description Get all posts
// @Accept json
// @Produce json
// @Success 200 {array} models.Post
// @Router /posts [get]
func GetPostsHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var posts []models.Post
		db.Find(&posts)

		// Determine the response format based on the Accept header
		if c.Request().Header.Get("Accept") == "application/xml" {
			return utils.XmlResponse(c, posts)
		} else {
			return utils.JsonResponse(c, posts)
		}
	}
}

// @tags Posts
// @Summary Get a post by ID
// @Description Get a post by its unique ID
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.Post
// @Router /posts/{id} [get]
func GetPostHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID parameter")
		}

		var post models.Post
		result := db.First(&post, id)
		if result.Error != nil {
			return c.String(http.StatusNotFound, "Post not found")
		}

		// Определяем формат ответа на основе заголовка Accept
		if c.Request().Header.Get("Accept") == "application/xml" {
			return utils.XmlResponse(c, post)
		} else {
			return utils.JsonResponse(c, post)
		}
	}
}

// @tags Posts
// @Summary Update a post by ID
// @Description Update a post with the given ID and data
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body models.Post true "Updated post data"
// @Success 200 {object} models.Post
// @Router /posts/{id} [put]
func UpdatePostHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID parameter")
		}

		var post models.Post
		result := db.First(&post, id)
		if result.Error != nil {
			return c.String(http.StatusNotFound, "Post not found")
		}

		var updatedPost models.Post
		err = c.Bind(&updatedPost)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		post.Title = updatedPost.Title
		post.Body = updatedPost.Body
		db.Save(&post)

		return utils.JsonResponse(c, post)
	}
}

// @tags Posts
// @Summary Delete a post by ID
// @Description Delete a post with the given ID
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 204 "No Content"
// @Router /posts/{id} [delete]
func DeletePostHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID parameter")
		}

		var post models.Post
		result := db.First(&post, id)
		if result.Error != nil {
			return c.String(http.StatusNotFound, "Post not found")
		}

		db.Delete(&post)
		return c.NoContent(http.StatusOK)
	}
}
