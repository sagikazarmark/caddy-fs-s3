# Caddy FS module for AWS S3

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/sagikazarmark/caddy-fs-s3/ci.yaml?style=flat-square)
![Caddy Version](https://img.shields.io/badge/caddy%20version-%3E=2.10.x-61CFDD.svg?style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sagikazarmark/caddy-fs-s3?style=flat-square&color=61CFDD)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/sagikazarmark/caddy-fs-s3/badge?style=flat-square)](https://deps.dev/go/github.com%252Fsagikazarmark%252Fcaddy-fs-s3)

## Installation

Build Caddy using [xcaddy](https://github.com/caddyserver/xcaddy):

```shell
xcaddy build --with github.com/sagikazarmark/caddy-fs-s3
```

## Usage

```caddyfile
{
	filesystem my-s3-fs s3 {
		bucket mybucket
		region us-east-1

		# endpoint <endpoint>
		# profile <profile>
		# use_path_style
	}
}

example.com {
    file_server {
        fs my-s3-fs
    }
}
```

> [!NOTE]
> For a full parameter reference, check out the module [documentation page](https://caddyserver.com/docs/modules/caddy.fs.s3).

### Authentication

The module uses the AWS SDK [default credential chain](https://docs.aws.amazon.com/sdkref/latest/guide/standardized-credentials.html) to find valid credentials.

The easiest way to try the module is setting [static credentials](https://docs.aws.amazon.com/sdkref/latest/guide/feature-static-credentials.html) either in your AWS credentials file or as environment variables:

```shell
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
```

Caddy will pick up the credentials automatically.

## Development

**For an optimal developer experience, it is recommended to install [Nix](https://nixos.org/download.html) and [direnv](https://direnv.net/docs/installation.html).**

Run Caddy with the following command:

```shell
just run up --ports 8080:80
```

When all coding and testing is done, please run the test suite:

```shell
just check
```

## License

The project is licensed under the [MIT License](LICENSE).
