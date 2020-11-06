/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"bytes"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/go-trellis/common/formats"
	"github.com/go-trellis/common/shell"

	"github.com/go-trellis/config"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	goos   = build.Default.GOOS
	goarch = build.Default.GOARCH

	mainFile = "main.go"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: `build trellis project`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		stop := make(chan bool, 1)

		go func() {
			if err := buildRun(); err != nil {
				os.Exit(1)
			}
			stop <- true
		}()

		select {
		case <-ch:
		case <-stop:
		}
		close(ch)
		close(stop)

		if buildConfig.Project.Build.DelBuildFile {
			buildConfig.RemovePath()
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVar(&cfgFile, "config", ".tr_complier.yaml", "config file (default is .tr_complier.yaml)")
	buildCmd.Flags().BoolVar(&buildConfig.verbose, "verbose", false, "print debug information")
}

// BuildConfig build config
type BuildConfig struct {
	Project Project `json:"project" yaml:"project"`

	verbose   bool
	buildPath string
	mainPath  string
}

var buildConfig = &BuildConfig{}

func buildRun() error {

	r, err := config.NewSuffixReader(config.ReaderOptionFilename(cfgFile))
	if err != nil {
		return err
	}

	if err = r.Read(buildConfig); err != nil {
		return err
	}

	if buildConfig.Project.Build.Path == "" {
		buildConfig.Project.Build.Path = "."
	}

	if buildConfig.Project.Build.Type == "origin" {
		buildConfig.originMain()
	} else {

		if err := buildConfig.writeMainFile(); err != nil {
			return err
		}
	}

	return buildConfig.build()
}

func (p *BuildConfig) originMain() {
	p.buildPath = p.Project.Build.Path
	p.mainPath = filepath.Join(p.buildPath, mainFile)
}

func (p *BuildConfig) writeMainFile() error {

	mapImports := map[string]bool{}
	tpl, err := template.New("main").Funcs(template.FuncMap{
		"imports": func(url string) string {
			if _, ok := mapImports[url]; ok {
				return ""
			}
			mapImports[url] = true
			return fmt.Sprintf("import _ %q", url)
		},
	}).Parse(mainTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, p.Project); err != nil {
		return err
	}

	p.buildPath = filepath.Join(p.Project.Build.Path, fmt.Sprintf(".%s", uuid.Must(uuid.NewUUID()).String()))

	if err = os.MkdirAll(p.buildPath, 0755); err != nil {
		return err
	}

	p.mainPath = filepath.Join(p.buildPath, mainFile)

	if p.verbose {
		fmt.Println(buf.String())
	}

	err = ioutil.WriteFile(p.mainPath, []byte(buf.String()), 0644)
	if err != nil {
		err = fmt.Errorf("write %s failure to temp dir: %s", mainFile, err)
		return err
	}

	return nil
}

func (p *BuildConfig) build() error {

	outputFile := p.Project.Name
	if goos == "windows" {
		outputFile += ".exe"
	}

	params := []string{"build", "-o", outputFile}

	if len(p.Project.Build.Flags) > 0 {
		params = append(params, strings.Split(p.Project.Build.Flags, " ")...)
	}

	ldFlags, err := p.getLdflags()
	if err != nil {
		return err
	}

	if len(ldFlags) > 0 {
		params = append(params, "-ldflags", ldFlags)
	}

	params = append(params, p.mainPath)

	os.Setenv("CGO_ENABLED", "0")
	if p.Project.Go.CGo {
		os.Setenv("CGO_ENABLED", "1")
	}
	defer os.Unsetenv("CGO_ENABLED")

	if p.verbose {
		fmt.Println(params)
	}

	return shell.RunCommand("go", params...)
}

func (p *BuildConfig) getLdflags() (string, error) {
	var ldflags []string

	if len(strings.TrimSpace(p.Project.Build.Ldflags)) > 0 {
		tmpl, err := template.New("ldflags").Funcs(template.FuncMap{
			"compiler": func() string { return shellOutput("go version") },
			"datetime": time.Now().Format,
			"author":   AuthorFunc,
			"hostname": os.Hostname,
		}).Parse(p.Project.Build.Ldflags)
		if err != nil {
			return "", err
		}

		var tmplOutput bytes.Buffer
		if err := tmpl.Execute(&tmplOutput, p.Project); err != nil {
			return "", err
		}

		ldflags = append(ldflags, strings.Split(tmplOutput.String(), "\n")...)
	}

	extLDFlags := p.Project.Build.ExtLDFlags
	if p.Project.Build.Static && goos != "darwin" && goos != "solaris" && !formats.StringInSlice("-static", extLDFlags) {
		extLDFlags = append(p.Project.Build.ExtLDFlags, "-static")
	}

	if len(extLDFlags) > 0 {
		ldflags = append(ldflags, fmt.Sprintf("-extldflags '%s'", strings.Join(extLDFlags, " ")))
	}

	return strings.Join(ldflags[:], " "), nil
}

// AuthorFunc returns the current username.
func AuthorFunc() string {
	// os/user.Current() doesn't always work without CGO
	return shellOutput("whoami")
}

// shellOutput executes a shell command and returns the trimmed output
func shellOutput(cmd string) string {
	args := strings.Fields(cmd)
	out, _ := exec.Command(args[0], args[1:]...).Output()
	return strings.Trim(string(out), " \n\r")
}

// RemovePath remove build path
func (p *BuildConfig) RemovePath() error {
	if p.buildPath == "" {
		return nil
	}
	return os.RemoveAll(p.buildPath)
}
