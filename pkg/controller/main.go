package controller

import (
	"fmt"
	metamorphv1alpha1 "github.com/metamorph/pkg/apis/v1alpha1"
)

func  ControllerInit() {

	host := metamorphv1alpha1.BareMetalHost{}
	fmt.Println(host)
	fmt.Println("From Controller")
}
