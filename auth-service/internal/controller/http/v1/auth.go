package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gitverse.ru/volatex/backend/internal/controller/http/v1/request"
	"gitverse.ru/volatex/backend/internal/controller/http/v1/response"
	"gitverse.ru/volatex/backend/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

// @Summary     Register
// @Description Register a new user
// @ID          register
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body request.Register true "User registration details"
// @Success     201 {object} response.User
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /auth/register [post]
func (r *Auth) register(ctx *fiber.Ctx) error {
	var body request.Register

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - register")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "http - v1 - register")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		r.l.Error(err, "http - v1 - register")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to process password")
	}

	user, err := r.u.Register(
		ctx.UserContext(),
		entity.User{
			Email:    body.Email,
			Password: string(hashedPassword),
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - register")
		return errorResponse(ctx, http.StatusInternalServerError, "registration failed")
	}

	return ctx.Status(http.StatusCreated).JSON(response.User{
		Id:    user.Id,
		Email: user.Email,
	})
}
