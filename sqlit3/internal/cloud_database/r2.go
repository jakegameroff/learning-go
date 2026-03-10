import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type R2Client struct {
	client *s3.Client
	bucket string
}

func NewR2Client(bucket string) (*R2Client, error) {
	var accountID = os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	var accessKeyID = os.Getenv("R2_ACCESS_KEY_ID")
	var accessKeySecret = os.Getenv("R2_ACCESS_KEY_SECRET")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
	})

	return &R2Client{
		client: client,
		bucket: bucket,
	}, nil
}

func (c *R2Client) Download(key string) (io.ReadCloser, error) {}