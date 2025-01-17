package token

import (
	"testing"

	"github.com/ranger/ranger/tests/framework/clients/ranger"
	fv3 "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/kubeapi/tokens"
	"github.com/ranger/ranger/tests/framework/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	initialTokenDesc = "my-token"
	updatedTokenDesc = "changed-token"
)

type TokenTestSuite struct {
	suite.Suite
	client  *ranger.Client
	session *session.Session
	cluster *management.Cluster
}

func (t *TokenTestSuite) TearDownSuite() {
	t.session.Cleanup()
}

func (t *TokenTestSuite) SetupSuite() {
	testSession := session.NewSession()
	t.session = testSession

	client, err := ranger.NewClient("", t.session)
	require.NoError(t.T(), err)

	t.client = client
}

func (t *TokenTestSuite) TestPatchToken() {
	tokenToCreate := &fv3.Token{Description: initialTokenDesc}
	createdToken, err := t.client.Management.Token.Create(tokenToCreate)
	require.NoError(t.T(), err)

	assert.Equal(t.T(), initialTokenDesc, createdToken.Description)

	patchedToken, unstructuredRes, err := tokens.PatchToken(t.client, t.client.RangerConfig.ClusterName, createdToken.Name, "replace", "/description", updatedTokenDesc)
	require.NoError(t.T(), err)

	assert.Equal(t.T(), updatedTokenDesc, patchedToken.Description)

	uc := unstructuredRes.UnstructuredContent()
	if val, ok := uc["groupPrincipals"]; ok {
		assert.NotEmpty(t.T(), val)
	}
}

func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}
