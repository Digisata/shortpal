package main

import (
	"net/http"
	"net/url"

	"github.com/Digisata/shortpal/helper"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *Service
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}

func (h handler) ShortenUrl(c *gin.Context) {
	var input ShortenUrlInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Bad request", "error", http.StatusBadRequest, errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = url.ParseRequestURI(input.Url)
	if err != nil {
		response := helper.APIResponse("Invalid URL format", "error", http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newUrl, err := h.service.ShortenUrl(input, c.Request.Context())
	if err != nil {
		response := helper.APIResponse("Failed to shorten url", "error", http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	formatter := FormatUrl(newUrl)

	response := helper.APIResponse("Success to shorten url", "success", http.StatusOK, formatter)

	c.JSON(http.StatusOK, response)
}

func (h handler) Redirect(c *gin.Context) {
	var input GetUrlDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Bad request", "error", http.StatusBadRequest, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	url, err := h.service.Redirect(input.ShortUrl, c.Request.Context())
	if err != nil {
		response := helper.APIResponse("Failed to redirect", "error", http.StatusNotFound, err.Error())
		c.JSON(http.StatusNotFound, response)
		return
	}

	c.Redirect(http.StatusFound, url.Url)
}
