package awsssm

import (
	"encoding/base64"
	"errors"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// Pass handles secrets using Linux secret-service as a store.
type Awsssm struct{
	Svc *ssm.SSM
}

const credsSsmPathPrefix = "/infrastructure/docker/"

// Add adds new credentials to the keychain.
func (h Awsssm) Add(creds *credentials.Credentials) error {
	if creds == nil {
		return errors.New("missing credentials")
	}

	encoded := base64.URLEncoding.WithPadding('.').EncodeToString([]byte(creds.ServerURL))

	secretparam := (&ssm.PutParameterInput{}).
		SetAllowedPattern("^[a-zA-Z0-9\\-,/;:=]+$").
		SetDescription("Secret for docker hub").
		SetName(credsSsmPathPrefix + encoded + "/secret").
		SetType("SecureString").
		SetValue(creds.Secret)

	_, err := h.Svc.PutParameter(secretparam)

	// Only do user name if password was successful
	if err == nil {
		userparam := (&ssm.PutParameterInput{}).
			SetAllowedPattern("^[a-zA-Z0-9\\-.@_]+$").
			SetDescription("Username for docker hub").
			SetName(credsSsmPathPrefix + encoded + "/username").
			SetType("String").
			SetValue(creds.Username)

		_, err = h.Svc.PutParameter(userparam)
	}

	return err
}

// Delete removes credentials from the store.
func (h Awsssm) Delete(serverURL string) error {
	if serverURL == "" {
		return errors.New("missing server url")
	}

  encoded := base64.URLEncoding.WithPadding('.').EncodeToString([]byte(serverURL))
	delsecretparam := (&ssm.DeleteParameterInput{}).
		SetName(credsSsmPathPrefix+encoded+"/secret")

	_, err := h.Svc.DeleteParameter(delsecretparam)

	if err == nil {
		deluserparam := (&ssm.DeleteParameterInput{}).
		  SetName(credsSsmPathPrefix+encoded+"/username")
		_, err = h.Svc.DeleteParameter(deluserparam)
	}
	return err
}


// Get returns the username and secret to use for a given registry server URL.
func (h Awsssm) Get(serverURL string) (string, string, error) {
	if serverURL == "" {
		return "", "", errors.New("missing server url")
	}

	var username, password string

	encoded := base64.URLEncoding.WithPadding('.').EncodeToString([]byte(serverURL))

	getsecretparam := (&ssm.GetParameterInput{}).
		SetName(credsSsmPathPrefix+encoded+"/secret").
		SetWithDecryption(true)

	secretparam, err := h.Svc.GetParameter(getsecretparam)

	if err != nil {
		return "", "", err
	} else {
		getuserparam := (&ssm.GetParameterInput{}).
			SetName(credsSsmPathPrefix + encoded + "/username").
			SetWithDecryption(false)
		userparam, err := h.Svc.GetParameter(getuserparam)
		if err != nil {
			return "", "", err
		} else {
			username = *userparam.Parameter.Value
			password = *secretparam.Parameter.Value
		}
	}

	return username, password, nil
}

// List not implemented
func (h Awsssm) List() (map[string]string, error) {
	
	return nil, nil
}
