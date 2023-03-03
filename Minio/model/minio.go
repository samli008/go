package model

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	end    = "xxx:9000"
	id     = "admin"
	Secret = "liyang@008"
	ssl    = false
)

type ob struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

func CreateBk(name string) {
	ctx := context.Background()
	endpoint := end
	accessKeyID := id
	secretAccessKey := Secret
	useSSL := ssl
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = client.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := client.BucketExists(ctx, name)
		if err == nil && exists {
			log.Printf("We already own %s\n", name)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", name)
	}
}

func UpOb(bucketName string, object string, file *multipart.FileHeader, contentType string) {
	ctx := context.Background()
	endpoint := end
	accessKeyID := id
	secretAccessKey := Secret
	useSSL := ssl
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	src, err1 := file.Open()
	if err1 != nil {
		log.Fatalln(err)
	}
	defer src.Close()

	info, err := client.PutObject(ctx, bucketName, object, src, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", object, info.Size)
}

func ListOb(bucketName string) []ob {
	ctx := context.Background()
	endpoint := end
	accessKeyID := id
	secretAccessKey := Secret
	useSSL := ssl
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	opts := minio.ListObjectsOptions{
		UseV1:     true,
		Prefix:    "",
		Recursive: false,
	}

	lists := make([]ob, 0)
	// List all objects from a bucket-name with a matching prefix.
	for object := range client.ListObjects(ctx, bucketName, opts) {
		if object.Err != nil {
			fmt.Println(object.Err)
		}
		list := ob{}
		list.Name = object.Key
		list.Size = FormatFileSize(object.Size)
		lists = append(lists, list)
	}
	return lists
}

func DeleteOb(bucketName string, objectName string) error {
	ctx := context.Background()
	endpoint := end
	accessKeyID := id
	secretAccessKey := Secret
	useSSL := ssl
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}

	err = client.RemoveObject(ctx, bucketName, objectName, opts)
	return err
}

// 字节的单位转换 保留两位小数
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fPB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

func DownOb(bucketName string, objectName string) {
	ctx := context.Background()
	endpoint := end
	accessKeyID := id
	secretAccessKey := Secret
	useSSL := ssl
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	reader, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()

	localFile, err := os.Create(objectName)
	if err != nil {
		log.Fatalln(err)
	}
	defer localFile.Close()

	stat, err := reader.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := io.CopyN(localFile, reader, stat.Size); err != nil {
		log.Fatalln(err)
	}
}
