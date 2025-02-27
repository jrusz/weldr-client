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
	metadataCmd = &cobra.Command{
		Use:   "metadata UUID",
		Short: "Get a tar of the the metadata for the compose",
		Long:  "Get a tar of the the metadata for the compose",
		RunE:  getMetadata,
		Args:  cobra.ExactArgs(1),
	}
)

func init() {
	composeCmd.AddCommand(metadataCmd)
}

func getMetadata(cmd *cobra.Command, args []string) error {
	fn, resp, err := root.Client.ComposeMetadata(args[0])
	if err != nil {
		return root.ExecutionError(cmd, "Metadata error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionErrors(cmd, resp.Errors)
	}

	fmt.Println(fn)

	return nil
}
