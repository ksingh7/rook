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

Some of the code below came from https://github.com/coreos/etcd-operator
which also has the apache 2.0 license.
*/
package cluster

import (
	"fmt"

	rookclient "github.com/rook/rook/pkg/rook/client"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/rook/rook/pkg/model"
)

const (
	replicatedType  = "replicated"
	erasureCodeType = "erasure-coded"
)

type Pool struct {
	Metadata v1.ObjectMeta `json:"metadata,omitempty"`
	PoolSpec `json:"spec"`
}

// Instantiate a new pool
func NewPool(spec PoolSpec) *Pool {
	return &Pool{PoolSpec: spec}
}

// Create the pool
func (p *Pool) Create(rclient rookclient.RookRestClient) error {
	// validate the pool settings
	if err := p.validate(); err != nil {
		return fmt.Errorf("invalid pool %s arguments. %+v", p.PoolSpec.Name, err)
	}

	// check if the pool already exists
	exists, err := p.exists(rclient)
	if err == nil && exists {
		logger.Infof("pool %s already exists in namespace %s", p.PoolSpec.Name, p.PoolSpec.Namespace)
		return nil
	}

	// create the pool
	pool := model.Pool{Name: p.PoolSpec.Name}
	r := p.replication()
	if r != nil {
		logger.Infof("creating pool %s in namespace %s with replicas %d", p.PoolSpec.Name, p.Namespace, r.Count)
		pool.ReplicationConfig.Size = r.Count
		pool.Type = model.Replicated
	} else {
		ec := p.erasureCode()
		logger.Infof("creating pool %s in namespace %s. coding chunks = %d, data chunks = %d", p.PoolSpec.Name, p.Namespace, ec.CodingChunks, ec.DataChunks)
		pool.ErasureCodedConfig.CodingChunkCount = ec.CodingChunks
		pool.ErasureCodedConfig.DataChunkCount = ec.DataChunks
		pool.Type = model.ErasureCoded
	}

	info, err := rclient.CreatePool(pool)
	if err != nil {
		return fmt.Errorf("failed to create pool %s. %+v", p.PoolSpec.Name, err)
	}

	logger.Infof("created pool %s. %s", p.PoolSpec.Name, info)
	return nil
}

// Delete the pool
func (p *Pool) Delete(rclient rookclient.RookRestClient) error {
	// check if the pool  exists
	exists, err := p.exists(rclient)
	if err == nil && !exists {
		return nil
	}

	logger.Infof("TODO: delete pool %s from namespace %s", p.PoolSpec.Name, p.PoolSpec.Namespace)
	//return p.client.DeletePool(p.PoolSpec.Name)
	return nil
}

// Check if the pool exists
func (p *Pool) exists(rclient rookclient.RookRestClient) (bool, error) {
	pools, err := rclient.GetPools()
	if err != nil {
		return false, err
	}
	for _, pool := range pools {
		if pool.Name == p.PoolSpec.Name {
			return true, nil
		}
	}
	return false, nil
}

// Validate the pool arguments
func (p *Pool) validate() error {
	if p.PoolSpec.Name == "" {
		return fmt.Errorf("missing name")
	}
	if p.PoolSpec.Namespace == "" {
		return fmt.Errorf("missing namespace")
	}
	if p.replication() != nil && p.erasureCode() != nil {
		return fmt.Errorf("both replication and erasure code settings cannot be specified")
	}
	if p.replication() == nil && p.erasureCode() == nil {
		return fmt.Errorf("neither replication nor erasure code settings were specified")
	}
	return nil
}

func (p *Pool) replication() *ReplicationSpec {
	if p.PoolSpec.Replication.Count > 0 {
		return &p.PoolSpec.Replication
	}
	return nil
}

func (p *Pool) erasureCode() *ErasureCodeSpec {
	ec := &p.PoolSpec.ErasureCoding
	if ec.CodingChunks > 0 || ec.DataChunks > 0 {
		return ec
	}
	return nil
}
