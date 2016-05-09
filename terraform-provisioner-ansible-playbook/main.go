package main

import (
	"fmt"
	"os"
	"os/exec"
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

	ansibleArgs := []string{"ansible-playbook"}
	playbook, _ := c.Get("playbook")
	ansibleArgs = append(ansibleArgs, playbook.(string))
	inventory, _ := c.Get("inventory")
	ansibleArgs = append(ansibleArgs, fmt.Sprintf("-i %s", inventory.(string)))
	user, _ := c.Get("user")
	ansibleArgs = append(ansibleArgs, fmt.Sprintf("-u %s", user.(string)))
	command := strings.Join(ansibleArgs, " ")

	cwd, _ := os.Getwd()

	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Dir = cwd
	cmd.Stdout = Output{o}
	cmd.Stderr = Output{o}
	return cmd.Run()

}

type Output struct {
	o terraform.UIOutput
}

func (out Output) Write(p []byte) (n int, err error) {
	out.o.Output(string(p))
	return len(p), nil
}

func (p *ResourceProvisioner) Validate(c *terraform.ResourceConfig) (ws []string, es []error) {
	return ws, es
}
