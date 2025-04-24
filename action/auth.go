package action

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	"github.com/briskt/go-htmx-app/api"
	"github.com/briskt/go-htmx-app/app"
	"github.com/briskt/go-htmx-app/core"
	"github.com/briskt/go-htmx-app/email"
	"github.com/briskt/go-htmx-app/log"
)

const (
	// http cookie access token
	AccessTokenSessionKey = "AccessToken"

	// http param and session key for ReturnTo
	ReturnToParam      = "return-to"
	ReturnToSessionKey = "ReturnTo"
)

// swagger:operation GET /auth/login Authentication AuthLogin
// AuthLogin
//
// Start the SAML login process
// ---
//
//	responses:
//	  '200':
//	    description: returns a "RedirectURL" key with the saml idp url that has a saml request
func (a *App) authLogin(c echo.Context) error {
	if err := clearSession(c); err != nil {
		return api.NewAppError(err, api.ErrorClearingSession, http.StatusInternalServerError)
	}

	setReturnToInSession(c)

	redirectURL, err := a.samlProvider.BuildAuthURL("")
	if err != nil {
		err = fmt.Errorf("failed to determine what the saml authentication url should be: %w", err)
		return api.NewAppError(err, api.ErrorGettingAuthURL, http.StatusInternalServerError)
	}

	// Reply with a 302 redirect to the IdP
	return c.Redirect(http.StatusFound, redirectURL)
}

// swagger:operation POST /auth/callback Authentication callback
// AuthCallback
//
// Complete the SAML login process
// ---
//
//	responses:
//	  '302':
//	    description: redirects to the home page
func (a *App) authCallback(c echo.Context) error {
	err := clearSession(c)
	if err != nil {
		return api.NewAppError(err, api.ErrorClearingSession, http.StatusInternalServerError)
	}

	staffID, err := a.samlProvider.GetUser(c)
	if err != nil {
		err = fmt.Errorf("auth response error: %w", err)
		return api.NewAppError(err, api.ErrorAuthProvidersCallback, http.StatusInternalServerError)
	}

	token, err := core.NewToken(toCtx(c), Tx(c), staffID)
	if err != nil {
		return err
	}

	user, err := core.FindUserByToken(toCtx(c), Tx(c), token)
	if err != nil {
		return err
	}

	// set person on log context
	log.SetUser(toCtx(c), user.EmployeeID, email.MaskString(user.GetDisplayName()), email.MaskEmail(user.GetEmail()))

	// Set the authentication token in the session, which by default is stored in a cookie
	err = sessionSetValue(c, AccessTokenSessionKey, token)
	if err != nil {
		return api.NewAppError(err, api.ErrorStoringAccessToken, http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusFound, getLoginSuccessRedirectURL(c))
}

// swagger:operation GET /auth/logout Authentication AuthLogout
// AuthLogout
//
// Logout of application
// ---
//
//	responses:
//	  '302':
//	    description: redirect to UI
func (a *App) authLogout(c echo.Context) error {
	err := clearSession(c)
	if err != nil {
		return api.NewAppError(err, api.ErrorClearingSession, http.StatusInternalServerError)
	}

	logoutURL := a.samlProvider.IdentityProviderSLOURL
	if logoutURL == "" {
		err = errors.New("IdentityProviderSLOURL value is empty")
		return err
	}

	redirectURL := fmt.Sprintf("%s?ReturnTo=%s", logoutURL, app.Env.AppURL)
	return c.Redirect(http.StatusFound, redirectURL)
}

// swagger:operation GET /auth/logout-callback Authentication AuthLogoutCallback
// AuthLogoutCallback
//
// Receive redirect from IdP after logout
// ---
//
//	responses:
//	  '302':
//	    description: redirect to UI
func (a *App) authLogoutCallback(c echo.Context) error {
	err := clearSession(c)
	if err != nil {
		return api.NewAppError(err, api.ErrorClearingSession, http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusFound, "/auth/login")
}

// getLoginSuccessRedirectURL generates the URL for redirection after a successful login
func getLoginSuccessRedirectURL(c echo.Context) string {
	profileURL := app.Env.AppURL
	params := ""

	returnTo := getReturnTo(c)
	if len(returnTo) > 0 && returnTo != profileURL {
		params = "?" + ReturnToParam + "=" + url.QueryEscape(returnTo)
	}

	return profileURL + params
}
