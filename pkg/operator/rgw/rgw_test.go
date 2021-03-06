/*
Copyright 2016 The Rook Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package rgw

import (
	"fmt"
	"os"
	"strings"
	"testing"

	testceph "github.com/rook/rook/pkg/cephmgr/client/test"
	testclient "github.com/rook/rook/pkg/cephmgr/client/test"
	cephrgw "github.com/rook/rook/pkg/cephmgr/rgw"
	testop "github.com/rook/rook/pkg/operator/test"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
)

func TestStartRGW(t *testing.T) {
	clientset := testop.New(3)
	info := testop.CreateClusterInfo(1)
	factory := &testclient.MockConnectionFactory{Fsid: "fsid", SecretKey: "mysecret"}
	conn, _ := factory.NewConnWithClusterAndUser(info.Name, "user")
	conn.(*testceph.MockConnection).MockMonCommand = func(buf []byte) (buffer []byte, info string, err error) {
		response := "{\"key\":\"mysecurekey\"}"
		return []byte(response), "", nil
	}

	c := New(clientset, factory, "ns", "version")
	c.dataDir = "/tmp/rgwtest"
	defer os.RemoveAll(c.dataDir)

	// start a basic cluster
	err := c.Start(info)
	assert.Nil(t, err)

	validateStart(t, c, clientset)

	// starting again should be a no-op
	err = c.Start(info)
	assert.Nil(t, err)

	validateStart(t, c, clientset)
}

func validateStart(t *testing.T, c *Cluster, clientset *fake.Clientset) {

	r, err := clientset.ExtensionsV1beta1().Deployments(c.Namespace).Get("rgw")
	assert.Nil(t, err)
	assert.Equal(t, "rgw", r.Name)

	s, err := clientset.CoreV1().Services(c.Namespace).Get("rgw")
	assert.Nil(t, err)
	assert.Equal(t, "rgw", s.Name)

	secret, err := clientset.CoreV1().Secrets(c.Namespace).Get("rgw")
	assert.Nil(t, err)
	assert.Equal(t, "rgw", secret.Name)
	assert.Equal(t, 1, len(secret.StringData))

}

func TestPodSpecs(t *testing.T) {
	c := New(nil, nil, "ns", "myversion")

	d := c.makeDeployment()
	assert.NotNil(t, d)
	assert.Equal(t, "rgw", d.Name)
	assert.Equal(t, v1.RestartPolicyAlways, d.Spec.Template.Spec.RestartPolicy)
	assert.Equal(t, 1, len(d.Spec.Template.Spec.Volumes))
	assert.Equal(t, "rook-data", d.Spec.Template.Spec.Volumes[0].Name)

	assert.Equal(t, "rgw", d.ObjectMeta.Name)
	assert.Equal(t, "rgw", d.Spec.Template.ObjectMeta.Labels["app"])
	assert.Equal(t, c.Namespace, d.Spec.Template.ObjectMeta.Labels["rook_cluster"])
	assert.Equal(t, 0, len(d.ObjectMeta.Annotations))

	cont := d.Spec.Template.Spec.Containers[0]
	assert.Equal(t, "quay.io/rook/rookd:myversion", cont.Image)
	assert.Equal(t, 1, len(cont.VolumeMounts))
	assert.Equal(t, 5, len(cont.Env))

	expectedCommand := fmt.Sprintf("/usr/bin/rookd rgw --data-dir=/var/lib/rook --rgw-port=%d --rgw-host=%s",
		cephrgw.RGWPort, cephrgw.DNSName)

	assert.NotEqual(t, -1, strings.Index(cont.Command[2], expectedCommand), cont.Command[2])
}
