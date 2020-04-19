package kubernetes

import (
   "fmt"
   apiv1 "k8s.io/api/core/v1"
   corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
   metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMapClient struct {
    c corev1.ConfigMapInterface
}

func NewConfigMap(namespace string) *ConfigMapClient{
    clientset := New().kube
    cm := ConfigMapClient{c: clientset.CoreV1().ConfigMaps("metalkube")}
    return &cm
}

func (c *ConfigMapClient) CreateConfigMap(ConfigMap *apiv1.ConfigMap) bool {
    result, _ := c.c.Create(ConfigMap)
    if result.GetObjectMeta().GetName() != ""{
        return true
    } else{
            return false
    }
    return false
}

func (c *ConfigMapClient) DeleteConfigMap(name, label_selector string) bool {
    fmt.Println(label_selector)
    _ = c.c.Delete(name, &metav1.DeleteOptions{})
    return true
}

func (c *ConfigMapClient) GetConfigMaps(name string, label_selector map[string]string) *apiv1.ConfigMap {
  result, _ := c.c.Get(name, metav1.GetOptions{})
  return result
}

func (c *ConfigMapClient) GetConfigMapDetails(name string) bool {
  return true
}
