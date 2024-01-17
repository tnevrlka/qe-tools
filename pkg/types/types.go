package types

const (
	ArtifactDirEnv string = "ARTIFACT_DIR"
	GithubTokenEnv string = "GITHUB_TOKEN" // #nosec G101
	ProwJobIDEnv   string = "PROW_JOB_ID"

	ArtifactDirParamName string = "artifact-dir"
	ProwJobIDParamName   string = "prow-job-id"

	JunitFilename string = `/(j?unit|e2e).*\.xml`
)

type CmdParameter[T any] struct {
	Name         string
	Env          string
	DefaultValue T
	Value        T
	Usage        string
}
