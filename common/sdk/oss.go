package sdk

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zhi-miao/qiansi/common/config"
)

type ossClient struct {
	c *oss.Client
	b *oss.Bucket
}

// NewOSSClient 初始化一个OSS客户端
func NewOSSClient() (*ossClient, error) {
	conf := config.GetConfig().Aliyun
	// 创建OSSClient实例。
	client, err := oss.New(conf.OSS.Endpoint, conf.AccessKey, conf.AccessSecret)
	if err != nil {
		return nil, err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(conf.OSS.BucketName)
	if err != nil {
		return nil, err
	}
	return &ossClient{c: client, b: bucket}, err
}

// ListFile 列文件
func (o *ossClient) ListFile(path string) ([]string, error) {
	// 列举包含指定前缀的文件。默认列举100个文件。
	// lsRes, err := o.b.ListObjects(oss.Prefix("my-object-"))
	lsRes, err := o.b.ListObjects(oss.Prefix(""))
	if err != nil {
		return nil, err
	}
	data := make([]string, 0)
	for _, obj := range lsRes.Objects {
		data = append(data, obj.Key)
	}
	return data, nil
}
