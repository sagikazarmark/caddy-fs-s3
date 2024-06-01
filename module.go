package caddyfss3

import (
	"errors"
	"io/fs"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/jszwec/s3fs/v2"
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

	// Set this to `true` to enable the client to use path-style addressing.
	UsePathStyle bool `json:"use_path_style,omitempty"`

	// DEPRECATED: please use 'use_path_style' instead.
	// Set this to `true` to force the request to use path-style addressing.
	S3ForcePathStyle bool `json:"force_path_style,omitempty"`
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

	if fs.S3ForcePathStyle {
		ctx.Logger().Warn("force_path_style is deprecated, please use use_path_style instead")
	}

	var configOpts []func(*config.LoadOptions) error

	if fs.Region != "" {
		configOpts = append(configOpts, config.WithRegion(fs.Region))
	}

	if fs.Profile != "" {
		configOpts = append(configOpts, config.WithSharedConfigProfile(fs.Profile))
	}

	cfg, err := config.LoadDefaultConfig(ctx.Context, configOpts...)
	if err != nil {
		ctx.Slogger().Error("could not create AWS config", slog.String("error", err.Error()))

		return err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if fs.Endpoint != "" {
			o.BaseEndpoint = aws.String(fs.Endpoint)
		}

		o.UsePathStyle = fs.UsePathStyle || fs.S3ForcePathStyle
	})

	// ReadSeeker is required by Caddy
	fs.StatFS = s3fs.New(client, fs.Bucket, s3fs.WithReadSeeker)

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
		case "use_path_style":
			fs.UsePathStyle = true
		case "force_path_style":
			fs.S3ForcePathStyle = true
		default:
			return d.Errf("%s not a valid caddy.fs.s3 option", d.Val())
		}
	}

	return nil
}
