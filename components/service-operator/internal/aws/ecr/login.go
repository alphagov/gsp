package ecr

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/alphagov/gsp/components/service-operator/internal/aws/sdk"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type ECRCredentials struct {
	Username string
	Password string
	Endpoint string
}

func GetECRCredentials(ctx context.Context, c sdk.Client) (*ECRCredentials, error) {
	res, err := c.GetAuthorizationTokenWithContext(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return nil, err
	}
	if len(res.AuthorizationData) == 0 {
		return nil, fmt.Errorf("GetECRCredentials: no credentials returned")
	} else if len(res.AuthorizationData) != 1 {
		return nil, fmt.Errorf("GetECRCredentials: unexpected number of credentials returns")
	}
	data := res.AuthorizationData[0]
	token := aws.StringValue(data.AuthorizationToken)
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	decodedTokenParts := strings.SplitN(string(decodedToken), ":", 2)
	if len(decodedTokenParts) != 2 {
		return nil, errors.New("GetECRCredentials: invalid credential data")
	}
	creds := &ECRCredentials{
		Username: "AWS",
		Password: decodedTokenParts[1],
		Endpoint: aws.StringValue(data.ProxyEndpoint),
	}
	return creds, nil

}
