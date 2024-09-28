package errorhandling

var GenericErrorsSource ErrorSource = ErrorSource{
	Name: "Generic",
	ParseError: func(err error) (CustomError, error) {
		return CustomError{}, nil
	},
}

var FailedToGenerateTraceId = CustomErrorPreset{
	Code:        401,
	UserMessage: "",
	DevMessage:  "Error while generating trace id.",
	LogMessage:  "Error while generating trace id.",
	Source:      GenericErrorsSource,
	Level:       ErrorWARN,
	HttpCode:    200,
}

// Preset for 500 Internal Server Error
var GenericInternalServerError = CustomErrorPreset{
	Code:        500,
	UserMessage: "Internal server error",
	DevMessage:  "An unexpected error occurred on the server",
	LogMessage:  "Internal server error encountered",
	Source:      GenericErrorsSource,
	Level:       ErrorMedium,
	HttpCode:    500,
}
