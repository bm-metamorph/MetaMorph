package kubernetes

import (
  //"github.com/redfishProvisioner/kubernetes/base"
   // "fmt"
   metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
   corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
   apiv1 "k8s.io/api/core/v1"
)

type PodClient struct {
    p corev1.PodInterface
}

func NewPod(namespace string) *PodClient{
    clientset := New().kube
    return &PodClient{p: clientset.CoreV1().Pods("metalkube")}
}

func (p *PodClient) CreatePod(Pod *apiv1.Pod) bool {
    result, _ := p.p.Create(Pod)
    if result.GetObjectMeta().GetName() != ""{
        return true
    } else{
            return false
    }
}

func (p *PodClient) DeletePod(name string, label_selector map[string]string) bool {
    _ = p.p.Delete(name, &metav1.DeleteOptions{})
    return true
}

func (p *PodClient) GetPods(name string, label_selector map[string]string) bool {
    return true
}

func (p *PodClient) GetPodDetails(name string) bool {
    return true
}
