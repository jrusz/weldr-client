// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	infoCmd = &cobra.Command{
		Use:   "info UUID",
		Short: "Show detailed information on the compose",
		Long:  "List basic information about composes",
		RunE:  info,
		Args:  cobra.ExactArgs(1),
	}
)

func init() {
	composeCmd.AddCommand(infoCmd)
}

func info(cmd *cobra.Command, args []string) error {
	info, resp, err := root.Client.ComposeInfo(args[0])
	if err != nil {
		return root.ExecutionError(cmd, "Info Error: %s", err)
	}
	if resp != nil {
		return root.ExecutionErrors(cmd, resp.Errors)
	}

	var imageSize string
	if info.ImageSize > 0 {
		imageSize = fmt.Sprintf("%d", info.ImageSize)
	}
	fmt.Printf("%s %-8s %-15s %s %-16s %s\n",
		info.ID,
		info.QueueStatus,
		info.Blueprint.Name,
		info.Blueprint.Version,
		info.ComposeType,
		imageSize)

	fmt.Println("Packages:")
	for _, p := range info.Blueprint.Packages {
		fmt.Printf("    %s-%s\n", p.Name, p.Version)
	}

	fmt.Println("Modules:")
	for _, m := range info.Blueprint.Modules {
		fmt.Printf("    %s-%s\n", m.Name, m.Version)
	}

	fmt.Println("Dependencies:")
	for _, d := range info.Deps.Packages {
		fmt.Printf("    %s\n", d)
	}

	return nil
}
