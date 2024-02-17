package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sayar/go-streaming/pkg/utils"
)

func main(){
	godotenv.Load(".env")
	file,err:=os.Open("test.jpg")
	if err != nil {
		panic(err)
	}
	key:=utils.UploadToSupabase(&utils.AudioFile{BucketId: "audio", Key: "test",Folder:"test", Extension: "jpg", File:file })
	fmt.Print(key)
}

