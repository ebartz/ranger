package prime

import (
	"testing"

	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"github.com/ranger/ranger/tests/framework/extensions/clusters"
	prime "github.com/ranger/ranger/tests/framework/extensions/prime"
	"github.com/ranger/ranger/tests/framework/extensions/rangerversion"
	"github.com/ranger/ranger/tests/framework/pkg/config"
	"github.com/ranger/ranger/tests/framework/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	systemRegistry = "system-default-registry"
	localCluster   = "local"
	uiBrand        = "ui-brand"
)

type PrimeTestSuite struct {
	suite.Suite
	session        *session.Session
	client         *ranger.Client
	brand          string
	isPrime        bool
	rangerVersion string
	primeRegistry  string
}

func (t *PrimeTestSuite) TearDownSuite() {
	t.session.Cleanup()
}

func (t *PrimeTestSuite) SetupSuite() {
	testSession := session.NewSession()
	t.session = testSession

	primeConfig := new(rangerversion.Config)
	config.LoadConfig(rangerversion.ConfigurationFileKey, primeConfig)

	t.brand = primeConfig.Brand
	t.isPrime = primeConfig.IsPrime
	t.rangerVersion = primeConfig.RangerVersion
	t.primeRegistry = primeConfig.Registry

	client, err := ranger.NewClient("", t.session)
	assert.NoError(t.T(), err)

	t.client = client
}

func (t *PrimeTestSuite) TestPrimeUIBrand() {
	rangerBrand, err := t.client.Management.Setting.ByID(uiBrand)
	require.NoError(t.T(), err)

	checkBrand := prime.CheckUIBrand(t.client, t.isPrime, rangerBrand, t.brand)
	assert.NoError(t.T(), checkBrand)
}

func (t *PrimeTestSuite) TestPrimeVersion() {
	serverConfig, err := rangerversion.RequestRangerVersion(t.client.RangerConfig.Host)
	require.NoError(t.T(), err)

	checkVersion := prime.CheckVersion(t.isPrime, t.rangerVersion, serverConfig)
	assert.NoError(t.T(), checkVersion)
}

func (t *PrimeTestSuite) TestSystemDefaultRegistry() {
	registry, err := t.client.Management.Setting.ByID(systemRegistry)
	require.NoError(t.T(), err)

	checkRegistry := prime.CheckSystemDefaultRegistry(t.isPrime, t.primeRegistry, registry)
	assert.NoError(t.T(), checkRegistry)
}

func (t *PrimeTestSuite) TestLocalClusterRangerImages() {
	adminClient, err := ranger.NewClient(t.client.RangerConfig.AdminToken, t.client.Session)
	require.NoError(t.T(), err)

	clusterID, err := clusters.GetClusterIDByName(adminClient, localCluster)
	require.NoError(t.T(), err)

	imageResults, imageErrors := prime.CheckLocalClusterRangerImages(t.client, t.isPrime, t.rangerVersion, t.primeRegistry, clusterID)
	assert.NotEmpty(t.T(), imageResults)
	assert.Empty(t.T(), imageErrors)
}

func TestPrimeTestSuite(t *testing.T) {
	suite.Run(t, new(PrimeTestSuite))
}
