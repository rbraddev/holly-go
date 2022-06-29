/*
Copyright © 2022 Ryan Bradshaw <ryan@rbradshaw.dev>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/rbraddev/holly/internal/hosts"
	"github.com/spf13/cobra"
)

var site string

var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "Site specific commands",

	Run: func(cmd *cobra.Command, args []string) {
		sp := hosts.SearchParams{
			Site:   []string{"123", "124"},
			Device: []string{"sw"},
		}

		hosts.LoadSolarwindsHosts(sp)

	},
}

func init() {
	rootCmd.AddCommand(sitesCmd)

	sitesCmd.Flags().Bool("enable", false, "Enable site")
	sitesCmd.Flags().Bool("disable", false, "Disable site")
	sitesCmd.Flags().Bool("check", false, "Check site current status")
	sitesCmd.Flags().StringVar(&site, "site", "", "Site number")

	sitesCmd.MarkFlagRequired("site")
}
