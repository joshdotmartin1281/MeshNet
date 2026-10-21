package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"MeshNet/internal/app"
	"MeshNet/uauth"
)

func UserCreate(ctx context.Context, port app.AdminPort, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: user-create <name>")
	}

	user, err := port.CreateUser(ctx, args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", user.Name)

	return nil
}

func UserList(ctx context.Context, port app.AdminPort, args []string) error {
	if len(args) != 0 {
		return errors.New("usage: user-list")
	}

	users, err := port.ListUsers(ctx)
	if err != nil {
		return err
	}

	for _, u := range users {
		grants := make([]string, 0, len(u.Grants))

		for _, g := range u.Grants {
			grants = append(grants, g.Collection+"="+g.Level.String())
		}

		fmt.Printf("%s\t%s\n", u.Name, strings.Join(grants, ","))
	}

	return nil
}
func TokenCreate(ctx context.Context, port app.AdminPort, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: token-create <user>")
	}

	token, err := port.IssueToken(ctx, args[0])
	if err != nil {
		return err
	}

	fmt.Println(token)
	fmt.Fprintln(os.Stderr, "store this token now; it cannot be shown again")

	return nil
}

func TokenRevoke(ctx context.Context, port app.AdminPort, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: token-revoke <user>")
	}

	return port.RevokeTokens(ctx, args[0])
}

func Grant(ctx context.Context, port app.AdminPort, args []string) error {
	fs := flag.NewFlagSet("grant", flag.ContinueOnError)

	collection := fs.String("c", "", "collection to grant access to")

	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*collection) == "" {
		return errors.New("collection is required")
	}

	if fs.NArg() != 2 {
		return errors.New("usage: grant -c <collection> <user> <read|write>")
	}

	level, err := auth.ParseLevel(fs.Arg(1))
	if err != nil {
		return err
	}

	return port.Grant(ctx, fs.Arg(0), *collection, level)
}

func Ungrant(ctx context.Context, port app.AdminPort, args []string) error {
	fs := flag.NewFlagSet("ungrant", flag.ContinueOnError)

	collection := fs.String("c", "", "collection to revoke access to")

	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*collection) == "" {
		return errors.New("collection is required")
	}

	if fs.NArg() != 1 {
		return errors.New("usage: ungrant -c <collection> <user>")
	}

	return port.Ungrant(ctx, fs.Arg(0), *collection)
}
