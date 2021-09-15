package service

import (
	"errors"
	"hal9000/types"
)

type interfaceStoreConcrete struct {
	interfaces map[string][]types.Interface
}

func InitInterfaceStore() types.InterfaceStore {
	return &interfaceStoreConcrete{make(map[string][]types.Interface)}
}

func (is *interfaceStoreConcrete) Register(person types.Person, iface types.Interface) {
	if _, ok := is.interfaces[person.GetID()]; !ok {
		is.interfaces[person.GetID()] = make([]types.Interface, 0)
	}
	is.interfaces[person.GetID()] = append(is.interfaces[person.GetID()], iface)
}

func (is *interfaceStoreConcrete) GetInterfacesForPerson(p types.Person, id string) []types.Interface {
	interfaces := make([]types.Interface, 0)
	if ifaces, ok := is.interfaces[p.GetID()]; ok {
		removeSet := make([]int, 0)
		for i, iface := range ifaces {
			if iface.IsStillValid() {
				if id == "" || (id != "" && id == iface.ID()) {
					interfaces = append(interfaces, iface)
				}
			} else {
				removeSet = append(removeSet, i)
			}
		}
		// if len(removeSet) > 0 {
		// 	for _, i := range removeSet {
		// 		ifaces = append(ifaces[:i], ifaces[i+1:]...)
		// 	}
		// 	is.interfaces[p.GetID()] = ifaces
		// }TODO
	}
	return interfaces
}

func (is *interfaceStoreConcrete) GetVisualInterfacesForPerson(p types.Person) []types.Interface {
	interfaces := make([]types.Interface, 0)
	for _, iface := range is.GetInterfacesForPerson(p, "") {
		if iface.SupportsVisuals() {
			interfaces = append(interfaces, iface)
		}
	}
	return interfaces
}

func (is *interfaceStoreConcrete) DetermineInterfaceOwner(runtime types.Runtime, iface types.Interface) (types.Person, error) {
	for owner, ifaces := range is.interfaces {
		for _, _iface := range ifaces {
			if iface.ID() == _iface.ID() {
				return runtime.People().GetPersonByID(owner)
			}
		}
	}
	return nil, errors.New("no owner for interface")
}
