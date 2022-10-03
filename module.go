package caddyfss3

import (
	"errors"
	"io/fs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/jszwec/s3fs"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(FS{})
}

// Interface guards
var (
	_ fs.StatFS             = (*FS)(nil)
	_ caddyfile.Unmarshaler = (*FS)(nil)
)

// FS is a Caddy virtual filesystem module for AWS S3 (and compatible) object store.
type FS struct {
	fs.StatFS `json:"-"`

	// The name of the S3 bucket.
	Bucket string `json:"bucket,omitempty"`

	// The AWS region the bucket is hosted in.
	Region string `json:"region,omitempty"`

	// The AWS profile to use if mulitple profiles are specified.
	Profile string `json:"profile,omitempty"`

	// Use non-standard endpoint for S3.
	Endpoint string `json:"endpoint,omitempty"`

	// Set this to `true` to force the request to use path-style addressing.
	S3ForcePathStyle bool `json:"force_path_style,omitempty"`

	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (FS) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.fs.s3",
		New: func() caddy.Module { return new(FS) },
	}
}

func (fs *FS) Provision(ctx caddy.Context) error {
	if fs.Bucket == "" {
		return errors.New("bucket must be set")
	}

	var config aws.Config

	if fs.Region != "" {
		config.Region = aws.String(fs.Region)
	}

	if fs.Endpoint != "" {
		config.Endpoint = aws.String(fs.Endpoint)
	}

	if fs.S3ForcePathStyle {
		config.S3ForcePathStyle = aws.Bool(fs.S3ForcePathStyle)
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  config,
		Profile: fs.Profile,
	})
	if err != nil {
		fs.logger.Error("could not create AWS session", zap.Error(err))
		return err
	}

	// ReadSeeker is required by Caddy
	fs.StatFS = s3fs.New(s3.New(sess), fs.Bucket, s3fs.WithReadSeeker)

	return nil
}

// UnmarshalCaddyfile unmarshals a caddyfile.
func (fs *FS) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if !d.Next() { // skip block beginning
		return d.ArgErr()
	}

	for nesting := d.Nesting(); d.NextBlock(nesting); {
		switch d.Val() {
		case "bucket":
			if !d.AllArgs(&fs.Bucket) {
				return d.ArgErr()
			}
		case "region":
			if !d.AllArgs(&fs.Region) {
				return d.ArgErr()
			}
		case "profile":
			if !d.AllArgs(&fs.Profile) {
				return d.ArgErr()
			}
		case "endpoint":
			if !d.AllArgs(&fs.Endpoint) {
				return d.ArgErr()
			}
		case "force_path_style":
			fs.S3ForcePathStyle = true
		default:
			return d.Errf("%s not a valid caddy.fs.s3 option", d.Val())
		}
	}

	return nil
}
