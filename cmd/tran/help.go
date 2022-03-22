package tran

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/scmn-dev/tran/ios"
	"github.com/scmn-dev/tran/tools"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func rootUsageFunc(command *cobra.Command) error {
	command.Printf("Usage:  %s", command.UseLine())

	subcommands := command.Commands()

	if len(subcommands) > 0 {
		command.Print("\n\nCommands:\n")
		for _, c := range subcommands {
			if c.Hidden {
				continue
			}

			command.Printf("  %s\n", c.Name())
		}

		return nil
	}

	flagUsages := command.LocalFlags().FlagUsages()

	if flagUsages != "" {
		command.Println("\n\nFlags:")
		command.Print(tools.Indent(dedent(flagUsages), "  "))
	}

	return nil
}

func rootFlagErrorFunc(cmd *cobra.Command, err error) error {
	if err == pflag.ErrHelp {
		return err
	}

	return &tools.FlagError{Err: err}
}

var hasFailed bool

func HasFailed() bool {
	return hasFailed
}

func nestedSuggestFunc(command *cobra.Command, arg string) {
	command.Printf("unknown command %q for %q\n", arg, command.CommandPath())

	var candidates []string
	if arg == "help" {
		candidates = []string{"--help"}
	} else {
		if command.SuggestionsMinimumDistance <= 0 {
			command.SuggestionsMinimumDistance = 2
		}

		candidates = command.SuggestionsFor(arg)
	}

	if len(candidates) > 0 {
		command.Print("\nDid you mean this?\n")

		for _, c := range candidates {
			command.Printf("\t%s\n", c)
		}
	}

	command.Print("\n")
	_ = rootUsageFunc(command)
}

func isRootCmd(command *cobra.Command) bool {
	return command != nil && !command.HasParent()
}

func rootHelpFunc(cs *ios.ColorScheme, command *cobra.Command, args []string) {
	if isRootCmd(command.Parent()) && len(args) >= 2 && args[1] != "--help" && args[1] != "-h" {
		nestedSuggestFunc(command, args[1])
		hasFailed = true
		return
	}

	commands := []string{}

	for _, c := range command.Commands() {
		if c.Short == "" {
			continue
		}
		if c.Hidden {
			continue
		}

		s := rpad(c.Name()+":", c.NamePadding()) + c.Short

		commands = append(commands, s)
	}

	if len(commands) == 0 {
		commands = []string{}
	}

	type helpEntry struct {
		Title string
		Body  string
	}

	helpEntries := []helpEntry{}

	if command.Long != "" {
		helpEntries = append(helpEntries, helpEntry{"", command.Long})
	} else if command.Short != "" {
		helpEntries = append(helpEntries, helpEntry{"", command.Short})
	}

	helpEntries = append(helpEntries, helpEntry{"USAGE", command.UseLine()})

	if len(commands) > 0 {
		helpEntries = append(helpEntries, helpEntry{"COMMANDS", strings.Join(commands, "\n")})
	}

	flagUsages := command.LocalFlags().FlagUsages()

	if flagUsages != "" {
		helpEntries = append(helpEntries, helpEntry{"FLAGS", dedent(flagUsages)})
	}

	if _, ok := command.Annotations["help:arguments"]; ok {
		helpEntries = append(helpEntries, helpEntry{"ARGUMENTS", command.Annotations["help:arguments"]})
	}

	if command.Example != "" {
		helpEntries = append(helpEntries, helpEntry{"EXAMPLES", command.Example})
	}

	helpEntries = append(helpEntries, helpEntry{"LEARN MORE", `
Use 'tran <command> <subcommand> --help' for more information about a command.`})
	if _, ok := command.Annotations["help:tellus"]; ok {
		helpEntries = append(helpEntries, helpEntry{"TELL US", command.Annotations["help:tellus"]})
	}

	out := command.OutOrStdout()
	for _, e := range helpEntries {
		if e.Title != "" {
			fmt.Fprintln(out, cs.Bold(e.Title))
			fmt.Fprintln(out, tools.Indent(strings.Trim(e.Body, "\r\n"), "  "))
		} else {
			fmt.Fprintln(out, e.Body)
		}

		fmt.Fprintln(out)
	}
}

func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds ", padding)
	return fmt.Sprintf(template, s)
}

func dedent(s string) string {
	lines := strings.Split(s, "\n")
	minIndent := -1

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}

		indent := len(l) - len(strings.TrimLeft(l, " "))

		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return s
	}

	var buf bytes.Buffer

	for _, l := range lines {
		fmt.Fprintln(&buf, strings.TrimPrefix(l, strings.Repeat(" ", minIndent)))
	}

	return strings.TrimSuffix(buf.String(), "\n")
}
