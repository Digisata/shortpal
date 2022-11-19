package main

import (
	"net/http"
	"net/url"

	"github.com/Digisata/shortpal/helper"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) ShortenUrl(c *gin.Context) {
	var input ShortenUrlInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Bad request", "error", http.StatusBadRequest, errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = url.ParseRequestURI(input.OriginUrl)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Invalid URL format", "error", http.StatusNotAcceptable, errorMessage)
		c.JSON(http.StatusNotAcceptable, response)
		return
	}

	newUrl, err := h.service.ShortenUrl(input)
	if err != nil {
		response := helper.APIResponse("Failed to shorten url", "error", http.StatusInternalServerError, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	formatter := FormatUrl(newUrl)

	response := helper.APIResponse("Success to shorten url", "success", http.StatusOK, formatter)

	c.JSON(http.StatusOK, response)
}

func (h *handler) RedirectToOriginUrl(c *gin.Context) {
	var input GetUrlDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail url", "error", http.StatusBadRequest, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	url, err := h.service.GetUrlByShortUrl(input.ShortUrl)
	if err != nil {
		response := helper.APIResponse("Failed to get detail url", "error", http.StatusInternalServerError, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.Redirect(http.StatusFound, url.OriginUrl)
}
