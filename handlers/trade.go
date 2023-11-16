package handlers

import (
	"net/http"
	// "qexchange/models"
	"qexchange/models"
	"qexchange/models/trade"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

func OpenTrade(service services.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(trade.OpenTradeRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(http.StatusBadRequest, response)
		}

		user,_ := c.Get("User").(models.User)
		statusCode, err := service.OpenTrade(*request, user)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, models.NewRespone("the trade successfully opened"))
	}
}

func CloseTrade(service services.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(trade.ClosedTradeRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(http.StatusBadRequest, response)
		}

		user,_ := c.Get("User").(models.User)
		statusCode, err := service.CloseTrade(*request, user)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, models.NewRespone("the closing request successfully processed"))
	}
}

func GetAllOpenTrades(service services.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(trade.OpenTradeRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(http.StatusBadRequest, response)
		}

		user,_ := c.Get("User").(models.User)
		allOpenTrades,statusCode, err := service.GetAllOpenTrades(user)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, allOpenTrades)
	}
}

func GetAllClosedTrades(service services.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(trade.ClosedTradeRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(http.StatusBadRequest, response)
		}

		user,_ := c.Get("User").(models.User)
		allClosedTrades,statusCode, err := service.GetAllClosedTrades(user)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, allClosedTrades)
	}
}