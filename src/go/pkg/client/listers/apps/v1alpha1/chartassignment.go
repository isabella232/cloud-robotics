// Copyright 2021 The Cloud Robotics Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/SAP/cloud-robotics/src/go/pkg/apis/apps/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ChartAssignmentLister helps list ChartAssignments.
// All objects returned here must be treated as read-only.
type ChartAssignmentLister interface {
	// List lists all ChartAssignments in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ChartAssignment, err error)
	// ChartAssignments returns an object that can list and get ChartAssignments.
	ChartAssignments(namespace string) ChartAssignmentNamespaceLister
	ChartAssignmentListerExpansion
}

// chartAssignmentLister implements the ChartAssignmentLister interface.
type chartAssignmentLister struct {
	indexer cache.Indexer
}

// NewChartAssignmentLister returns a new ChartAssignmentLister.
func NewChartAssignmentLister(indexer cache.Indexer) ChartAssignmentLister {
	return &chartAssignmentLister{indexer: indexer}
}

// List lists all ChartAssignments in the indexer.
func (s *chartAssignmentLister) List(selector labels.Selector) (ret []*v1alpha1.ChartAssignment, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ChartAssignment))
	})
	return ret, err
}

// ChartAssignments returns an object that can list and get ChartAssignments.
func (s *chartAssignmentLister) ChartAssignments(namespace string) ChartAssignmentNamespaceLister {
	return chartAssignmentNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ChartAssignmentNamespaceLister helps list and get ChartAssignments.
// All objects returned here must be treated as read-only.
type ChartAssignmentNamespaceLister interface {
	// List lists all ChartAssignments in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ChartAssignment, err error)
	// Get retrieves the ChartAssignment from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ChartAssignment, error)
	ChartAssignmentNamespaceListerExpansion
}

// chartAssignmentNamespaceLister implements the ChartAssignmentNamespaceLister
// interface.
type chartAssignmentNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ChartAssignments in the indexer for a given namespace.
func (s chartAssignmentNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ChartAssignment, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ChartAssignment))
	})
	return ret, err
}

// Get retrieves the ChartAssignment from the indexer for a given namespace and name.
func (s chartAssignmentNamespaceLister) Get(name string) (*v1alpha1.ChartAssignment, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("chartassignment"), name)
	}
	return obj.(*v1alpha1.ChartAssignment), nil
}
