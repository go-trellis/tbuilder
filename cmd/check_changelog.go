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
	"os"

	"github.com/go-trellis/tbuilder/utils/changelog"
	"github.com/go-trellis/tbuilder/utils/repository"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var projInfo repository.Info

//  checkChangeLogCmd represents the  check licenses command
var checkChangeLogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "checking project's changelog",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		projInfo, err = repository.NewInfo(warn)
		if err != nil {
			fatal(err)
		}
		err = checkConfig.runCheckChangeLog()
		if err != nil {
			fatal(err)
		}
	},
}

func init() {

	checkChangeLogCmd.Flags().StringSliceVar(&checkConfig.Extensions, "extensions", []string{".go"},
		"Comma separated list of valid source code extensions (default is .go)")
	checkChangeLogCmd.Flags().IntVar(&checkConfig.Length, "length", 10,
		"The number of lines to read from the head of the file")
	checkChangeLogCmd.Flags().StringVar(&checkConfig.Location, "location", "CHANGELOG.md", "Directory path to check changelog")

	checkCmd.AddCommand(checkChangeLogCmd)
}

func (p *CheckConfig) runCheckChangeLog() error {

	if p.Version == "" {
		_, err := projInfo.ToSemver()
		if err != nil {
			return errors.Wrap(err, "invalid semver version")
		}

		p.Version = projInfo.Version
	}

	f, err := os.Open(p.Location)
	if err != nil {
		return err
	}
	defer f.Close()

	entry, err := changelog.ReadEntry(f, p.Version)
	if err != nil {
		return errors.Wrapf(err, "%s:", p.Location)
	}

	// Check that the changes are ordered correctly.
	err = entry.Changes.Sorted()
	if err != nil {
		return errors.Wrap(err, "invalid changelog entry")
	}

	return nil
}
