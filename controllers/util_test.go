//go:build exclude

package controllers

import (
	"encoding/base64"
	"errors"
	"github.com/stretchr/testify/assert"
	. "github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	"github.com/vshn/stardog-userrole-operator/stardogrest/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"
	"time"
)

func Test_getData(t *testing.T) {

	tests := []struct {
		name           string
		secret         v1.Secret
		credentialType string
		expectValue    string
		expectErr      error
	}{
		{
			name: "GivenASecret_WhenUsernameIsPresent_ThenReturnIt",
			secret: v1.Secret{
				StringData: map[string]string{
					"fake":     "data1",
					"username": "dark",
					"password": "1234",
				},
				Data: map[string][]byte{
					"fake":     []byte("fake"),
					"username": []byte("dark"),
					"password": []byte("1234"),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test-namespace",
					Name:      "test-secret",
				},
			},
			credentialType: "username",
			expectValue:    "dark",
			expectErr:      nil,
		},
		{
			name: "GivenASecret_WhenPasswordIsPresent_ThenReturnIt",
			secret: v1.Secret{
				StringData: map[string]string{
					"fake":     "data1",
					"username": "dark",
					"password": "1234",
				},
				Data: map[string][]byte{
					"fake":     []byte("fake"),
					"username": []byte("dark"),
					"password": []byte("1234"),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test-namespace",
					Name:      "test-secret",
				},
			},
			credentialType: "password",
			expectValue:    "1234",
			expectErr:      nil,
		},
		{
			name: "GivenASecret_WhenDataIsMissing_ThenGetFromStringData",
			secret: v1.Secret{
				StringData: map[string]string{
					"fake":     "data1",
					"username": "dark",
					"password": "2345",
				},
				Data: map[string][]byte{
					"fake": []byte(base64.StdEncoding.EncodeToString([]byte("data1"))),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test-namespace",
					Name:      "test-secret",
				},
			},
			credentialType: "password",
			expectValue:    "2345",
			expectErr:      nil,
		},
		{
			name: "GivenASecret_WhenUsernameIsEmpty_ThenRaiseError",
			secret: v1.Secret{
				StringData: map[string]string{
					"fake":     "data1",
					"username": "",
					"password": "1234",
				},
				Data: map[string][]byte{
					"fake":     []byte(base64.StdEncoding.EncodeToString([]byte("data1"))),
					"username": []byte(base64.StdEncoding.EncodeToString([]byte(""))),
					"password": []byte(base64.StdEncoding.EncodeToString([]byte("1234"))),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test-namespace",
					Name:      "test-secret",
				},
			},
			credentialType: "username",
			expectValue:    "",
			expectErr:      errors.New(".data.username in the Secret test-namespace/test-secret is required"),
		},
		{
			name: "GivenASecret_WhenPasswordIsEmpty_ThenRaiseError",
			secret: v1.Secret{
				StringData: map[string]string{
					"fake":     "data1",
					"username": "dark",
					"password": "",
				},
				Data: map[string][]byte{
					"fake":     []byte(base64.StdEncoding.EncodeToString([]byte("data1"))),
					"username": []byte(base64.StdEncoding.EncodeToString([]byte("dark"))),
					"password": []byte(base64.StdEncoding.EncodeToString([]byte(""))),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test-namespace",
					Name:      "test-secret",
				},
			},
			credentialType: "password",
			expectValue:    "",
			expectErr:      errors.New(".data.password in the Secret test-namespace/test-secret is required"),
		},
		{
			name: "GivenASecret_WhenCredentialsNotPresent_ThenRaiseError",
			secret: v1.Secret{
				StringData: map[string]string{},
				Data:       map[string][]byte{},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test-namespace",
					Name:      "test-secret",
				},
			},
			credentialType: "password",
			expectValue:    "",
			expectErr:      errors.New(".data.password in the Secret test-namespace/test-secret is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := getSecretData(tt.secret, tt.credentialType)
			if tt.expectErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err, tt.expectErr.Error())
			}

			assert.Equal(t, tt.expectValue, data)
		})
	}
}

