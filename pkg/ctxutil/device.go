package ctxutil

import (
	"context"
	"strings"

	"github.com/ua-parser/uap-go/uaparser"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Device struct {
	IP             string
	UserAgent      string
	OS             string
	OSVersion      string
	Browser        string
	BrowserVersion string
}

// InjectDevice injects device information into the context
func InjectDeviceInfo(ctx context.Context, ip, userAgent string) context.Context {
	parser := uaparser.NewFromSaved()
	client := parser.Parse(userAgent)

	device := &Device{
		IP:             ip,
		UserAgent:      userAgent,
		OS:             client.Os.Family,
		OSVersion:      strings.Join([]string{client.Os.Major, client.Os.Minor, client.Os.Patch}, "."),
		Browser:        client.UserAgent.Family,
		BrowserVersion: strings.Join([]string{client.UserAgent.Major, client.UserAgent.Minor, client.UserAgent.Patch}, "."),
	}

	return context.WithValue(ctx, "device", device)
}

// ExtractDeviceInfo extracts device information from the context
func ExtractDeviceInfo(ctx context.Context) (*Device, error) {
	val, ok := ctx.Value("device").(*Device)
	if !ok {
		return nil, status.Error(codes.Internal, "missing device info")
	}
	return val, nil
}
