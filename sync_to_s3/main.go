package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

// 硬编码的 AWS 配置
const (
    awsRegion  = "ap-east-1"           // AWS 区域
    bucketName = "" // S3 Bucket 名称
    accessKey  = ""  // 替换为你的 AK
    secretKey  = ""     // 替换为你的 SK
)

// 上传文件到 S3
func uploadFileToS3(sess *session.Session, bucketName, filePath, keyName string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("无法打开文件 %q, %v", filePath, err)
    }
    defer file.Close()

    svc := s3.New(sess)
    _, err = svc.PutObject(&s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(keyName),
        Body:   file,
    })
    if err != nil {
        return fmt.Errorf("无法上传文件 %q 到 S3 bucket %q, %v", keyName, bucketName, err)
    }

    fmt.Printf("文件 %q 已成功上传至 S3 bucket %q\n", keyName, bucketName)
    return nil
}

func main() {
    // 定义命令行标志，用于传递文件路径
    filePaths := flag.String("files", "", "要上传到 S3 的文件路径，多个文件用逗号分隔")
    flag.Parse()

    if *filePaths == "" {
        log.Fatalf("必须提供至少一个文件路径，请使用 -files 参数")
    }

    // 将传递的文件路径分隔成多个文件
    files := filepath.SplitList(*filePaths)

    // 创建 AWS session，使用显式的 AK/SK
    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String(awsRegion),
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
    })
    if err != nil {
        log.Fatalf("无法创建 AWS session, %v", err)
    }

    // 遍历所有文件并上传到 S3
    for _, filePath := range files {
        keyName := filepath.Base(filePath)
        if err := uploadFileToS3(sess, bucketName, filePath, keyName); err != nil {
            log.Printf("上传文件 %q 失败, %v", filePath, err)
        }
    }
}
