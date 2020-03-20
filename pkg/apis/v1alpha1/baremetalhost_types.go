package v1aplha1


//import (
//	"fmt"
//)

type BareMetalHostBmcCred struct {
}

type BareMetalHostImage struct {
}

type BareMetalHostSpec struct {
	BmcCred BareMetalHostBmcCred `json:"spec,omitempty"`
	Image BareMetalHostImage `json:"spec,omitempty"`

}

type BareMetalHost struct {

	Spec   BareMetalHostSpec   `json:"spec,omitempty"`
//	Status BareMetalHostStatus `json:"status,omitempty"`
}
