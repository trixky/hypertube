//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Build() error {
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "hypertube", ".")
	return cmd.Run()
}

func Run() error {
	fmt.Println("Run...")
	cmd := exec.Command("go", "run", "main.go")

	return cmd.Run()
}

func TestApiAuth() error {
	// --------------------- REDIS
	if redis_ctx, redis_container, err := createRedisContainer(); err != nil {
		return err
	} else {
		defer redis_container.Terminate(redis_ctx)
	}

	// --------------------- POSTGRES
	if postgres_ctx, postgres_container, err := createPostgresContainer(); err != nil {
		return err
	} else {
		defer postgres_container.Terminate(postgres_ctx)
	}

	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = "../../api-auth"

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func TestApiUser() error {
	// --------------------- REDIS
	if redis_ctx, redis_container, err := createRedisContainer(); err != nil {
		return err
	} else {
		defer redis_container.Terminate(redis_ctx)
	}

	// --------------------- POSTGRES
	if postgres_ctx, postgres_container, err := createPostgresContainer(); err != nil {
		return err
	} else {
		defer postgres_container.Terminate(postgres_ctx)
	}

	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = "../../api-user"

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ---------------------------------------------
// ---------------------------------------------
// ---------------------------------------------

func createRedisContainer() (context.Context, testcontainers.Container, error) {
	// https://golang.testcontainers.org/quickstart/gotest/

	ctx := context.Background()

	request := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})

	if err != nil {
		return ctx, nil, err
	}

	// Give the container ip to the tests
	if container_ip, err := container.ContainerIP(ctx); err != nil {
		return ctx, nil, err
	} else {
		os.Setenv("REDIS_HOST", container_ip)
	}

	return ctx, container, err
}

func createPostgresContainer() (context.Context, testcontainers.Container, error) {
	// https://brietsparks.com/testcontainers-golang-db-access/

	ctx := context.Background()

	environment := map[string]string{
		"POSTGRES_HOST":      "postgres",
		"POSTGRES_USER":      "admin_test",
		"POSTGRES_PASSWORD":  "1234",
		"POSTGRES_DB":        "hypertube",
		"POSTGRES_PORT":      "5432",
		"POSTGRES_PORT_FULL": "5432/tcp",
	}

	dbURL := func(port nat.Port) string {
		return fmt.Sprintf("%s://%s:%s@localhost:%s/%s?sslmode=disable", environment["POSTGRES_HOST"], environment["POSTGRES_USER"], environment["POSTGRES_PASSWORD"], port.Port(), environment["POSTGRES_DB"])
	}

	for key, value := range environment {
		os.Setenv(key, value)
	}

	abs_init_path, err := filepath.Abs("../../postgres/init.sql")
	if err != nil {
		return ctx, nil, err
	}

	request := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{environment["POSTGRES_PORT_FULL"]},
		Env:          environment,
		WaitingFor:   wait.ForSQL(nat.Port(environment["POSTGRES_PORT_FULL"]), "postgres", dbURL).Timeout(time.Second * 15),
		Mounts: testcontainers.ContainerMounts{
			testcontainers.BindMount(abs_init_path, "/docker-entrypoint-initdb.d/init.sql"),
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})

	if err != nil {
		return ctx, nil, err
	}

	// Give the container ip to the tests
	if container_ip, err := container.ContainerIP(ctx); err != nil {
		return ctx, nil, err
	} else {
		os.Setenv("POSTGRES_HOST", container_ip)
	}

	return ctx, container, err
}
