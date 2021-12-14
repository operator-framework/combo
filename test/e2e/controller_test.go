package e2e

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/combo/api/v1alpha1"
	comboConditions "github.com/operator-framework/combo/pkg/conditions"
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

			// Create and defer deletion of a valid template
			err := kubeclient.Create(ctx, validTemplateCRCopy)
			Expect(err).To(BeNil(), "failed to create template CR")

			// Create and defer deletion of a valid combination
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

				g.Expect(conditionReasons).To(ContainElement(comboConditions.ProccessedCombinationsCondition.Reason))
				g.Expect(retrievedCombination.Status.Evaluations).To(ContainElements(expectedEvaluations))

				return err
			}).Should(BeZero())
		})

		It("should reevaluate whenever its associated template gets updated", func() {
			// Give up to a minute (default eventually timeout) for the combination to update after a template is updated
			Eventually(func(g Gomega) error {
				// Create and defer deletion of a valid template
				validTemplateCRCopy.Spec.Body = validUpdatedTemplateCR.Spec.Body
				err := kubeclient.Update(ctx, validTemplateCRCopy)
				g.Expect(err).To(BeNil(), "failed to create template CR")

				var retrievedCombination v1alpha1.Combination
				err = kubeclient.Get(ctx, types.NamespacedName{Name: validCombinationCRCopy.Name}, &retrievedCombination)

				var conditionReasons []string
				for _, condition := range retrievedCombination.Status.Conditions {
					conditionReasons = append(conditionReasons, condition.Reason)
				}

				g.Expect(conditionReasons).To(ContainElement(comboConditions.ProccessedCombinationsCondition.Reason))
				g.Expect(retrievedCombination.Status.Evaluations).To(ContainElements(expectedUpdatedEvaluations))

				return err
			}).Should(BeZero())
		})

	})

	When("given healthy input and a non-existent template", func() {
		var ctx context.Context
		var validCombinationCRCopy *v1alpha1.Combination

		BeforeEach(func() {
			ctx = context.Background()

			// Create copies of valid CR
			validCombinationCRCopy = validCombinationCR.DeepCopy()

			// Create and defer deletion of a valid combination
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
			// Give up to a minute (default eventually timeout) for the combination to process properly
			Eventually(func() ([]string, error) {
				var retrievedCombination v1alpha1.Combination
				err := kubeclient.Get(ctx, types.NamespacedName{Name: validCombinationCRCopy.Name}, &retrievedCombination)

				var conditionReasons []string
				for _, condition := range retrievedCombination.Status.Conditions {
					conditionReasons = append(conditionReasons, condition.Reason)
				}

				return conditionReasons, err
			}).Should(ContainElement(comboConditions.TemplateNotFoundCondition.Reason))
		})
	})
})
