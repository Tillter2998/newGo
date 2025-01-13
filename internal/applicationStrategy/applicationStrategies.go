package applicationStrategy

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	templates "github.com/Tillter2998/newGo/internal"
)

type ApplicationStrategy interface {
	Execute(name string, directory string) error
}

type Empty struct{}

func (e *Empty) Execute(name string, directory string) error {
	fmt.Println(name)
	fmt.Println(directory)

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

	relPath := filepath.Join("templates", "go", "{{empty}}/")
	err = fs.WalkDir(template, relPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(relPath, path)
		if err != nil {
			return err
		}

		err = copyFiles(template, "{{empty}}", path, basePath, relativePath, name, d)
		if err != nil {
			fmt.Println("Error copying files", err)
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println("Error walking templates:", err)
		return err
	}

	// Snipit of characters for future project idea
	// `project/
	// ├── main.go
	// ├── files/
	// │   └── example.txt`

	return nil
}

func init() {
	GetRegistry().RegisterStrategy("empty", &Empty{})
}

// copyFiles copies the specified projectType to a destination built from basePath and relativePath
func copyFiles(
	template embed.FS,
	projectType string,
	path string,
	basePath string,
	relativePath string,
	name string,
	dirEntry fs.DirEntry,
) error {
	updatedPath := strings.ReplaceAll(relativePath, projectType, name)
	updatedPath = strings.ReplaceAll(updatedPath, "{{name}}", name)
	updatedPath = strings.ReplaceAll(updatedPath, ".template", "")
	destPath := filepath.Join(basePath, updatedPath)

	if dirEntry.IsDir() {
		os.MkdirAll(destPath, 0o755)
	} else {
		file, err := template.Open(path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		writer, err := os.Create(destPath)
		io.Copy(writer, file)
	}

	// TODO: Investigate running commands like go mod init or other language equivalents after copying

	return nil
}
