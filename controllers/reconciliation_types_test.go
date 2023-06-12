//go:build exclude

package controllers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	stardogv1alpha1 "github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	stardogrestapi "github.com/vshn/stardog-userrole-operator/stardogrest/mocks"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_initStardogClientFromRef(t *testing.T) {
	namespace := "namespace-test"
	stardogInstanceName := "instance-test"
	secretName := "secret-test"
	serverURL := "url"
	username := "admin"
	password := "1234"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogClient := stardogrestapi.NewMockExtendedBaseClientAPI(mockCtrl)

	tests := []struct {
		name            string
		namespace       v1.Namespace
		stardogInstance stardogv1alpha1.StardogInstance
		secret          v1.Secret
		err             error
	}{
		{
			name:            "GivenCorrectSetup_WhenUrlAndCredentialsPresent_ThenSetStardogClientConnection",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			err:             nil,
		},
		{
			name:            "GivenCorrectSetup_WhenStardogInstanceNotFound_ThenRaiseError",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, "non-existing-instance", secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			err: errors.New(fmt.Sprintf("cannot retrieve stardogInstanceRef %s/%s: "+
				"stardoginstances.stardog.vshn.ch \"%s\" not found", namespace, stardogInstanceName, stardogInstanceName)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeKubeClient, err := createKubeFakeClient(&tt.namespace, &tt.stardogInstance, &tt.secret)
			assert.NoError(t, err)

			rc := &ReconciliationContext{
				context:       context.Background(),
				conditions:    make(v1alpha1.StardogConditionMap),
				namespace:     namespace,
				stardogClient: stardogClient,
			}

			base64.StdEncoding.EncodeToString([]byte(username))
			if tt.err == nil {
				stardogClient.EXPECT().SetConnection(
					serverURL,
					base64.StdEncoding.EncodeToString([]byte(username)),
					base64.StdEncoding.EncodeToString([]byte(password)),
				)
			}

			err = rc.initStardogClientFromRef(fakeKubeClient, stardogInstanceName)
			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_getCredentials(t *testing.T) {
	namespace := "namespace-test"
	alternativeNamespace := "namespace-alternative-test"
	secretName := "secret-test"
	alternativeSecretName := "secret-alternative-test"
	username := "admin"
	password := "1234"

	tests := []struct {
		name        string
		namespace   v1.Namespace
		credentials stardogv1alpha1.StardogUserCredentialsSpec
		secret      v1.Secret
		user        string
		pass        string
		err         error
	}{
		{
			name:      "GivenCorrectSetup_WhenCorrectCredentials_ThenGetUsernameAndPassword",
			namespace: *createNamespace(namespace),
			credentials: stardogv1alpha1.StardogUserCredentialsSpec{
				Namespace: namespace,
				SecretRef: secretName,
			},
			secret: *createFullSecret(namespace, secretName, username, password),
			user:   base64.StdEncoding.EncodeToString([]byte(username)),
			pass:   base64.StdEncoding.EncodeToString([]byte(password)),
			err:    nil,
		},
		{
			name:      "GivenCorrectSetup_WhenSecretDoesNotExist_ThenRaiseError",
			namespace: *createNamespace(namespace),
			credentials: stardogv1alpha1.StardogUserCredentialsSpec{
				Namespace: namespace,
				SecretRef: alternativeSecretName,
			},
			secret: *createPartialSecret(namespace, secretName, username, password),
			user:   "",
			pass:   "",
			err: errors.New(fmt.Sprintf("cannot retrieve credentials from Secret %s/%s: "+
				"secrets \"%s\" not found", namespace, alternativeSecretName, alternativeSecretName)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeKubeClient, err := createKubeFakeClient(&tt.namespace, &tt.secret)
			assert.NoError(t, err)

			rc := &ReconciliationContext{
				context:       context.Background(),
				conditions:    make(v1alpha1.StardogConditionMap),
				namespace:     namespace,
				stardogClient: nil,
			}

			user, pass, err := rc.getCredentials(fakeKubeClient, tt.credentials, alternativeNamespace)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.user, user)
			assert.Equal(t, tt.pass, pass)
		})
	}
}

func createNamespace(name string) *v1.Namespace {
	namespace := &v1.Namespace{
		TypeMeta:   metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}
	return namespace
}

func createFullSecret(namespace, name, username, password string) *v1.Secret {
	secret := &v1.Secret{
		TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		StringData: map[string]string{
			"username": username,
			"password": password},
		Data: map[string][]byte{
			"username": []byte(base64.StdEncoding.EncodeToString([]byte(username))),
			"password": []byte(base64.StdEncoding.EncodeToString([]byte(password)))},
	}
	return secret
}

func createPartialSecret(namespace, name, username, password string) *v1.Secret {
	secret := &v1.Secret{
		TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		StringData: map[string]string{
			"username": username,
			"password": password},
	}
	return secret
}

func createKubeFakeClient(initObjs ...runtime.Object) (client.Client, error) {
	err := stardogv1alpha1.AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, err
	}
	kubeClient := fake.NewFakeClientWithScheme(scheme.Scheme, initObjs...)
	return kubeClient, nil
}
