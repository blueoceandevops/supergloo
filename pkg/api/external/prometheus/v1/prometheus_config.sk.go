// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"sort"

	github_com_solo_io_supergloo_api_external_prometheus "github.com/solo-io/supergloo/api/external/prometheus"

	"github.com/solo-io/go-utils/hashutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
)

func NewPrometheusConfig(namespace, name string) *PrometheusConfig {
	prometheusconfig := &PrometheusConfig{}
	prometheusconfig.PrometheusConfig.SetMetadata(core.Metadata{
		Name:      name,
		Namespace: namespace,
	})
	return prometheusconfig
}

// require custom resource to implement Clone() as well as resources.Resource interface

type CloneablePrometheusConfig interface {
	resources.Resource
	Clone() *github_com_solo_io_supergloo_api_external_prometheus.PrometheusConfig
}

var _ CloneablePrometheusConfig = &github_com_solo_io_supergloo_api_external_prometheus.PrometheusConfig{}

type PrometheusConfig struct {
	github_com_solo_io_supergloo_api_external_prometheus.PrometheusConfig
}

func (r *PrometheusConfig) Clone() resources.Resource {
	return &PrometheusConfig{PrometheusConfig: *r.PrometheusConfig.Clone()}
}

func (r *PrometheusConfig) Hash() uint64 {
	clone := r.PrometheusConfig.Clone()

	resources.UpdateMetadata(clone, func(meta *core.Metadata) {
		meta.ResourceVersion = ""
	})

	return hashutils.HashAll(clone)
}

type PrometheusConfigList []*PrometheusConfig

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list PrometheusConfigList) Find(namespace, name string) (*PrometheusConfig, error) {
	for _, prometheusConfig := range list {
		if prometheusConfig.GetMetadata().Name == name {
			if namespace == "" || prometheusConfig.GetMetadata().Namespace == namespace {
				return prometheusConfig, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find prometheusConfig %v.%v", namespace, name)
}

func (list PrometheusConfigList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, prometheusConfig := range list {
		ress = append(ress, prometheusConfig)
	}
	return ress
}

func (list PrometheusConfigList) Names() []string {
	var names []string
	for _, prometheusConfig := range list {
		names = append(names, prometheusConfig.GetMetadata().Name)
	}
	return names
}

func (list PrometheusConfigList) NamespacesDotNames() []string {
	var names []string
	for _, prometheusConfig := range list {
		names = append(names, prometheusConfig.GetMetadata().Namespace+"."+prometheusConfig.GetMetadata().Name)
	}
	return names
}

func (list PrometheusConfigList) Sort() PrometheusConfigList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].GetMetadata().Less(list[j].GetMetadata())
	})
	return list
}

func (list PrometheusConfigList) Clone() PrometheusConfigList {
	var prometheusConfigList PrometheusConfigList
	for _, prometheusConfig := range list {
		prometheusConfigList = append(prometheusConfigList, resources.Clone(prometheusConfig).(*PrometheusConfig))
	}
	return prometheusConfigList
}

func (list PrometheusConfigList) Each(f func(element *PrometheusConfig)) {
	for _, prometheusConfig := range list {
		f(prometheusConfig)
	}
}

func (list PrometheusConfigList) EachResource(f func(element resources.Resource)) {
	for _, prometheusConfig := range list {
		f(prometheusConfig)
	}
}

func (list PrometheusConfigList) AsInterfaces() []interface{} {
	var asInterfaces []interface{}
	list.Each(func(element *PrometheusConfig) {
		asInterfaces = append(asInterfaces, element)
	})
	return asInterfaces
}
