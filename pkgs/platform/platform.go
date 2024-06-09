package platform

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"slices"

	"github.com/go-enry/go-enry/v2"
)

// We want to support possible future container types
type Platform interface {
	BuildImage(file *os.File, labels *Labels) error
	Exec() error
	ListImages() ([]*Image, error)
	RemoveImage(name string) error
	Rollback(hash string) error
	Test(name string) error
}

type Image struct {
	Name    string
	Runtime string
	Hash    string
	Meta    *Labels
}

type Labels struct {
	GuildID   string
	OwnerID   string
	UserID    string
	Timestamp time.Time
}

//go:embed runtimes/*
var RuntimeFiles embed.FS

func FunctionTemplate(name string, fullfunc bool, runtime string) error {
	runtimeList, err := ListRuntimes()
	if err != nil {
		return fmt.Errorf("error listing runtimes: %v", err)
	}

	if !slices.Contains(runtimeList, runtime) {
		return fmt.Errorf("invalid runtime: %s", runtime)
	}
	if runtime == "" {
		runtime, _ = enry.GetLanguageByExtension(name) //Lookup by extension
	}

	fp, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	runtimeFp := filepath.Join("runtimes", runtime)

	runFiles, err := fs.ReadDir(RuntimeFiles, runtimeFp)
	if err != nil {
		return fmt.Errorf("error reading runtime files: %v", err)
	}

	funcFile, err := fs.ReadDir(RuntimeFiles, runtimeFp+"/function")
	if err != nil {
		return fmt.Errorf("error reading runtime files: %v", err)
	}

	render := func(entry string) error {
		data, err := RuntimeFiles.ReadFile(entry)
		if err != nil {
			return fmt.Errorf("error reading runtime file: %v", err)
		}
		fmt.Print(fp)
		err = os.WriteFile(filepath.Join(fp, filepath.Base(entry)), data, 0644)
		if err != nil {
			return fmt.Errorf("error writing file: %v", err)
		}
		return nil
	}

	if !fullfunc {
		for _, entry := range funcFile {
			err := render(filepath.Join(runtimeFp, "/function/", entry.Name()))
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, entry := range runFiles {
		if entry.IsDir() {
			fp := filepath.Join(fp, entry.Name())
			if err := os.MkdirAll(fp, 0755); err != nil {
				return fmt.Errorf("error creating directory: %v", err)
			}
		} else {
			err := render(runtimeFp)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ListRuntimes() ([]string, error) {
	runtimes, err := RuntimeFiles.ReadDir("runtimes")
	if err != nil {
		return nil, fmt.Errorf("error reading runtime files: %v", err)
	}
	var runtimeNames []string
	for _, entry := range runtimes {
		if entry.IsDir() {
			runtimeNames = append(runtimeNames, entry.Name())
		}
	}
	return runtimeNames, nil
}
