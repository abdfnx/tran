package tran

import (
	"github.com/spf13/cobra"
)

var HelpTopics = map[string]map[string]string{}

func NewHelpTopic(topic string) *cobra.Command {
	cmd := &cobra.Command{
		Use:    topic,
		Short:  HelpTopics[topic]["short"],
		Long:   HelpTopics[topic]["long"],
		Hidden: true,
		Annotations: map[string]string{
			"markdown:generate": "true",
			"markdown:basename": "tran_help_" + topic,
		},
	}

	cmd.SetHelpFunc(helpTopicHelpFunc)
	cmd.SetUsageFunc(helpTopicUsageFunc)

	return cmd
}

func helpTopicHelpFunc(command *cobra.Command, args []string) {
	command.Print(command.Long)
}

func helpTopicUsageFunc(command *cobra.Command) error {
	command.Printf("Usage: tran help %s", command.Use)
	return nil
}
