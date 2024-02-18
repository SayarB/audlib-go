package utils

import (
	"fmt"
	"io"
	"os"

	storage_go "github.com/supabase-community/storage-go"
)

type AudioFile struct{
	BucketId string
	Folder string
	Key string
	Extension string
	File io.Reader
	MIMEType string
}

func UploadToSupabase(file *AudioFile) string {

	storageClient := storage_go.NewClient(fmt.Sprintf("https://%s.supabase.co/storage/v1", os.Getenv("SUPABASE_PROJECT_REFERENCE_ID")), os.Getenv("SUPABASE_API_SECRET_KEY"), nil)

	

	upsert:=true
	contentType:=file.MIMEType

	_, err := storageClient.UploadFile(file.BucketId, fmt.Sprintf("%s/%s%s",file.Folder,file.Key, file.Extension), file.File, storage_go.FileOptions{
		Upsert: &upsert,
		ContentType:&contentType,
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fmt.Sprintf("%s%s",file.Key, file.Extension))

	return fmt.Sprintf("%s%s",file.Key, file.Extension)
}