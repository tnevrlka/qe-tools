package webhook

import (
	"fmt"
	"github.com/redhat-appstudio/qe-tools/pkg/prow"
	"github.com/redhat-appstudio/qe-tools/pkg/types"
	"github.com/redhat-appstudio/qe-tools/pkg/webhook"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
	"os"
)

// AppStudio QE webhook configuration values will be used by default (if none are provided via env vars)
const (
	appstudioQESaltSecret       = "123456789"
	appstudioQEWebhookTargetURL = "https://hook.pipelinesascode.com/EyFYTakxEgEy"
)

var (
	openshiftJobSpec   *prow.OpenshiftJobSpec
	requiredParameters = []types.CmdParameter[string]{jobType, jobName, repoOwner, repoName, prNumber, saltSecret, webhookTargetUrl, jobSpec}
	jobType            = types.CmdParameter[string]{
		Name:  "job-type",
		Env:   "JOB_TYPE",
		Usage: "Type of the job",
	}
	jobName = types.CmdParameter[string]{
		Name:  "job-name",
		Env:   "JOB_NAME",
		Usage: "Name of the job",
	}
	repoOwner = types.CmdParameter[string]{
		Name:  "repo-owner",
		Env:   "REPO_OWNER",
		Usage: "Owner of the repository",
	}
	repoName = types.CmdParameter[string]{
		Name:  "repo-name",
		Env:   "REPO_OWNER",
		Usage: "Name of the repository",
	}
	prNumber = types.CmdParameter[string]{
		Name:  "pr-number",
		Env:   "PR_NUMBER",
		Usage: "Number of the pull request",
	}
	saltSecret = types.CmdParameter[string]{
		Name:         "salt-secret",
		Env:          "SALT_SECRET",
		DefaultValue: appstudioQESaltSecret,
		Usage:        "Salt for webhook config",
	}
	webhookTargetUrl = types.CmdParameter[string]{
		Name:         "target-url",
		Env:          "TARGET_URL",
		DefaultValue: appstudioQEWebhookTargetURL,
		Usage:        "Target URL for webhook",
	}
	jobSpec = types.CmdParameter[string]{
		Name:  "job-spec",
		Env:   "JOB_SPEC",
		Usage: "Job spec",
	}
)

var reportPortalWebhookCmd = &cobra.Command{
	Use: "report-portal",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		for _, reqParam := range requiredParameters {
			if reqParam.Value == "" {
				return fmt.Errorf("parameter '%s' and env var '%s' is empty", reqParam.Name, reqParam.Env)
			}
		}
		var err error
		openshiftJobSpec, err = prow.ParseJobSpec(jobSpec.Value)
		if err != nil {
			return fmt.Errorf("error parsing openshift job spec: %+v", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var repoURL string
		if jobType.Value == "periodic" {
			repoURL = "https://github.com/redhat-appstudio/infra-deployments"
			repoOwner.Value = "redhat-appstudio"
			repoName.Value = "infra-deployments"
			prNumber.Value = "periodic"
		} else if repoName.Value == "e2e-tests" || repoName.Value == "infra-deployments" {
			repoURL = openshiftJobSpec.Refs.RepoLink
		} else {
			klog.Infof("sending webhook for jobType %s, jobName %s is not supported", jobType, jobName)
			return nil
		}

		path, err := os.Executable()
		if err != nil {
			return fmt.Errorf("error when sending webhook: %+v", err)
		}

		wh := webhook.Webhook{
			Path: path,
			Repository: webhook.Repository{
				FullName:   fmt.Sprintf("%s/%s", repoOwner, repoName),
				PullNumber: prNumber.Value,
			},
			RepositoryURL: repoURL,
		}
		resp, err := wh.CreateAndSend(saltSecret.Value, webhookTargetUrl.Value)
		if err != nil {
			return fmt.Errorf("error sending webhook: %+v", err)
		}
		klog.Infof("webhook response: %+v", resp)

		return nil
	},
}

func init() {
	for _, reqParam := range requiredParameters {
		reportPortalWebhookCmd.Flags().StringVar(&reqParam.Value, reqParam.Name, reqParam.DefaultValue, reqParam.Usage)
		_ = viper.BindEnv(reqParam.Name, reqParam.Env)
		reqParam.Value = viper.GetString(reqParam.Name)
	}
}
