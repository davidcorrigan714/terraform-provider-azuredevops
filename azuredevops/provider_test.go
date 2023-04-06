package azuredevops_test

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops"
	mock_azuredevops "github.com/microsoft/terraform-provider-azuredevops/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider_HasChildResources(t *testing.T) {
	expectedResources := []string{
		"azuredevops_resource_authorization",
		"azuredevops_build_definition",
		"azuredevops_build_definition_permissions",
		"azuredevops_branch_policy_build_validation",
		"azuredevops_branch_policy_min_reviewers",
		"azuredevops_branch_policy_auto_reviewers",
		"azuredevops_branch_policy_work_item_linking",
		"azuredevops_branch_policy_comment_resolution",
		"azuredevops_branch_policy_merge_types",
		"azuredevops_branch_policy_status_check",
		"azuredevops_library_permissions",
		"azuredevops_project",
		"azuredevops_project_features",
		"azuredevops_project_pipeline_settings",
		"azuredevops_check_branch_control",
		"azuredevops_check_business_hours",
		"azuredevops_serviceendpoint_github",
		"azuredevops_serviceendpoint_github_enterprise",
		"azuredevops_serviceendpoint_dockerregistry",
		"azuredevops_serviceendpoint_azuredevops",
		"azuredevops_serviceendpoint_azurerm",
		"azuredevops_serviceendpoint_azurecr",
		"azuredevops_serviceendpoint_runpipeline",
		"azuredevops_serviceendpoint_bitbucket",
		"azuredevops_serviceendpoint_kubernetes",
		"azuredevops_serviceendpoint_servicefabric",
		"azuredevops_serviceendpoint_argocd",
		"azuredevops_serviceendpoint_aws",
		"azuredevops_serviceendpoint_artifactory",
		"azuredevops_serviceendpoint_sonarqube",
		"azuredevops_serviceendpoint_sonarcloud",
		"azuredevops_serviceendpoint_ssh",
		"azuredevops_serviceendpoint_npm",
		"azuredevops_serviceendpoint_nuget",
		"azuredevops_serviceendpoint_generic",
		"azuredevops_serviceendpoint_generic_git",
		"azuredevops_serviceendpoint_octopusdeploy",
		"azuredevops_serviceendpoint_incomingwebhook",
		"azuredevops_serviceendpoint_jfrog_artifactory_v2",
		"azuredevops_serviceendpoint_jfrog_distribution_v2",
		"azuredevops_serviceendpoint_jfrog_platform_v2",
		"azuredevops_serviceendpoint_jfrog_xray_v2",
		"azuredevops_serviceendpoint_externaltfs",
		"azuredevops_variable_group",
		"azuredevops_repository_policy_author_email_pattern",
		"azuredevops_repository_policy_case_enforcement",
		"azuredevops_repository_policy_file_path_pattern",
		"azuredevops_repository_policy_max_file_size",
		"azuredevops_repository_policy_max_path_length",
		"azuredevops_repository_policy_reserved_names",
		"azuredevops_repository_policy_check_credentials",
		"azuredevops_git_repository",
		"azuredevops_git_repository_branch",
		"azuredevops_git_repository_file",
		"azuredevops_user_entitlement",
		"azuredevops_group_membership",
		"azuredevops_group",
		"azuredevops_agent_pool",
		"azuredevops_agent_queue",
		"azuredevops_project_permissions",
		"azuredevops_git_permissions",
		"azuredevops_workitemquery_permissions",
		"azuredevops_area_permissions",
		"azuredevops_iteration_permissions",
		"azuredevops_team",
		"azuredevops_team_members",
		"azuredevops_team_administrators",
		"azuredevops_serviceendpoint_permissions",
		"azuredevops_servicehook_permissions",
		"azuredevops_variable_group_permissions",
		"azuredevops_tagging_permissions",
		"azuredevops_environment",
		"azuredevops_build_folder",
		"azuredevops_build_folder_permissions",
		"azuredevops_workitem",
	}

	resources := azuredevops.Provider().ResourcesMap
	require.Equal(t, len(expectedResources), len(resources), "There are an unexpected number of registered resources")

	for _, resource := range expectedResources {
		require.Contains(t, resources, resource, "An expected resource was not registered")
		require.NotNil(t, resources[resource], "A resource cannot have a nil schema")
	}
}

