package prowjob

import (
	"github.com/spf13/cobra"
)

const (
	artifactDirEnv       string = "ARTIFACT_DIR"
	artifactDirParamName string = "artifact-dir"

	prowJobIDEnv       string = "PROW_JOB_ID"
	prowJobIDParamName string = "prow-job-id"
)

var (
	artifactDir string
	prowJobID   string
)

// ProwjobCmd represents the prowjob command
var ProwjobCmd = &cobra.Command{
	Use:   "prowjob",
	Short: "Commands for processing Prow jobs",
}

func init() {
	ProwjobCmd.AddCommand(periodicSlackReportCmd)
	ProwjobCmd.AddCommand(createReportCmd)
}
