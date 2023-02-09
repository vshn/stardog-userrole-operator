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

var _ = Describe("DatabaseSet controller", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Second

		DatabasesetName      = "foobar"
		InstanceName         = "dev-instance"
		ExpectedDatabaseName = "foobar-dev-instance"
		Namespace            = "default"
		NamedGraphPrefix     = "http://foobar"
	)
	var databaseSet *stardogv1beta1.DatabaseSet
	var instance *stardogv1beta1.Instance

	BeforeEach(func() {
		databaseSet = &stardogv1beta1.DatabaseSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       "DatabaseSet",
				APIVersion: "stardog.vshn.ch/v1beta1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      DatabasesetName,
				Namespace: Namespace,
			},
			Spec: stardogv1beta1.DatabaseSetSpec{
				NamedGraphPrefix: NamedGraphPrefix,
				Instances: []stardogv1beta1.StardogInstanceRef{
					{
						Name: InstanceName,
					},
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
		Expect(k8sClient.Create(ctx, databaseSet)).Should(Succeed())
		Expect(k8sClient.Create(ctx, instance)).Should(Succeed())
	})

	AfterEach(func() {
		ctx := context.Background()
		Expect(k8sClient.Delete(ctx, databaseSet)).Should(Succeed())
		Expect(k8sClient.Delete(ctx, instance)).Should(Succeed())
	})

	Context("When creating a DatabaseSet", func() {
		It("Should create correct Database objects", func() {
			By("By creating a new DatabaseSet")
			ctx := context.Background()

			createdDatabase := &stardogv1beta1.Database{}
			Eventually(
				k8sClient.Get(ctx, types.NamespacedName{Name: ExpectedDatabaseName, Namespace: Namespace}, createdDatabase),
			).WithContext(ctx).WithTimeout(timeout).WithPolling(interval).Should(Succeed())
			Expect(createdDatabase.Spec.DatabaseName).Should(Equal(DatabasesetName))
		})
	})
})
