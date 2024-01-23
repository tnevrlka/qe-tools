package webhook

import "github.com/spf13/cobra"

var WebhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Command for triggering webhooks",
}

func init() {
	WebhookCmd.AddCommand(reportPortalWebhookCmd)
}
