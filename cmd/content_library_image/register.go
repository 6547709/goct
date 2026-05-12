package content_library_image

import "github.com/spf13/cobra"

const groupID = "content_library_image"

func Register(root *cobra.Command) {
	root.AddCommand(newLs(), newDelete(), newImport(), newDistribute())
}
