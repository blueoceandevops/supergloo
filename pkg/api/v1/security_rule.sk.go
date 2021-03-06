// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"sort"

	"github.com/solo-io/go-utils/hashutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewSecurityRule(namespace, name string) *SecurityRule {
	securityrule := &SecurityRule{}
	securityrule.SetMetadata(core.Metadata{
		Name:      name,
		Namespace: namespace,
	})
	return securityrule
}

func (r *SecurityRule) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

func (r *SecurityRule) SetStatus(status core.Status) {
	r.Status = status
}

func (r *SecurityRule) Hash() uint64 {
	metaCopy := r.GetMetadata()
	metaCopy.ResourceVersion = ""
	return hashutils.HashAll(
		metaCopy,
		r.TargetMesh,
		r.SourceSelector,
		r.DestinationSelector,
		r.AllowedPaths,
		r.AllowedMethods,
	)
}

type SecurityRuleList []*SecurityRule

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list SecurityRuleList) Find(namespace, name string) (*SecurityRule, error) {
	for _, securityRule := range list {
		if securityRule.GetMetadata().Name == name {
			if namespace == "" || securityRule.GetMetadata().Namespace == namespace {
				return securityRule, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find securityRule %v.%v", namespace, name)
}

func (list SecurityRuleList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, securityRule := range list {
		ress = append(ress, securityRule)
	}
	return ress
}

func (list SecurityRuleList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, securityRule := range list {
		ress = append(ress, securityRule)
	}
	return ress
}

func (list SecurityRuleList) Names() []string {
	var names []string
	for _, securityRule := range list {
		names = append(names, securityRule.GetMetadata().Name)
	}
	return names
}

func (list SecurityRuleList) NamespacesDotNames() []string {
	var names []string
	for _, securityRule := range list {
		names = append(names, securityRule.GetMetadata().Namespace+"."+securityRule.GetMetadata().Name)
	}
	return names
}

func (list SecurityRuleList) Sort() SecurityRuleList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].GetMetadata().Less(list[j].GetMetadata())
	})
	return list
}

func (list SecurityRuleList) Clone() SecurityRuleList {
	var securityRuleList SecurityRuleList
	for _, securityRule := range list {
		securityRuleList = append(securityRuleList, resources.Clone(securityRule).(*SecurityRule))
	}
	return securityRuleList
}

func (list SecurityRuleList) Each(f func(element *SecurityRule)) {
	for _, securityRule := range list {
		f(securityRule)
	}
}

func (list SecurityRuleList) EachResource(f func(element resources.Resource)) {
	for _, securityRule := range list {
		f(securityRule)
	}
}

func (list SecurityRuleList) AsInterfaces() []interface{} {
	var asInterfaces []interface{}
	list.Each(func(element *SecurityRule) {
		asInterfaces = append(asInterfaces, element)
	})
	return asInterfaces
}

var _ resources.Resource = &SecurityRule{}

// Kubernetes Adapter for SecurityRule

func (o *SecurityRule) GetObjectKind() schema.ObjectKind {
	t := SecurityRuleCrd.TypeMeta()
	return &t
}

func (o *SecurityRule) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*SecurityRule)
}

var SecurityRuleCrd = crd.NewCrd("supergloo.solo.io",
	"securityrules",
	"supergloo.solo.io",
	"v1",
	"SecurityRule",
	"sr",
	false,
	&SecurityRule{})
