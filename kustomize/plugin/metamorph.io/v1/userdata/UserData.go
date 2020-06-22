// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"text/template"
	"encoding/base64"
	"os"
	"io/ioutil"
	"fmt"


	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/yaml"
)

// A simple generator example.  Makes one service.
type plugin struct {
	rf               *resmap.Factory
	types.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Namespace        string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	NetworkConfig    string `json:"network_config,omitempty" yaml:"network_config,omitempty"`
	Resources        []string `json:"resources,omitempty" yaml:"resources,omitempty"`
	UserData         string
}

//nolint: golint
//noinspection GoUnusedGlobalVariable
var KustomizePlugin plugin


const tmpl = `

kind: Secret
apiVersion: v1
type: Opaque
metadata:
  namespace: {{.Namespace}}
  name: {{.Name}}
data:
  userData: {{.UserData}}
`


func check(e error) {
    if e != nil {
        panic(e)
    }
}

func (p *plugin) Config(h *resmap.PluginHelpers, config []byte) error {
	p.rf = h.ResmapFactory()
	return yaml.Unmarshal(config, p)
}

func (p *plugin) Generate() (resmap.ResMap, error){
	path, err := os.Getwd()
        check(err)
	network_config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s",path, p.NetworkConfig))
        check(err)

        var globalUserData map[string]interface{}
	gud, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, p.Resources[0]))
        check(err)
	err = yaml.Unmarshal(gud, &globalUserData)
        check(err)


        var nodeUserData map[string]interface{}
	nud, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, p.Resources[1]))
        check(err)
	err = yaml.Unmarshal(nud, &nodeUserData)
        check(err)

	for k, v := range nodeUserData {
          globalUserData[k] = v
        }

	globalUserData["NetworkConfig"] = string(base64.StdEncoding.EncodeToString(network_config))

	mergedUserData, err := yaml.Marshal(globalUserData)
        check(err)

        mergedUserDataJSON, err := yaml.YAMLToJSON(mergedUserData)

	p.UserData= string(base64.StdEncoding.EncodeToString(mergedUserDataJSON))

	var buf bytes.Buffer
	temp := template.Must(template.New("tmpl").Parse(tmpl))
	err = temp.Execute(&buf, p)
	if err != nil {
		return nil, err
	}
	return p.rf.NewResMapFromBytes(buf.Bytes())
}
