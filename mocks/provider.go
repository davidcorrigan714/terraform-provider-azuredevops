// Code generated by MockGen. DO NOT EDIT.
// Source: azuredevops/provider.go

// Package mock_azuredevops is a generated GoMock package.
package mock_azuredevops

import (
	context "context"
	crypto "crypto"
	x509 "crypto/x509"
	reflect "reflect"

	azcore "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	policy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	gomock "github.com/golang/mock/gomock"
	azuredevops "github.com/microsoft/terraform-provider-azuredevops/azuredevops"
)

// MockTokenGetter is a mock of TokenGetter interface.
type MockTokenGetter struct {
	ctrl     *gomock.Controller
	recorder *MockTokenGetterMockRecorder
}

// MockTokenGetterMockRecorder is the mock recorder for MockTokenGetter.
type MockTokenGetterMockRecorder struct {
	mock *MockTokenGetter
}

// NewMockTokenGetter creates a new mock instance.
func NewMockTokenGetter(ctrl *gomock.Controller) *MockTokenGetter {
	mock := &MockTokenGetter{ctrl: ctrl}
	mock.recorder = &MockTokenGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenGetter) EXPECT() *MockTokenGetterMockRecorder {
	return m.recorder
}

// GetToken mocks base method.
func (m *MockTokenGetter) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetToken", ctx, opts)
	ret0, _ := ret[0].(azcore.AccessToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetToken indicates an expected call of GetToken.
func (mr *MockTokenGetterMockRecorder) GetToken(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetToken", reflect.TypeOf((*MockTokenGetter)(nil).GetToken), ctx, opts)
}

// MockAzIdentityFuncs is a mock of AzIdentityFuncs interface.
type MockAzIdentityFuncs struct {
	ctrl     *gomock.Controller
	recorder *MockAzIdentityFuncsMockRecorder
}

// MockAzIdentityFuncsMockRecorder is the mock recorder for MockAzIdentityFuncs.
type MockAzIdentityFuncsMockRecorder struct {
	mock *MockAzIdentityFuncs
}

// NewMockAzIdentityFuncs creates a new mock instance.
func NewMockAzIdentityFuncs(ctrl *gomock.Controller) *MockAzIdentityFuncs {
	mock := &MockAzIdentityFuncs{ctrl: ctrl}
	mock.recorder = &MockAzIdentityFuncsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAzIdentityFuncs) EXPECT() *MockAzIdentityFuncsMockRecorder {
	return m.recorder
}

// NewClientAssertionCredential mocks base method.
func (m *MockAzIdentityFuncs) NewClientAssertionCredential(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (azuredevops.TokenGetter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewClientAssertionCredential", tenantID, clientID, getAssertion, options)
	ret0, _ := ret[0].(azuredevops.TokenGetter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewClientAssertionCredential indicates an expected call of NewClientAssertionCredential.
func (mr *MockAzIdentityFuncsMockRecorder) NewClientAssertionCredential(tenantID, clientID, getAssertion, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewClientAssertionCredential", reflect.TypeOf((*MockAzIdentityFuncs)(nil).NewClientAssertionCredential), tenantID, clientID, getAssertion, options)
}

// NewClientCertificateCredential mocks base method.
func (m *MockAzIdentityFuncs) NewClientCertificateCredential(tenantID, clientID string, certs []*x509.Certificate, key crypto.PrivateKey, options *azidentity.ClientCertificateCredentialOptions) (azuredevops.TokenGetter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewClientCertificateCredential", tenantID, clientID, certs, key, options)
	ret0, _ := ret[0].(azuredevops.TokenGetter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewClientCertificateCredential indicates an expected call of NewClientCertificateCredential.
func (mr *MockAzIdentityFuncsMockRecorder) NewClientCertificateCredential(tenantID, clientID, certs, key, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewClientCertificateCredential", reflect.TypeOf((*MockAzIdentityFuncs)(nil).NewClientCertificateCredential), tenantID, clientID, certs, key, options)
}

// NewClientSecretCredential mocks base method.
func (m *MockAzIdentityFuncs) NewClientSecretCredential(tenantID, clientID, clientSecret string, options *azidentity.ClientSecretCredentialOptions) (azuredevops.TokenGetter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewClientSecretCredential", tenantID, clientID, clientSecret, options)
	ret0, _ := ret[0].(azuredevops.TokenGetter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewClientSecretCredential indicates an expected call of NewClientSecretCredential.
func (mr *MockAzIdentityFuncsMockRecorder) NewClientSecretCredential(tenantID, clientID, clientSecret, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewClientSecretCredential", reflect.TypeOf((*MockAzIdentityFuncs)(nil).NewClientSecretCredential), tenantID, clientID, clientSecret, options)
}
