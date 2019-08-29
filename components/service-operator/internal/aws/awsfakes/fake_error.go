package awsfakes

import "github.com/aws/aws-sdk-go/aws/awserr"

var _ awserr.Error = &MockAWSError{}

type MockAWSError struct {
	C string
	M string
	O error
}

func (err *MockAWSError) Code() string {
	return err.C
}
func (err *MockAWSError) Error() string {
	return err.C
}
func (err *MockAWSError) Message() string {
	return err.M
}
func (err *MockAWSError) OrigErr() error {
	return err.O
}

var ResourceNotFoundException awserr.Error = &MockAWSError{
	C: "ResourceNotFoundException",
	M: "fake version of error returned when no stack",
}

var NoUpdateRequiredException awserr.Error = &MockAWSError{
	C: "No updates",
	M: "No updates",
}
