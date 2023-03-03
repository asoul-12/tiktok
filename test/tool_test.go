package test

import (
	"fmt"
	"log"
	"testing"
	"tiktok/tools"
)

func TestSnowFlake(t *testing.T) {
	flake := tools.SnowFlake{}
	println(flake.GenSnowID())
}
func TestToken(t *testing.T) {
	token, err := tools.GenerateToken("asoul")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)
	registeredClaims, err := tools.ParseToken(token)
	fmt.Println(registeredClaims.Audience[0])
}
func TestAESScode(t *testing.T) {
	aes, _ := tools.EncryptByAes([]byte("asoul"))
	fmt.Printf("%v\n", aes)
	decryptByAes, _ := tools.DecryptByAes(aes)
	fmt.Printf("%s", decryptByAes)
}

func TestIdenticon(t *testing.T) {
	tools.GenerateAvatar()
}
