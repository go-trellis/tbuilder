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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "trellis",
	Short:   "build trellis project with config file",
	Version: "v0.1.0",
	Long:    ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "print debug information")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// info prints the given message only if running in verbose mode
func info(message string) {
	if verbose {
		fmt.Println(message)
	}
}

// warn prints a non-fatal error
func warn(err error) {
	if verbose {
		fmt.Fprintf(os.Stderr, `/!\ %+v\n`, err)
	} else {
		fmt.Fprintln(os.Stderr, `/!\`, err)
	}
}

// fatal prints a error and exit
func fatal(err error) {
	printErr(err)
	os.Exit(1)
}

// printErr prints a error
func printErr(err error) {
	if verbose {
		fmt.Fprintf(os.Stderr, "!! %+v\n", err)
	} else {
		fmt.Fprintln(os.Stderr, "!!", err)
	}
}