func Test_contains(t *testing.T) {

	tests := []struct {
		name        string
		value       string
		slice       []string
		expectValue bool
	}{
		{
			name:        "GivenAList_FindExistingValueString_ThenReturnTrue",
			value:       "StrB",
			slice:       []string{"StrS", "StrA", "StrC", "StrB", "StrI"},
			expectValue: true,
		},
		{
			name:        "GivenAList_FindNonExistingValueString_ThenReturnFalse",
			value:       "StrMissing",
			slice:       []string{"StrS", "StrA", "StrC", "StrB", "StrI"},
			expectValue: false,
		},
		{
			name:        "GivenAList_FindEmptyString_ThenReturnFalse",
			value:       "",
			slice:       []string{"StrS", "StrA", "StrC", "StrB", "StrI"},
			expectValue: false,
		},
		{
			name:        "GivenAnEmptyList_FindAString_ThenReturnFalse",
			value:       "StrA",
			slice:       []string{},
			expectValue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.slice, tt.value)
			assert.Equal(t, tt.expectValue, result)
		})
	}
}

func Test_missingAtLeastOne(t *testing.T) {

	tests := []struct {
		name        string
		values      []string
		slice       []string
		expectValue bool
	}{
		{
			name:        "GivenAList_FindMissingValueString_ThenReturnTrue",
			values:      []string{"StrB", "StrZ"},
			slice:       []string{"StrS", "StrA", "StrC", "StrB", "StrI"},
			expectValue: true,
		},
		{
			name:        "GivenAList_FindNonExistingValuesString_ThenReturnTrue",
			values:      []string{"StrT", "StrZ"},
			slice:       []string{"StrS", "StrA", "StrC", "StrB", "StrI"},
			expectValue: true,
		},
		{
			name:        "GivenAList_FindAllExistingValuesString_ThenReturnFalse",
			values:      []string{"StrS", "StrI"},
			slice:       []string{"StrS", "StrA", "StrC", "StrB", "StrI"},
			expectValue: false,
		},
		{
			name:        "GivenAList_FindNonExistingValuesString_ThenReturnTrue",
			values:      []string{""},
			slice:       []string{"StrS", "StrA", "StrC", "StrB", "StrI"},
			expectValue: true,
		},
		{
			name:        "GivenAEmptyList_FindNonExistingValuesString_ThenReturnTrue",
			values:      []string{"StrS"},
			slice:       []string{},
			expectValue: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := missingAtLeastOne(tt.slice, tt.values...)
			assert.Equal(t, tt.expectValue, result)
		})
	}
}

func Test_equals(t *testing.T) {

	actionRead := "READ"
	actionWrite := "WRITE"
	resourceTypeAll := "*"
	resourcesX := []string{"resourceA", "resourceB"}

	tests := []struct {
		name            string
		permissionTypeA models.Permission
		permissionTypeB StardogPermissionSpec
		expectValue     bool
	}{
		{
			name: "GivenEqualTypes_ThenReturnTrue",
			permissionTypeA: models.Permission{
				Action:       &actionRead,
				ResourceType: &resourceTypeAll,
				Resource:     resourcesX,
			},
			permissionTypeB: StardogPermissionSpec{
				Action:       actionRead,
				ResourceType: resourceTypeAll,
				Resources:    resourcesX,
			},
			expectValue: true,
		},
		{
			name: "GivenNonEqualTypes1_ThenReturnFalse",
			permissionTypeA: models.Permission{
				Action:       &actionRead,
				ResourceType: &resourceTypeAll,
				Resource:     resourcesX,
			},
			permissionTypeB: StardogPermissionSpec{
				Action:       actionRead,
				ResourceType: resourceTypeAll,
				Resources:    []string{"resourceB"},
			},
			expectValue: false,
		},
		{
			name: "GivenNonEqualTypes2_ThenReturnFalse",
			permissionTypeA: models.Permission{
				Action:       &actionRead,
				ResourceType: &resourceTypeAll,
				Resource:     resourcesX,
			},
			permissionTypeB: StardogPermissionSpec{
				Action:       actionWrite,
				ResourceType: resourceTypeAll,
				Resources:    resourcesX,
			},
			expectValue: false,
		},
		{
			name: "GivenMissingAttribute1_ThenReturnFalse",
			permissionTypeA: models.Permission{
				ResourceType: &resourceTypeAll,
				Resource:     resourcesX,
			},
			permissionTypeB: StardogPermissionSpec{
				Action:       actionWrite,
				ResourceType: resourceTypeAll,
				Resources:    resourcesX,
			},
			expectValue: false,
		},
		{
			name: "GivenMissingAttribute2_ThenReturnFalse",
			permissionTypeA: models.Permission{
				Action:       &actionRead,
				ResourceType: &resourceTypeAll,
				Resource:     resourcesX,
			},
			permissionTypeB: StardogPermissionSpec{
				Action:       actionWrite,
				ResourceType: resourceTypeAll,
			},
			expectValue: false,
		},
		{
			name: "GivenMissingAttribute3_ThenReturnFalse",
			permissionTypeA: models.Permission{
				Action:       &actionWrite,
				ResourceType: &resourceTypeAll,
			},
			permissionTypeB: StardogPermissionSpec{
				Action:       actionWrite,
				ResourceType: resourceTypeAll,
			},
			expectValue: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := equals(tt.permissionTypeA, tt.permissionTypeB)
			assert.Equal(t, tt.expectValue, result)
		})
	}
}

