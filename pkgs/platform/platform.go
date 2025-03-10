package platform

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"slices"

	"golang.org/x/exp/maps"
)

var UserLangDir = map[string][]string{
	"golang": {"function/"},
	"ruby":   {"function.rb", "content_pb.rb"},
}

var ConstTempl = []string{"README.md", "dfaas.yaml"}

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
				err := render(dir, filepath.Join("templates", runtime, embedfp.Name()))
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	for _, usrfiles := range UserLangDir[runtime] {
		if err := render(dir, "templates/"+runtime+"/"+usrfiles); err != nil {
			return err
		}
	}

	for _, file := range ConstTempl {
		if err := renderConst(dir, file, runtime); err != nil {
			return err
		}
	}

	return nil
}

func renderConst(dir string, file string, runtime string) error {
	readme, err := RuntimeFiles.ReadFile(filepath.Join("templates", file))
	if err != nil {
		return fmt.Errorf("error reading runtime file: %v", err)
	}

	t, err := template.New(file).Parse(string(readme))
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	f, err := os.Create(filepath.Join(dir, file))
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer f.Close()

	err = t.Execute(f, struct {
		Runtime string
		Name    string
	}{
		Runtime: runtime,
		Name:    filepath.Base(dir),
	})
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}

func render(dir string, efp string) error {
	if strings.HasSuffix(efp, "/") {
		efp = strings.TrimSuffix(efp, "/")
		dirFiles, err := RuntimeFiles.ReadDir(efp)
		if err != nil {
			return fmt.Errorf("error reading runtime directory: %v", err)
		}
		funcdir := filepath.Join(dir, filepath.Base(efp))
		if _, err := os.Stat(funcdir); os.IsNotExist(err) {
			err = os.Mkdir(funcdir, 0755)
			if err != nil {
				return fmt.Errorf("error creating directory: %v", err)
			}
		} else if err != nil {
			return fmt.Errorf("error checking directory: %v", err)
		}

		for _, embedfp := range dirFiles {
			err := render(funcdir, filepath.Join(efp, embedfp.Name()))
			if err != nil {
				return err
			}
		}
		return nil
	}

	data, err := RuntimeFiles.ReadFile(efp)
	if err != nil {
		return fmt.Errorf("error reading runtime file: %v", err)
	}
	err = os.WriteFile(filepath.Join(dir, filepath.Base(efp)), data, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}
	return nil
}
