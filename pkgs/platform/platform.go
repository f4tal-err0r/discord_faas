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

// We want to support possible future container types
type Platform interface {
	BuildImage(*os.File, *Image) error
	Exec(string) error
	ListImages() ([]*Image, error)
	RemoveImage(string) error
	Rollback(string) error
	Test(string) error
}

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
	"golang": []string{"function/*"},
	"ruby":   []string{"function.rb", "content_pb.rb"},
}

//go:embed templates/*
var RuntimeFiles embed.FS

func FunctionTemplate(name string, build bool, runtime string) error {
	runtimeList := maps.Keys(UserLangDir)

	if !slices.Contains(runtimeList, runtime) {
		return fmt.Errorf("invalid runtime: %s", runtime)
	}

	cfp, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	if err = os.MkdirAll(filepath.Join(cfp, name), 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	} else {
		fmt.Printf("./%s/\n", name)
	}

	var render func(string) error

	render = func(fp string) error {
		if strings.HasSuffix(fp, "/*") {
			fp = strings.TrimSuffix(fp, "/*")
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
		err = os.WriteFile(filepath.Join(cfp, name, filepath.Base(fp)), data, 0644)
		if err != nil {
			return fmt.Errorf("error writing file: %v", err)
		}
		fmt.Print("./" + filepath.Join(name, filepath.Base(fp)) + "\n")
		return nil
	}

	for _, usrfiles := range UserLangDir[runtime] {
		if err := render(filepath.Join("templates", runtime, usrfiles)); err != nil {
			return err
		}
	}

	return nil
}
