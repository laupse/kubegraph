package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"dagger.io/dagger"
	"github.com/dagger/dagger/engine"
	"github.com/dagger/dagger/router"
	"github.com/spf13/pflag"
	viper "github.com/spf13/viper"
)

func initDagger(ctx context.Context) (*dagger.Client, *dagger.Directory, error) {
	pflag.String("registry-url", "proxy-registry:5002", "")
	pflag.String("image-name", "kubegraph", "")
	pflag.String("image-tag", "latest", "")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.SetEnvPrefix("CI") // will be uppercased automatically
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	client, err := dagger.Connect(
		ctx,
		dagger.WithLogOutput(os.Stdout),
		dagger.WithConfigPath("dagger.json"),
	)
	if err != nil {
		return nil, nil, err
	}

	workdir := client.Host().Workdir()

	return client, workdir, nil
}

func innerBuild(
	ctx context.Context,
	client *dagger.Client,
	workdir *dagger.Directory,
) (*dagger.Container, error) {

	build := client.
		Container().
		Build(workdir)
	build.Publish(ctx, "")

	return build, nil
}

func Dev() error {
	err := engine.Start(
		context.Background(),
		&engine.Config{
			LogOutput: os.Stdout,
		},
		func(ctx context.Context, r *router.Router) error {
			srv := http.Server{
				Addr:              fmt.Sprintf(":%d", 8080),
				Handler:           r,
				ReadHeaderTimeout: 30 * time.Second,
			}
			fmt.Fprintf(os.Stderr, "==> dev server listening on http://localhost:%d", 8080)
			return srv.ListenAndServe()
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func Build() error {
	ctx := context.Background()

	client, workdir, err := initDagger(ctx)
	if err != nil {
		return err
	}

	_, err = innerBuild(ctx, client, workdir)
	if err != nil {
		return err
	}

	return nil
}

func Deploy() error {
	ctx := context.Background()

	client, _, err := initDagger(ctx)
	if err != nil {
		return err
	}

	_, err = client.Git("https://github.com/laupse/dagger-kapp").
		Branch("main").
		Tree().
		LoadProject("dagger.json").
		Install(ctx)
	if err != nil {
		return err
	}

	return nil
}

func IntTest() {

}

func UnitTest() error {
	ctx := context.Background()

	client, workdir, err := initDagger(ctx)
	if err != nil {
		return err
	}

	testContainerId, err := client.Container().
		From("golang:1.19-alpine3.16").
		WithMountedDirectory("/go/src/hello", workdir).
		WithWorkdir("/go/src/hello").
		Exec(dagger.ContainerExecOpts{
			Args: []string{"apk", "update"},
		}).
		Exec(dagger.ContainerExecOpts{
			Args: []string{"apk", "add", "build-base"},
		}).
		Exec(dagger.ContainerExecOpts{
			Args: []string{"/usr/local/go/bin/go", "mod", "download"},
		}).
		ID(ctx)
	if err != nil {
		return err
	}

	_, err = client.Container(dagger.ContainerOpts{ID: testContainerId}).
		Exec(dagger.ContainerExecOpts{
			Args: []string{"/usr/local/go/bin/go", "test", "./adapter/k8s", "-v"},
		}).
		Stdout().
		Contents(ctx)
	if err != nil {
		return err
	}

	return nil
}

func Push() error {
	ctx := context.Background()

	client, workdir, err := initDagger(ctx)
	if err != nil {
		return err
	}

	build, err := innerBuild(ctx, client, workdir)
	if err != nil {
		return err
	}
	ref := fmt.Sprintf(
		"%s/%s:%s",
		viper.GetString("registry-url"),
		viper.GetString("image-name"),
		viper.GetString("image-tag"),
	)
	_, err = build.Publish(ctx, ref)
	if err != nil {
		return err
	}

	return nil
}
