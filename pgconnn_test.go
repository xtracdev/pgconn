package pgconn

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func cleanEnv() {
	os.Unsetenv(DBUser)
	os.Unsetenv(DBPassword)
	os.Unsetenv(DBHost)
	os.Unsetenv(DBPort)
	os.Unsetenv(DBName)
}

func setEnv(config map[string]string) {
	for k, v := range config {
		os.Setenv(k, v)
	}
}

func TestConfig(t *testing.T) {
	var configTests = []struct {
		testName            string
		config              map[string]string
		errorComponents     []string
		expectError         bool
		connectString       string
		maskedConnectString string
	}{
		{
			"all environment present",
			map[string]string{DBUser: "user", DBPassword: "password", DBHost: "host", DBPort: "port", DBName: "svc"},
			[]string{},
			false,
			"user=user password=password dbname=svc host=host port=port sslmode=disable", "user=user password=XXX dbname=svc host=host port=port sslmode=disable",
		},
		{
			"no environment present",
			map[string]string{},
			[]string{DBUser, DBPassword, DBHost, DBPort, DBName},
			true,
			"", "",
		},
		{
			"some environment present",
			map[string]string{DBUser: "user", DBPort: "port", DBName: "svc"},
			[]string{DBPassword, DBHost},
			true,
			"", "",
		},
	}

	for _, test := range configTests {
		t.Run(test.testName, func(t *testing.T) {
			cleanEnv()
			setEnv(test.config)
			config, err := NewEnvConfig()
			if test.expectError {
				assert.NotNil(t, err, "expected error")
				errString := err.Error()
				for _, ec := range test.errorComponents {
					assert.True(t, strings.Contains(errString, ec))
				}
			} else {
				assert.Equal(t, test.connectString, config.ConnectString())
			}
		})
	}
}
