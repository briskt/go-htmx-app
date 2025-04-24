package api

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
)

type ErrorKey struct {
	name string
}

func (e ErrorKey) String() string {
	return e.name
}

// AppError holds information that is helpful in logging and reporting api errors
// swagger:model
type AppError struct {
	Err error `json:"-"`

	// Don't change the value of these Key entries without making a corresponding change on the UI,
	// since these will be converted to human-friendly texts for presentation to the user
	Key ErrorKey `json:"key"`

	HttpStatus int `json:"status"`

	// detailed error message for debugging, only provided in development environment
	DebugMsg string `json:"debug_msg,omitempty"`

	// user-facing error message
	Message string `json:"message"`

	// Extra data providing detail about the error condition, only provided in development environment
	Extras map[string]any `json:"extras,omitempty"`

	// URL to redirect, if HttpStatus is in 300-series
	RedirectURL string `json:"-"`
}

func (a *AppError) Error() string {
	if a.Err == nil {
		return ""
	}
	return a.Err.Error()
}

func (a *AppError) Unwrap() error {
	return a.Err
}

// Is tests the provided error against the AppError by comparing the Key. If they match, Is returns true.
// This can be used like `if errors.Is(err, &api.AppError{Key: ErrorUserNotFound})`
// See https://pkg.go.dev/errors#Is for more information.
func (a *AppError) Is(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Key == a.Key
}

// NewAppError returns a new AppError with its Err, Key and HttpStatus set
func NewAppError(err error, key ErrorKey, status int) *AppError {
	a := &AppError{
		Err:        err,
		Key:        key,
		HttpStatus: status,
	}
	a.Extras = map[string]any{"function": getFunctionName(2)}
	return a
}

// getFunctionName provides the filename, line number, and function name of the caller, skipping the top `skip`
// functions on the stack.
func getFunctionName(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d %s", file, line, fn.Name())
}

// LoadTranslatedMessage assigns the error message by translating the Key into a user-friendly string, either
// from a list of translated strings (see errors.en) or by breaking down the Key into individual words
func (a *AppError) LoadTranslatedMessage(c echo.Context) {
	key := a.Key

	if a.HttpStatus == http.StatusInternalServerError {
		key = ErrorInternal
	}

	a.Message = keyToReadableString(key.String())
}

// keyToReadableString takes a key like ErrorSomethingSomethingOther and returns Error something something other
// although it will lose initial lowercase letters, if it has a non-initial uppercase letter
func keyToReadableString(key string) string {
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	words := re.FindAllString(key, -1)

	if len(words) == 0 {
		return key
	}

	if len(words) > 1 && words[0] == "Error" {
		words = words[1:]
	}

	count := len(words)
	newWords := []string{}

	// Lowercase all but first word.
	for i := 0; i < count; i++ {
		// If a word is longer than one character, just use it as is
		if len(words[i]) > 1 {
			newWords = append(newWords, strings.ToLower(words[i]))
			continue
		}

		// Combine single character words
		next := words[i]
		for j := i + 1; j < count; j++ {
			if len(words[j]) == 1 {
				next += words[j]
				i++ // avoid reprocessing the same word
			} else {
				break
			}
		}
		newWords = append(newWords, strings.ToLower(next))
	}

	firstUpper := strings.ToUpper(newWords[0][0:1])
	newWords[0] = firstUpper + newWords[0][1:]

	return strings.Join(newWords, " ")
}
