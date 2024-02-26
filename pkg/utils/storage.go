package utils

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
)

// func UploadToSupabase(file *models.AudioFile) string {

// 	storageClient := storage_go.NewClient(fmt.Sprintf("https://%s.supabase.co/storage/v1", os.Getenv("SUPABASE_PROJECT_REFERENCE_ID")), os.Getenv("SUPABASE_API_SECRET_KEY"), nil)

// 	upsert:=true
// 	contentType:=file.MIMEType

// 	_, err := storageClient.UploadFile(file.BucketId, fmt.Sprintf("%s/%s%s",file.Folder,file.Key, file.Extension), file.File, storage_go.FileOptions{
// 		Upsert: &upsert,
// 		ContentType:&contentType,
// 	})

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(fmt.Sprintf("%s%s",file.Key, file.Extension))

// 	return fmt.Sprintf("%s%s",file.Key, file.Extension)
// }

// func DownloadFileFromSupabase(file *models.AudioFile) []byte{
// 	storageClient := storage_go.NewClient(fmt.Sprintf("https://%s.supabase.co/storage/v1", os.Getenv("SUPABASE_PROJECT_REFERENCE_ID")), os.Getenv("SUPABASE_API_SECRET_KEY"), nil)

// 	result, err := storageClient.DownloadFile(file.BucketId, fmt.Sprintf("%s/%s%s",file.Folder,file.Key, file.Extension))

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	return result
// }






func UploadToS3(file *models.AudioFile) (string, error) {


	fileName:=file.Key + file.Extension
	fmt.Println(fileName)
	_, err := config.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(file.BucketId),
		Key:    aws.String(fileName),
		Body:   file.File,
	})
	if err != nil {
		log.Printf("Couldn't upload file to %v:%v. Here's why: %v\n", file.BucketId, file.Key, err)
	}
	return file.Key + file.Extension, err
}

type S3DownloadInput struct{
	Key string
	BucketId string
	Extension string
}

func DownloadFileFromS3(file *S3DownloadInput, chunkStart int, chunkEnd int ) (string, error) {

	fileName:=file.Key + file.Extension

	result, err := config.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(file.BucketId),
		Key:    aws.String(fileName),
		Range: aws.String(fmt.Sprintf("bytes=%d-%d", chunkStart, chunkEnd)), // 1KB
	})


	if err != nil {
		log.Printf("Couldn't download file from %v:%v. Here's why: %v\n", file.BucketId, file.Key, err)
	}

	defer result.Body.Close()

	
	res, err:=io.ReadAll(result.Body)
	
	return string(res), err
}