package sync

import "testing"

type User struct {
	Name string
}

func TestLockable(t *testing.T) {
	// 1 directly use
	var safeUser Lockable[User]
	safeUser.Set(User{Name: "test"})

	// 2 new lockable type
	type LockableUser Lockable[User]

}
