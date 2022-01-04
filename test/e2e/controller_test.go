package e2e

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/combo/api/v1alpha1"
	"github.com/operator-framework/combo/pkg/controller"
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

		It("should have the ReferencedTemplate label correctly applied on the combination", func() {
			Eventually(func(g Gomega) error {
				var retrievedCombination v1alpha1.Combination
				err := kubeclient.Get(ctx, types.NamespacedName{Name: validCombinationCRCopy.Name}, &retrievedCombination)

				g.Expect(retrievedCombination.Labels).To(HaveKeyWithValue(controller.ReferencedTemplateLabel, validTemplateCRCopy.Name))

				return err
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
})
