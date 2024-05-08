// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package mongolocal

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	testenv "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
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
	// Get a random port.
	port := RandomPort()

	// Get DB path.
	dbPath := filepath.Join(testenv.ExecRoot(), "tmp", "mongodb", fmt.Sprintf("db-%d", port))

	return Config{
		DbPath: dbPath,
		Port:   port,
	}
}

type MongoLocal struct {
	proc    *exec.Cmd
	watcher *exec.Cmd
	config  Config
	logger  *zap.Logger
}

func New(t *testing.T, config Config) (*MongoLocal, error) {
	// Create a command to start MongoDB.
	cmd := exec.Command("mongod", "--dbpath", config.DbPath, "--port", fmt.Sprintf("%d", config.Port))

	// Create a logger.
	level := zap.InfoLevel
	if testenv.IsDebug() {
		level = zap.DebugLevel
	}
	logger := zaptest.NewLogger(t, zaptest.Level(level))

	return &MongoLocal{
		proc:   cmd,
		config: config,
		logger: logger,
	}, nil
}

func RunWithServer(t *testing.T, callback func(*MongoLocal)) {
	// Create a new MongoDB server.
	server, err := New(t, DefaultConfig())
	logger := server.logger
	logger.Info("Starting MongoDB server...")
	if err != nil {
		logger.Error("Failed to create MongoDB server", zap.Error(err))
		return
	}

	// Attach logger to the server.
	reader, err := server.Process().StdoutPipe()
	if err != nil {
		logger.Error("Failed to attach logger to MongoDB server", zap.Error(err))
		return
	}
	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			logger.Sugar().Debugf("[mongod] %s", scanner.Text())
		}
	}()

	// Start the server.
	if err := server.Start(); err != nil {
		logger.Error("Failed to start MongoDB server", zap.Error(err))
		if err := server.Stop(); err != nil {
			logger.Error("Failed to stop MongoDB server", zap.Error(err))
		}
		return
	}
	logger.Sugar().Infof("MongoDB server started on port %d", server.Config().Port)

	defer func() {
		logger.Info("Stopping MongoDB server...")
		if err := server.Stop(); err != nil {
			logger.Error("Failed to stop MongoDB server", zap.Error(err))
			return
		}
	}()
	callback(server)
}

func (m *MongoLocal) createStorage() error {
	return os.MkdirAll(m.config.DbPath, 0755)
}

func (m *MongoLocal) clearStorage() error {
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

func (m *MongoLocal) Logger() *zap.Logger {
	return m.logger
}

func (m *MongoLocal) Start() error {
	if err := m.clearStorage(); err != nil {
		return err
	}
	if err := m.createStorage(); err != nil {
		return err
	}
	if err := m.proc.Start(); err != nil {
		return err
	}

	// Create a watcher to kill the process if the parent process dies.
	m.watcher = exec.Command("/bin/sh", "-c", watcherScript(os.Getpid(), m.proc.Process.Pid))

	if err := m.watcher.Start(); err != nil {
		return err
	}

	return nil
}

func (m *MongoLocal) Stop() error {
	if m.proc.Process != nil {
		if err := m.proc.Process.Kill(); err != nil {
			m.logger.Warn("Failed to kill MongoDB process", zap.Error(err))
		}
		if err := m.proc.Wait(); err != nil {
			m.logger.Warn("Failed to wait for MongoDB process", zap.Error(err))
		}
	}
	if m.watcher.Process != nil {
		if err := m.watcher.Process.Kill(); err != nil {
			m.logger.Warn("Failed to kill MongoDB watcher process", zap.Error(err))
		}
		if err := m.watcher.Wait(); err != nil {
			m.logger.Warn("Failed to wait for MongoDB watcher process", zap.Error(err))
		}
	}
	return m.clearStorage()
}
