package action

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/briskt/go-htmx-app/api"
	"github.com/briskt/go-htmx-app/app"
	"github.com/briskt/go-htmx-app/log"
)

// customHTTPErrorHandler adds details to an error and renders the error with echo.Render.
// If the HTTP status code provided is in the 300 family, echo.Redirect is used instead.
func customHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	var appErr *api.AppError
	if !errors.As(err, &appErr) {
		appErr = api.NewAppError(err, api.ErrorInternal, http.StatusInternalServerError)
	}
	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		log.Debug("converting echo Error to api.AppError")
		appErr.HttpStatus = echoErr.Code
		appErr.Message = echoErr.Error()
		appErr.Err = echoErr.Internal
	}

	if appErr.Extras == nil {
		appErr.Extras = map[string]any{}
	}

	appErr.Extras["key"] = appErr.Key
	appErr.Extras["status"] = appErr.HttpStatus
	appErr.Extras["redirectURL"] = appErr.RedirectURL
	appErr.Extras["method"] = c.Request().Method
	appErr.Extras["URI"] = c.Request().RequestURI
	appErr.Extras["employeeID"] = CurrentUser(c).EmployeeID

	address, _ := getClientIPAddress(c.Request())
	appErr.Extras["IP"] = address

	appErr.LoadTranslatedMessage(c)

	// clear out debugging info if not in development or test
	if app.Env.AppEnv == app.EnvDevelopment || app.Env.AppEnv == app.EnvTest {
		appErr.DebugMsg = err.Error()
	} else {
		appErr.Extras = map[string]any{}
	}

	if appErr.HttpStatus >= 300 && appErr.HttpStatus <= 399 {
		if appErr.RedirectURL == "" {
			appErr.RedirectURL = app.Env.AppURL + "/logged-out?appError=" + appErr.Message
		}
		err = c.Redirect(appErr.HttpStatus, appErr.RedirectURL)
		if err != nil {
			log.Errorf("c.Redirect returned error: %v", err)
		}
		return
	}

	err = c.JSON(appErr.HttpStatus, appErr)
	if err != nil {
		appErr.Extras = map[string]any{}
		appErr.Err = fmt.Errorf("unable to encode extras for error (%s): %w", err, appErr.Err)
		_ = c.JSON(http.StatusInternalServerError, appErr)
	}
}

// getClientIPAddress gets the client IP address from CF-Connecting-IP or RemoteAddr
func getClientIPAddress(req *http.Request) (net.IP, error) {
	// https://developers.cloudflare.com/fundamentals/get-started/reference/http-request-headers/#cf-connecting-ip
	if cf := req.Header.Get("CF-Connecting-IP"); cf != "" {
		return net.ParseIP(cf), nil
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("userip: %q is not IP:port, %w", req.RemoteAddr, err)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not a valid IP address, %w", req.RemoteAddr, err)
	}

	return userIP, nil
}
