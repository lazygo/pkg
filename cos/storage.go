package cos

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var client func() *cos.Client

var ErrUninitialized = errors.New("cos 未初始化")

func Init(conf Config) error {
	u, err := url.Parse(conf.BucketURL)
	if err != nil {
		return err
	}
	b := &cos.BaseURL{BucketURL: u}
	hc := &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.SecretID,
			SecretKey: conf.SecretKey,
		},
	}
	client = func() *cos.Client {
		return cos.NewClient(b, hc)
	}
	return nil
}

func Upload(name string, r io.Reader) error {
	if client == nil {
		return ErrUninitialized
	}

	_, err := client().Object.Put(context.Background(), "upload/"+name, r, nil)
	if err != nil {
		return err
	}

	return nil
}

func Download(name string) (io.ReadCloser, error) {
	if client == nil {
		return nil, ErrUninitialized
	}

	resp, err := client().Object.Get(context.Background(), name, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
