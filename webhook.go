package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"os/exec"

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
	Name    string   `yaml:"name"`
	Command string `yaml:"command"`
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

	workflow, err := readWorkflowFile(clonePath)
	if err != nil {
		return err
	}

	for _, job := range workflow.Jobs {
		filePath := "/tmp/" + job.Name + ".sh"

		// prepend command with `#!/bin/bash` to make it executable
		job.Command = "#!/bin/bash\n" + job.Command
		
		err := os.WriteFile(filePath, []byte(job.Command), 0755)
		if err != nil {
			return err
		}

		_, err = Execute(filePath, []string{})
		if err != nil {
			return err
		}
	}

	return nil
}

func CheckoutRepository(p Payload) (*git.Repository, string, error) {
	clonePath := "/tmp/repo/" + p.Repository.Name
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

	slog.Info("Workflow file content", "workflow", workflow)
	return &workflow, nil
}

func Execute(script string, command []string) (bool, error) {
	cmd := &exec.Cmd{
		Path:   script,
		Args:   command,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	// slog.Info("Executing command ", "command", cmd.String())

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
