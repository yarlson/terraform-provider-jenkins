package jenkins

import (
	"bytes"
	"context"
	"testing"

	jenkins "github.com/bndr/gojenkins"
)

type mockJenkinsClient struct {
	mockCreateJobInFolder func(ctx context.Context, config string, jobName string, parentIDs ...string) (*jenkins.Job, error)
	mockDeleteJobInFolder func(ctx context.Context, name string, parentIDs ...string) (bool, error)
	mockGetJob            func(ctx context.Context, id string, parentIDs ...string) (*jenkins.Job, error)
	mockGetFolder         func(ctx context.Context, id string, parentIDs ...string) (*jenkins.Folder, error)
	mockRegisterNode      func(ctx context.Context, name string, numExecutors int, description, remoteFS, credentials string, labels []string, ip string, port int) (*jenkins.Node, error)
	mockGetNode           func(ctx context.Context, name string) (*jenkins.Node, error)
	mockDeleteNode        func(ctx context.Context, name string) (bool, error)
}

func (m *mockJenkinsClient) CreateJobInFolder(ctx context.Context, config string, jobName string, parentIDs ...string) (*jenkins.Job, error) {
	return m.mockCreateJobInFolder(ctx, config, jobName, parentIDs...)
}

func (m *mockJenkinsClient) Credentials() *jenkins.CredentialsManager {
	return &jenkins.CredentialsManager{}
}

func (m *mockJenkinsClient) DeleteJobInFolder(ctx context.Context, name string, parentIDs ...string) (bool, error) {
	return m.mockDeleteJobInFolder(ctx, name, parentIDs...)
}

func (m *mockJenkinsClient) GetJob(ctx context.Context, id string, parentIDs ...string) (*jenkins.Job, error) {
	return m.mockGetJob(ctx, id, parentIDs...)
}

func (m *mockJenkinsClient) GetFolder(ctx context.Context, id string, parentIDs ...string) (*jenkins.Folder, error) {
	return m.mockGetFolder(ctx, id, parentIDs...)
}

func (m *mockJenkinsClient) RegisterNode(ctx context.Context, name string, numExecutors int, description, remoteFS, credentials string, labels []string, ip string, port int) (*jenkins.Node, error) {
	return m.mockRegisterNode(ctx, name, numExecutors, description, remoteFS, credentials, labels, ip, port)
}

func (m *mockJenkinsClient) GetNode(ctx context.Context, name string) (*jenkins.Node, error) {
	return m.mockGetNode(ctx, name)
}

func (m *mockJenkinsClient) DeleteNode(ctx context.Context, name string) (bool, error) {
	return m.mockDeleteNode(ctx, name)
}

func TestNewJenkinsClient(t *testing.T) {
	c := newJenkinsClient(&Config{})
	if c == nil {
		t.Errorf("Expected populated client")
	}

	c = newJenkinsClient(&Config{
		CACert: bytes.NewBufferString("certificate"),
	})
	if string(c.Requester.CACert) != "certificate" {
		t.Errorf("Initialization did not extract certificate data")
	}
}

func TestJenkinsAdapter_Credentials(t *testing.T) {
	c := newJenkinsClient(&Config{})
	cm := c.Credentials()

	if cm == nil {
		t.Errorf("Expected populated client")
	} else if cm.J != c.Jenkins {
		t.Error("Expected credentials client to match client")
	}
}
