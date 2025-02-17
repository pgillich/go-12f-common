package timer

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Config_GetConfigFlagSet(t *testing.T) {
	const EXPECTED_TIME_STEP_FROM_ENV_VAR = "env_10s"
	const EXPECTED_TIME_STEP_FROM_CLI_ARG = "cli_10s"

	envVars := map[string]string{
		"TIME_STEP": EXPECTED_TIME_STEP_FROM_ENV_VAR,
	}
	cliArgs := []string{
		fmt.Sprintf("--%v=%v", TIME_STEP_ARG_NAME, EXPECTED_TIME_STEP_FROM_CLI_ARG),
	}
	testCases := map[string]struct {
		expectedConfig Config
		envVars        map[string]string
		cliArgs        []string
	}{
		"default values": {
			expectedConfig: Config{
				TimeStep: TIME_STEP_DEFAULT,
			},
		},
		"from environment variables": {
			expectedConfig: Config{
				TimeStep: EXPECTED_TIME_STEP_FROM_ENV_VAR,
			},
			envVars: envVars,
		},
		"from cli args": {
			expectedConfig: Config{
				TimeStep: EXPECTED_TIME_STEP_FROM_CLI_ARG,
			},
			cliArgs: cliArgs,
		},
		"prefer cli args over env vars": {
			expectedConfig: Config{
				TimeStep: EXPECTED_TIME_STEP_FROM_CLI_ARG,
			},
			envVars: envVars,
			cliArgs: cliArgs,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			// given
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			cfg := &Config{}

			for k, v := range testCase.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// when
			cfg.GetConfigFlagSet(fs)
			require.NoError(t, fs.Parse(testCase.cliArgs))

			err := cfg.LoadConfig(fs)

			// then
			assert := assert.New(t)
			assert.NoError(err)
			assert.Equal(testCase.expectedConfig, *cfg)
		})
	}
}
