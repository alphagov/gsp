// Code generated by counterfeiter. DO NOT EDIT.
package sdkfakes

import (
	"context"
	"sync"

	"github.com/alphagov/gsp/components/service-operator/internal/aws/sdk"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type FakeClient struct {
	AssumeRoleStub        func(string) sdk.Client
	assumeRoleMutex       sync.RWMutex
	assumeRoleArgsForCall []struct {
		arg1 string
	}
	assumeRoleReturns struct {
		result1 sdk.Client
	}
	assumeRoleReturnsOnCall map[int]struct {
		result1 sdk.Client
	}
	CreateStackWithContextStub        func(context.Context, *cloudformation.CreateStackInput, ...request.Option) (*cloudformation.CreateStackOutput, error)
	createStackWithContextMutex       sync.RWMutex
	createStackWithContextArgsForCall []struct {
		arg1 context.Context
		arg2 *cloudformation.CreateStackInput
		arg3 []request.Option
	}
	createStackWithContextReturns struct {
		result1 *cloudformation.CreateStackOutput
		result2 error
	}
	createStackWithContextReturnsOnCall map[int]struct {
		result1 *cloudformation.CreateStackOutput
		result2 error
	}
	DeleteStackWithContextStub        func(context.Context, *cloudformation.DeleteStackInput, ...request.Option) (*cloudformation.DeleteStackOutput, error)
	deleteStackWithContextMutex       sync.RWMutex
	deleteStackWithContextArgsForCall []struct {
		arg1 context.Context
		arg2 *cloudformation.DeleteStackInput
		arg3 []request.Option
	}
	deleteStackWithContextReturns struct {
		result1 *cloudformation.DeleteStackOutput
		result2 error
	}
	deleteStackWithContextReturnsOnCall map[int]struct {
		result1 *cloudformation.DeleteStackOutput
		result2 error
	}
	DescribeStackEventsWithContextStub        func(context.Context, *cloudformation.DescribeStackEventsInput, ...request.Option) (*cloudformation.DescribeStackEventsOutput, error)
	describeStackEventsWithContextMutex       sync.RWMutex
	describeStackEventsWithContextArgsForCall []struct {
		arg1 context.Context
		arg2 *cloudformation.DescribeStackEventsInput
		arg3 []request.Option
	}
	describeStackEventsWithContextReturns struct {
		result1 *cloudformation.DescribeStackEventsOutput
		result2 error
	}
	describeStackEventsWithContextReturnsOnCall map[int]struct {
		result1 *cloudformation.DescribeStackEventsOutput
		result2 error
	}
	DescribeStacksWithContextStub        func(context.Context, *cloudformation.DescribeStacksInput, ...request.Option) (*cloudformation.DescribeStacksOutput, error)
	describeStacksWithContextMutex       sync.RWMutex
	describeStacksWithContextArgsForCall []struct {
		arg1 context.Context
		arg2 *cloudformation.DescribeStacksInput
		arg3 []request.Option
	}
	describeStacksWithContextReturns struct {
		result1 *cloudformation.DescribeStacksOutput
		result2 error
	}
	describeStacksWithContextReturnsOnCall map[int]struct {
		result1 *cloudformation.DescribeStacksOutput
		result2 error
	}
	GetAuthorizationTokenWithContextStub        func(context.Context, *ecr.GetAuthorizationTokenInput, ...request.Option) (*ecr.GetAuthorizationTokenOutput, error)
	getAuthorizationTokenWithContextMutex       sync.RWMutex
	getAuthorizationTokenWithContextArgsForCall []struct {
		arg1 context.Context
		arg2 *ecr.GetAuthorizationTokenInput
		arg3 []request.Option
	}
	getAuthorizationTokenWithContextReturns struct {
		result1 *ecr.GetAuthorizationTokenOutput
		result2 error
	}
	getAuthorizationTokenWithContextReturnsOnCall map[int]struct {
		result1 *ecr.GetAuthorizationTokenOutput
		result2 error
	}
	GetRoleCredentialsStub        func(string) *credentials.Credentials
	getRoleCredentialsMutex       sync.RWMutex
	getRoleCredentialsArgsForCall []struct {
		arg1 string
	}
	getRoleCredentialsReturns struct {
		result1 *credentials.Credentials
	}
	getRoleCredentialsReturnsOnCall map[int]struct {
		result1 *credentials.Credentials
	}
	GetSecretValueWithContextStub        func(context.Context, *secretsmanager.GetSecretValueInput, ...request.Option) (*secretsmanager.GetSecretValueOutput, error)
	getSecretValueWithContextMutex       sync.RWMutex
	getSecretValueWithContextArgsForCall []struct {
		arg1 context.Context
		arg2 *secretsmanager.GetSecretValueInput
		arg3 []request.Option
	}
	getSecretValueWithContextReturns struct {
		result1 *secretsmanager.GetSecretValueOutput
		result2 error
	}
	getSecretValueWithContextReturnsOnCall map[int]struct {
		result1 *secretsmanager.GetSecretValueOutput
		result2 error
	}
	UpdateStackWithContextStub        func(context.Context, *cloudformation.UpdateStackInput, ...request.Option) (*cloudformation.UpdateStackOutput, error)
	updateStackWithContextMutex       sync.RWMutex
	updateStackWithContextArgsForCall []struct {
		arg1 context.Context
		arg2 *cloudformation.UpdateStackInput
		arg3 []request.Option
	}
	updateStackWithContextReturns struct {
		result1 *cloudformation.UpdateStackOutput
		result2 error
	}
	updateStackWithContextReturnsOnCall map[int]struct {
		result1 *cloudformation.UpdateStackOutput
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) AssumeRole(arg1 string) sdk.Client {
	fake.assumeRoleMutex.Lock()
	ret, specificReturn := fake.assumeRoleReturnsOnCall[len(fake.assumeRoleArgsForCall)]
	fake.assumeRoleArgsForCall = append(fake.assumeRoleArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("AssumeRole", []interface{}{arg1})
	fake.assumeRoleMutex.Unlock()
	if fake.AssumeRoleStub != nil {
		return fake.AssumeRoleStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.assumeRoleReturns
	return fakeReturns.result1
}

func (fake *FakeClient) AssumeRoleCallCount() int {
	fake.assumeRoleMutex.RLock()
	defer fake.assumeRoleMutex.RUnlock()
	return len(fake.assumeRoleArgsForCall)
}

func (fake *FakeClient) AssumeRoleCalls(stub func(string) sdk.Client) {
	fake.assumeRoleMutex.Lock()
	defer fake.assumeRoleMutex.Unlock()
	fake.AssumeRoleStub = stub
}

func (fake *FakeClient) AssumeRoleArgsForCall(i int) string {
	fake.assumeRoleMutex.RLock()
	defer fake.assumeRoleMutex.RUnlock()
	argsForCall := fake.assumeRoleArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) AssumeRoleReturns(result1 sdk.Client) {
	fake.assumeRoleMutex.Lock()
	defer fake.assumeRoleMutex.Unlock()
	fake.AssumeRoleStub = nil
	fake.assumeRoleReturns = struct {
		result1 sdk.Client
	}{result1}
}

func (fake *FakeClient) AssumeRoleReturnsOnCall(i int, result1 sdk.Client) {
	fake.assumeRoleMutex.Lock()
	defer fake.assumeRoleMutex.Unlock()
	fake.AssumeRoleStub = nil
	if fake.assumeRoleReturnsOnCall == nil {
		fake.assumeRoleReturnsOnCall = make(map[int]struct {
			result1 sdk.Client
		})
	}
	fake.assumeRoleReturnsOnCall[i] = struct {
		result1 sdk.Client
	}{result1}
}

func (fake *FakeClient) CreateStackWithContext(arg1 context.Context, arg2 *cloudformation.CreateStackInput, arg3 ...request.Option) (*cloudformation.CreateStackOutput, error) {
	fake.createStackWithContextMutex.Lock()
	ret, specificReturn := fake.createStackWithContextReturnsOnCall[len(fake.createStackWithContextArgsForCall)]
	fake.createStackWithContextArgsForCall = append(fake.createStackWithContextArgsForCall, struct {
		arg1 context.Context
		arg2 *cloudformation.CreateStackInput
		arg3 []request.Option
	}{arg1, arg2, arg3})
	fake.recordInvocation("CreateStackWithContext", []interface{}{arg1, arg2, arg3})
	fake.createStackWithContextMutex.Unlock()
	if fake.CreateStackWithContextStub != nil {
		return fake.CreateStackWithContextStub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.createStackWithContextReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) CreateStackWithContextCallCount() int {
	fake.createStackWithContextMutex.RLock()
	defer fake.createStackWithContextMutex.RUnlock()
	return len(fake.createStackWithContextArgsForCall)
}

func (fake *FakeClient) CreateStackWithContextCalls(stub func(context.Context, *cloudformation.CreateStackInput, ...request.Option) (*cloudformation.CreateStackOutput, error)) {
	fake.createStackWithContextMutex.Lock()
	defer fake.createStackWithContextMutex.Unlock()
	fake.CreateStackWithContextStub = stub
}

func (fake *FakeClient) CreateStackWithContextArgsForCall(i int) (context.Context, *cloudformation.CreateStackInput, []request.Option) {
	fake.createStackWithContextMutex.RLock()
	defer fake.createStackWithContextMutex.RUnlock()
	argsForCall := fake.createStackWithContextArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) CreateStackWithContextReturns(result1 *cloudformation.CreateStackOutput, result2 error) {
	fake.createStackWithContextMutex.Lock()
	defer fake.createStackWithContextMutex.Unlock()
	fake.CreateStackWithContextStub = nil
	fake.createStackWithContextReturns = struct {
		result1 *cloudformation.CreateStackOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) CreateStackWithContextReturnsOnCall(i int, result1 *cloudformation.CreateStackOutput, result2 error) {
	fake.createStackWithContextMutex.Lock()
	defer fake.createStackWithContextMutex.Unlock()
	fake.CreateStackWithContextStub = nil
	if fake.createStackWithContextReturnsOnCall == nil {
		fake.createStackWithContextReturnsOnCall = make(map[int]struct {
			result1 *cloudformation.CreateStackOutput
			result2 error
		})
	}
	fake.createStackWithContextReturnsOnCall[i] = struct {
		result1 *cloudformation.CreateStackOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) DeleteStackWithContext(arg1 context.Context, arg2 *cloudformation.DeleteStackInput, arg3 ...request.Option) (*cloudformation.DeleteStackOutput, error) {
	fake.deleteStackWithContextMutex.Lock()
	ret, specificReturn := fake.deleteStackWithContextReturnsOnCall[len(fake.deleteStackWithContextArgsForCall)]
	fake.deleteStackWithContextArgsForCall = append(fake.deleteStackWithContextArgsForCall, struct {
		arg1 context.Context
		arg2 *cloudformation.DeleteStackInput
		arg3 []request.Option
	}{arg1, arg2, arg3})
	fake.recordInvocation("DeleteStackWithContext", []interface{}{arg1, arg2, arg3})
	fake.deleteStackWithContextMutex.Unlock()
	if fake.DeleteStackWithContextStub != nil {
		return fake.DeleteStackWithContextStub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.deleteStackWithContextReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) DeleteStackWithContextCallCount() int {
	fake.deleteStackWithContextMutex.RLock()
	defer fake.deleteStackWithContextMutex.RUnlock()
	return len(fake.deleteStackWithContextArgsForCall)
}

func (fake *FakeClient) DeleteStackWithContextCalls(stub func(context.Context, *cloudformation.DeleteStackInput, ...request.Option) (*cloudformation.DeleteStackOutput, error)) {
	fake.deleteStackWithContextMutex.Lock()
	defer fake.deleteStackWithContextMutex.Unlock()
	fake.DeleteStackWithContextStub = stub
}

func (fake *FakeClient) DeleteStackWithContextArgsForCall(i int) (context.Context, *cloudformation.DeleteStackInput, []request.Option) {
	fake.deleteStackWithContextMutex.RLock()
	defer fake.deleteStackWithContextMutex.RUnlock()
	argsForCall := fake.deleteStackWithContextArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) DeleteStackWithContextReturns(result1 *cloudformation.DeleteStackOutput, result2 error) {
	fake.deleteStackWithContextMutex.Lock()
	defer fake.deleteStackWithContextMutex.Unlock()
	fake.DeleteStackWithContextStub = nil
	fake.deleteStackWithContextReturns = struct {
		result1 *cloudformation.DeleteStackOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) DeleteStackWithContextReturnsOnCall(i int, result1 *cloudformation.DeleteStackOutput, result2 error) {
	fake.deleteStackWithContextMutex.Lock()
	defer fake.deleteStackWithContextMutex.Unlock()
	fake.DeleteStackWithContextStub = nil
	if fake.deleteStackWithContextReturnsOnCall == nil {
		fake.deleteStackWithContextReturnsOnCall = make(map[int]struct {
			result1 *cloudformation.DeleteStackOutput
			result2 error
		})
	}
	fake.deleteStackWithContextReturnsOnCall[i] = struct {
		result1 *cloudformation.DeleteStackOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) DescribeStackEventsWithContext(arg1 context.Context, arg2 *cloudformation.DescribeStackEventsInput, arg3 ...request.Option) (*cloudformation.DescribeStackEventsOutput, error) {
	fake.describeStackEventsWithContextMutex.Lock()
	ret, specificReturn := fake.describeStackEventsWithContextReturnsOnCall[len(fake.describeStackEventsWithContextArgsForCall)]
	fake.describeStackEventsWithContextArgsForCall = append(fake.describeStackEventsWithContextArgsForCall, struct {
		arg1 context.Context
		arg2 *cloudformation.DescribeStackEventsInput
		arg3 []request.Option
	}{arg1, arg2, arg3})
	fake.recordInvocation("DescribeStackEventsWithContext", []interface{}{arg1, arg2, arg3})
	fake.describeStackEventsWithContextMutex.Unlock()
	if fake.DescribeStackEventsWithContextStub != nil {
		return fake.DescribeStackEventsWithContextStub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.describeStackEventsWithContextReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) DescribeStackEventsWithContextCallCount() int {
	fake.describeStackEventsWithContextMutex.RLock()
	defer fake.describeStackEventsWithContextMutex.RUnlock()
	return len(fake.describeStackEventsWithContextArgsForCall)
}

func (fake *FakeClient) DescribeStackEventsWithContextCalls(stub func(context.Context, *cloudformation.DescribeStackEventsInput, ...request.Option) (*cloudformation.DescribeStackEventsOutput, error)) {
	fake.describeStackEventsWithContextMutex.Lock()
	defer fake.describeStackEventsWithContextMutex.Unlock()
	fake.DescribeStackEventsWithContextStub = stub
}

func (fake *FakeClient) DescribeStackEventsWithContextArgsForCall(i int) (context.Context, *cloudformation.DescribeStackEventsInput, []request.Option) {
	fake.describeStackEventsWithContextMutex.RLock()
	defer fake.describeStackEventsWithContextMutex.RUnlock()
	argsForCall := fake.describeStackEventsWithContextArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) DescribeStackEventsWithContextReturns(result1 *cloudformation.DescribeStackEventsOutput, result2 error) {
	fake.describeStackEventsWithContextMutex.Lock()
	defer fake.describeStackEventsWithContextMutex.Unlock()
	fake.DescribeStackEventsWithContextStub = nil
	fake.describeStackEventsWithContextReturns = struct {
		result1 *cloudformation.DescribeStackEventsOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) DescribeStackEventsWithContextReturnsOnCall(i int, result1 *cloudformation.DescribeStackEventsOutput, result2 error) {
	fake.describeStackEventsWithContextMutex.Lock()
	defer fake.describeStackEventsWithContextMutex.Unlock()
	fake.DescribeStackEventsWithContextStub = nil
	if fake.describeStackEventsWithContextReturnsOnCall == nil {
		fake.describeStackEventsWithContextReturnsOnCall = make(map[int]struct {
			result1 *cloudformation.DescribeStackEventsOutput
			result2 error
		})
	}
	fake.describeStackEventsWithContextReturnsOnCall[i] = struct {
		result1 *cloudformation.DescribeStackEventsOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) DescribeStacksWithContext(arg1 context.Context, arg2 *cloudformation.DescribeStacksInput, arg3 ...request.Option) (*cloudformation.DescribeStacksOutput, error) {
	fake.describeStacksWithContextMutex.Lock()
	ret, specificReturn := fake.describeStacksWithContextReturnsOnCall[len(fake.describeStacksWithContextArgsForCall)]
	fake.describeStacksWithContextArgsForCall = append(fake.describeStacksWithContextArgsForCall, struct {
		arg1 context.Context
		arg2 *cloudformation.DescribeStacksInput
		arg3 []request.Option
	}{arg1, arg2, arg3})
	fake.recordInvocation("DescribeStacksWithContext", []interface{}{arg1, arg2, arg3})
	fake.describeStacksWithContextMutex.Unlock()
	if fake.DescribeStacksWithContextStub != nil {
		return fake.DescribeStacksWithContextStub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.describeStacksWithContextReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) DescribeStacksWithContextCallCount() int {
	fake.describeStacksWithContextMutex.RLock()
	defer fake.describeStacksWithContextMutex.RUnlock()
	return len(fake.describeStacksWithContextArgsForCall)
}

func (fake *FakeClient) DescribeStacksWithContextCalls(stub func(context.Context, *cloudformation.DescribeStacksInput, ...request.Option) (*cloudformation.DescribeStacksOutput, error)) {
	fake.describeStacksWithContextMutex.Lock()
	defer fake.describeStacksWithContextMutex.Unlock()
	fake.DescribeStacksWithContextStub = stub
}

func (fake *FakeClient) DescribeStacksWithContextArgsForCall(i int) (context.Context, *cloudformation.DescribeStacksInput, []request.Option) {
	fake.describeStacksWithContextMutex.RLock()
	defer fake.describeStacksWithContextMutex.RUnlock()
	argsForCall := fake.describeStacksWithContextArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) DescribeStacksWithContextReturns(result1 *cloudformation.DescribeStacksOutput, result2 error) {
	fake.describeStacksWithContextMutex.Lock()
	defer fake.describeStacksWithContextMutex.Unlock()
	fake.DescribeStacksWithContextStub = nil
	fake.describeStacksWithContextReturns = struct {
		result1 *cloudformation.DescribeStacksOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) DescribeStacksWithContextReturnsOnCall(i int, result1 *cloudformation.DescribeStacksOutput, result2 error) {
	fake.describeStacksWithContextMutex.Lock()
	defer fake.describeStacksWithContextMutex.Unlock()
	fake.DescribeStacksWithContextStub = nil
	if fake.describeStacksWithContextReturnsOnCall == nil {
		fake.describeStacksWithContextReturnsOnCall = make(map[int]struct {
			result1 *cloudformation.DescribeStacksOutput
			result2 error
		})
	}
	fake.describeStacksWithContextReturnsOnCall[i] = struct {
		result1 *cloudformation.DescribeStacksOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetAuthorizationTokenWithContext(arg1 context.Context, arg2 *ecr.GetAuthorizationTokenInput, arg3 ...request.Option) (*ecr.GetAuthorizationTokenOutput, error) {
	fake.getAuthorizationTokenWithContextMutex.Lock()
	ret, specificReturn := fake.getAuthorizationTokenWithContextReturnsOnCall[len(fake.getAuthorizationTokenWithContextArgsForCall)]
	fake.getAuthorizationTokenWithContextArgsForCall = append(fake.getAuthorizationTokenWithContextArgsForCall, struct {
		arg1 context.Context
		arg2 *ecr.GetAuthorizationTokenInput
		arg3 []request.Option
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetAuthorizationTokenWithContext", []interface{}{arg1, arg2, arg3})
	fake.getAuthorizationTokenWithContextMutex.Unlock()
	if fake.GetAuthorizationTokenWithContextStub != nil {
		return fake.GetAuthorizationTokenWithContextStub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getAuthorizationTokenWithContextReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetAuthorizationTokenWithContextCallCount() int {
	fake.getAuthorizationTokenWithContextMutex.RLock()
	defer fake.getAuthorizationTokenWithContextMutex.RUnlock()
	return len(fake.getAuthorizationTokenWithContextArgsForCall)
}

func (fake *FakeClient) GetAuthorizationTokenWithContextCalls(stub func(context.Context, *ecr.GetAuthorizationTokenInput, ...request.Option) (*ecr.GetAuthorizationTokenOutput, error)) {
	fake.getAuthorizationTokenWithContextMutex.Lock()
	defer fake.getAuthorizationTokenWithContextMutex.Unlock()
	fake.GetAuthorizationTokenWithContextStub = stub
}

func (fake *FakeClient) GetAuthorizationTokenWithContextArgsForCall(i int) (context.Context, *ecr.GetAuthorizationTokenInput, []request.Option) {
	fake.getAuthorizationTokenWithContextMutex.RLock()
	defer fake.getAuthorizationTokenWithContextMutex.RUnlock()
	argsForCall := fake.getAuthorizationTokenWithContextArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) GetAuthorizationTokenWithContextReturns(result1 *ecr.GetAuthorizationTokenOutput, result2 error) {
	fake.getAuthorizationTokenWithContextMutex.Lock()
	defer fake.getAuthorizationTokenWithContextMutex.Unlock()
	fake.GetAuthorizationTokenWithContextStub = nil
	fake.getAuthorizationTokenWithContextReturns = struct {
		result1 *ecr.GetAuthorizationTokenOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetAuthorizationTokenWithContextReturnsOnCall(i int, result1 *ecr.GetAuthorizationTokenOutput, result2 error) {
	fake.getAuthorizationTokenWithContextMutex.Lock()
	defer fake.getAuthorizationTokenWithContextMutex.Unlock()
	fake.GetAuthorizationTokenWithContextStub = nil
	if fake.getAuthorizationTokenWithContextReturnsOnCall == nil {
		fake.getAuthorizationTokenWithContextReturnsOnCall = make(map[int]struct {
			result1 *ecr.GetAuthorizationTokenOutput
			result2 error
		})
	}
	fake.getAuthorizationTokenWithContextReturnsOnCall[i] = struct {
		result1 *ecr.GetAuthorizationTokenOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetRoleCredentials(arg1 string) *credentials.Credentials {
	fake.getRoleCredentialsMutex.Lock()
	ret, specificReturn := fake.getRoleCredentialsReturnsOnCall[len(fake.getRoleCredentialsArgsForCall)]
	fake.getRoleCredentialsArgsForCall = append(fake.getRoleCredentialsArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetRoleCredentials", []interface{}{arg1})
	fake.getRoleCredentialsMutex.Unlock()
	if fake.GetRoleCredentialsStub != nil {
		return fake.GetRoleCredentialsStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getRoleCredentialsReturns
	return fakeReturns.result1
}

func (fake *FakeClient) GetRoleCredentialsCallCount() int {
	fake.getRoleCredentialsMutex.RLock()
	defer fake.getRoleCredentialsMutex.RUnlock()
	return len(fake.getRoleCredentialsArgsForCall)
}

func (fake *FakeClient) GetRoleCredentialsCalls(stub func(string) *credentials.Credentials) {
	fake.getRoleCredentialsMutex.Lock()
	defer fake.getRoleCredentialsMutex.Unlock()
	fake.GetRoleCredentialsStub = stub
}

func (fake *FakeClient) GetRoleCredentialsArgsForCall(i int) string {
	fake.getRoleCredentialsMutex.RLock()
	defer fake.getRoleCredentialsMutex.RUnlock()
	argsForCall := fake.getRoleCredentialsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) GetRoleCredentialsReturns(result1 *credentials.Credentials) {
	fake.getRoleCredentialsMutex.Lock()
	defer fake.getRoleCredentialsMutex.Unlock()
	fake.GetRoleCredentialsStub = nil
	fake.getRoleCredentialsReturns = struct {
		result1 *credentials.Credentials
	}{result1}
}

func (fake *FakeClient) GetRoleCredentialsReturnsOnCall(i int, result1 *credentials.Credentials) {
	fake.getRoleCredentialsMutex.Lock()
	defer fake.getRoleCredentialsMutex.Unlock()
	fake.GetRoleCredentialsStub = nil
	if fake.getRoleCredentialsReturnsOnCall == nil {
		fake.getRoleCredentialsReturnsOnCall = make(map[int]struct {
			result1 *credentials.Credentials
		})
	}
	fake.getRoleCredentialsReturnsOnCall[i] = struct {
		result1 *credentials.Credentials
	}{result1}
}

func (fake *FakeClient) GetSecretValueWithContext(arg1 context.Context, arg2 *secretsmanager.GetSecretValueInput, arg3 ...request.Option) (*secretsmanager.GetSecretValueOutput, error) {
	fake.getSecretValueWithContextMutex.Lock()
	ret, specificReturn := fake.getSecretValueWithContextReturnsOnCall[len(fake.getSecretValueWithContextArgsForCall)]
	fake.getSecretValueWithContextArgsForCall = append(fake.getSecretValueWithContextArgsForCall, struct {
		arg1 context.Context
		arg2 *secretsmanager.GetSecretValueInput
		arg3 []request.Option
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetSecretValueWithContext", []interface{}{arg1, arg2, arg3})
	fake.getSecretValueWithContextMutex.Unlock()
	if fake.GetSecretValueWithContextStub != nil {
		return fake.GetSecretValueWithContextStub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getSecretValueWithContextReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetSecretValueWithContextCallCount() int {
	fake.getSecretValueWithContextMutex.RLock()
	defer fake.getSecretValueWithContextMutex.RUnlock()
	return len(fake.getSecretValueWithContextArgsForCall)
}

func (fake *FakeClient) GetSecretValueWithContextCalls(stub func(context.Context, *secretsmanager.GetSecretValueInput, ...request.Option) (*secretsmanager.GetSecretValueOutput, error)) {
	fake.getSecretValueWithContextMutex.Lock()
	defer fake.getSecretValueWithContextMutex.Unlock()
	fake.GetSecretValueWithContextStub = stub
}

func (fake *FakeClient) GetSecretValueWithContextArgsForCall(i int) (context.Context, *secretsmanager.GetSecretValueInput, []request.Option) {
	fake.getSecretValueWithContextMutex.RLock()
	defer fake.getSecretValueWithContextMutex.RUnlock()
	argsForCall := fake.getSecretValueWithContextArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) GetSecretValueWithContextReturns(result1 *secretsmanager.GetSecretValueOutput, result2 error) {
	fake.getSecretValueWithContextMutex.Lock()
	defer fake.getSecretValueWithContextMutex.Unlock()
	fake.GetSecretValueWithContextStub = nil
	fake.getSecretValueWithContextReturns = struct {
		result1 *secretsmanager.GetSecretValueOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetSecretValueWithContextReturnsOnCall(i int, result1 *secretsmanager.GetSecretValueOutput, result2 error) {
	fake.getSecretValueWithContextMutex.Lock()
	defer fake.getSecretValueWithContextMutex.Unlock()
	fake.GetSecretValueWithContextStub = nil
	if fake.getSecretValueWithContextReturnsOnCall == nil {
		fake.getSecretValueWithContextReturnsOnCall = make(map[int]struct {
			result1 *secretsmanager.GetSecretValueOutput
			result2 error
		})
	}
	fake.getSecretValueWithContextReturnsOnCall[i] = struct {
		result1 *secretsmanager.GetSecretValueOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) UpdateStackWithContext(arg1 context.Context, arg2 *cloudformation.UpdateStackInput, arg3 ...request.Option) (*cloudformation.UpdateStackOutput, error) {
	fake.updateStackWithContextMutex.Lock()
	ret, specificReturn := fake.updateStackWithContextReturnsOnCall[len(fake.updateStackWithContextArgsForCall)]
	fake.updateStackWithContextArgsForCall = append(fake.updateStackWithContextArgsForCall, struct {
		arg1 context.Context
		arg2 *cloudformation.UpdateStackInput
		arg3 []request.Option
	}{arg1, arg2, arg3})
	fake.recordInvocation("UpdateStackWithContext", []interface{}{arg1, arg2, arg3})
	fake.updateStackWithContextMutex.Unlock()
	if fake.UpdateStackWithContextStub != nil {
		return fake.UpdateStackWithContextStub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.updateStackWithContextReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) UpdateStackWithContextCallCount() int {
	fake.updateStackWithContextMutex.RLock()
	defer fake.updateStackWithContextMutex.RUnlock()
	return len(fake.updateStackWithContextArgsForCall)
}

func (fake *FakeClient) UpdateStackWithContextCalls(stub func(context.Context, *cloudformation.UpdateStackInput, ...request.Option) (*cloudformation.UpdateStackOutput, error)) {
	fake.updateStackWithContextMutex.Lock()
	defer fake.updateStackWithContextMutex.Unlock()
	fake.UpdateStackWithContextStub = stub
}

func (fake *FakeClient) UpdateStackWithContextArgsForCall(i int) (context.Context, *cloudformation.UpdateStackInput, []request.Option) {
	fake.updateStackWithContextMutex.RLock()
	defer fake.updateStackWithContextMutex.RUnlock()
	argsForCall := fake.updateStackWithContextArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) UpdateStackWithContextReturns(result1 *cloudformation.UpdateStackOutput, result2 error) {
	fake.updateStackWithContextMutex.Lock()
	defer fake.updateStackWithContextMutex.Unlock()
	fake.UpdateStackWithContextStub = nil
	fake.updateStackWithContextReturns = struct {
		result1 *cloudformation.UpdateStackOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) UpdateStackWithContextReturnsOnCall(i int, result1 *cloudformation.UpdateStackOutput, result2 error) {
	fake.updateStackWithContextMutex.Lock()
	defer fake.updateStackWithContextMutex.Unlock()
	fake.UpdateStackWithContextStub = nil
	if fake.updateStackWithContextReturnsOnCall == nil {
		fake.updateStackWithContextReturnsOnCall = make(map[int]struct {
			result1 *cloudformation.UpdateStackOutput
			result2 error
		})
	}
	fake.updateStackWithContextReturnsOnCall[i] = struct {
		result1 *cloudformation.UpdateStackOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.assumeRoleMutex.RLock()
	defer fake.assumeRoleMutex.RUnlock()
	fake.createStackWithContextMutex.RLock()
	defer fake.createStackWithContextMutex.RUnlock()
	fake.deleteStackWithContextMutex.RLock()
	defer fake.deleteStackWithContextMutex.RUnlock()
	fake.describeStackEventsWithContextMutex.RLock()
	defer fake.describeStackEventsWithContextMutex.RUnlock()
	fake.describeStacksWithContextMutex.RLock()
	defer fake.describeStacksWithContextMutex.RUnlock()
	fake.getAuthorizationTokenWithContextMutex.RLock()
	defer fake.getAuthorizationTokenWithContextMutex.RUnlock()
	fake.getRoleCredentialsMutex.RLock()
	defer fake.getRoleCredentialsMutex.RUnlock()
	fake.getSecretValueWithContextMutex.RLock()
	defer fake.getSecretValueWithContextMutex.RUnlock()
	fake.updateStackWithContextMutex.RLock()
	defer fake.updateStackWithContextMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ sdk.Client = new(FakeClient)
