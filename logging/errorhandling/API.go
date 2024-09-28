package errorhandling

var APIErrorsSource ErrorSource = ErrorSource{
	Name: "API",
	ParseError: func(err error) (CustomError, error) {
		return CustomError{}, nil
	},
}

// Preset for 401 Unauthorized
var UnauthorizedError = CustomErrorPreset{
	Code:        401,
	UserMessage: "Unauthorized access",
	DevMessage:  "Authentication required to access this resource",
	LogMessage:  "Unauthorized request made to the server",
	Source:      APIErrorsSource,
	Level:       ErrorWrongUsage,
	HttpCode:    401,
}

// Preset for 403 Forbidden
var ForbiddenError = CustomErrorPreset{
	Code:        403,
	UserMessage: "Forbidden",
	DevMessage:  "Access denied to this resource",
	LogMessage:  "Forbidden request attempted",
	Source:      APIErrorsSource,
	Level:       ErrorWrongUsage,
	HttpCode:    403,
}

// Preset for 404 Not Found
var NotFoundError = CustomErrorPreset{
	Code:        404,
	UserMessage: "Resource not found",
	DevMessage:  "The requested resource could not be found on the server",
	LogMessage:  "Resource not found",
	Source:      APIErrorsSource,
	Level:       ErrorWrongUsage,
	HttpCode:    404,
}

// Preset for 405 Method Not Allowed
var MethodNotAllowedError = CustomErrorPreset{
	Code:        405,
	UserMessage: "Method not allowed",
	DevMessage:  "The HTTP method used is not allowed for this resource",
	LogMessage:  "Method not allowed for this endpoint",
	Source:      APIErrorsSource,
	Level:       ErrorWrongUsage,
	HttpCode:    405,
}
