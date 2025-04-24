package api

var (
	TODO = ErrorKey{"unknown error"}

	// General

	ErrorInvalidRequestBody = ErrorKey{"ErrorInvalidRequestBody"}
	ErrorClearingSession    = ErrorKey{"ErrorClearingSession"}
	ErrorRenderingTemplate  = ErrorKey{"ErrorRenderingTemplate"}

	// HTTP error types

	ErrorInternal         = ErrorKey{"ErrorInternal"}
	ErrorNotFound         = ErrorKey{"ErrorNotFound"}
	ErrorNotAuthenticated = ErrorKey{"ErrorNotAuthenticated"}

	// Authentication

	ErrorAuthProvidersCallback = ErrorKey{"ErrorAuthProvidersCallback"}
	ErrorGeneratingRandomToken = ErrorKey{"ErrorGeneratingRandomToken"}
	ErrorCreatingAccessToken   = ErrorKey{"ErrorCreatingAccessToken"}
	ErrorStoringAccessToken    = ErrorKey{"ErrorStoringAccessToken"}
	ErrorGettingAuthURL        = ErrorKey{"ErrorGettingAuthURL"}

	// User

	ErrorUserNotFound      = ErrorKey{"ErrorUserNotFound"}
	ErrorPasswordSetFailed = ErrorKey{"ErrorPasswordSetFailed"}
)
