package storage

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSProvider 阿里云 OSS 存储（生产环境）
type OSSProvider struct {
	client   *oss.Client
	bucket   *oss.Bucket
	endpoint string
	bucketNm string
}

func NewOSSProvider(endpoint, accessKeyID, accessKeySecret, bucketName string) (*OSSProvider, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建 OSS 客户端失败: %w", err)
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("获取 OSS bucket 失败: %w", err)
	}

	return &OSSProvider{
		client:   client,
		bucket:   bucket,
		endpoint: endpoint,
		bucketNm: bucketName,
	}, nil
}

func (p *OSSProvider) Upload(data []byte, filename string) (string, error) {
	if err := p.bucket.PutObject(filename, bytes.NewReader(data)); err != nil {
		return "", fmt.Errorf("上传 OSS 失败: %w", err)
	}

	// 拼接公网访问 URL: https://{bucket}.{endpoint}/{filename}
	host := strings.TrimPrefix(p.endpoint, "https://")
	host = strings.TrimPrefix(host, "http://")
	return fmt.Sprintf("https://%s.%s/%s", p.bucketNm, host, filename), nil
}

func (p *OSSProvider) Delete(url string) error {
	// 从 URL 提取 object key
	idx := strings.Index(url, ".com/")
	if idx == -1 {
		return fmt.Errorf("无效的 OSS URL")
	}
	objectKey := url[idx+5:]

	if err := p.bucket.DeleteObject(objectKey); err != nil {
		return fmt.Errorf("删除 OSS 对象失败: %w", err)
	}
	return nil
}
