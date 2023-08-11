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
// @tags Comments
// @Summary Create a new comment
// @Description Create a new comment with the given data
// @Accept json
// @Produce json
// @Param Authorization header string true "your_token_here"
// @Param comment body models.Comment true "Comment data"
// @Success 200 {object} models.Comment
// @Router /api/comments [post]
func CreateCommentHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var comment models.Comment
		err := c.Bind(&comment)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		// Get the userID from the context (added by JWT middleware)
		userID := c.Get("userID").(uint)
		comment.UserID = int(userID)

		err = models.InsertComment(db, comment)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return utils.JsonResponse(c, comment)
	}
}

// @tags Comments
// @Summary Get all comments
// @Description Get all comments
// @Accept json
// @Produce json
// @Success 200 {array} models.Comment
// @Router /comments [get]
func GetCommentsHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var comments []models.Comment
		db.Find(&comments)

		if c.Request().Header.Get("Accept") == "application/xml" {
			return utils.XmlResponse(c, comments)
		} else {
			return utils.JsonResponse(c, comments)
		}
	}
}

// @tags Comments
// @Summary Get a comment by ID
// @Description Get a comment by its unique ID
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} models.Comment
// @Router /comments/{id} [get]
func GetCommentHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID parameter")
		}

		var comment models.Comment
		result := db.First(&comment, id)
		if result.Error != nil {
			return c.String(http.StatusNotFound, "Comment not found")
		}

		if c.Request().Header.Get("Accept") == "application/xml" {
			return utils.XmlResponse(c, comment)
		} else {
			return utils.JsonResponse(c, comment)
		}
	}
}

// @tags Comments
// @Summary Update a comment by ID
// @Description Update a comment with the given ID and data
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param comment body models.Comment true "Updated comment data"
// @Success 200 {object} models.Comment
// @Router /comments/{id} [put]
func UpdateCommentHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID parameter")
		}

		var comment models.Comment
		result := db.First(&comment, id)
		if result.Error != nil {
			return c.String(http.StatusNotFound, "Comment not found")
		}

		var updatedComment models.Comment
		err = c.Bind(&updatedComment)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		comment.Name = updatedComment.Name
		comment.Email = updatedComment.Email
		comment.Body = updatedComment.Body
		db.Save(&comment)

		return utils.JsonResponse(c, comment)
	}
}

// @tags Comments
// @Summary Delete a comment by ID
// @Description Delete a comment with the given ID
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Success 204 "No Content"
// @Router /comments/{id} [delete]
func DeleteCommentHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID parameter")
		}

		var comment models.Comment
		result := db.First(&comment, id)
		if result.Error != nil {
			return c.String(http.StatusNotFound, "Comment not found")
		}

		db.Delete(&comment)
		return c.NoContent(http.StatusOK)
	}
}
