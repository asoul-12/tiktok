package test

import (
	"fmt"
	"testing"
	"tiktok/repo"
	"tiktok/serverInit"
)

func TestUser(t *testing.T) {
	serverInit.InitDatabase()
	userRepo := repo.UserRepo{}
	name := userRepo.FindUserByUserName("123456")
	fmt.Println(name)
}
