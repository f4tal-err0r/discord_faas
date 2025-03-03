package platform

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"slices"

	"golang.org/x/exp/maps"
)

type Image struct {
	Name      string
	Runtime   string
	Hash      string
	Meta      *Labels
	Timestamp time.Time
}

type Labels struct {
	GuildID string
	OwnerID string
	UserID  string
}

var UserLangDir = map[string][]string{
	"golang": {"function/"},
	"ruby":   {"function.rb", "content_pb.rb"},
}

//go:embed templates/*
var RuntimeFiles embed.FS

func FunctionTemplate(dir string, build bool, runtime string) error {
	runtimeList := maps.Keys(UserLangDir)

	if !slices.Contains(runtimeList, runtime) {
		return fmt.Errorf("invalid runtime: %s", runtime)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	var render func(string) error

	render = func(fp string) error {
		if strings.HasSuffix(fp, "/") {
			fp = strings.TrimSuffix(fp, "/")
			dirFiles, err := RuntimeFiles.ReadDir(fp)
			if err != nil {
				return fmt.Errorf("error reading runtime directory: %v", err)
			}
			for _, embedfp := range dirFiles {
				err := render(fp + "/" + embedfp.Name())
				if err != nil {
					return err
				}
			}
			return nil
		}

		data, err := RuntimeFiles.ReadFile(fp)
		if err != nil {
			return fmt.Errorf("error reading runtime file: %v", err)
		}
		err = os.WriteFile(filepath.Join(dir, filepath.Base(fp)), data, 0644)
		if err != nil {
			return fmt.Errorf("error writing file: %v", err)
		}
		return nil
	}

	if build {
		//only render filepaths not in UserLangDir
		dirFiles, err := RuntimeFiles.ReadDir("templates/" + runtime)
		if err != nil {
			return fmt.Errorf("error reading runtime directory: %v", err)
		}

		for _, embedfp := range dirFiles {
			if !slices.ContainsFunc(UserLangDir[runtime], func(s string) bool {
				return embedfp.Name() == strings.TrimSuffix(s, "/")
			}) {
				err := render("templates/" + runtime + "/" + embedfp.Name())
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	for _, usrfiles := range UserLangDir[runtime] {
		if err := render(filepath.Join("templates", runtime, usrfiles)); err != nil {
			return err
		}
	}

	return nil
}
