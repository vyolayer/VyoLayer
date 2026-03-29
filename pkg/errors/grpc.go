package errors

import (
	"net/http"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPC converts an *AppError into a gRPC status error.
// If appErr is nil, it returns nil.
func ToGRPC(appErr *AppError) error {
	if appErr == nil {
		return nil
	}
	return status.Error(httpStatusToGRPCCode(appErr.HTTPStatus), appErr.Message)
}

// httpStatusToGRPCCode maps HTTP status codes to gRPC codes.
func httpStatusToGRPCCode(httpStatus int) codes.Code {
	switch httpStatus {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusUnprocessableEntity:
		return codes.InvalidArgument
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusGatewayTimeout, http.StatusRequestTimeout:
		return codes.DeadlineExceeded
	case http.StatusInternalServerError:
		return codes.Internal
	default:
		return codes.Internal
	}
}

// FromGRPC converts any gRPC error to AppError.
func FromGRPC(err error) *AppError {
	if err == nil {
		return nil
	}

	if appErr, ok := As(err); ok {
		return appErr
	}

	st, ok := status.FromError(err)
	if !ok {
		return Wrap(err, ErrInternalUnexpected, "Unexpected upstream error")
	}

	code := grpcCodeToErrorCode(st.Code())
	appErr := NewWithMessage(code, st.Message())
	appErr.Wrapped = err
	appErr.WithMetadata("grpc_code", st.Code().String())

	for _, detail := range st.Details() {
		switch d := detail.(type) {
		case *errdetails.BadRequest:
			violations := make([]map[string]string, 0, len(d.GetFieldViolations()))
			for _, v := range d.GetFieldViolations() {
				violations = append(violations, map[string]string{
					"field":       v.GetField(),
					"description": v.GetDescription(),
				})
			}
			if len(violations) > 0 {
				appErr.WithMetadata("validation_errors", violations)
			}
		case *errdetails.ErrorInfo:
			appErr.WithMetadata("grpc_error_reason", d.GetReason())
			// if len(d.GetMetadata()) > 0 {
			// 	appErr.WithMetadata("grpc_error_metadata", d.GetMetadata())
			// }
		case *errdetails.RetryInfo:
			if d.GetRetryDelay() != nil {
				appErr.WithMetadata("retry_delay", d.GetRetryDelay().AsDuration().String())
			}
		}
	}

	return appErr
}

func grpcCodeToErrorCode(code codes.Code) ErrorCode {
	switch code {
	case codes.InvalidArgument:
		return ErrRequestInvalidParams
	case codes.NotFound:
		return ErrResourceNotFound
	case codes.AlreadyExists:
		return ErrResourceAlreadyExists
	case codes.PermissionDenied:
		return ErrAuthForbidden
	case codes.Unauthenticated:
		return ErrAuthUnauthorized
	case codes.ResourceExhausted:
		return ErrRequestRateLimited
	case codes.FailedPrecondition:
		return ErrBusinessRuleViolation
	case codes.Aborted:
		return ErrResourceConflict
	case codes.OutOfRange:
		return ErrValidationOutOfRange
	case codes.Unimplemented:
		return ErrInternalNotImplemented
	case codes.Internal:
		return ErrInternalUnexpected
	case codes.Unavailable:
		return ErrExternalServiceUnavailable
	case codes.DeadlineExceeded:
		return ErrRequestTimeout
	case codes.Canceled:
		return ErrRequestTimeout
	default:
		return ErrInternalUnexpected
	}
}
