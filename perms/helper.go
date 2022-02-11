package perms

import (
	"context"

	"github.com/contextgg/pkg/identity"
)

func CheckCurrent(ctx context.Context, manager Manager, relationship string, object string) error {
	user, ok := identity.FromContext(ctx)
	if !ok || user == nil {
		return manager.Check("Anonymous", relationship, object)
	}

	for _, role := range user.Roles {
		// check if the role is ok!
		if err := manager.Check("Role:"+role, relationship, object); err == nil {
			return nil
		}
	}

	for _, identity := range user.Identities {
		if err := manager.Check("Identity:"+identity.Connection+"@"+identity.UserId, relationship, object); err == nil {
			return nil
		}
	}

	return ErrCheckFailed
}
