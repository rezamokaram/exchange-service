package handlers

import (
	"errors"
	"net/http"
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

		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorRespone("bad user data", nil)
			return c.JSON(http.StatusBadRequest, response)
		}

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

		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorRespone("bad user data", nil)
			return c.JSON(http.StatusBadRequest, response)
		}

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

		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorRespone("bad user data", nil)
			return c.JSON(http.StatusBadRequest, response)
		}

		allOpenTrades, statusCode, err := service.GetAllOpenTrades(user)
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

		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorRespone("", errors.New("bad user data"))
			return c.JSON(http.StatusBadRequest, response)
		}

		allClosedTrades, statusCode, err := service.GetAllClosedTrades(user)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}
		return c.JSON(http.StatusOK, allClosedTrades)
	}
}

func SetFutureOrder(service services.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(trade.FutureOrderRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(http.StatusBadRequest, response)
		}

		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorRespone("", errors.New("bad user data"))
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.SetFutureOrder(*request, user)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, models.NewRespone("the future order successfully set"))
	}
}

func GetAllFutureOrders(service services.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {

		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorRespone("", errors.New("bad user data"))
			return c.JSON(http.StatusBadRequest, response)
		}

		allFutureOrders, statusCode, err := service.GetAllFutureOrders(user)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}
		return c.JSON(http.StatusOK, allFutureOrders)
	}
}