func TestProvider_HasChildDataSources(t *testing.T) {
	expectedDataSources := []string{
		"azuredevops_build_definition",
		"azuredevops_client_config",
		"azuredevops_group",
		"azuredevops_project",
		"azuredevops_projects",
		"azuredevops_git_repositories",
		"azuredevops_git_repository",
		"azuredevops_users",
		"azuredevops_agent_pool",
		"azuredevops_agent_pools",
		"azuredevops_agent_queue",
		"azuredevops_area",
		"azuredevops_iteration",
		"azuredevops_team",
		"azuredevops_teams",
		"azuredevops_groups",
		"azuredevops_variable_group",
		"azuredevops_serviceendpoint_azurerm",
		"azuredevops_serviceendpoint_github",
	}

	dataSources := azuredevops.Provider().DataSourcesMap
	require.Equal(t, len(expectedDataSources), len(dataSources), "There are an unexpected number of registered data sources")

	for _, resource := range expectedDataSources {
		require.Contains(t, dataSources, resource, "An expected data source was not registered")
		require.NotNil(t, dataSources[resource], "A data source cannot have a nil schema")
	}
}

func TestProvider_SchemaIsValid(t *testing.T) {
	type testParams struct {
		name          string
		required      bool
		defaultEnvVar string
		sensitive     bool
	}

	tests := []testParams{
		{"org_service_url", false, "AZDO_ORG_SERVICE_URL", false},
		{"personal_access_token", false, "AZDO_PERSONAL_ACCESS_TOKEN", true},
		{"sp_client_id", false, "AZDO_SP_CLIENT_ID", false},
		{"sp_tenant_id", false, "AZDO_SP_TENANT_ID", false},
		{"sp_client_id_plan", false, "AZDO_SP_CLIENT_ID_PLAN", false},
		{"sp_tenant_id_plan", false, "AZDO_SP_TENANT_ID_PLAN", false},
		{"sp_client_id_apply", false, "AZDO_SP_CLIENT_ID_APPLY", false},
		{"sp_tenant_id_apply", false, "AZDO_SP_TENANT_ID_APPLY", false},
		{"sp_client_secret", false, "AZDO_SP_CLIENT_SECRET", true},
		{"sp_client_secret_path", false, "AZDO__SP_CLIENT_SECRET_PATH", false},
		{"sp_oidc_token", false, "AZDO_SP_OIDC_TOKEN", true},
		{"sp_oidc_token_path", false, "AZDO_SP_OIDC_TOKEN_PATH", false},
		{"sp_oidc_github_actions", false, "AZDO_SP_OIDC_GITHUB_ACTIONS", false},
		{"sp_oidc_github_actions_audience", false, "AZDO_SP_OIDC_GITHUB_ACTIONS_AUDIENCE", false},
		{"sp_oidc_hcp", false, "AZDO_SP_OIDC_HCP", false},
		{"sp_client_certificate_path", false, "AZDO_SP_CLIENT_CERTIFICATE_PATH", false},
		{"sp_client_certificate", false, "AZDO_SP_CLIENT_CERTIFICATE", true},
		{"sp_client_certificate_password", false, "AZDO_SP_CLIENT_CERTIFICATE_PASSWORD", true},
	}

	schema := azuredevops.Provider().Schema
	require.Equal(t, len(tests), len(schema), "There are an unexpected number of properties in the schema")

	for _, test := range tests {
		require.Contains(t, schema, test.name, "An expected property was not found in the schema")
		require.NotNil(t, schema[test.name], "A property in the schema cannot have a nil value")
		require.Equal(t, test.sensitive, schema[test.name].Sensitive, "A property in the schema has an incorrect sensitivity value")
		require.Equal(t, test.required, schema[test.name].Required, "A property in the schema has an incorrect required value")

		if test.defaultEnvVar != "" {
			expectedValue := os.Getenv(test.defaultEnvVar)

			actualValue, err := schema[test.name].DefaultFunc()
			if actualValue == nil {
				actualValue = ""
			}

			require.Nil(t, err, "An error occurred when getting the default value from the environment")
			require.Equal(t, expectedValue, actualValue, "The default value pulled from the environment has the wrong value")
		}
	}
}

