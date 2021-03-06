package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sula0/ibanator/v2/cmd/ibanator/docs"
	"github.com/sula0/ibanator/v2/pkg/iban"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

type IBANValidationResponse struct {
	Valid         bool   `json:"valid:"`
	InvalidReason string `json:"invalidReason,omitempty"`
}

type IBANValidationRequest struct {
	IBAN string `json:"iban" binding:"required"`
}

// @BasePath /
// @Summary Validates IBAN
// @Tags IBAN validation
// @Description Runs validation on IBAN. The response gives the validation result, and, if the validation failed, the reason for the failure.
// @Accept json
// @Produce json
// @Param IBANValidationRequest body IBANValidationRequest{iban=string} true "{ "iban": "AL35202111090000000001234567" }"
// @Success 200 {object} Response{data=IBANValidationResponse} "Response wrapper"
// @Failure 400  {object} Response{error=string} "Response wrapper containing only an error string"
// @Router /iban/validate [post]
func validateIBAN(c *gin.Context) {
	var request IBANValidationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// Not sure if error is the correct value to use
	// as a "bad" response.
	valid, err := iban.ValidateIBAN(request.IBAN)
	validationResponse := IBANValidationResponse{Valid: valid}
	if err != nil {
		validationResponse.InvalidReason = err.Error()
	}

	c.JSON(http.StatusOK, Response{Data: validationResponse})
}

// @title IBANator
// @description Service for validating IBAN.
// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"

	r.POST("/iban/validate", validateIBAN)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := r.Run()
	if err != nil {
		panic(err)
	}
}
