package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v2"
)

type Payload struct {
	Repository struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"repository"`
	HeadCommit struct {
		Id      string `json:"id"`
		Message string `json:"message"`
		Url     string `json:"url"`
	} `json:"head_commit"`
}

type Worfkflow struct {
	Name string `yaml:"name"`
	Jobs []Job  `yaml:"jobs"`
}

type Job struct {
	Name   string   `yaml:"name"`
	Script string   `yaml:"script"`
	Env    map[string]string `yaml:"env"`
}

func HandlePayload(body io.ReadCloser) error {
	var payload Payload
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		return err
	}

	slog.Info("Received payload", "payload", payload)
	_, clonePath, err := CheckoutRepository(payload)
	if err != nil {
		return err
	}
	defer os.RemoveAll(filepath.Dir(clonePath)) // clean up

	workflow, err := readWorkflowFile(clonePath)
	if err != nil {
		return err
	}

	for _, job := range workflow.Jobs {
		// prepend command with `#!/bin/bash` to make it executable
		job.Script = "#!/usr/bin/env bash --noprofile --norc -eo pipefail\n" + job.Script
		
		filePath := "/tmp/" + job.Name + ".sh"

		err := os.WriteFile(filePath, []byte(job.Script), 0755)
		if err != nil {
			return err
		}

		_, err = Execute(filePath, []string{}, job.Env)
		if err != nil {
			return err
		}
	}

	return nil
}

func CheckoutRepository(p Payload) (*git.Repository, string, error) {
	tmpdir, err := os.MkdirTemp("", "repos")
	if err != nil {
		return nil, "", err
	}

	clonePath := filepath.Join(tmpdir, p.Repository.Name)
	r, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL: p.Repository.Url,
	})
	if err != nil {
		return nil, "", err
	}
	return r, clonePath, nil
}

func readWorkflowFile(clonePath string) (*Worfkflow, error) {
	filePath := clonePath + "/.popci.yaml"
	if os.Getenv("IS_LOCAL") == "true" {
		slog.Info("Running locally, using local workflow file")
		filePath = ".popci.yaml"
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var workflow Worfkflow
	err = yaml.Unmarshal(content, &workflow)
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func Execute(script string, command []string, env map[string]string) (bool, error) {
	// convert env to format expected by exec.Cmd
	envArr := []string{}
	for k, v := range env {
		envArr = append(envArr, k+"="+v)
	}

	stdout := &OutputWriter{w: os.Stdout}

	cmd := &exec.Cmd{
		Path:   script,
		Args:   command,
		Stdout: stdout,
		Stderr: os.Stderr,
		Env:    append(os.Environ(), envArr...),
	}

	err := cmd.Start()
	if err != nil {
		return false, err
	}

	err = cmd.Wait()
	if err != nil {
		return false, err
	}

	return true, nil
}
