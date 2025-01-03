# Caddy FS module for AWS S3

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/sagikazarmark/caddy-fs-s3/ci.yaml?style=flat-square)
![Caddy Version](https://img.shields.io/badge/caddy%20version-%3E=2.8.x-61CFDD.svg?style=flat-square)

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

Run Caddy with the following command:

```shell
task run
```

When all coding and testing is done, please run the test suite:

```shell
task check
```

For the best developer experience, install [Nix](https://builtwithnix.org/) and [direnv](https://direnv.net/).

Alternatively, install Go, xcaddy and the rest of the dependencies manually or using a package manager.

## License

The project is licensed under the [MIT License](LICENSE).
