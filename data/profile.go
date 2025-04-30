package data

import "github.com/a-h/templ"

type ProfileView struct {
	DisplayName   string
	Enabled       bool
	HelpCenterURL templ.SafeURL
	AppName       string
	LastLogin     string
	UserID        string
	Username      string
}