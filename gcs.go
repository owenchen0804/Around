package main

import (
    "context"
    "fmt"
    "io"

    "cloud.google.com/go/storage"
)

const (
    BUCKET_NAME = "owen-chen-bucket"
)

func saveToGCS(r io.Reader, objectName string) (string, error) { // reader是发过来的文件
    ctx := context.Background()     // 提供空白context(运行HTTP的参数等info），.Background()是不specify anything

    client, err := storage.NewClient(ctx)   // GCS提供的NewClient
    if err != nil {
        return "", err
    }

object := client.Bucket(BUCKET_NAME).Object(objectName)     // 创建bucket as a placeholder给一个名字
    wc := object.NewWriter(ctx)     // new a NewWriter
    if _, err := io.Copy(wc, r); err != nil {   // Use io package copy (destination, 要写的文件)
        // it is a blocking call that may take a few seconds
        // 不需要结果所以是_
        return "", err
    }

    if err := wc.Close(); err != nil { // 没问题就close writer
        return "", err
    }

    if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
        // ACL access control lists, 给object设定权限. 把delete权限给server，交给代码，提出delete要求
        // 的可以是上传的人
        return "", err
    }

    attrs, err := object.Attrs(ctx)
    if err != nil {
        return "", err
    }

    fmt.Printf("Image is saved to GCS: %s\n", attrs.MediaLink)
    return attrs.MediaLink, nil     // 成功了返回MediaLink as String
} 