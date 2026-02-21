package apperror

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	arrayReplaceRegexString = `(\[.+?\])`
	tagSplitLimit           = 2
)

var arrayReplaceRegex = regexp.MustCompile(arrayReplaceRegexString)
var (
	errIsRequired          = errors.New("is required")
	errShouldBeAnUUID      = errors.New("should be a valid UUID")
	errShouldBeAnEmail     = errors.New("should be a valid email")
	errInvalidValue        = errors.New("has invalid value")
	errMinThree            = errors.New("must be at least 3 characters")
	errMaxFifty            = errors.New("cannot exceed 50 characters")
	errMinSix              = errors.New("must be at least 6 characters")
	errInvalidGender       = errors.New("must be one of: M, F, O")
	errInvalidCookingLevel = errors.New("must be one of: Novice, Intermediate, Proficient, Expert")
	errInvalidDate         = errors.New("must be a valid date")
	errInvalidNumeric      = errors.New("invalid numeric value")
	errInvalidBoolean      = errors.New("invalid boolean value")
)

//
// --------------------
// Custom Error Map
// --------------------
//

var customErrors = map[string]error{
	"name.required": errIsRequired,
	"name.min":      errMinThree,
	"name.max":      errMaxFifty,

	"email.required": errIsRequired,
	"email.email":    errShouldBeAnEmail,

	"password.required": errIsRequired,
	"password.min":      errMinSix,
	"password.max":      errMaxFifty,

	"gender.required": errIsRequired,
	"gender.gender":   errInvalidGender,

	"dob.required": errIsRequired,

	"levelOfCooking.required":       errIsRequired,
	"levelOfCooking.levelOfCooking": errInvalidCookingLevel,
}

func CustomValidationError(err error) []map[string]string {
	errs := make([]map[string]string, 0)

	var (
		validationErrors   validator.ValidationErrors
		numError           *strconv.NumError
		unmarshalTypeError *json.UnmarshalTypeError
	)

	switch {
	case errors.As(err, &validationErrors):
		for _, e := range validationErrors {
			errorMap := make(map[string]string)

			key := e.Field() + "." + e.Tag()
			newKey := arrayReplaceRegex.ReplaceAllString(key, ".dive")

			if v, ok := customErrors[newKey]; ok {
				errorMap[e.Field()] = v.Error()
			} else {
				errorMap[e.Field()] = fmt.Sprintf("validation failed on %s", e.Tag())
			}

			errs = append(errs, errorMap)
		}

		return errs

	case errors.As(err, &numError):
		switch numError.Func {
		case "ParseBool":
			errs = append(errs, map[string]string{"body": errInvalidBoolean.Error()})
		case "ParseInt", "ParseUint", "ParseFloat", "Atoi":
			errs = append(errs, map[string]string{"body": errInvalidNumeric.Error()})
		default:
			errs = append(errs, map[string]string{"body": errInvalidValue.Error()})
		}
		return errs

	case errors.As(err, &unmarshalTypeError):
		errs = append(errs, map[string]string{
			unmarshalTypeError.Field: fmt.Sprintf("%v cannot be %v", unmarshalTypeError.Field, unmarshalTypeError.Value),
		})
		return errs
	}

	if errors.Is(err, io.EOF) {
		errs = append(errs, map[string]string{"body": "request body cannot be empty"})
	} else {
		errs = append(errs, map[string]string{"unknown": fmt.Sprintf("unsupported error: %v", err)})
	}

	return errs
}

func RegisterTags(v *validator.Validate) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tags := []string{"json", "uri", "form"}

		for _, key := range tags {
			tag := fld.Tag.Get(key)
			name := strings.SplitN(tag, ",", tagSplitLimit)[0]

			if name == "-" {
				return ""
			} else if len(name) != 0 {
				return name
			}
		}
		return ""
	})
}
