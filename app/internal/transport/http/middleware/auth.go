package middleware

import (
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type AuthContext struct {
	Permissions []string `json:"permissions"`
	Role        string   `json:"role"`
	RequestorID string   `json:"requestorID"`
}

var (
	ErrNotAuthorized = errors.New("not authorized to access resource")
)

func GetAuthContext(event events.APIGatewayV2HTTPRequest) (*AuthContext, error) {
	requestorID, ok := event.RequestContext.Authorizer.Lambda["requestorID"].(string)
	if !ok {
		return nil, errors.New("err no sub")
	}

	role, ok := event.RequestContext.Authorizer.Lambda["role"].(string)
	if !ok {
		return nil, errors.New("err no role")
	}

	permissions, ok := event.RequestContext.Authorizer.Lambda["permissions"].([]interface{})
	if !ok {
		return nil, errors.New("err no permissions")
	}

	typedPermissions := make([]string, len(permissions))
	for i, p := range permissions {
		str, ok := p.(string)
		if !ok {
			return nil, errors.New("err invalid permission type")
		}
		typedPermissions[i] = str
	}

	return &AuthContext{
		Permissions: typedPermissions,
		Role:        role,
		RequestorID: requestorID,
	}, nil
}

func HasPermission(permissions []string, requestorID, resourceID, requiredResourceType, requiredResourceAction string) bool {
	for _, permission := range permissions {
		permissionParts := strings.Split(permission, ":")

		if len(permissionParts) != 3 {
			// log and skip invalid permissions
			continue
		}

		scope, resourceType, resourceAction := permissionParts[0], permissionParts[1], permissionParts[2]

		if resourceType != requiredResourceType || resourceAction != requiredResourceAction {
			continue
		}

		switch scope {
		case "all":
			return true
		case "own":
			return requestorID == resourceID
		}
	}

	return false
}
