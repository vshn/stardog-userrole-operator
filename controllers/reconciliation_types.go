package controllers

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime"
	auth "github.com/go-openapi/runtime/client"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/vshn/stardog-userrole-operator/api/v1beta1"
	stardog "github.com/vshn/stardog-userrole-operator/stardogrest/client"
	"net/url"

	. "github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ReconciliationContext struct {
	context       context.Context
	conditions    map[StardogConditionType]StardogCondition
	stardogClient *stardog.Stardog
	namespace     string
}

type DatabaseReconciliation struct {
	resource              *v1beta1.Database
	reconciliationContext *ReconciliationContext
}

type StardogInstanceReconciliation struct {
	resource              *StardogInstance
	reconciliationContext *ReconciliationContext
}

type StardogRoleReconciliation struct {
	resource              *StardogRole
	reconciliationContext *ReconciliationContext
}

type StardogUserReconciliation struct {
	resource              *StardogUser
	reconciliationContext *ReconciliationContext
}

// SetStatusCondition adds the given condition to the status condition of the Stardog CRDs. Overwrites existing conditions
// of the same type.
func (rc *ReconciliationContext) SetStatusCondition(condition StardogCondition) {
	rc.conditions[condition.Type] = condition
}

// SetStatusIfExisting sets the condition of the given type to the given status, if the condition already exists, otherwise noop
func (rc *ReconciliationContext) SetStatusIfExisting(conditionType StardogConditionType, status v1.ConditionStatus) {
	if condition, found := rc.conditions[conditionType]; found {
		condition.Status = status
		rc.conditions[conditionType] = condition
	}
}

func (rc *ReconciliationContext) initStardogClient(kubeClient client.Client, stardogInstance StardogInstance) (runtime.ClientAuthInfoWriter, error) {
	adminCredentials := stardogInstance.Spec.AdminCredentials
	serverUrl := stardogInstance.Spec.ServerUrl
	adminUsername, adminPassword, err := rc.getCredentials(kubeClient, adminCredentials, rc.namespace)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(serverUrl)
	if err != nil || u.Host == "" {
		return nil, fmt.Errorf("invalid url from stardoginstance %s: %s", stardogInstance.Name, serverUrl)
	}

	rc.stardogClient.SetTransport(httptransport.New(u.Host, stardog.DefaultBasePath, stardog.DefaultSchemes))
	return auth.BasicAuth(adminUsername, adminPassword), nil
}

func (rc *ReconciliationContext) initStardogClientFromRef(kubeClient client.Client, instance v1beta1.StardogInstanceRef) (runtime.ClientAuthInfoWriter, error) {
	stardogInstance := &StardogInstance{}
	err := kubeClient.Get(rc.context, types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name}, stardogInstance)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve stardogInstanceRef %s/%s: %v", instance.Namespace, instance.Name, err)
	}
	rc.namespace = stardogInstance.Namespace
	return rc.initStardogClient(kubeClient, *stardogInstance)
}

func (rc *ReconciliationContext) getCredentials(kubeClient client.Client, credentials StardogUserCredentialsSpec, alternativeNamespace string) (username, password string, err error) {
	secret := &v1.Secret{}
	namespace := credentials.Namespace
	if namespace == "" {
		namespace = alternativeNamespace
	}
	err = kubeClient.Get(rc.context, types.NamespacedName{Namespace: namespace, Name: credentials.SecretRef}, secret)
	if err != nil {
		return "", "", fmt.Errorf("cannot retrieve credentials from Secret %s/%s: %v", namespace, credentials.SecretRef, err)
	}

	username, err = getSecretData(*secret, "username")
	if err != nil {
		return "", "", err
	}

	password, err = getSecretData(*secret, "password")
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}
