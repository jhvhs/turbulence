package director

import (
	boshdir "github.com/cloudfoundry/bosh-cli/director"
)

type DirectorImpl struct {
	director boshdir.Director
}

func (d DirectorImpl) AllInstances() ([]Instance, error) {
	deps, err := d.director.Deployments()
	if err != nil {
		return nil, err
	}

	var instances []Instance

	for _, dep := range deps {
		insts, err := dep.Instances()
		if err != nil {
			return nil, err
		}

		for _, inst := range insts {
			instances = append(instances, InstanceImpl{
				id:         inst.ID,
				group:      inst.Group,
				deployment: dep,
				az:         inst.AZ,

				cid:     inst.VMID,
				agentID: inst.AgentID,
			})
		}
	}

	return instances, nil
}
