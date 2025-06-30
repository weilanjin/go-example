package repo

import (
	"context"
	"fmt"
)

type Repos[T any] interface {
	Create(context.Context, *T) error
	Update(context.Context, *T) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*T, error)
	Select(ctx context.Context, query string, args ...any) ([]*T, error)
}

type UserGroupAggregate struct {
	Users   Repos[User]
	Groups  Repos[UserGroup]
	Members Repos[UserGroupMember]
	Rules   Repos[UserGroupRule]
}

func (r *UserGroupAggregate) IsMember(ctx context.Context, groupID, userID string) (bool, error) {
	g, err := r.Groups.Get(ctx, groupID)
	if err != nil {
		return false, fmt.Errorf("r.Groups.Get %w", err)
	}
	_ = g
	u, err := r.Users.Get(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("r.Users.Get %w", err)
	}
	_ = u
	m, err := r.GetMembers(ctx, groupID, userID)
	if err != nil {
		return false, fmt.Errorf("r.GetMembers %w", err)
	}
	return len(m) > 0, nil
}

func (r *UserGroupAggregate) GetMembers(ctx context.Context, groupId, userId string) ([]*UserGroupMember, error) {
	return r.Members.Select(ctx, "group_id = ? AND user_id = ?", groupId, userId)
}
