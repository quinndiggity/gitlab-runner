package commands_helpers

import (
	"io"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"errors"
	"gitlab.com/gitlab-org/gitlab-ci-multi-runner/common"
	"gitlab.com/gitlab-org/gitlab-ci-multi-runner/helpers/archives"
	"gitlab.com/gitlab-org/gitlab-ci-multi-runner/helpers/formatter"
	"gitlab.com/gitlab-org/gitlab-ci-multi-runner/network"
)

type ArtifactsUploaderCommand struct {
	common.BuildCredentials
	fileArchiver
	retryHelper
	network common.Network
}

func (c *ArtifactsUploaderCommand) createAndUpload() (bool, error) {
	pr, pw := io.Pipe()
	defer pr.Close()

	// Create the archive
	go func() {
		err := archives.CreateZipArchive(pw, c.sortedFiles())
		pw.CloseWithError(err)
	}()

	// Upload the data
	switch c.network.UploadRawArtifacts(c.BuildCredentials, pr, "artifacts.zip") {
	case common.UploadSucceeded:
		return false, nil
	case common.UploadForbidden:
		return false, os.ErrPermission
	case common.UploadTooLarge:
		return false, errors.New("Too large")
	case common.UploadFailed:
		return true, os.ErrInvalid
	default:
		return false, os.ErrInvalid
	}
}

func (c *ArtifactsUploaderCommand) Execute(*cli.Context) {
	formatter.SetRunnerFormatter()

	if len(c.URL) == 0 || len(c.Token) == 0 {
		logrus.Fatalln("Missing runner credentials")
	}
	if c.ID <= 0 {
		logrus.Fatalln("Missing build ID")
	}

	// Enumerate files
	err := c.enumerate()
	if err != nil {
		logrus.Fatalln(err)
	}

	// If the upload fails, exit with a non-zero exit code to indicate an issue?
	err = c.doRetry(c.createAndUpload)
	if err != nil {
		logrus.Fatalln(err)
	}
}

func init() {
	common.RegisterCommand2("artifacts-uploader", "create and upload build artifacts (internal)", &ArtifactsUploaderCommand{
		network: &network.GitLabClient{},
		retryHelper: retryHelper{
			Retry:     2,
			RetryTime: time.Second,
		},
	})
}
