package terraform

import (
	"github.com/gruntwork-io/terratest/modules/shell"
	"testing"
	"github.com/gruntwork-io/terratest/modules/retry"
	"strings"
	"github.com/gruntwork-io/terratest/modules/logger"
	"fmt"
)

// Run terraform with the given arguments and options and return stdout/stderr
func RunTerraformCommand(t *testing.T, options *Options, args ... string) string {
	out, err := RunTerraformCommandE(t, options, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// Run terraform with the given arguments and options and return stdout/stderr
func RunTerraformCommandE(t *testing.T, options *Options, args ... string) (string, error) {
	description := fmt.Sprintf("Running terraform %v", args)
	return retry.DoWithRetryE(t, description, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		cmd := shell.Command{
			Command:    "terraform",
			Args:       args,
			WorkingDir: options.TerraformDir,
			Env:        options.EnvVars,
		}

		out, err := shell.RunCommandAndGetOutputE(t, cmd)
		if err == nil {
			return out, nil
		}

		for errorText, errorMessage := range options.RetryableTerraformErrors {
			if strings.Contains(err.Error(), errorText) {
				logger.Logf(t, "terraform failed with the error '%s' but this error was expected and warrants a retry. Further details: %s\n", errorText, errorMessage)
				return "", err
			}
		}

		return "", retry.FatalError{Underlying: err}
	})
}
