package controllers

import (
	"fmt"
	. "github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	"github.com/vshn/stardog-userrole-operator/stardogrest"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"strings"
)

// createStatusConditionReady is a shortcut for adding a StardogReady condition.
func createStatusConditionReady(isReady bool, message string) StardogCondition {
	readyCondition := StardogCondition{
		Status:             v1.ConditionFalse,
		Type:               StardogReady,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonFailed,
		Message:            message,
	}
	if isReady {
		readyCondition.Status = v1.ConditionTrue
		readyCondition.Reason = ReasonSucceeded
		readyCondition.Message = message
	}
	return readyCondition
}

// createStatusConditionErrored is a shortcut for adding a StardogErrored condition with the given error message.
func createStatusConditionErrored(err error) StardogCondition {
	return StardogCondition{
		Status:             v1.ConditionTrue,
		Type:               StardogErrored,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonFailed,
		Message:            err.Error(),
	}
}

// createStatusConditionInvalid is a shortcut for adding a StardogInvalid condition with the given error message.
func createStatusConditionInvalid(err error) StardogCondition {
	return StardogCondition{
		Status:             v1.ConditionTrue,
		Type:               StardogInvalid,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonSpecInvalid,
		Message:            err.Error(),
	}
}

// createStatusConditionTerminating is a shortcut for adding a StardogTerminating condition with the given error message.
func createStatusConditionTerminating(err error) StardogCondition {
	return StardogCondition{
		Status:             v1.ConditionTrue,
		Type:               StardogTerminating,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonTerminating,
		Message:            err.Error(),
	}
}

func getSecretData(secret v1.Secret, value string) (string, error) {
	if len(secret.Data[value]) == 0 {
		if secret.StringData[value] == "" {
			return "", fmt.Errorf(".data.%s in the Secret %s/%s is required", value, secret.Namespace, secret.Name)
		}
		return secret.StringData[value], nil
	}

	return string(secret.Data[value]), nil
}

func mergeWithExistingConditions(existing []StardogCondition, new StardogConditionMap) (merged []StardogCondition) {
	exMap := mapConditionsToTypeAndDisableStatus(existing)
	for _, condition := range new {
		exMap[condition.Type] = condition
	}
	for _, condition := range exMap {
		merged = append(merged, condition)
	}
	return merged
}

func mapConditionsToTypeAndDisableStatus(conditions []StardogCondition) (m StardogConditionMap) {
	m = make(StardogConditionMap)
	for _, condition := range conditions {
		condition.Status = v1.ConditionFalse
		m[condition.Type] = condition
	}
	return m
}

func missingAtLeastOne(list []string, strings ...string) bool {
	for _, s := range strings {
		if !contains(list, s) {
			return true
		}
	}
	return false
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func containsStardogPermission(permissionsTypeB []StardogPermissionSpec, permissionTypeA stardogrest.Permission) bool {
	for _, permissionTypeB := range permissionsTypeB {
		if equals(permissionTypeA, permissionTypeB) {
			return true
		}
	}
	return false
}

func containsOperatorPermission(permissionsTypeA []stardogrest.Permission, permissionTypeB StardogPermissionSpec) bool {
	for _, permissionTypeA := range permissionsTypeA {
		if equals(permissionTypeA, permissionTypeB) {
			return true
		}
	}
	return false
}

func equals(permissionTypeA stardogrest.Permission, permissionTypeB StardogPermissionSpec) bool {
	var action bool
	if permissionTypeA.Action == nil {
		if permissionTypeB.Action != "" {
			return false
		}
		action = true
	} else {
		action = strings.EqualFold(*permissionTypeA.Action, permissionTypeB.Action)
	}

	var resourceType bool
	if permissionTypeA.ResourceType == nil {
		if permissionTypeB.ResourceType != "" {
			return false
		}
		resourceType = true
	} else {
		resourceType = strings.EqualFold(*permissionTypeA.ResourceType, permissionTypeB.ResourceType)
	}

	var resources bool
	if permissionTypeA.Resource == nil {
		if permissionTypeB.Resources != nil {
			return false
		}
		resources = true
	} else {
		resources = reflect.DeepEqual(*permissionTypeA.Resource, permissionTypeB.Resources)
	}

	return action && resourceType && resources
}
