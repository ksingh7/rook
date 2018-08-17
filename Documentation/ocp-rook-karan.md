
```
oc login -u system:admin

git clone https://github.com/rook/rook.git
cd rook/cluster/examples/kubernetes/ceph

oc create -f scc.yaml


- Edit https://github.com/rook/rook/blob/master/cluster/examples/kubernetes/ceph/operator.yaml

        - name: ROOK_HOSTPATH_REQUIRES_PRIVILEGED
          value: "true"

oc create -f operator.yaml
oc -n rook-ceph-system get pod -o wide

oc get nodes --show-labels

oc label node node04.internal.aws.testdrive.openshift.com role=storage-node
oc label node node05.internal.aws.testdrive.openshift.com role=storage-node
oc label node node06.internal.aws.testdrive.openshift.com role=storage-node

oc create -f cluster-v1.yaml

oc adm pod-network join-projects rook-ceph-system --to=rook-ceph
oc get netnamespaces rook-ceph rook-ceph-system

oc -n rook-ceph get pods
oc -n rook-ceph get service
oc expose service rook-ceph-mgr-dashboard
oc get route

oc create -f toolbox.yaml
oc get pods -n rook-ceph

oc rsh -n rook-ceph rook-ceph-tools
```
### Edit object.yaml , set port to 8080
```
oc create -f object.yaml
oc -n rook-ceph get pod -l app=rook-ceph-rgw
oc -n rook-ceph get service rook-ceph-rgw-my-store
oc expose service rook-ceph-rgw-my-store
oc get route


oc scale --replicas=2 deploy/rook-ceph-rgw-my-store
oc describe svc/rook-ceph-rgw-my-store -n rook-ceph
```
