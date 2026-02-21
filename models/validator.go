package models

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	Validator *validator.Validate

	errRegisterValidationFailed = errors.New("failed to register validation")
	nonZeroExitCode             = 1
)

func init() {
	var ok bool

	if Validator, ok = binding.Validator.Engine().(*validator.Validate); ok {

		if err := registerValidationTags(context.Background(), Validator); err != nil {
			os.Exit(nonZeroExitCode)
		}
	}
}

func registerValidationTags(ctx context.Context, v *validator.Validate) error {

	if err := v.RegisterValidation("gender", genderValidator); err != nil {
		slog.ErrorContext(ctx, errRegisterValidationFailed.Error(), slog.Any("error", err))
		return err
	}

	if err := v.RegisterValidation("levelOfCooking", levelOfCookingValidator); err != nil {
		slog.ErrorContext(ctx, errRegisterValidationFailed.Error(), slog.Any("error", err))
		return err
	}

	if err := v.RegisterValidation("strongPassword", strongPasswordValidator); err != nil {
		slog.ErrorContext(ctx, errRegisterValidationFailed.Error(), slog.Any("error", err))
		return err
	}

	if err := v.RegisterValidation("date", dateValidator); err != nil {
		slog.ErrorContext(ctx, errRegisterValidationFailed.Error(), slog.Any("error", err))
		return err
	}

	return nil
}

var genderValidator validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(Gender)
	if !ok {
		return false
	}

	return value.IsValid()
}

var levelOfCookingValidator validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(LevelOfCooking)
	if !ok {
		return false
	}

	return value.IsValid()
}

var strongPasswordValidator validator.Func = func(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	// min 8, max 32
	if len(password) < 8 || len(password) > 32 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#~$%^&*()+|_]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}

var dateValidator validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	// expecting YYYY-MM-DD
	_, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	return err == nil
}
