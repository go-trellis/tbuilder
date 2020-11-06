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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-trellis/common/formats"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//  checkLicensesCmd represents the  check licenses command
var checkLicensesCmd = &cobra.Command{
	Use:   "licenses",
	Short: "checking files' licenses",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		checkConfig.runCheckLicenses()
	},
}

var (
	validHeaderStrings = []string{"copyright", "generated"}
)

func init() {

	checkLicensesCmd.Flags().StringSliceVar(&checkConfig.Extensions, "extensions", []string{".go"},
		"Comma separated list of valid source code extensions (default is .go)")
	checkLicensesCmd.Flags().IntVar(&checkConfig.Length, "length", 10,
		"The number of lines to read from the head of the file")
	checkLicensesCmd.Flags().StringVar(&checkConfig.Location, "location", ".", "Directory path to check licenses")

	checkCmd.AddCommand(checkLicensesCmd)
}

func (p *CheckConfig) runCheckLicenses() {
	path := fmt.Sprintf("%s%c", filepath.Clean(p.Location), filepath.Separator)

	filesMissingHeaders, err := checkLicenses(path, p.Length, p.Extensions)
	if err != nil {
		fatal(errors.Wrap(err, "Failed to check files for license header"))
	}

	if len(filesMissingHeaders) == 0 {
		return
	}

	for _, file := range filesMissingHeaders {
		fmt.Println(file)
	}
}

func checkLicenses(path string, n int, extensions []string) ([]string, error) {
	var missingHeaders []string
	walkFunc := func(filepath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if strings.HasPrefix(filepath, "vendor/") {
			return nil
		}

		if !formats.SuffixStringInSlice(f.Name(), extensions) {
			return nil
		}

		file, err := os.Open(filepath)
		if err != nil {
			return err
		}

		defer file.Close()

		pass := false
		scanner := bufio.NewScanner(file)
		for i := 0; i < n; i++ {
			scanner.Scan()

			if err = scanner.Err(); err != nil {
				return err
			}

			if formats.StringContainedInSlice(strings.ToLower(scanner.Text()), validHeaderStrings) {
				pass = true
			}
		}

		if !pass {
			missingHeaders = append(missingHeaders, filepath)
		}

		return nil
	}

	err := filepath.Walk(path, walkFunc)
	if err != nil {
		return nil, err
	}

	return missingHeaders, nil
}
