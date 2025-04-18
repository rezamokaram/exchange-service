package handlers

import (
	"net/http"

	"github.com/rezamokaram/exchange-service/internal"
	"github.com/rezamokaram/exchange-service/models"
	"github.com/rezamokaram/exchange-service/models/trade"
	userModels "github.com/rezamokaram/exchange-service/models/user"

	"github.com/labstack/echo/v4"
)

// OpenTrade opens a new trade
// @Summary Open a trade
// @Description Opens a new trade
// @Tags Trade
// @Accept  json
// @Produce  json
// @Param   body  body      trade.OpenTradeRequest  true  "Open Trade"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /open-trade [post]
func OpenTrade(service internal.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(trade.OpenTradeRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorResponse("trade failed", err.Error())
			return c.JSON(http.StatusBadRequest, response)
		}

		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("trade failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.OpenTrade(*request, user)
		if err != nil {
			response := models.NewErrorResponse("trade failed", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, models.NewResponse("the trade successfully opened"))
	}
}

// CloseTrade closes an existing trade
// @Summary Close a trade
// @Description Closes an existing trade
// @Tags Trade
// @Accept  json
// @Produce  json
// @Param   body  body      trade.ClosedTradeRequest  true  "Close Trade"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /close-trade [post]
func CloseTrade(service internal.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(trade.ClosedTradeRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorResponse("trade failed", err.Error())
			return c.JSON(http.StatusBadRequest, response)
		}

		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("trade failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.CloseTrade(*request, user)
		if err != nil {
			response := models.NewErrorResponse("trade failed", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, models.NewResponse("the closing request successfully processed"))
	}
}

// GetAllOpenTrades retrieves all open trades for the authenticated user
// @Summary Get all open trades
// @Description Retrieves all open trades for the authenticated user
// @Tags Trade
// @Produce  json
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200   {array}   trade.OpenTrade
// @Failure 400   {object}  models.Response
// @Router /open-trade/get-all [get]
func GetAllOpenTrades(service internal.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {

		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("trade failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		allOpenTrades, statusCode, err := service.GetAllOpenTrades(user)
		if err != nil {
			response := models.NewErrorResponse("trade failed", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, allOpenTrades)
	}
}

// GetAllClosedTrades retrieves all closed trades for the authenticated user
// @Summary Get all closed trades
// @Description Retrieves all closed trades for the authenticated user
// @Tags Trade
// @Produce  json
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200   {array}   trade.ClosedTrade
// @Failure 400   {object}  models.Response
// @Router /close-trade/get-all [get]
func GetAllClosedTrades(service internal.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {

		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("trade failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		allClosedTrades, statusCode, err := service.GetAllClosedTrades(user)
		if err != nil {
			response := models.NewErrorResponse("trade failed", err.Error())
			return c.JSON(statusCode, response)
		}
		return c.JSON(http.StatusOK, allClosedTrades)
	}
}

// SetFutureOrder sets a future order for the authenticated user
// @Summary Set a future order
// @Description Sets a future order
// @Tags Trade
// @Accept  json
// @Produce  json
// @Param   body  body      trade.FutureOrderRequest  true  "Future Order"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /future-order [post]
func SetFutureOrder(service internal.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(trade.FutureOrderRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorResponse("trade failed", err.Error())
			return c.JSON(http.StatusBadRequest, response)
		}

		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("trade failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.SetFutureOrder(*request, user)
		if err != nil {
			response := models.NewErrorResponse("trade failed", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, models.NewResponse("the future order successfully set"))
	}
}

// DeleteFutureOrder handles the deletion of a future order
// @Summary Delete Future Order
// @Description Deletes a specified future order belonging to the authenticated user
// @Tags Trade
// @Accept  json
// @Produce  json
// @Param   order_id  body    trade.DeleteFutureOrderRequest  true  "Delete Future Order"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200  {object}  models.Response  "the future order successfully deleted"
// @Failure 400  {object}  models.Response  "Bad Request"
// @Failure 403  {object}  models.Response  "Forbidden - Order belongs to another user"
// @Failure 500  {object}  models.Response  "Internal Server Error"
// @Router /future-order [delete]
func DeleteFutureOrder(service internal.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(trade.DeleteFutureOrderRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorResponse("trade failed", err.Error())
			return c.JSON(http.StatusBadRequest, response)
		}

		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("trade failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.DeleteFutureOrder(*request, user)
		if err != nil {
			response := models.NewErrorResponse("trade failed", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, models.NewResponse("the future order successfully deleted"))
	}
}

// GetAllFutureOrders retrieves all future orders for the authenticated user
// @Summary Get all future orders
// @Description Retrieves all future orders
// @Tags Trade
// @Produce  json
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200   {array}   trade.FutureOrder
// @Failure 400   {object}  models.Response
// @Router /future-order/get-all [get]
func GetAllFutureOrders(service internal.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {

		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("trade failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		allFutureOrders, statusCode, err := service.GetAllFutureOrders(user)
		if err != nil {
			response := models.NewErrorResponse("trade failed", err.Error())
			return c.JSON(statusCode, response)
		}
		return c.JSON(http.StatusOK, allFutureOrders)
	}
}

// FilterTrades filters closed trades based on criteria for the authenticated user
// @Summary Filter closed trades
// @Description Filters closed trades based on given criteria
// @Tags Trade
// @Accept  json
// @Produce  json
// @Param   body  body      trade.FilterTradesRequest  true  "Filter Criteria"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200   {object}  trade.FilterTradesResponse
// @Failure 400   {object}  models.Response
// @Router /close-trade/filter-all [get]
func FilterTrades(service internal.TradeService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(trade.FilterTradesRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorResponse("trade failed", err.Error())
			return c.JSON(http.StatusBadRequest, response)
		}

		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("trade failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		allTargetTrades, statusCode, err := service.FilterClosedTrades(user, *request)
		if err != nil {
			response := models.NewErrorResponse("trade failed", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, allTargetTrades)
	}
}
