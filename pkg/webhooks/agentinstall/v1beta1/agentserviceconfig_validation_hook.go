package v1beta1

import (
	"fmt"
	"net/http"

	"github.com/openshift/assisted-service/api/v1beta1"
	log "github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	agentServiceConfigResource = "agentserviceconfigs"
	agentServiceConfigGroup    = "admission.agentinstall.openshift.io"
	agentServiceConfigVersion  = "v1beta1"
)

type AgentServiceConfigValidatingAdmissionHook struct {
	decoder *admission.Decoder
}

// NewAgentServiceConfigValidatingAdmissionHook constructs a new AgentServiceConfigValidatingAdmissionHook
func NewAgentServiceConfigValidatingAdmissionHook(decoder *admission.Decoder) *AgentServiceConfigValidatingAdmissionHook {
	return &AgentServiceConfigValidatingAdmissionHook{decoder: decoder}
}

func (a *AgentServiceConfigValidatingAdmissionHook) ValidatingResource() (plural schema.GroupVersionResource, singular string) {
	log.WithFields(log.Fields{
		"group":    agentServiceConfigGroup,
		"version":  agentServiceConfigVersion,
		"resource": agentServiceConfigResource,
	}).Info("Registering validation REST resource")
	return schema.GroupVersionResource{
			Group:    agentServiceConfigGroup,
			Version:  agentServiceConfigVersion,
			Resource: agentServiceConfigResource,
		},
		"agentserviceconfigvalidator"
}

// Initialize is called by generic-admission-server on startup to setup any special initialization that your webhook needs.
func (a *AgentServiceConfigValidatingAdmissionHook) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	log.WithFields(log.Fields{
		"group":    agentServiceConfigGroup,
		"version":  agentServiceConfigVersion,
		"resource": agentServiceConfigResource,
	}).Info("Initializing validation REST resource")
	return nil // No initialization needed right now.
}

// Validate is called by generic-admission-server when the registered REST resource above is called with an admission request.
// Usually it's the kube apiserver that is making the admission validation request.
func (a *AgentServiceConfigValidatingAdmissionHook) Validate(admissionSpec *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	contextLogger := log.WithFields(log.Fields{
		"operation": admissionSpec.Operation,
		"group":     admissionSpec.Resource.Group,
		"version":   admissionSpec.Resource.Version,
		"resource":  admissionSpec.Resource.Resource,
		"method":    "Validate",
	})

	if !a.shouldValidate(admissionSpec) {
		contextLogger.Info("Skipping validation for request")
		// The request object isn't something that this validator should validate.
		// Therefore, we say that it's allowed.
		return &admissionv1.AdmissionResponse{
			Allowed: true,
		}
	}

	contextLogger.Info("Validating request")

	switch admissionSpec.Operation {
	case admissionv1.Create, admissionv1.Update:
		agentserviceconfig := &v1beta1.AgentServiceConfig{}
		if err := a.decoder.DecodeRaw(admissionSpec.Object, agentserviceconfig); err != nil {
			message := ""
			contextLogger.Error(message)
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Status: metav1.StatusFailure, Code: http.StatusBadRequest, Reason: metav1.StatusReasonBadRequest,
					Message: message,
				},
			}
		}

		seen := make(map[string]bool)
		for _, osImage := range agentserviceconfig.Spec.OSImages {
			if seen[osImage.Version] {
				message := fmt.Sprintf("Found duplicate OSImages entry: %s", osImage.Version)
				contextLogger.Infof("Failed validation: %v", message)
				contextLogger.Error(message)
				return &admissionv1.AdmissionResponse{
					Allowed: false,
					Result: &metav1.Status{
						Status: metav1.StatusFailure, Code: http.StatusBadRequest, Reason: metav1.StatusReasonBadRequest,
						Message: message,
					},
				}
			}
			seen[osImage.Version] = true
		}
	}

	return &admissionv1.AdmissionResponse{
		Allowed: true,
	}
}

// shouldValidate explicitly checks if the request should be validated. For example, this webhook may have accidentally been registered to check
// the validity of some other type of object with a different GVR.
func (a *AgentServiceConfigValidatingAdmissionHook) shouldValidate(admissionSpec *admissionv1.AdmissionRequest) bool {
	contextLogger := log.WithFields(log.Fields{
		"operation": admissionSpec.Operation,
		"group":     admissionSpec.Resource.Group,
		"version":   admissionSpec.Resource.Version,
		"resource":  admissionSpec.Resource.Resource,
		"method":    "shouldValidate",
	})

	if admissionSpec.Resource.Group != v1beta1.Group {
		contextLogger.Debug("Returning False, not our group")
		return false
	}

	if admissionSpec.Resource.Version != v1beta1.Version {
		contextLogger.Debug("Returning False, it's our group, but not the right version")
		return false
	}

	if admissionSpec.Resource.Resource != agentServiceConfigResource {
		contextLogger.Debug("Returning False, it's our group and version, but not the right resource")
		return false
	}

	// If we get here, then we're supposed to validate the object.
	contextLogger.Debug("Returning True, passed all prerequisites.")
	return true
}
