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
		DatabaseName       = "foo"
		InstanceName       = "test-instance"
		CombinedName       = "foo-test-instance"
		ExpectedSecretName = "foo-test-instance-credentials"
		Namespace          = "default"
		NamedGraphPrefix   = "http://foobar"
	)
	var database *stardogv1beta1.Database

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
		Expect(k8sClient.Create(ctx, database)).Should(Succeed())
	})

	AfterEach(func() {
		Expect(k8sClient.Delete(context.Background(), database)).Should(Succeed())
	})

	Context("When creating a Database", func() {
		It("Should create Secret objects", func() {
			By("By creating a new Database")
			ctx := context.Background()

			createdDatabase := &stardogv1beta1.Database{}
			Eventually(ctx, func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: CombinedName, Namespace: Namespace}, createdDatabase)
				return err == nil
			}).WithTimeout(60 * time.Second).Should(BeTrue())
			Expect(createdDatabase.Spec.DatabaseName).Should(Equal(DatabaseName))

			createdCredentialsSecret := &corev1.Secret{}
			Eventually(ctx, func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: ExpectedSecretName, Namespace: Namespace}, createdCredentialsSecret)
				return err == nil
			}).WithTimeout(60 * time.Second).Should(BeTrue())
		})
	})
})
