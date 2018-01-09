package caddytlss3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mholt/caddy/caddytls"
)

func init() {
	caddytls.RegisterStorageProvider("s3", NewS3Storage)
}

type S3Storage struct {
	client *s3.S3
	bucket *string
	prefix string
}

func NewS3Storage(caURL *url.URL) (caddytls.Storage, error) {
	b := os.Getenv("CADDY_S3_BUCKET")
	if b == "" {
		return nil, fmt.Errorf("CADDY_S3_BUCKET environment variable is not set")
	}

	cfg := aws.NewConfig().WithCredentialsChainVerboseErrors(true)

	prefix := caURL.Hostname()
	if p := os.Getenv("CADDY_S3_PREFIX"); p != "" {
		prefix = path.Join(prefix, p)
	}

	if r := os.Getenv("CADDY_S3_REGION"); r != "" {
		cfg = cfg.WithRegion(r)
	}

	return &S3Storage{
		client: s3.New(session.Must(session.NewSession()), cfg),
		bucket: aws.String(b),
		prefix: prefix,
	}, nil
}

func (s *S3Storage) siteKey(domain string) *string {
	return aws.String(path.Join(s.prefix, "sites", domain))
}

func (s *S3Storage) userKey(email string) *string {
	return aws.String(path.Join(s.prefix, "users", email))
}

func (s *S3Storage) DeleteSite(domain string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.bucket,
		Key:    s.siteKey(domain),
	})
	return err
}

func (s *S3Storage) LoadSite(domain string) (*caddytls.SiteData, error) {
	res, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: s.bucket,
		Key:    s.siteKey(domain),
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to load site: %v", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read site data: %v", err)
	}

	data := &caddytls.SiteData{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, fmt.Errorf("Cannot unmarshal site data: %v", err)
	}

	return data, nil
}

func (s *S3Storage) LoadUser(email string) (*caddytls.UserData, error) {
	res, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: s.bucket,
		Key:    s.userKey(email),
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch user: %v", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read site data: %v", err)
	}

	data := &caddytls.UserData{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, fmt.Errorf("Cannot unmarshal user data: %v", err)
	}

	return data, nil
}

func (s *S3Storage) MostRecentUserEmail() string {
	return ""
}

func (s *S3Storage) SiteExists(domain string) (bool, error) {
	_, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: s.bucket,
		Key:    s.siteKey(domain),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *S3Storage) StoreSite(domain string, data *caddytls.SiteData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Unable to marshal: %v", err)
	}

	_, err = s.client.PutObject(&s3.PutObjectInput{
		Body:                 aws.ReadSeekCloser(bytes.NewReader(b)),
		Bucket:               s.bucket,
		Key:                  s.siteKey(domain),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return fmt.Errorf("Unable to store site data: %v", err)
	}

	return nil
}

func (s *S3Storage) StoreUser(email string, data *caddytls.UserData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Unable to marshal: %v", err)
	}

	_, err = s.client.PutObject(&s3.PutObjectInput{
		Body:                 aws.ReadSeekCloser(bytes.NewReader(b)),
		Bucket:               s.bucket,
		Key:                  s.userKey(email),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return fmt.Errorf("Unable to store user data: %v", err)
	}

	return nil
}
