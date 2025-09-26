package main

import (
	"context"

	"github.com/sagikazarmark/dagx/pipeline"

	"github.com/sagikazarmark/caddy-fs-s3/.dagger/internal/dagger"
)

type CaddyFsS3 struct {
	// Project source directory
	//
	// +private
	Source *dagger.Directory
}

func New(
	// Project source directory.
	//
	// +defaultPath="/"
	// +ignore=[".devenv", ".direnv", ".github"]
	source *dagger.Directory,
) *CaddyFsS3 {
	return &CaddyFsS3{
		Source: source,
	}
}

var supportedGoVersions = []string{"1.24", "1.25"}

func (m *CaddyFsS3) Check(ctx context.Context) error {
	p := pipeline.New(ctx)

	for _, goVersion := range supportedGoVersions {
		pipeline.AddSyncStep(p, m.Build(goVersion))
	}

	pipeline.AddSyncStep(p, m.Test())
	pipeline.AddSyncStep(p, m.Lint())

	return pipeline.Run(p)
}

func (m *CaddyFsS3) Build(
	// Go version to use.
	//
	// +optional
	goVersion string,
) *dagger.Container {
	if goVersion == "" {
		goVersion = defaultGoVersion
	}

	return dag.Xcaddy(dagger.XcaddyOpts{
		GoVersion: goVersion,
	}).Build().
		Plugin("github.com/sagikazarmark/caddy-fs-s3", dagger.XcaddyBuildPluginOpts{Replacement: m.Source}).
		Container()
}

func (m *CaddyFsS3) Test() *dagger.Container {
	return dag.Go(dagger.GoOpts{
		Version: defaultGoVersion,
	}).
		WithSource(m.Source).
		Exec([]string{"go", "test", "-race", "-v", "./..."})
}

// TODO: add e2e test
func (m *CaddyFsS3) Run() *dagger.Service {
	return m.Build("").
		WithEnvVariable("AWS_ACCESS_KEY_ID", "admin").
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", "admin1234").
		WithServiceBinding("minio", minio()).
		WithMountedFile("/etc/caddy/Caddyfile", m.Source.File("Caddyfile")).
		WithExposedPort(80).
		AsService()
}

func minio() *dagger.Service {
	return dag.Container().
		From("quay.io/minio/minio").
		WithEnvVariable("MINIO_ROOT_USER", "admin").
		WithEnvVariable("MINIO_ROOT_PASSWORD", "admin1234").
		WithExposedPort(9000).
		WithExposedPort(9090).
		AsService(dagger.ContainerAsServiceOpts{
			Args:          []string{"server", "/data", "--console-address", ":9090"},
			UseEntrypoint: true,
		})
}

func (m *CaddyFsS3) Lint() *dagger.Container {
	return dag.GolangciLint(dagger.GolangciLintOpts{
		Version:   golangciLintVersion,
		GoVersion: defaultGoVersion,
	}).
		Run(m.Source, dagger.GolangciLintRunOpts{Verbose: true})
}