func Test_containsStardogPermission(t *testing.T) {

	actionRead := "READ"
	actionWrite := "WRITE"
	actionAll := "ALL"
	resourceTypeAll := "*"
	resourceTypeDB := "DB"
	resourceTypeGraph := "Graph"
	resourcesX := []string{"resourceA", "resourceB"}
	resourcesY := []string{"resourceC", "resourceD", "resourceS"}

	tests := []struct {
		name             string
		permissionTypeA  models.Permission
		permissionTypesB []StardogPermissionSpec
		expectValue      bool
	}{
		{
			name: "GivenAListOfStardogPermissionSpec_WhenPermissionExists_ThenReturnTrue",
			permissionTypeA: models.Permission{
				Action:       &actionRead,
				ResourceType: &resourceTypeAll,
				Resource:     resourcesX,
			},
			permissionTypesB: []StardogPermissionSpec{
				{
					Action:       actionRead,
					ResourceType: resourceTypeGraph,
					Resources:    resourcesX,
				},
				{
					Action:       actionAll,
					ResourceType: resourceTypeDB,
					Resources:    resourcesY,
				},
				{
					Action:       actionRead,
					ResourceType: resourceTypeAll,
					Resources:    resourcesX,
				},
				{
					Action:       actionAll,
					ResourceType: resourceTypeAll,
					Resources:    resourcesY,
				},
				{
					Action:       actionWrite,
					ResourceType: resourceTypeGraph,
					Resources:    resourcesY,
				},
			},
			expectValue: true,
		},
		{
			name:            "GivenAListOfStardogPermissionSpec_WhenPermissionIsEmpty_ThenReturnFalse",
			permissionTypeA: models.Permission{},
			permissionTypesB: []StardogPermissionSpec{
				{
					Action:       actionRead,
					ResourceType: resourceTypeGraph,
					Resources:    resourcesX,
				},
				{
					Action:       actionAll,
					ResourceType: resourceTypeDB,
					Resources:    resourcesY,
				},
				{
					Action:       actionRead,
					ResourceType: resourceTypeAll,
					Resources:    resourcesX,
				},
				{
					Action:       actionAll,
					ResourceType: resourceTypeAll,
					Resources:    resourcesY,
				},
				{
					Action:       actionWrite,
					ResourceType: resourceTypeGraph,
					Resources:    resourcesY,
				},
			},
			expectValue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsStardogPermission(tt.permissionTypesB, tt.permissionTypeA)
			assert.Equal(t, tt.expectValue, result)
		})
	}
}

func Test_containsOperatorPermission(t *testing.T) {

	actionRead := "READ"
	actionWrite := "WRITE"
	actionAll := "ALL"
	resourceTypeAll := "*"
	resourceTypeDB := "DB"
	resourceTypeGraph := "Graph"
	resourcesX := []string{"resourceA", "resourceB"}
	resourcesY := []string{"resourceC", "resourceD", "resourceS"}

	tests := []struct {
		name             string
		permissionTypeA  StardogPermissionSpec
		permissionTypesB []*models.Permission
		expectValue      bool
	}{
		{
			name: "GivenAListOfPermission_WhenStardogPermissionSpecExists_ThenReturnTrue",
			permissionTypeA: StardogPermissionSpec{
				Action:       actionRead,
				ResourceType: resourceTypeAll,
				Resources:    resourcesX,
			},
			permissionTypesB: []*models.Permission{
				{
					Action:       &actionRead,
					ResourceType: &resourceTypeGraph,
					Resource:     resourcesX,
				},
				{
					Action:       &actionAll,
					ResourceType: &resourceTypeDB,
					Resource:     resourcesY,
				},
				{
					Action:       &actionRead,
					ResourceType: &resourceTypeAll,
					Resource:     resourcesX,
				},
				{
					Action:       &actionAll,
					ResourceType: &resourceTypeAll,
					Resource:     resourcesY,
				},
				{
					Action:       &actionWrite,
					ResourceType: &resourceTypeGraph,
					Resource:     resourcesY,
				},
			},
			expectValue: true,
		},
		{
			name:            "GivenAListOfPermission_WhenStardogPermissionSpecIsEmpty_ThenReturnFalse",
			permissionTypeA: StardogPermissionSpec{},
			permissionTypesB: []*models.Permission{
				{
					Action:       &actionRead,
					ResourceType: &resourceTypeGraph,
					Resource:     resourcesX,
				},
				{
					Action:       &actionAll,
					ResourceType: &resourceTypeDB,
					Resource:     resourcesY,
				},
				{
					Action:       &actionRead,
					ResourceType: &resourceTypeAll,
					Resource:     resourcesX,
				},
				{
					Action:       &actionAll,
					ResourceType: &resourceTypeAll,
					Resource:     resourcesY,
				},
				{
					Action:       &actionWrite,
					ResourceType: &resourceTypeGraph,
					Resource:     resourcesY,
				},
			},
			expectValue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsOperatorPermission(tt.permissionTypesB, tt.permissionTypeA)
			assert.Equal(t, tt.expectValue, result)
		})
	}
}

