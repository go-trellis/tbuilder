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

// Project project config info
type Project struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`

	Build Build `yaml:"build"`

	Go Go `yaml:"go"`

	Services map[string]Service `yaml:"services"`
}

// Build build configure
type Build struct {
	Type         string   `yaml:"type"` // trellis, origin(main.go). default trellis
	Path         string   `yaml:"path"`
	Static       bool     `yaml:"static"`
	Flags        string   `yaml:"flags"`
	Ldflags      string   `yaml:"ldflags"`
	ExtLDFlags   []string `yaml:"ext_ldflags"`
	DelBuildFile bool     `yaml:"delete_build_file"`
}

// Go go
type Go struct {
	CGo bool `yaml:"cgo"`
}

// Service service info
type Service struct {
	URL string `yaml:"url"`
}