func TestAuthPAT(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockAzIdentityFuncs(ctrl)

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("personal_access_token", "test123")

	resp, err := azuredevops.GetAuthToken(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	assert.Equal(t, "test123", resp)
}

type simpleTokenGetter struct {
	token string
}

func (s simpleTokenGetter) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     s.token,
		ExpiresOn: time.Now(),
	}, nil
}

func TestAuthOIDCToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockAzIdentityFuncs(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	oidcToken := "buffalo123"
	accessToken := "thepassword"

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("sp_client_id", clientId)
	resourceData.Set("sp_tenant_id", tenantId)
	resourceData.Set("sp_oidc_token", oidcToken)

	mockIdentityClient.EXPECT().NewClientAssertionCredential(tenantId, clientId, gomock.Any(), nil).DoAndReturn(
		func(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (*simpleTokenGetter, error) {
			assertion, err := getAssertion(context.Background())
			assert.Nil(t, err)
			assert.Equal(t, oidcToken, assertion)
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := azuredevops.GetAuthToken(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	assert.Equal(t, accessToken, resp)
}

func TestAuthOIDCTokenFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockAzIdentityFuncs(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	oidcToken := "buffalo123"
	tempFile := t.TempDir() + "/clientSecret.txt"
	err := os.WriteFile(tempFile, []byte(oidcToken), 0644)
	assert.Nil(t, err)

	accessToken := "thepassword"

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("sp_client_id", clientId)
	resourceData.Set("sp_tenant_id", tenantId)
	resourceData.Set("sp_oidc_token_path", tempFile)

	mockIdentityClient.EXPECT().NewClientAssertionCredential(tenantId, clientId, gomock.Any(), nil).DoAndReturn(
		func(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (*simpleTokenGetter, error) {
			assertion, err := getAssertion(context.Background())
			assert.Nil(t, err)
			assert.Equal(t, oidcToken, assertion)
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := azuredevops.GetAuthToken(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	assert.Equal(t, accessToken, resp)
}

func TestAuthClientSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockAzIdentityFuncs(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	clientSecret := "buffalo123"
	accessToken := "thepassword"

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("sp_client_id", clientId)
	resourceData.Set("sp_tenant_id", tenantId)
	resourceData.Set("sp_client_secret", clientSecret)

	mockIdentityClient.EXPECT().NewClientSecretCredential(tenantId, clientId, gomock.Any(), nil).DoAndReturn(
		func(tenantID, clientID, secret string, options *azidentity.ClientSecretCredentialOptions) (*simpleTokenGetter, error) {
			assert.Equal(t, clientSecret, secret)
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := azuredevops.GetAuthToken(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	assert.Equal(t, accessToken, resp)
}

func TestAuthClientSecretFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockAzIdentityFuncs(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	clientSecret := "buffalo123"
	tempFile := t.TempDir() + "/clientSecret.txt"
	err := os.WriteFile(tempFile, []byte(clientSecret), 0644)
	assert.Nil(t, err)

	accessToken := "thepassword"

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("sp_client_id", clientId)
	resourceData.Set("sp_tenant_id", tenantId)
	resourceData.Set("sp_client_secret_path", tempFile)

	mockIdentityClient.EXPECT().NewClientSecretCredential(tenantId, clientId, gomock.Any(), nil).DoAndReturn(
		func(tenantID, clientID, secret string, options *azidentity.ClientSecretCredentialOptions) (*simpleTokenGetter, error) {
			assert.Equal(t, clientSecret, secret)
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := azuredevops.GetAuthToken(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	assert.Equal(t, accessToken, resp)
}

func TestAuthTrfm(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockAzIdentityFuncs(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	fakeTokenValue := "tokenvalue"
	os.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN", fakeTokenValue)
	accessToken := "thepassword"

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("sp_client_id", clientId)
	resourceData.Set("sp_tenant_id", tenantId)
	resourceData.Set("sp_oidc_hcp", true)

	mockIdentityClient.EXPECT().NewClientAssertionCredential(tenantId, clientId, gomock.Any(), nil).DoAndReturn(
		func(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (*simpleTokenGetter, error) {
			assertion, err := getAssertion(context.Background())
			assert.Nil(t, err)
			assert.Equal(t, fakeTokenValue, assertion)
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := azuredevops.GetAuthToken(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	assert.Equal(t, accessToken, resp)
}

func TestAuthTrfmPlanApply(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockAzIdentityFuncs(ctrl)
	clientId_apply := "00000000-0000-0000-0000-000000000003"
	tenantId_apply := "00000000-0000-0000-0000-000000000004"
	clientId_plan := "00000000-0000-0000-0000-000000000005"
	tenantId_plan := "00000000-0000-0000-0000-000000000006"
	trfm_fake_token_plan := fmt.Sprintf("header.%s.signature", base64.StdEncoding.EncodeToString([]byte("{\"terraform_run_phase\":\"plan\"}")))
	trfm_fake_token_apply := fmt.Sprintf("header.%s.signature", base64.StdEncoding.EncodeToString([]byte("{\"terraform_run_phase\":\"apply\"}")))
	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	accessToken := "thepassword"
	resourceData.Set("sp_client_id_apply", clientId_apply)
	resourceData.Set("sp_tenant_id_apply", tenantId_apply)
	resourceData.Set("sp_client_id_plan", clientId_plan)
	resourceData.Set("sp_tenant_id_plan", tenantId_plan)
	resourceData.Set("sp_oidc_hcp", true)

	// Apply phase test
	os.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN", trfm_fake_token_apply)
	mockIdentityClient.EXPECT().NewClientAssertionCredential(tenantId_apply, clientId_apply, gomock.Any(), nil).DoAndReturn(
		func(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (*simpleTokenGetter, error) {
			assertion, err := getAssertion(context.Background())
			assert.Nil(t, err)
			assert.Equal(t, trfm_fake_token_apply, assertion)
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := azuredevops.GetAuthToken(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	assert.Equal(t, accessToken, resp)

	// Plan phase test
	os.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN", trfm_fake_token_plan)
	mockIdentityClient.EXPECT().NewClientAssertionCredential(tenantId_plan, clientId_plan, gomock.Any(), nil).DoAndReturn(
		func(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (*simpleTokenGetter, error) {
			assertion, err := getAssertion(context.Background())
			assert.Nil(t, err)
			assert.Equal(t, trfm_fake_token_plan, assertion)
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err = azuredevops.GetAuthToken(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	assert.Equal(t, accessToken, resp)
}

func generateCert() []byte {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)

	template := x509.Certificate{
		SerialNumber: new(big.Int).SetUint64(20),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		DNSNames:  []string{"localhost"},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Minute * 5),
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("Failed to create private key: %v", err)
	}

	publicBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	privateBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	return append(publicBytes[:], privateBytes[:]...)
}

func TestAuthClientCert(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockAzIdentityFuncs(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	cert := generateCert()
	accessToken := "thepassword"

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("sp_client_id", clientId)
	resourceData.Set("sp_tenant_id", tenantId)
	resourceData.Set("sp_client_certificate", base64.StdEncoding.EncodeToString(cert))

	theseCerts, theseKey, err := azidentity.ParseCertificates(cert, nil)
	assert.Nil(t, err)

	mockIdentityClient.EXPECT().NewClientCertificateCredential(tenantId, clientId, gomock.Any(), gomock.Any(), nil).DoAndReturn(
		func(tenantID string, clientID string, certs []*x509.Certificate, key crypto.PrivateKey, options *azidentity.ClientCertificateCredentialOptions) (*simpleTokenGetter, error) {
			assert.Equal(t, theseCerts, certs)
			assert.Equal(t, theseKey, key)
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := azuredevops.GetAuthToken(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	assert.Equal(t, accessToken, resp)
}

func TestAuthClientCertFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockAzIdentityFuncs(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	cert := generateCert()
	accessToken := "thepassword"
	tempFile := t.TempDir() + "/clientCerts.pem"
	err := os.WriteFile(tempFile, []byte(cert), 0644)

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("sp_client_id", clientId)
	resourceData.Set("sp_tenant_id", tenantId)
	resourceData.Set("sp_client_certificate_path", tempFile)

	theseCerts, theseKey, err := azidentity.ParseCertificates(cert, nil)
	assert.Nil(t, err)

	mockIdentityClient.EXPECT().NewClientCertificateCredential(tenantId, clientId, gomock.Any(), gomock.Any(), nil).DoAndReturn(
		func(tenantID string, clientID string, certs []*x509.Certificate, key crypto.PrivateKey, options *azidentity.ClientCertificateCredentialOptions) (*simpleTokenGetter, error) {
			assert.Equal(t, theseCerts, certs)
			assert.Equal(t, theseKey, key)
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := azuredevops.GetAuthToken(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	assert.Equal(t, accessToken, resp)
}

func TestGHActionsNoAudience(t *testing.T) {
	testCases := []struct {
		testAudience     string
		expectedAudience string
	}{
		{
			testAudience:     "",
			expectedAudience: "api://AzureADTokenExchange",
		},
		{
			testAudience:     "my-test-audience",
			expectedAudience: "my-test-audience",
		},
	}

	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockAzIdentityFuncs(ctrl)
	clientId := "00000000-0000-0000-0000-000000000003"
	tenantId := "00000000-0000-0000-0000-000000000004"
	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	accessToken := "thepassword"
	ghToken := "gh_oidc_token"
	resourceData.Set("sp_client_id", clientId)
	resourceData.Set("sp_tenant_id", tenantId)
	resourceData.Set("sp_oidc_github_actions", true)

	for _, testCase := range testCases {
		resourceData.Set("sp_oidc_github_actions_audience", testCase.testAudience)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, testCase.expectedAudience, r.URL.Query().Get("audience"))
			assert.Equal(t, int64(0), r.ContentLength)
			assert.Equal(t, "Bearer "+ghToken, r.Header.Get("Authorization"))
			w.Header().Add("content-type", "application/json")
			fmt.Fprintln(w, "{\"value\":\""+accessToken+"\"}")
		}))
		defer ts.Close()

		os.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", ts.URL)
		os.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", ghToken)

		mockIdentityClient.EXPECT().NewClientAssertionCredential(tenantId, clientId, gomock.Any(), nil).DoAndReturn(
			func(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (*simpleTokenGetter, error) {
				assertion, err := getAssertion(context.Background())
				assert.Nil(t, err)
				assert.Equal(t, accessToken, assertion)
				getter := simpleTokenGetter{token: accessToken}
				return &getter, nil
			}).Times(1)
		resp, err := azuredevops.GetAuthToken(context.Background(), resourceData, mockIdentityClient)
		assert.Nil(t, err)
		assert.Equal(t, accessToken, resp)
	}
}
