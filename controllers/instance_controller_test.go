package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	stardogv1beta1 "github.com/vshn/stardog-userrole-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Instance controller", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Second

		InstanceName = "test-instance"
		Namespace    = "default"
	)
	var instance *stardogv1beta1.Instance

	BeforeEach(func() {
		instance = &stardogv1beta1.Instance{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Instance",
				APIVersion: "stardog.vshn.ch/v1beta1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      InstanceName,
				Namespace: Namespace,
			},
			Spec: stardogv1beta1.InstanceSpec{
				URL: "http://localhost:8080",
				AdminCredentialRef: stardogv1beta1.SecretKeyRef{
					Name: "admin-credentials",
					Key:  "password",
				},
			},
		}
		Expect(k8sClient.Create(ctx, instance)).Should(Succeed())
	})

	AfterEach(func() {
		Expect(k8sClient.Delete(context.Background(), instance)).Should(Succeed())
	})

	Context("When creating an Instance", func() {
		It("should check if it is available", func() {
			ctx := context.Background()

			createdInstance := &stardogv1beta1.Instance{}
			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{Name: InstanceName, Namespace: Namespace}, createdInstance)
			}).WithContext(ctx).WithTimeout(timeout).WithPolling(interval).Should(Succeed())

			Eventually(func() int {
				k8sClient.Get(ctx, types.NamespacedName{Name: InstanceName, Namespace: Namespace}, createdInstance)
				return len(createdInstance.Status.Conditions)
			}).WithContext(ctx).WithTimeout(timeout).WithPolling(interval).ShouldNot(BeZero())
		})
	})
})
