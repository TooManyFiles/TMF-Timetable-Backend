package errorhandling

var DBErrorsSource ErrorSource = ErrorSource{
	Name: "DB",
	ParseError: func(err error) (CustomError, error) {
		return CustomError{}, nil
	},
}
