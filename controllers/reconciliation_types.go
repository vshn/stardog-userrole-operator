package controllers

import (
	"context"
	"fmt"
	. "github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	"github.com/vshn/stardog-userrole-operator/stardogrest/stardogrestapi"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ReconciliationContext struct {
	context       context.Context
	conditions    map[StardogConditionType]StardogCondition
	stardogClient stardogrestapi.ExtendedBaseClientAPI
	namespace     string
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

func (rc *ReconciliationContext) initStardogClientFromRef(kubeClient client.Client, stardogInstanceRef string) error {
	stardogInstance := &StardogInstance{}
	err := kubeClient.Get(rc.context, types.NamespacedName{Namespace: rc.namespace, Name: stardogInstanceRef}, stardogInstance)
	if err != nil {
		return fmt.Errorf("cannot retrieve stardogInstanceRef %s/%s: %v", rc.namespace, stardogInstanceRef, err)
	}

	adminCredentials := stardogInstance.Spec.AdminCredentials
	adminUsername, adminPassword, err := rc.getCredentials(kubeClient, adminCredentials, rc.namespace)
	if err != nil {
		return err
	}

	rc.stardogClient.SetConnection(stardogInstance.Spec.ServerUrl, adminUsername, adminPassword)
	return nil
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
