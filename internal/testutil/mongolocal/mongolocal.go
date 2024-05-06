// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package mongolocal

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type Config struct {
	DbPath string
	Port   int
}

func RandomPort() int {
	for {
		port := rand.Intn(10000) + 20000
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			listener.Close()
			return port
		}
	}
}

func DefaultConfig() Config {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	port := RandomPort()
	return Config{
		DbPath: filepath.Join(cwd, "tmp", "mongodb", fmt.Sprintf("db-%d", port)),
		Port:   RandomPort(),
	}
}

type MongoLocal struct {
	proc   *exec.Cmd
	config Config
}

func New(config Config) (*MongoLocal, error) {
	// Ensure the path exists.
	if err := os.MkdirAll(config.DbPath, 0755); err != nil {
		return nil, err
	}

	cmd := exec.Command("mongod", "--dbpath", config.DbPath, "--port", fmt.Sprintf("%d", config.Port))
	return &MongoLocal{
		proc:   cmd,
		config: config,
	}, nil
}

func WithMongoLocal(t *testing.T, callback func(*MongoLocal)) {
	server, err := New(DefaultConfig())
	t.Log("Starting MongoDB server...")
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := server.Start(); err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("MongoDB server started on port %d", server.Config().Port)

	callback(server)

	t.Log("Stopping MongoDB server...")
	if err := server.Stop(); err != nil {
		t.Fatal(err)
		return
	}
}

func (m *MongoLocal) clean() error {
	return os.RemoveAll(m.config.DbPath)
}

func (m *MongoLocal) Config() Config {
	return m.config
}

func (m *MongoLocal) Process() *exec.Cmd {
	return m.proc
}

func (m *MongoLocal) URI() string {
	return fmt.Sprintf("mongodb://localhost:%d", m.config.Port)
}

func (m *MongoLocal) Start() error {
	if err := m.clean(); err != nil {
		return err
	}
	return m.proc.Start()
}

func (m *MongoLocal) Stop() error {
	if err := m.proc.Process.Kill(); err != nil {
		return err
	}
	if err := m.proc.Wait(); err != nil {
		return err
	}
	return m.clean()
}
