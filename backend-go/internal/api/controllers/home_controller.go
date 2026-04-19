package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type HomeController struct{}

func (h *HomeController) Index(c *gin.Context) {
	utils.Success(c, nil, "Welcome to the Home route!", 200)
}
