// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package kubectl

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	kubectldrain "k8s.io/kubernetes/pkg/kubectl/cmd/drain"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var log logr.Logger = ctrl.Log.WithName("kubectl")

func getOptions(vmName string) (*kubectldrain.DrainCmdOptions, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to set up client config")
	}

	streams := genericclioptions.IOStreams{
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	f := cmdutil.NewFactory(&RESTConfigClientGetter{Config: cfg})
	drain := kubectldrain.NewCmdDrain(f, streams)
	options := kubectldrain.NewDrainCmdOptions(f, streams)
	err = drain.ParseFlags([]string{"--ignore-daemonsets", "--force", "--delete-local-data", "--grace-period=60"})
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing flags")
	}
	err = options.Complete(f, drain, []string{vmName})
	if err != nil {
		return nil, errors.Wrapf(err, "error setting up drain")
	}

	return options, nil
}

// Cordon cordons  vmname from kubernetes
func Cordon(vmName string) error {
	options, err := getOptions(vmName)
	if err != nil {
		return errors.Wrapf(err, "error getting options")
	}

	log.Info("Cordon", "VMName", vmName)
	err = options.RunCordonOrUncordon(true)
	if err != nil {
		return errors.Wrapf(err, "error cordoning node")
	}

	return nil
}

// Drain drains vmname from kubernetes
func Drain(vmName string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return errors.Wrapf(err, "unable to set up client config")
	}

	streams := genericclioptions.IOStreams{
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	f := cmdutil.NewFactory(&RESTConfigClientGetter{Config: cfg})
	drain := kubectldrain.NewCmdDrain(f, streams)
	drain.SetArgs([]string{vmName, "--ignore-daemonsets", "--force", "--delete-local-data", "--grace-period=60"})

	log.Info("Draining", "VMName", vmName)

	err = drain.Execute()
	if err != nil {
		return errors.Wrapf(err, "error draining node")
	}

	return nil
}

// Uncordon uncordons vmname from kubernetes
func Uncordon(vmName string) error {
	options, err := getOptions(vmName)
	if err != nil {
		return errors.Wrapf(err, "error getting options")
	}

	log.Info("Uncordon", "VMName", vmName)
	err = options.RunCordonOrUncordon(false)
	if err != nil {
		return errors.Wrapf(err, "error cordoning node")
	}

	return nil
}