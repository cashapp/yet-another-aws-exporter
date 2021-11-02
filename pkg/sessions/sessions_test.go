package sessions

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func restoreEnvFn(envVars map[string]string) {
	for k := range envVars {
		os.Unsetenv(k)
	}
}

// TestCreateAWSSession calls the function which initializes
// an AWS session from the local environment.
func TestCreateAWSSession(t *testing.T) {
	cases := map[string]struct {
		InEnvs map[string]string
	}{
		"with an env var that breaks new session creation": {
			InEnvs: map[string]string{
				"AWS_STS_REGIONAL_ENDPOINTS": "fake",
			},
		},
		"returns session with shared config": {
			InEnvs: map[string]string{},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			defer restoreEnvFn(c.InEnvs)

			for k, v := range c.InEnvs {
				os.Setenv(k, v)
			}

			sess, err := CreateAWSSession()
			if err != nil {
				assert.NotNil(t, err)
				return
			}
			assert.NotNil(t, sess)
		})
	}
}
