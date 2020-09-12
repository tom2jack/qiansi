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
	// 列举所有文件。
	marker := ""
	data := make([]string, 0)
	for {
		// oss.Prefix("") 可指定前缀
		lsRes, err := o.b.ListObjects(oss.Marker(marker))
		if err != nil {
			return nil, err
		}
		// 打印列举文件，默认情况下一次返回100条记录。
		for _, obj := range lsRes.Objects {
			data = append(data, obj.Key)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return data, nil
}
