package user

import (
	"fmt"
	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newCreate() *cobra.Command {
	var name, username, password, roleID, email string
	c := &cobra.Command{
		Use: "user.create", Short: "Create a user", GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			ref, err := service.NewUser(cli).Create(c.Context(), adapter.UserCreateSpec{
				Name: name, Username: username, Password: password, RoleID: roleID, Email: email,
			})
			if err != nil { return err }
			if ref.IsSync() { fmt.Fprintln(c.OutOrStdout(), "User created (sync)"); return nil }
			return task.New(cli, task.Options{Out: c.OutOrStderr()}).Watch(c.Context(), ref.ID)
		},
	}
	c.Flags().StringVar(&name, "name", "", "Display name (required)")
	c.Flags().StringVar(&username, "username", "", "Login username (required)")
	c.Flags().StringVar(&password, "password", "", "Password (required)")
	c.Flags().StringVar(&roleID, "role", "", "Role ID (required)")
	c.Flags().StringVar(&email, "email", "", "Email address")
	_ = c.MarkFlagRequired("name"); _ = c.MarkFlagRequired("username")
	_ = c.MarkFlagRequired("password"); _ = c.MarkFlagRequired("role")
	return c
}
