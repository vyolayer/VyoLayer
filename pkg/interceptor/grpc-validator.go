package interceptor

import (
	"context"
	"errors"

	"buf.build/go/protovalidate"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func GRPCValidationInterceptor(v protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// 1. Cast the generic request to a proto.Message
		if msg, ok := req.(proto.Message); ok {
			if v == nil {
				return nil, status.Error(codes.Internal, "validation system not configured")
			}

			// 2. Run the Buf validation rules
			if err := v.Validate(msg); err != nil {

				var valErr *protovalidate.ValidationError
				if errors.As(err, &valErr) {
					// 3. Format the standard gRPC BadRequest error
					st := status.New(codes.InvalidArgument, "Invalid request parameters")
					br := &errdetails.BadRequest{}

					for _, violation := range valErr.Violations {
						var field string
						var description string
						if violation != nil && violation.Proto != nil {
							field = protovalidate.FieldPathString(violation.Proto.GetField())
							description = violation.Proto.GetMessage()
						}

						br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
							Field:       field,
							Description: description,
						})
					}

					st, attachErr := st.WithDetails(br)
					if attachErr != nil {
						return nil, status.Error(codes.InvalidArgument, "Invalid request parameters")
					}

					// Return the error immediately; the handler is NEVER called
					return nil, st.Err()
				}
				return nil, status.Error(codes.Internal, "Validation system failed")
			}
		}

		// 4. If validation passes, proceed to the actual service handler
		return handler(ctx, req)
	}
}
