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
		timeout  = time.Second * 30
		interval = time.Second * 1

		DatabasesetName      = "foobar"
		InstanceName         = "test-instance"
		ExpectedDatabaseName = "foobar-test-instance"
		Namespace            = "default"
		NamedGraphPrefix     = "http://foobar"
	)
	var databaseSet *stardogv1beta1.DatabaseSet

	BeforeEach(func() {
		databaseSet = &stardogv1beta1.DatabaseSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Database",
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
		Expect(k8sClient.Create(ctx, databaseSet)).Should(Succeed())
	})

	AfterEach(func() {
		Expect(k8sClient.Delete(context.Background(), databaseSet)).Should(Succeed())
	})

	Context("When creating a DatabaseSet", func() {
		It("Should create correct Database objects", func() {
			By("By creating a new DatabaseSet")
			ctx := context.Background()

			createdDatabase := &stardogv1beta1.Database{}
			Eventually(ctx, func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: ExpectedDatabaseName, Namespace: Namespace}, createdDatabase)
				return err == nil
			}).WithTimeout(timeout).WithPolling(interval).Should(BeTrue())
			Expect(createdDatabase.Spec.DatabaseName).Should(Equal(DatabasesetName))
		})
	})
})
