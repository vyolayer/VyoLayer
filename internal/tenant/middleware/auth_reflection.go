package middleware

import (
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	// Import your generated proto code to access the extensions
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

// RouteAuthConfig holds all the security rules extracted from the Protobuf schema
type RouteAuthConfig struct {
	SkipOrgCheck          bool
	RequiredOrgPermission string
}

// routeAuthCache stores the reflection results so we only look them up once per route.
var routeAuthCache sync.Map

// getRouteAuthConfig reads the protobuf method options at runtime and caches them.
func getRouteAuthConfig(fullMethod string) RouteAuthConfig {
	// Check the cache first (O(1) lookup, extremely fast)
	if val, ok := routeAuthCache.Load(fullMethod); ok {
		return val.(RouteAuthConfig)
	}

	// Default safe configuration (fail closed)
	config := RouteAuthConfig{
		SkipOrgCheck:          false,
		RequiredOrgPermission: "",
	}

	// Convert gRPC path to Protobuf Full Name
	// From: "/tenant.v1.OrganizationService/CreateOrganization"
	// To:   "tenant.v1.OrganizationService.CreateOrganization"
	methodName := strings.TrimPrefix(fullMethod, "/")
	methodName = strings.ReplaceAll(methodName, "/", ".")

	// Look up the method in the global Protobuf registry
	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(methodName))
	if err != nil {
		routeAuthCache.Store(fullMethod, config)
		return config
	}

	// Ensure it's actually a method descriptor
	md, ok := desc.(protoreflect.MethodDescriptor)
	if !ok {
		routeAuthCache.Store(fullMethod, config)
		return config
	}

	// Extract the options and check for our custom extensions
	opts := md.Options().(*descriptorpb.MethodOptions)

	if opts != nil {
		// Check for skip_org_check (Bool)
		if proto.HasExtension(opts, tenantV1.E_SkipOrgCheck) {
			config.SkipOrgCheck = proto.GetExtension(opts, tenantV1.E_SkipOrgCheck).(bool)
		}

		// Check for required_org_permission (String)
		if proto.HasExtension(opts, tenantV1.E_RequiredOrgPermission) {
			config.RequiredOrgPermission = proto.GetExtension(opts, tenantV1.E_RequiredOrgPermission).(string)
		}
	}

	// Save the result in the cache for all future requests
	routeAuthCache.Store(fullMethod, config)

	return config
}
