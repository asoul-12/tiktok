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

func TestAES(t *testing.T) {
	aes, _ := tools.EncryptByAes([]byte("asoul"))
	fmt.Printf("%v\n", aes)
	decryptByAes, _ := tools.DecryptByAes(aes)
	fmt.Printf("%s", decryptByAes)
}

func TestFFmpeg(t *testing.T) {
	var ffmpeg = &tools.Bind{
		FFMpeg:         "D:\\tools\\ffmpeg\\bin\\ffmpeg.exe",
		FFProbe:        "D:\\tools\\ffmpeg\\bin\\ffprobe.exe",
		CommandTimeout: 5000,
	}
	err := ffmpeg.Thumbnail("../assets/bear.mp4", "../assets/testdata/bear.jpg", 0, true)
	if err != nil {
		fmt.Println(err)
	}
	//_, err = ffmpeg.Transcoding("../assets/bear.mp4", "../assets/testdata/bear.mp4", true)
	//if err != nil {
	//	fmt.Println( err)
	//}
}
