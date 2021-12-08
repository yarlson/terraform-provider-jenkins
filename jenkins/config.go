package jenkins

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	jenkins "github.com/bndr/gojenkins"
)

type jenkinsClient interface {
	CreateJobInFolder(ctx context.Context, config string, jobName string, parentIDs ...string) (*jenkins.Job, error)
	Credentials() *jenkins.CredentialsManager
	DeleteJobInFolder(ctx context.Context, name string, parentIDs ...string) (bool, error)
	GetJob(ctx context.Context, id string, parentIDs ...string) (*jenkins.Job, error)
	GetFolder(ctx context.Context, id string, parents ...string) (*jenkins.Folder, error)
	RegisterNode(ctx context.Context, name string, numExecutors int, description, remoteFS, credentials string, labels []string, ip string, port int) (*jenkins.Node, error)
	GetNode(ctx context.Context, name string) (*jenkins.Node, error)
	DeleteNode(ctx context.Context, name string) (bool, error)
}

// jenkinsAdapter wraps the Jenkins client, enabling additional functionality
type jenkinsAdapter struct {
	*jenkins.Jenkins
}

// Config is the set of parameters needed to configure the Jenkins provider.
type Config struct {
	ServerURL string
	CACert    io.Reader
	Username  string
	Password  string
}

func newJenkinsClient(c *Config) *jenkinsAdapter {
	client := jenkins.CreateJenkins(nil, c.ServerURL, c.Username, c.Password)
	if c.CACert != nil {
		// provide CA certificate if server is using self-signed certificate
		client.Requester.CACert, _ = ioutil.ReadAll(c.CACert)
	}

	// return the Jenkins API client
	return &jenkinsAdapter{Jenkins: client}
}

func (j *jenkinsAdapter) Credentials() *jenkins.CredentialsManager {
	return &jenkins.CredentialsManager{
		J: j.Jenkins,
	}
}

// DeleteJobInFolder assists in running DeleteJob funcs, as DeleteJob is not folder aware
// and cannot take a canonical job ID without mishandling it.
func (j *jenkinsAdapter) DeleteJobInFolder(ctx context.Context, name string, parentIDs ...string) (bool, error) {
	return j.DeleteJob(ctx, strings.Join(append(parentIDs, name), "/job/"))
}

func (j *jenkinsAdapter) RegisterNode(ctx context.Context, name string, numExecutors int, description, remoteFS, credentials string, labels []string, ip string, port int) (*jenkins.Node, error) {
	return j.CreateNode(
		ctx,
		name,
		numExecutors,
		description,
		remoteFS,
		strings.Join(labels, ", "),
		map[string]string{
			"method":                         "SSHLauncher",
			"host":                           ip,
			"sshHostKeyVerificationStrategy": jenkins.NonVerifyingKey,
			"javaPath":                       "/usr/bin/java",
			"credentialsId":                  credentials,
			"port":                           fmt.Sprintf("%d", port),
			"mode":                           "EXCLUSIVE",
		},
		map[string]string{
			"mode": "EXCLUSIVE",
		},
	)
}
