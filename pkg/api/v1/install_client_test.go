// Code generated by solo-kit. DO NOT EDIT.

// +build solokit

package v1

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/tests/typed"
)

var _ = Describe("InstallClient", func() {
	var (
		namespace string
	)
	for _, test := range []typed.ResourceClientTester{
		&typed.KubeRcTester{Crd: InstallCrd},
		&typed.ConsulRcTester{},
		&typed.FileRcTester{},
		&typed.MemoryRcTester{},
		&typed.VaultRcTester{},
		&typed.KubeSecretRcTester{},
		&typed.KubeConfigMapRcTester{},
	} {
		Context("resource client backed by "+test.Description(), func() {
			var (
				client              InstallClient
				err                 error
				name1, name2, name3 = "foo" + helpers.RandString(3), "boo" + helpers.RandString(3), "goo" + helpers.RandString(3)
			)

			BeforeEach(func() {
				namespace = helpers.RandString(6)
				factory := test.Setup(namespace)
				client, err = NewInstallClient(factory)
				Expect(err).NotTo(HaveOccurred())
			})
			AfterEach(func() {
				test.Teardown(namespace)
			})
			It("CRUDs Installs "+test.Description(), func() {
				InstallClientTest(namespace, client, name1, name2, name3)
			})
		})
	}
})

func InstallClientTest(namespace string, client InstallClient, name1, name2, name3 string) {
	err := client.Register()
	Expect(err).NotTo(HaveOccurred())

	name := name1
	input := NewInstall(namespace, name)

	r1, err := client.Write(input, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	_, err = client.Write(input, clients.WriteOpts{})
	Expect(err).To(HaveOccurred())
	Expect(errors.IsExist(err)).To(BeTrue())

	Expect(r1).To(BeAssignableToTypeOf(&Install{}))
	Expect(r1.GetMetadata().Name).To(Equal(name))
	Expect(r1.GetMetadata().Namespace).To(Equal(namespace))
	Expect(r1.GetMetadata().ResourceVersion).NotTo(Equal(input.GetMetadata().ResourceVersion))
	Expect(r1.GetMetadata().Ref()).To(Equal(input.GetMetadata().Ref()))
	Expect(r1.Status).To(Equal(input.Status))
	Expect(r1.Disabled).To(Equal(input.Disabled))
	Expect(r1.InstallationNamespace).To(Equal(input.InstallationNamespace))

	_, err = client.Write(input, clients.WriteOpts{
		OverwriteExisting: true,
	})
	Expect(err).To(HaveOccurred())

	resources.UpdateMetadata(input, func(meta *core.Metadata) {
		meta.ResourceVersion = r1.GetMetadata().ResourceVersion
	})
	r1, err = client.Write(input, clients.WriteOpts{
		OverwriteExisting: true,
	})
	Expect(err).NotTo(HaveOccurred())
	read, err := client.Read(namespace, name, clients.ReadOpts{})
	Expect(err).NotTo(HaveOccurred())
	Expect(read).To(Equal(r1))
	_, err = client.Read("doesntexist", name, clients.ReadOpts{})
	Expect(err).To(HaveOccurred())
	Expect(errors.IsNotExist(err)).To(BeTrue())

	name = name2
	input = &Install{}

	input.SetMetadata(core.Metadata{
		Name:      name,
		Namespace: namespace,
	})

	r2, err := client.Write(input, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())
	list, err := client.List(namespace, clients.ListOpts{})
	Expect(err).NotTo(HaveOccurred())
	Expect(list).To(ContainElement(r1))
	Expect(list).To(ContainElement(r2))
	err = client.Delete(namespace, "adsfw", clients.DeleteOpts{})
	Expect(err).To(HaveOccurred())
	Expect(errors.IsNotExist(err)).To(BeTrue())
	err = client.Delete(namespace, "adsfw", clients.DeleteOpts{
		IgnoreNotExist: true,
	})
	Expect(err).NotTo(HaveOccurred())
	err = client.Delete(namespace, r2.GetMetadata().Name, clients.DeleteOpts{})
	Expect(err).NotTo(HaveOccurred())

	Eventually(func() InstallList {
		list, err = client.List(namespace, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		return list
	}, time.Second*10).Should(ContainElement(r1))
	Eventually(func() InstallList {
		list, err = client.List(namespace, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		return list
	}, time.Second*10).ShouldNot(ContainElement(r2))
	w, errs, err := client.Watch(namespace, clients.WatchOpts{
		RefreshRate: time.Hour,
	})
	Expect(err).NotTo(HaveOccurred())

	var r3 resources.Resource
	wait := make(chan struct{})
	go func() {
		defer close(wait)
		defer GinkgoRecover()

		resources.UpdateMetadata(r2, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		r2, err = client.Write(r2, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		name = name3
		input = &Install{}
		Expect(err).NotTo(HaveOccurred())
		input.SetMetadata(core.Metadata{
			Name:      name,
			Namespace: namespace,
		})

		r3, err = client.Write(input, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
	}()
	<-wait

	select {
	case err := <-errs:
		Expect(err).NotTo(HaveOccurred())
	case list = <-w:
	case <-time.After(time.Millisecond * 5):
		Fail("expected a message in channel")
	}

	go func() {
		defer GinkgoRecover()
		for {
			select {
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Second / 4):
				return
			}
		}
	}()

	Eventually(w, time.Second*5, time.Second/10).Should(Receive(And(ContainElement(r1), ContainElement(r3), ContainElement(r3))))
}
