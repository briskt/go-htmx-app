package action

import (
	"fmt"
	"net/http"
)

func (s *Suite) TestApp_authLogin() {
	const returnToPath = "/foo"
	path := fmt.Sprintf("/auth/login?%s=%s", ReturnToParam, returnToPath)
	response := s.requestResponse("GET", path, "", nil)
	s.Equal(http.StatusFound, response.Code)
	s.Contains(response.Header().Get("Location"), "http://localhost:8106/module.php/saml/idp/singleSignOnService?SAMLRequest=")
	s.Equal(returnToPath, s.session.Values[ReturnToSessionKey])
}

func (s *Suite) TestApp_authLogout() {
	response := s.requestResponse("GET", fmt.Sprintf("/auth/logout"), "", nil)
	s.Equal(http.StatusFound, response.Code)
	s.Contains(response.Header().Get("Location"), "http://localhost:8106/module.php/saml/idp/singleLogout")
	s.Len(s.session.Values, 0)
}

func (s *Suite) TestApp_authLogoutCallback() {
	response := s.requestResponse("GET", fmt.Sprintf("/auth/logout-callback"), "", nil)
	s.Equal(http.StatusFound, response.Code)
	s.Contains(response.Header().Get("Location"), "/auth/login")
	s.Len(s.session.Values, 0)
}
