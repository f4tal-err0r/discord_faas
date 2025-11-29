package runtimes

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"
)

var runtimeList = []string{
	"python",
	"golang",
}

type RuntimeConfig struct {
	Name        string   `yaml:"name"`
	Lang        string   `yaml:"lang"`
	Version     string   `yaml:"version"`
	LockedFiles []string `yaml:"lockedFiles"`
	Build       struct {
		Image    string   `yaml:"image"`
		Commands []string `yaml:"commands"`
	} `yaml:"build"`
	Run struct {
		Commands []string `yaml:"commands"`
	} `yaml:"run"`
}

//go:embed _templates/*
var RuntimeFiles embed.FS

func GenerateFunc(name string, runtime string) error {
	if !slices.Contains(runtimeList, runtime) {
		return fmt.Errorf("invalid runtime: %s", runtime)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %v", err)
	}

	funcDir := filepath.Join(cwd, name)
	if err := os.MkdirAll(funcDir, 0755); err != nil {
		return fmt.Errorf("error creating function directory: %v", err)
	}

	if err := renderDir(filepath.Join("_templates", runtime), funcDir); err != nil {
		return err
	}

	return nil
}

func ProtectedFilesFunc(fp string, runtime string) error {
	rc, err := GetRuntimeConfig(runtime)
	if err != nil {
		return err
	}

	for _, lockedFile := range rc.LockedFiles {
		srcPath := filepath.Join(runtime, lockedFile)
		dstPath := filepath.Join(fp, lockedFile)

		data, err := RuntimeFiles.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("error reading locked file %s: %v", srcPath, err)
		}
		if err := os.WriteFile(dstPath, data, 0644); err != nil {
			return fmt.Errorf("error writing locked file %s: %v", dstPath, err)
		}
	}
	return nil
}

func GetRuntimeConfig(runtime string) (*RuntimeConfig, error) {
	if !slices.Contains(runtimeList, runtime) {
		return nil, fmt.Errorf("invalid runtime: %s", runtime)
	}

	data, err := RuntimeFiles.ReadFile(filepath.Join(runtime, "runtime.yaml"))
	if err != nil {
		return nil, fmt.Errorf("error reading runtime config: %v", err)
	}

	var rc RuntimeConfig
	if err := yaml.Unmarshal(data, &rc); err != nil {
		return nil, fmt.Errorf("error unmarshaling runtime config: %v", err)
	}

	return &rc, nil
}

func renderDir(runtimeDir, targetDir string) error {
	entries, err := RuntimeFiles.ReadDir(runtimeDir)
	if err != nil {
		return fmt.Errorf("error reading embedded directory %s: %v", runtimeDir, err)
	}

	for _, entry := range entries {
		if entry.Name() == "runtime.yaml" {
			continue
		}
		srcPath := filepath.Join(runtimeDir, entry.Name())
		dstPath := filepath.Join(targetDir, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return fmt.Errorf("error creating directory %s: %v", dstPath, err)
			}
			if err := renderDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := RuntimeFiles.ReadFile(srcPath)
			if err != nil {
				return fmt.Errorf("error reading embedded file %s: %v", srcPath, err)
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return fmt.Errorf("error writing file %s: %v", dstPath, err)
			}
			fmt.Printf("%s\n", dstPath)
		}
	}
	return nil
}
