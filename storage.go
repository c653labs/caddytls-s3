package caddytlss3

import (
	"net/url"

	"github.com/mholt/caddy/caddytls"
)

func init() {
	caddytls.RegisterStorageProvider("s3", NewS3Storage)
}

type S3Storage struct {
}

func NewS3Storage(caURL *url.URL) (caddytls.Storage, error) {

	return &S3Storage{}, nil
}

func (s *S3Storage) DeleteSite(domain string) error {
	return nil
}

func (s *S3Storage) LoadSite(domain string) (*caddytls.SiteData, error) {
	return nil, nil
}

func (s *S3Storage) LoadUser(email string) (*caddytls.UserData, error) {
	return nil, nil
}

func (s *S3Storage) MostRecentUserEmail() string {
	return ""
}

func (s *S3Storage) SiteExists(domain string) (bool, error) {
	return true, nil
}

func (s *S3Storage) StoreSite(domain string, data *caddytls.SiteData) error {
	return nil
}

func (s *S3Storage) StoreUser(email string, data *caddytls.SiteUser) error {
	return nil
}
