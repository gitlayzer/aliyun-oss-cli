package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

var (
	endpoint     = os.Getenv("OSS_ENDPOINT")    // OSS endpoint
	accessKey    = os.Getenv("OSS_ACCESS_KEY")  // OSS accessKey
	secretSecret = os.Getenv("OSS_SECRET_KEY")  // OSS secretKey
	bucketName   = os.Getenv("OSS_BUCKET_NAME") // OSS bucketName
	uploadFile   = ""
)

func loadParams() {
	flag.StringVar(&uploadFile, "f", "", "set oss upload file")
	flag.Parse()
}

func validate() error {
	if endpoint == "" || accessKey == "" || secretSecret == "" {
		return fmt.Errorf("endpoint, accessKey, secretSecret must be set")
	}
	if uploadFile == "" {
		return fmt.Errorf("upload file path required")
	}
	return nil
}

func upload(filepath string) error {
	client, err := oss.New(endpoint, accessKey, secretSecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	return bucket.PutObjectFromFile(uploadFile, filepath)
}

func main() {
	loadParams()
	if err := validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := upload(uploadFile); err != nil {
		fmt.Println("upload file failed, err:", err)
		os.Exit(1)
	} else {
		fmt.Printf("%s Upload Success！！！\n", uploadFile)
		fmt.Printf("Download URL: %s\n", "https"+"://"+bucketName+"."+endpoint+"/"+uploadFile)
	}
}
