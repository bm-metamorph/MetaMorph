package kubernetes

import (
  // "github.com/redfishProvisioner/kubernetes/base"
  "fmt"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  corev1 "k8s.io/client-go/kubernetes/typed/batch/v1"
  batchv1 "k8s.io/api/batch/v1"
)

type JobClient struct {
    j corev1.JobInterface
}

func NewJob(namespace string) *JobClient{
    clientset := New().kube
    return &JobClient{j: clientset.BatchV1().Jobs("metalkube")}
}

func (j *JobClient) CreateJob(job *batchv1.Job) bool {
    result, _ := j.j.Create(job)
    status := false
    if result.GetObjectMeta().GetName() != ""{
      for {
          job, err := j.j.Get(result.GetObjectMeta().GetName(), metav1.GetOptions{})
          if err != nil {
              fmt.Println("Unable to fetch job")
              break
          }
          if job.Status.Failed > 0 {
            fmt.Println("job failed")
            break
          }
          if job.Status.Succeeded > 0 {
            fmt.Println("job success")
            status = true
            break
          }
      }
    }
    return status
}

func (j *JobClient) DeleteJob(name, label_selector string) bool {
    _ = j.j.Delete(name, &metav1.DeleteOptions{})
    return true
}

func (j *JobClient) GetJobs(name string, label_selector map[string]string) bool {
    return true
}

func (j *JobClient) GetJobDetails(name string) bool {
    return true
}