func Test_mergeWithExistingConditions(t *testing.T) {

	now := metav1.Now()

	tests := []struct {
		name                      string
		existingStardogConditions []StardogCondition
		newStardogConditions      StardogConditionMap
		expectConditions          []StardogCondition
	}{
		{
			name: "GivenAListOfStardogConditions_WithNewStardogConditions_ThenReturnMergedConditions",
			existingStardogConditions: []StardogCondition{
				{
					Type:               StardogErrored,
					Status:             v1.ConditionTrue,
					Reason:             ReasonTerminating,
					LastTransitionTime: now,
				},
				{
					Type:               StardogReady,
					Status:             v1.ConditionTrue,
					Reason:             ReasonSpecInvalid,
					LastTransitionTime: now,
				},
				{
					Type:               StardogInvalid,
					Status:             v1.ConditionFalse,
					Reason:             ReasonSpecInvalid,
					LastTransitionTime: now,
				},
			},
			newStardogConditions: map[StardogConditionType]StardogCondition{
				StardogReady: {
					Type:               StardogReady,
					Status:             v1.ConditionTrue,
					Reason:             ReasonSucceeded,
					LastTransitionTime: now,
				},
				StardogErrored: {
					Type:               StardogErrored,
					Status:             v1.ConditionTrue,
					Reason:             ReasonSucceeded,
					LastTransitionTime: now,
				},
			},
			expectConditions: []StardogCondition{
				{
					Type:               StardogErrored,
					Status:             v1.ConditionTrue,
					Reason:             ReasonSucceeded,
					LastTransitionTime: now,
				},
				{
					Type:               StardogReady,
					Status:             v1.ConditionTrue,
					Reason:             ReasonSucceeded,
					LastTransitionTime: now,
				},
				{
					Type:               StardogInvalid,
					Status:             v1.ConditionFalse,
					Reason:             ReasonSpecInvalid,
					LastTransitionTime: now,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeWithExistingConditions(tt.existingStardogConditions, tt.newStardogConditions)
			assert.ElementsMatch(t, tt.expectConditions, result)
		})
	}
}

func Test_init(t *testing.T) {

	tests := []struct {
		name                    string
		reconFreqErr            string
		reconFreq               string
		expectedReconFreqErrDur time.Duration
		expectedReconFreqDur    time.Duration
	}{
		{
			name:                    "GiveReconFreq_WhenIsCorrectPopulated_ThenReturnDuration",
			reconFreqErr:            "1s",
			reconFreq:               "1h",
			expectedReconFreqErrDur: time.Second,
			expectedReconFreqDur:    time.Hour,
		},
		{
			name:                    "GiveReconFreq_WhenIsNotParsable_ThenReturn0Duration",
			reconFreqErr:            "1asd",
			reconFreq:               "1d",
			expectedReconFreqErrDur: 0,
			expectedReconFreqDur:    0,
		},
		{
			name:                    "GiveReconFreq_WhenIsNegativeValue_ThenReturn0Duration",
			reconFreqErr:            "-24s",
			reconFreq:               "-1h",
			expectedReconFreqErrDur: 0,
			expectedReconFreqDur:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("RECONCILIATION_FREQUENCY_ON_ERROR", tt.reconFreqErr)
			_ = os.Setenv("RECONCILIATION_FREQUENCY", tt.reconFreq)
			InitEnv()
			assert.Equal(t, tt.expectedReconFreqDur, ReconFreq)
			assert.Equal(t, tt.expectedReconFreqErrDur, ReconFreqErr)
		})
	}
}
