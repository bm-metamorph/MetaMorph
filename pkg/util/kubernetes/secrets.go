package kubernetes 

import (
  // "github.com/redfishProvisioner/kubernetes/base"
  // "fmt"
   metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
   corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
   apiv1 "k8s.io/api/core/v1"
)

type SecretClient struct {
    s corev1.SecretInterface
}

func NewSecret(namespace string) *SecretClient{
    clientset := New().kube
    return &SecretClient{s: clientset.CoreV1().Secrets("metalkube")}
}

func (s *SecretClient) CreateSecret(Secret *apiv1.Secret) bool {
    result, _ := s.s.Create(Secret)
    if result.GetObjectMeta().GetName() != ""{
        return true
    } else{
            return false
    }
}

func (s *SecretClient) DeleteSecret(name string, label_selector map[string]string) bool {
    _ = s.s.Delete(name, &metav1.DeleteOptions{})
    return true
}

func (s *SecretClient) GetSecrets(name string, label_selector map[string]string) bool {
    return true
}

func (s *SecretClient) GetSecretDetails(name string) bool {
    return true
}
