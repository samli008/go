package model

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

var end = "oss-cn-chengdu.aliyuncs.com"
var id = "LTAI5tAYEnLLsapWhEKdAUeH"
var Secret = "6McLsMtl1MdTAay9KfJibq9mQ0g0ZK"

type ob struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

func CreateBk(name string) {
	endpoint := end
	accessKeyId := id
	accessKeySecret := Secret
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		handleError(err)
	}
	err = client.CreateBucket(name)
	if err != nil {
		handleError(err)
	}
}

type OssProgressListener struct {
}

// 定义进度变更事件处理函数。
func (listener *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		fmt.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferDataEvent:
		fmt.Printf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
			event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
	case oss.TransferCompletedEvent:
		fmt.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferFailedEvent:
		fmt.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}

func UpOb(bucketName string, object string, fileByte []byte) {
	endpoint := end
	accessKeyId := id
	accessKeySecret := Secret
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		handleError(err)
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	err = bucket.PutObject(object, bytes.NewReader([]byte(fileByte)), oss.Progress(&OssProgressListener{}))
	if err != nil {
		handleError(err)
	}
}

func ListOb(bucketName string) []ob {
	endpoint := end
	accessKeyId := id
	accessKeySecret := Secret
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		handleError(err)
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	marker := ""
	lists := make([]ob, 0)
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			handleError(err)
		}

		for _, object := range lsRes.Objects {
			list := ob{}
			list.Name = object.Key
			list.Size = FormatFileSize(object.Size)
			lists = append(lists, list)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return lists
}

func DeleteOb(bucketName string, objectName string) error {
	endpoint := end
	accessKeyId := id
	accessKeySecret := Secret
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}
	err = bucket.DeleteObject(objectName)
	fmt.Println(objectName)
	if err != nil {
		return err
	}
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
	endpoint := end
	accessKeyId := id
	accessKeySecret := Secret
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		handleError(err)
	}
	// 获取存储空间
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	// 下载文件
	err = bucket.GetObjectToFile(objectName, "./files/"+objectName)
	if err != nil {
		handleError(err)
	}
}
