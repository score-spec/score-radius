// Copyright 2025 The Score Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"sort"
	"strings"

	"github.com/score-spec/score-radius/internal/provisioners"
	"github.com/score-spec/score-radius/internal/provisioners/loader"
	"github.com/score-spec/score-radius/internal/state"

	"github.com/score-spec/score-go/formatter"
)

var (
	provisionersGroup = &cobra.Command{
		Use:   "provisioners",
		Short: "Subcommands related to resources provisioners",
	}
	provisionersList = &cobra.Command{
		Use:   "list [--format table|json]",
		Short: "List the resources provisioners",
		Long: `The list command will list out the resources provisioners. This requires an active score-cloudrun state
after 'init' has been run. The list of resources provisioners will be empty if no provisioners are defined.
`,
		Args:          cobra.ArbitraryArgs,
		SilenceErrors: true,
		RunE:          listProvisioners,
	}
)

func listProvisioners(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	sd, ok, err := state.LoadStateDirectory(".")
	if err != nil {
		return fmt.Errorf("failed to load existing state directory: %w", err)
	} else if !ok {
		return fmt.Errorf("no state directory found, run 'score-cloudrun init' first")
	}
	slog.Debug(fmt.Sprintln("Listing resources provisioners"))

	provisioners, err := loader.LoadProvisionersFromDirectory(sd.Path, loader.ProvisionersFileSuffix)
	if err != nil {
		return fmt.Errorf("failed to load resources provisioners in %s: %w", sd.Path, err)
	}

	if len(provisioners) == 0 {
		slog.Info("No resources provisioners found")
		return nil
	}

	outputFormat := cmd.Flag("format").Value.String()
	return displayProvisioners(provisioners, outputFormat)
}

func displayProvisioners(provisioners []provisioners.Provisioner, outputFormat string) error {
	var outputFormatter formatter.OutputFormatter
	sortedProvisioners := sortProvisionersByType(provisioners)

	switch outputFormat {
	case "json":
		type jsonData struct {
			Type        string
			Class       string
			Format      string
			Params      []string
			Outputs     []string
			Description string
		}
		var outputs []jsonData
		for _, provisioner := range sortedProvisioners {
			outputs = append(outputs, jsonData{
				Type:        provisioner.ResType,
				Class:       provisioner.Class,
				Format:      provisioner.Format,
				Params:      provisioner.Params,
				Outputs:     provisioner.Outputs,
				Description: provisioner.Description,
			})
		}
		outputFormatter = &formatter.JSONOutputFormatter[[]jsonData]{Data: outputs}
	default:
		rows := [][]string{}
		for _, provisioner := range sortedProvisioners {
			rows = append(rows, []string{provisioner.ResType, provisioner.Class, provisioner.Format, strings.Join(provisioner.Params, ", "), strings.Join(provisioner.Outputs, ", "), provisioner.Description})
		}
		headers := []string{"Type", "Class", "Format", "Params", "Outputs", "Description"}
		outputFormatter = &formatter.TableOutputFormatter{
			Headers: headers,
			Rows:    rows,
		}
	}
	return outputFormatter.Display()
}

func sortProvisionersByType(provisioners []provisioners.Provisioner) []provisioners.Provisioner {
	sort.Slice(provisioners, func(i, j int) bool {
		return provisioners[i].ResType < provisioners[j].ResType
	})
	return provisioners
}

func init() {
	provisionersList.Flags().StringP("format", "f", "table", "Format of the output: table (default), json")
	provisionersGroup.AddCommand(provisionersList)
	rootCmd.AddCommand(provisionersGroup)
}