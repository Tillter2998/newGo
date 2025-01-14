package applicationStrategy

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	templates "github.com/Tillter2998/newGo/internal"
)

type ApplicationStrategy interface {
	Execute(githubUserName string, name string, directory string) error
}

type (
	Empty   struct{}
	RESTApi struct{}
)

func (e *Empty) Execute(githubUserName string, name string, directory string) error {
	homeDir, err := os.UserHomeDir()
	var basePath string
	if err != nil {
		return err
	}

	if len(directory) > 1 && directory[:2] == "~/" {
		basePath = filepath.Join(homeDir, directory[2:], name)
	} else {
		basePath = filepath.Join(directory, name)
	}

	template := templates.Templates

	relPath := filepath.Join("templates", "{{empty}}/")
	err = fs.WalkDir(template, relPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(relPath, path)
		if err != nil {
			return err
		}

		updatedPath := strings.ReplaceAll(relativePath, "{{empty}}", name)
		updatedPath = strings.ReplaceAll(updatedPath, "{{name}}", name)
		updatedPath = strings.ReplaceAll(updatedPath, ".template", "")
		destPath := filepath.Join(basePath, updatedPath)

		if d.IsDir() {
			os.MkdirAll(destPath, 0o755)
		} else {
			file, err := template.Open(path)
			if err != nil {
				return err
			}
			writer, err := os.Create(destPath)
			if err != nil {
				return err
			}
			io.Copy(writer, file)
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = goModInit(githubUserName, name, basePath)
	if err != nil {
		return err
	}

	return nil
}

func (r *RESTApi) Execute(githubUserName string, name string, directory string) error {
	homeDir, err := os.UserHomeDir()
	var basePath string
	if err != nil {
		return err
	}

	if len(directory) > 1 && directory[:2] == "~/" {
		basePath = filepath.Join(homeDir, directory[2:], name)
	} else {
		basePath = filepath.Join(directory, name)
	}

	template := templates.Templates

	relPath := filepath.Join("templates", "{{restApi}}/")
	err = fs.WalkDir(template, relPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(relPath, path)
		if err != nil {
			return err
		}

		updatedPath := strings.ReplaceAll(relativePath, "{{restApi}}", name)
		updatedPath = strings.ReplaceAll(updatedPath, "{{name}}", name)
		updatedPath = strings.ReplaceAll(updatedPath, ".template", "")
		destPath := filepath.Join(basePath, updatedPath)

		if d.IsDir() {
			os.MkdirAll(destPath, 0o755)
		} else {
			file, err := template.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			writer, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer writer.Close()

			io.Copy(writer, file)
		}

		return nil
	})
	if err != nil {
		return err
	}

	content, err := os.ReadFile(filepath.Join(basePath, "cmd", name, "main.go"))
	if err != nil {
		return err
	}

	text := string(content)

	replacements := map[string]string{
		"{{githubAccount}}": githubUserName,
		"{{name}}":          name,
	}

	for token, replacement := range replacements {
		text = strings.ReplaceAll(text, token, replacement)
	}

	err = os.WriteFile(filepath.Join(basePath, "cmd", name, "main.go"), []byte(text), 0o755)
	if err != nil {
		return err
	}

	err = goModInit(githubUserName, name, basePath)
	if err != nil {
		return err
	}

	return nil
}

func goModInit(githubUserName string, name string, path string) error {
	cmd := exec.Command("go", "mod", "init", fmt.Sprintf("github.com/%s/%s", githubUserName, name))
	cmd.Dir = path

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	GetRegistry().RegisterStrategy("empty", &Empty{})
	GetRegistry().RegisterStrategy("restApi", &RESTApi{})
}
