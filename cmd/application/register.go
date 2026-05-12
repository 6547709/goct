package application

import "github.com/spf13/cobra"

const groupID = "application"

func Register(root *cobra.Command) {
	root.AddCommand(
		newGetApplications(),
		newGetPackages(),
		newUploadPackage(),
		newDeletePackage(),
		newDeploy(),
	)
}
