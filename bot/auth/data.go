package auth

import "fmt"

var Ids = make(map[int]string)

func InitUsers(getUsersFunc func() map[int]string) error {
	users := getUsersFunc()
	if users == nil {
		return fmt.Errorf("failed to initialize users: user map is nil")
	}
	
	Ids = users
	return nil
}
