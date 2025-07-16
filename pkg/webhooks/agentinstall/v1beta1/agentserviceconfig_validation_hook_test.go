package v1beta1

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admission/v1"
)

var _ = Describe("agent serivce config validate", func() {
	tests := []struct {
		name      string
		allowed   bool
		operation admissionv1.Operation
	}{}

	for i := range tests {
		tc := tests[i]
		It(tc.name, func() {
			hook := NewAgentServiceConfigValidatingAdmissionHook(createDecoder())

			request := &admissionv1.AdmissionRequest{
				Operation: tc.operation,
			}
			response := hook.Validate(request)
			Expect(response.Allowed).To(Equal(tc.allowed))
		})
	}
})
