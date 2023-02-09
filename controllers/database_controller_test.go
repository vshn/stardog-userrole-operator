package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	stardogv1beta1 "github.com/vshn/stardog-userrole-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Database controller", func() {
	const (
		timeout  = time.Second * 15
		interval = time.Second

		DatabaseName       = "foo"
		InstanceName       = "prod-instance"
		CombinedName       = "foo-prod-instance"
		ExpectedSecretName = "foo-prod-instance-credentials"
		Namespace          = "default"
		NamedGraphPrefix   = "http://foobar"
	)
	var database *stardogv1beta1.Database
	var instance *stardogv1beta1.Instance

	BeforeEach(func() {
		database = &stardogv1beta1.Database{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Database",
				APIVersion: "stardog.vshn.ch/v1beta1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      CombinedName,
				Namespace: Namespace,
			},
			Spec: stardogv1beta1.DatabaseSpec{
				NamedGraphPrefix: NamedGraphPrefix,
				DatabaseName:     DatabaseName,
				InstanceRef: stardogv1beta1.StardogInstanceRef{
					Name: InstanceName,
				},
			},
		}
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
		Expect(k8sClient.Create(ctx, database)).Should(Succeed())
		Expect(k8sClient.Create(ctx, instance)).Should(Succeed())
	})

	AfterEach(func() {
		Expect(k8sClient.Delete(context.Background(), database)).Should(Succeed())
		Expect(k8sClient.Delete(context.Background(), instance)).Should(Succeed())
	})

	Context("When creating a Database", func() {
		It("Should create Secret objects", func() {
			By("By creating a new Database")
			ctx := context.Background()

			createdDatabase := &stardogv1beta1.Database{}
			Eventually(
				k8sClient.Get(ctx, types.NamespacedName{Name: CombinedName, Namespace: Namespace}, createdDatabase),
			).WithContext(ctx).WithTimeout(timeout).WithPolling(interval).Should(Succeed())
			Expect(createdDatabase.Spec.DatabaseName).Should(Equal(DatabaseName))

			createdCredentialsSecret := &corev1.Secret{}
			Eventually(
				k8sClient.Get(ctx, types.NamespacedName{Name: ExpectedSecretName, Namespace: Namespace}, createdCredentialsSecret),
			).WithContext(ctx).WithTimeout(timeout).WithPolling(interval).Should(Succeed())
		})
	})
})
