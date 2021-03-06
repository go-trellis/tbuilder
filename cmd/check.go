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
	"github.com/spf13/cobra"
)

//  check represents the check licenses command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "checking project's information, such as: licenses, changelog",
	Long:  ``,
}

var (
	checkConfig = CheckConfig{}
)

// CheckConfig check licences & change logs' config
type CheckConfig struct {
	Extensions []string
	Length     int
	Location   string

	Version string
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
