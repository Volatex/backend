package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitverse.ru/volatex/backend/market-service/internal/controller/http/v1/request"
	"gitverse.ru/volatex/backend/market-service/internal/entity"
	"net/http"
)

// @Summary     Save strategy
// @Description Save user's trading strategy
// @ID          save-strategy
// @Tags        strategy
// @Accept      json
// @Produce     json
// @Param       request body request.SaveStrategy true "Strategy details"
// @Success     200
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /strategy/add [post]
// @Security ApiKeyAuth
func (r *Market) saveStrategy(ctx *fiber.Ctx) error {
	var body request.SaveStrategy

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - saveStrategy")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "http - v1 - saveStrategy")
		return errorResponse(ctx, http.StatusBadRequest, "validation failed")
	}

	userIDRaw := ctx.Locals("user_id")
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		r.l.Error(nil, "http - v1 - saveStrategy - user_id missing or invalid")
		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		r.l.Error(err, "http - v1 - saveStrategy - invalid UUID")
		return errorResponse(ctx, http.StatusBadRequest, "invalid user ID")
	}

	strategy := &entity.Strategy{
		UserID:       userUUID,
		Figi:         body.Figi,
		BuyPrice:     body.BuyPrice,
		BuyQuantity:  body.BuyQuantity,
		SellPrice:    body.SellPrice,
		SellQuantity: body.SellQuantity,
		//TinkoffToken: body.TinkoffToken,
	}

	if err := r.u.Create(ctx.UserContext(), strategy); err != nil {
		r.l.Error(err, "http - v1 - saveStrategy")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to save strategy")
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

// @Summary     Save user's token
// @Description Save user's token
// @ID          save-user-token
// @Tags        strategy
// @Accept      json
// @Produce     json
// @Param       request body request.SaveUserToken true "Strategy details"
// @Success     200
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /strategy/add_token [post]
// @Security ApiKeyAuth
func (r *Market) saveUserToken(ctx *fiber.Ctx) error {
	var body request.SaveUserToken
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "validation failed")
	}

	userIDRaw := ctx.Locals("user_id")
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid user ID")
	}

	token := &entity.UserToken{
		UserID:       userUUID,
		TinkoffToken: body.TinkoffToken,
	}

	if err := r.u.SaveUserToken(ctx.UserContext(), token); err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, "failed to save token")
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
