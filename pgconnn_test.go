package pgconn

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/xtracdev/envinject"
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
			map[string]string{DBUser: "user", DBPassword: "secretpassword", DBHost: "host", DBPort: "port", DBName: "svc"},
			[]string{},
			false,
			"user=user password=secretpassword dbname=svc host=host port=port sslmode=disable", "user=user password=XXX dbname=svc host=host port=port sslmode=disable",
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
			//config, err := NewEnvConfig()
			env, _ := envinject.NewInjectedEnv()
			connectString, err := ConnectStringFromInjectedEnv(env)
			masked, maskedErr := MaskedConnectStringFromInjectedEnv(env)
			if test.expectError {
				assert.NotNil(t, err, "expected error")
				assert.NotNil(t, maskedErr, "expected masked error")
				errString := err.Error()
				for _, ec := range test.errorComponents {
					assert.True(t, strings.Contains(errString, ec))
				}
			} else {
				assert.Equal(t, test.connectString, connectString)
				assert.False(t, strings.Contains(masked, "secret"))
			}
		})
	}
}

func TestNilEnvProducesError(t *testing.T) {
	_, err := ConnectStringFromInjectedEnv(nil)
	assert.NotNil(t, err)

	_, err = MaskedConnectStringFromInjectedEnv(nil)
	assert.NotNil(t, err)
}

func TestGetIntFromEnv(t *testing.T) {
	var readIntTests = []struct {
		testName   string
		varName    string
		varValue   string
		defaultVal int
		expected   int
	}{
		{
			"read from environment",
			"FOO",
			"123",
			1,
			123,
		},
		{
			"no value",
			"FOO",
			"",
			1,
			1,
		},
		{
			"malformed value",
			"FOO",
			"not an integer",
			456,
			456,
		},
	}

	for _, test := range readIntTests {
		t.Run(test.testName, func(t *testing.T) {
			os.Setenv(test.varName, test.varValue)
			env, _ := envinject.NewInjectedEnv()
			pgdbShell := PostgresDB{nil, env}
			v := pgdbShell.getIntFromEnv(test.varName, test.defaultVal)

			assert.Equal(t, v, test.expected)
		})
	}
}

func TestGetDefaultMaxConnections(t *testing.T) {
	os.Unsetenv(maxConns)
	env, _ := envinject.NewInjectedEnv()
	pgdbShell := PostgresDB{nil, env}
	max := pgdbShell.getMaxConns()
	assert.Equal(t, defaultMaxConns, max)
}

func TestGetDefaultIdleConnections(t *testing.T) {
	os.Unsetenv(idleConns)
	env, _ := envinject.NewInjectedEnv()
	pgdbShell := PostgresDB{nil, env}
	max := pgdbShell.getIdleConns()
	assert.Equal(t, defaultIdleConns, max)
}

func TestConnectionErrorUmTest(t *testing.T) {
	assert.True(t, IsConnectionError(errors.New("connection refused")))
	assert.False(t, IsConnectionError(errors.New("something went wrong")))
}
