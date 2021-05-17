package errors

import "errors"

var (
	UserNotFoundError                         = errors.New("user not found")
	UserReturnWrongLenError                   = errors.New("search of users wrong len")
	EventAddressIsRequiredError               = errors.New("event topic is required")
	DatabaseDsnIsRequiredError                = errors.New("database dsn is required")
	DatabaseTransactionAtContextNotFoundError = errors.New("database transaction not found in context")
	DatabaseUniqueEntryViolationError         = errors.New("database unique entry violation error")
	RepositoryProjectionsIsRequiredError      = errors.New("at least one projection are required")
	RepositoryIdIsRequiredError               = errors.New("id is required for repository")
	ObjsIsNotSliceValidatorError              = errors.New("objs arg is not a slice")
)
