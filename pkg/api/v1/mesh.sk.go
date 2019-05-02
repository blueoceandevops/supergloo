// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"sort"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/hashutils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewMesh(namespace, name string) *Mesh {
	mesh := &Mesh{}
	mesh.SetMetadata(core.Metadata{
		Name:      name,
		Namespace: namespace,
	})
	return mesh
}

func (r *Mesh) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

func (r *Mesh) SetStatus(status core.Status) {
	r.Status = status
}

func (r *Mesh) Hash() uint64 {
	metaCopy := r.GetMetadata()
	metaCopy.ResourceVersion = ""
	return hashutils.HashAll(
		metaCopy,
		r.DiscoveryMetadata,
		r.MeshType,
	)
}

type MeshList []*Mesh
type MeshesByNamespace map[string]MeshList

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list MeshList) Find(namespace, name string) (*Mesh, error) {
	for _, mesh := range list {
		if mesh.GetMetadata().Name == name {
			if namespace == "" || mesh.GetMetadata().Namespace == namespace {
				return mesh, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find mesh %v.%v", namespace, name)
}

func (list MeshList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, mesh := range list {
		ress = append(ress, mesh)
	}
	return ress
}

func (list MeshList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, mesh := range list {
		ress = append(ress, mesh)
	}
	return ress
}

func (list MeshList) Names() []string {
	var names []string
	for _, mesh := range list {
		names = append(names, mesh.GetMetadata().Name)
	}
	return names
}

func (list MeshList) NamespacesDotNames() []string {
	var names []string
	for _, mesh := range list {
		names = append(names, mesh.GetMetadata().Namespace+"."+mesh.GetMetadata().Name)
	}
	return names
}

func (list MeshList) Sort() MeshList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].GetMetadata().Less(list[j].GetMetadata())
	})
	return list
}

func (list MeshList) Clone() MeshList {
	var meshList MeshList
	for _, mesh := range list {
		meshList = append(meshList, resources.Clone(mesh).(*Mesh))
	}
	return meshList
}

func (list MeshList) Each(f func(element *Mesh)) {
	for _, mesh := range list {
		f(mesh)
	}
}

func (list MeshList) AsInterfaces() []interface{} {
	var asInterfaces []interface{}
	list.Each(func(element *Mesh) {
		asInterfaces = append(asInterfaces, element)
	})
	return asInterfaces
}

func (byNamespace MeshesByNamespace) Add(mesh ...*Mesh) {
	for _, item := range mesh {
		byNamespace[item.GetMetadata().Namespace] = append(byNamespace[item.GetMetadata().Namespace], item)
	}
}

func (byNamespace MeshesByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace MeshesByNamespace) List() MeshList {
	var list MeshList
	for _, meshList := range byNamespace {
		list = append(list, meshList...)
	}
	return list.Sort()
}

func (byNamespace MeshesByNamespace) Clone() MeshesByNamespace {
	cloned := make(MeshesByNamespace)
	for ns, list := range byNamespace {
		cloned[ns] = list.Clone()
	}
	return cloned
}

var _ resources.Resource = &Mesh{}

// Kubernetes Adapter for Mesh

func (o *Mesh) GetObjectKind() schema.ObjectKind {
	t := MeshCrd.TypeMeta()
	return &t
}

func (o *Mesh) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*Mesh)
}

var MeshCrd = crd.NewCrd("supergloo.solo.io",
	"meshes",
	"supergloo.solo.io",
	"v1",
	"Mesh",
	"m",
	false,
	&Mesh{})
