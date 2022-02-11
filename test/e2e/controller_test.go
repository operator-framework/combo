package e2e

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/combo/api/v1alpha1"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Combination controller", func() {
	validTemplateCR := v1alpha1.Template{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "validtemplate",
		},
		Spec: v1alpha1.TemplateSpec{
			Body:       "---\nFIRSTNAME: LASTNAME",
			Parameters: []string{"FIRSTNAME", "LASTNAME"},
		},
	}

	validUpdatedTemplateCR := v1alpha1.Template{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "validtemplate",
		},
		Spec: v1alpha1.TemplateSpec{
			Body:       "---\nFIRSTNAME: foo\nLASTNAME: bar",
			Parameters: []string{"FIRSTNAME", "LASTNAME"},
		},
	}

	validCombinationCR := v1alpha1.Combination{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "validcombination",
		},
		Spec: v1alpha1.CombinationSpec{
			Arguments: []v1alpha1.Argument{
				{
					Key: "FIRSTNAME",
					Values: []string{
						"John",
						"Luke",
					},
				},
				{
					Key: "LASTNAME",
					Values: []string{
						"Snow",
						"Skywalker",
					},
				},
			},
		},
	}

	expectedEvaluations := []string{
		"John: Snow",
		"John: Skywalker",
		"Luke: Snow",
		"Luke: Skywalker",
	}

	expectedUpdatedEvaluations := []string{
		"John: foo\nSnow: bar",
		"John: foo\nSkywalker: bar",
		"Luke: foo\nSkywalker: bar",
		"Luke: foo\nSnow: bar",
	}

	When("given healthy input and a healthy template", func() {
		var ctx context.Context
		var validTemplateCRCopy *v1alpha1.Template
		var validCombinationCRCopy *v1alpha1.Combination

		BeforeEach(func() {
			ctx = context.Background()

			// Create copies of valid CR's
			validTemplateCRCopy = validTemplateCR.DeepCopy()
			validCombinationCRCopy = validCombinationCR.DeepCopy()

			// Create a valid template
			err := kubeclient.Create(ctx, validTemplateCRCopy)
			Expect(err).To(BeNil(), "failed to create template CR")

			// Create a valid combination
			validCombinationCRCopy.Spec.Template = validTemplateCRCopy.Name
			err = kubeclient.Create(ctx, validCombinationCRCopy)
			Expect(err).To(BeNil(), "failed to create combination CR")
		})

		AfterEach(func() {
			err := kubeclient.Delete(ctx, validCombinationCRCopy)
			Expect(err).To(BeNil(), "failed to clean-up combination CR after test")

			err = kubeclient.Delete(ctx, validTemplateCRCopy)
			Expect(err).To(BeNil(), "failed to clean-up template CR after test")

			ctx.Done()
		})

		It("should get the correct evaluations and a Processed status", func() {
			// Give up to a minute (default eventually timeout) for the combination to process properly
			Eventually(func(g Gomega) error {
				var retrievedCombination v1alpha1.Combination
				err := kubeclient.Get(ctx, types.NamespacedName{Name: validCombinationCRCopy.Name}, &retrievedCombination)

				var conditionReasons []string
				for _, condition := range retrievedCombination.Status.Conditions {
					conditionReasons = append(conditionReasons, condition.Reason)
				}

				g.Expect(conditionReasons).To(ContainElement(v1alpha1.ReasonProcessed))
				g.Expect(retrievedCombination.Status.Evaluations).To(ContainElements(expectedEvaluations))

				return err
			}).Should(Succeed())
		})

		It("should reevaluate whenever its associated template gets updated", func() {
			validTemplateCRCopy.Spec.Body = validUpdatedTemplateCR.Spec.Body
			Eventually(func(g Gomega) error {
				var retrievedCombination v1alpha1.Combination
				if err := kubeclient.Get(ctx, types.NamespacedName{Name: validCombinationCRCopy.Name}, &retrievedCombination); err != nil {
					return err
				}

				if err := kubeclient.Update(ctx, validTemplateCRCopy); err != nil {
					return err
				}

				var conditionReasons []string
				for _, condition := range retrievedCombination.Status.Conditions {
					conditionReasons = append(conditionReasons, condition.Reason)
				}

				g.Expect(conditionReasons).To(ContainElement(v1alpha1.ReasonProcessed))
				g.Expect(retrievedCombination.Status.Evaluations).To(ContainElements(expectedUpdatedEvaluations))

				return nil
			}).Should(Succeed())
		})

		It("should not requeue combinations that are not related to an updated template", func() {
			// Create and defer deletion of an alternative template
			altTemplate := validTemplateCR.DeepCopy()
			altTemplate.Name = "alttemplate"
			err := kubeclient.Create(ctx, altTemplate)
			Expect(err).To(BeNil(), "failed to create altTemplate CR")

			// Create and defer deletion of an alternative combination with the alternative template referenced
			altCombination := validCombinationCR.DeepCopy()
			altCombination.Spec.Template = "alttemplate"
			altCombination.Name = "altcombination"
			err = kubeclient.Create(ctx, altCombination)
			Expect(err).To(BeNil(), "failed to create altCombination CR")

			// Update template to trigger a requeue event
			validTemplateCRCopy.Spec.Body = validUpdatedTemplateCR.Spec.Body
			err = kubeclient.Update(ctx, validTemplateCRCopy)
			Expect(err).To(BeNil(), "failed to update original template CR")

			Eventually(func(g Gomega) error {
				var retrievedCombination v1alpha1.Combination
				err = kubeclient.Get(ctx, types.NamespacedName{Name: altCombination.Name}, &retrievedCombination)

				var conditionReasons []string
				for _, condition := range retrievedCombination.Status.Conditions {
					conditionReasons = append(conditionReasons, condition.Reason)
				}

				g.Expect(conditionReasons).To(ContainElement(v1alpha1.ReasonProcessed))
				g.Expect(retrievedCombination.Status.Evaluations).To(ContainElements(expectedEvaluations))

				return err
			}).Should(Succeed())

			Expect(kubeclient.Delete(ctx, altTemplate)).To(BeNil())
			Expect(kubeclient.Delete(ctx, altCombination)).To(BeNil())
		})
	})

	When("given healthy input and a non-existent template", func() {
		var ctx context.Context
		var validCombinationCRCopy *v1alpha1.Combination

		BeforeEach(func() {
			ctx = context.Background()

			// Create a valid CR
			validCombinationCRCopy = validCombinationCR.DeepCopy()

			// Create a valid combination
			validCombinationCRCopy.Spec.Template = "doesnotexist"
			err := kubeclient.Create(ctx, validCombinationCRCopy)
			Expect(err).To(BeNil(), "failed to create combination CR")
		})

		AfterEach(func() {
			err := kubeclient.Delete(ctx, validCombinationCRCopy)
			Expect(err).To(BeNil(), "failed to clean-up combination CR after test")
			ctx.Done()
		})

		It("should fail and output a TemplateNotFound status", func() {
			Eventually(func() ([]string, error) {
				var retrievedCombination v1alpha1.Combination
				err := kubeclient.Get(ctx, types.NamespacedName{Name: validCombinationCRCopy.Name}, &retrievedCombination)

				var conditionReasons []string
				for _, condition := range retrievedCombination.Status.Conditions {
					conditionReasons = append(conditionReasons, condition.Reason)
				}

				return conditionReasons, err
			}).Should(ContainElement(v1alpha1.ReasonTemplateNotFound))
		})
	})

	rawConfigMap :=
		`
---
apiVersion: v1
data: 
  key: VALUE
kind: ConfigMap
metadata: 
  name: cm
  namespace: default
`

	configMapTemplate := v1alpha1.Template{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "cm-template",
		},
		Spec: v1alpha1.TemplateSpec{
			Body:       rawConfigMap,
			Parameters: []string{"VALUE"},
		},
	}

	configMapCombo := v1alpha1.Combination{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "cm-combination",
		},
		Spec: v1alpha1.CombinationSpec{
			Arguments: []v1alpha1.Argument{
				{
					Key: "VALUE",
					Values: []string{
						"mango",
						"fox",
					},
				},
			},
			Apply: true,
		},
	}

	When("given a configmap template with a combination that has apply set", func() {
		var ctx context.Context

		BeforeEach(func() {
			ctx = context.Background()

			err := kubeclient.Create(ctx, &configMapTemplate)
			Expect(err).To(BeNil(), "failed to create template CR")

			configMapCombo.Spec.Template = configMapTemplate.Name
			err = kubeclient.Create(ctx, &configMapCombo)
			Expect(err).To(BeNil(), "failed to create combo CR")
		})

		AfterEach(func() {
			err := kubeclient.Delete(ctx, &configMapTemplate)
			Expect(err).To(BeNil(), "failed to cleanup template CR")

			err = kubeclient.Delete(ctx, &configMapCombo)
			Expect(err).To(BeNil(), "failed to cleanup combo CR")
		})

		It("should fail to create the resource on-cluster due to missing RBAC", func() {
			Eventually(func(g Gomega) error {
				var retrievedCombination v1alpha1.Combination
				err := kubeclient.Get(ctx, types.NamespacedName{Name: configMapCombo.Name}, &retrievedCombination)

				var conditionReasons []string
				for _, condition := range retrievedCombination.Status.Conditions {
					conditionReasons = append(conditionReasons, condition.Reason)
				}
				g.Expect(conditionReasons).To(ContainElement(v1alpha1.ReasonProcessed))
				g.Expect(conditionReasons).To(ContainElement(v1alpha1.ReasonApplyFailed))

				var conditionMessages []string
				for _, condition := range retrievedCombination.Status.Conditions {
					conditionMessages = append(conditionMessages, condition.Message)
				}
				g.Expect(conditionMessages).To(ContainElement("failed to apply manifest cm to cluster: configmaps is forbidden: User \"system:serviceaccount:combo:combo-operator\" cannot create resource \"configmaps\" in API group \"\" in the namespace \"default\""))

				return err
			}).Should(Succeed())
		})
	})

	rawTemplate :=
		`
--- 
apiVersion: combo.io/v1alpha1
kind: Template
metadata: 
  labels: 
    environment: VALUE
  name: embedded-template-VALUE
spec: 
  body: |
      ---
      apiVersion: rbac.authorization.k8s.io/v1
      kind: RoleBinding
      metadata:
        name: feature-controller
        namespace: TARGET_NAMESPACE
  parameters: 
    - TARGET_NAMESPACE
`

	embeddedTemplate := v1alpha1.Template{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "em-template",
		},
		Spec: v1alpha1.TemplateSpec{
			Body:       rawTemplate,
			Parameters: []string{"VALUE"},
		},
	}

	embeddedCombo := v1alpha1.Combination{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "em-combination",
		},
		Spec: v1alpha1.CombinationSpec{
			Arguments: []v1alpha1.Argument{
				{
					Key: "VALUE",
					Values: []string{
						"dev",
						"prod",
					},
				},
			},
			Apply: true,
		},
	}

	// This test ensures that when combo has permissions to create objects in a certain API group, it can do so successfully.
	// Since by default combo only has permissions for its own APIs, this test uses a Template type in the template body
	// so there are adequate permissions to create the resource on-cluster.
	When("given a template with another template embedded in it that has apply set", func() {
		var ctx context.Context

		BeforeEach(func() {
			ctx = context.Background()

			// before creation, modify the combo operator cluster role to allow for the creation of objects in the combo API group
			// this simulates the cluster-admin modifying the permissions granted to combo before applying resources.
			var cr rbacv1.ClusterRole
			err := kubeclient.Get(ctx, types.NamespacedName{Name: "combo-operator"}, &cr)
			Expect(err).To(BeNil(), "failed to get combo clusterrole")
			cr.Rules[0].Verbs = append(cr.Rules[0].Verbs, "create")
			err = kubeclient.Update(ctx, &cr)
			Expect(err).To(BeNil(), "failed to update combo clusterrole")

			err = kubeclient.Create(ctx, &embeddedTemplate)
			Expect(err).To(BeNil(), "failed to create template CR")

			embeddedCombo.Spec.Template = embeddedTemplate.Name
			err = kubeclient.Create(ctx, &embeddedCombo)
			Expect(err).To(BeNil(), "failed to create combo CR")
		})

		AfterEach(func() {
			err := kubeclient.Delete(ctx, &embeddedTemplate)
			Expect(err).To(BeNil(), "failed to cleanup template CR")

			err = kubeclient.Delete(ctx, &embeddedCombo)
			Expect(err).To(BeNil(), "failed to cleanup combo CR")
		})

		It("should successfully create the templated evaluation on-cluster", func() {
			Eventually(func(g Gomega) error {
				var retrievedCombination v1alpha1.Combination
				err := kubeclient.Get(ctx, types.NamespacedName{Name: embeddedCombo.Name}, &retrievedCombination)
				Expect(err).To(BeNil(), "failed to get combo CR")

				var conditionReasons []string
				for _, condition := range retrievedCombination.Status.Conditions {
					conditionReasons = append(conditionReasons, condition.Reason)
				}
				g.Expect(conditionReasons).To(ContainElement(v1alpha1.ReasonProcessed))
				// g.Expect(conditionReasons).To(ContainElement(v1alpha1.ReasonApplySucceeded))

				// fetch created template on-cluster and ensure templating worked as expected
				var embeddedTemplate v1alpha1.Template
				err = kubeclient.Get(ctx, types.NamespacedName{Name: "embedded-template-dev"}, &embeddedTemplate)

				g.Expect(embeddedTemplate.GetLabels()["environment"]).To(Equal("dev"))

				return err
			}).Should(Succeed())
		})

	})
})
