package platform

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"slices"

	"github.com/go-enry/go-enry/v2"
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
	"ruby":   []string{"!main.rb"},
}

//go:embed templates
var RuntimeFiles embed.FS

func FunctionTemplate(name string, fullfunc bool, runtime string) error {
	runtimeList := maps.Keys(UserLangDir)

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

	if err = os.MkdirAll(filepath.Join(fp, name), 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	var render func(string) error

	render = func(embedfp string) error {
		if strings.HasSuffix(embedfp, "/*") {
			dirFiles, err := RuntimeFiles.ReadDir(embedfp)
			if err != nil {
				return fmt.Errorf("error reading runtime directory: %v", err)
			}
			for _, embedfp := range dirFiles {
				err := render(embedfp.Name())
				if err != nil {
					return err
				}
			}
			return nil
		}
		data, err := RuntimeFiles.ReadFile(embedfp)
		if err != nil {
			return fmt.Errorf("error reading runtime file: %v", err)
		}
		err = os.WriteFile(filepath.Join(fp, filepath.Base(embedfp)), data, 0644)
		if err != nil {
			return fmt.Errorf("error writing file: %v", err)
		}
		return nil
	}

	return nil
}
