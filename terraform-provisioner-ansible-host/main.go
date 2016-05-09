package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func ResourceProvisionerBuilder() terraform.ResourceProvisioner {
	return &ResourceProvisioner{}
}

func main() {
	serveOpts := &plugin.ServeOpts{
		ProvisionerFunc: ResourceProvisionerBuilder,
	}

	plugin.Serve(serveOpts)
}

type ResourceProvisioner struct{}

func (p *ResourceProvisioner) Apply(o terraform.UIOutput,
	s *terraform.InstanceState,
	c *terraform.ResourceConfig) error {

	hostFile := ""
	inventoryPath, inventoryPathOk := c.Get("inventory_path")
	if !inventoryPathOk {
		return errors.New("inventory path is required")
	}
	host, hostOk := c.Get("host")
	if !hostOk {
		return errors.New("host is required")
	}
	if group, groupOk := c.Get("group"); groupOk {
		hostFile += fmt.Sprintf("[%s]\n", group.(string))
	}
	options := []string{host.(string)}
	for k, v := range c.Config {
		if k != "host" && k != "group" && k != "inventory_path" {
			options = append(options, fmt.Sprintf("%s=%s", k, v.(string)))
		}
	}
	hostFile += strings.Join(options, " ")
	cwd, _ := os.Getwd()

	path := filepath.Join(cwd, inventoryPath.(string))
	ioutil.WriteFile(path, []byte(hostFile), 0666)
	return nil
}

func (p *ResourceProvisioner) Validate(c *terraform.ResourceConfig) (ws []string, es []error) {
	if _, ok := c.Get("inventory_path"); !ok {
		es = append(es, errors.New("inventory path is required"))
	}
	if _, ok := c.Get("host"); !ok {
		es = append(es, errors.New("host is required"))
	}
	return ws, es
}
