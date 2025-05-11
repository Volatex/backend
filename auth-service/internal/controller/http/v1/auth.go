package v1

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gitverse.ru/volatex/backend/internal/controller/http/v1/request"
	"gitverse.ru/volatex/backend/internal/controller/http/v1/response"
	"gitverse.ru/volatex/backend/internal/entity"
	"golang.org/x/crypto/argon2"
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

	const (
		memory      = 64 * 1024 // 64 MB
		iterations  = 1
		parallelism = 4
		saltLength  = 16
		keyLength   = 32
	)

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		r.l.Error(err, "http - v1 - register")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to generate salt")
	}

	hash := argon2.IDKey([]byte(body.Password), salt, iterations, memory, parallelism, keyLength)

	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	hashedPassword := fmt.Sprintf("argon2id$%d$%d$%d$%s$%s", memory, iterations, parallelism, saltB64, hashB64)

	user, err := r.u.Register(
		ctx.UserContext(),
		entity.User{
			Email:    body.Email,
			Password: hashedPassword,
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

// @Summary     Verify email
// @Description Confirm user's email with a verification code
// @ID          verify-email
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body request.VerifyEmail true "Email and verification code"
// @Success     200
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /auth/verify-email [post]
func (r *Auth) verifyEmail(ctx *fiber.Ctx) error {
	var body request.VerifyEmail

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - verifyEmail")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "http - v1 - verifyEmail")
		return errorResponse(ctx, http.StatusBadRequest, "validation failed")
	}

	err := r.u.VerifyEmail(ctx.UserContext(), body.Email, body.Code)
	if err != nil {
		r.l.Error(err, "http - v1 - verifyEmail")
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.SendStatus(http.StatusOK)
}
