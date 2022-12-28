package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
)

var ErrStrategyBreak = errors.New("replica configuration complete")

// Configure applies one or more replica strategies to determine the replica
// configuration, ensuring that the replica is unique in the system. Each strategy is
// applied in the order specified and will continue until a ErrStrategyBreak is returned
// or all strategies have been applied. A new configuration is returned and the original
// configuration is unmodified.
func (c ReplicaConfig) Configure(strategies ...ReplicaStrategy) (_ ReplicaConfig, err error) {
	for _, strategy := range strategies {
		if c, err = strategy(c, err); err == ErrStrategyBreak {
			return c, nil
		}
	}
	return c, err
}

// A ReplicaStrategy defines alternate methods for populating a ReplicaConfig besides
// simply loading them from the environment (the default configuration method). The
// strategy will receive an input config value (not a pointer) and it must create a new
// configuration with the modifications so that the original config is not changed
// during processing. The error result returned by the function defines how the strategy
// continues if multiple strategies are applied. If the function returns a special
// StrategyBreak error, the replica configuration will stop processing without an error
// returned; all other non-nil errors are passed to the next strategy, so that they can
// continue to pass the error
type ReplicaStrategy func(in ReplicaConfig, err error) (ReplicaConfig, error)

// MultiError collects errors from the strategy into a single error for reporting.
func MultiError(f ReplicaStrategy) ReplicaStrategy {
	return func(in ReplicaConfig, err error) (out ReplicaConfig, nerr error) {
		if out, nerr = f(in, err); nerr != nil && nerr != err {
			return out, multierror.Append(err, nerr)
		}
		return out, err
	}
}

var hostRE = regexp.MustCompile(`^(?P<Name>[\w\d\_\-\.]+)-(?P<Pod>\d+)$`)

// HostnamePID attempts to process a pod-n hostname from a kubernetes pod created by a
// stateful set to determine the PID. If the hostname passed in is empty, then
// os.Hostname is used to determine the name. The PID (e.g. n from pod-n) is then added
// to the configured PID so that it can be used to specify a range of PIDs. E.g. if the
// environment PID is 40 and the hostname is trtl-2 then the PID will be 42.
// StatefulSets (or vms, or hardware) can then be configured with ranges of PID values.
// This strategy returns a MultiError for downstream processing if it cannot fetch the
// hostname or parse it; it does not return a StrategyBreak.
func HostnamePID(hostname string) ReplicaStrategy {
	return MultiError(func(in ReplicaConfig, err error) (_ ReplicaConfig, nerr error) {
		if hostname == "" {
			if hostname, nerr = os.Hostname(); nerr != nil {
				return in, fmt.Errorf("could not fetch hostname from kernel: %s", nerr)
			}
		}

		groups := hostRE.FindStringSubmatch(hostname)
		names := hostRE.SubexpNames()
		if len(names) != len(groups) {
			return in, fmt.Errorf("could not parse %q - does not match host-pid expression", hostname)
		}

		parts := make(map[string]string)
		for i, name := range names {
			parts[name] = groups[i]
		}

		// Parse the POD
		var pod uint64
		if pod, nerr = strconv.ParseUint(parts["Pod"], 10, 64); nerr != nil {
			return in, fmt.Errorf("could not parse %q - not a uint64", parts["Pod"])
		}

		// Update the PID and the hostname from the configuration
		in.PID = in.PID + pod
		in.Name = fmt.Sprintf("%s-%d", parts["Name"], in.PID)
		return in, nil
	})
}

// FilePID reads the PID from a file and is intended as a fast way to set the PID on a
// POD or in a development environment, very similar to how a /var/run/proc.pid works.
// This strategy returns a MultiError for downstream processing if it cannot fetch the
// hostname or parse it; it does not return a StrategyBreak.
func FilePID(path string) ReplicaStrategy {
	return MultiError(func(in ReplicaConfig, err error) (_ ReplicaConfig, nerr error) {
		// Read the file
		var data []byte
		if data, nerr = os.ReadFile(path); nerr != nil {
			return in, fmt.Errorf("could not read file pid path %q: %s", path, nerr)
		}

		// Parse the PID from the file
		var pid uint64
		if pid, nerr = strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64); nerr != nil {
			return in, fmt.Errorf("could not parse data in file - not a uint64")
		}

		// Update the PID and the name from the configuration
		in.PID = in.PID + pid
		in.Name = fmt.Sprintf("%s-%d", in.Name, in.PID)
		return in, nil
	})
}

// JSONConfig attempts to configure the Replica by loading a JSON file that maps PID to
// configuration data. If the previous error is not nil, then no processing will happen
// and the error will be returned. If the configuration is successful, then this will
// return StrategyBreak. The idea here is that a universal config map can be loaded into
// a kubernetes cluster and a HostnamePID or FilePID applied before JSONConfig (if they
// fail, we shouldn't try to lookup the default PID). If successfully configured, the
// processing can stop, otherwise we can continue to other configuration strategies.
func JSONConfig(path string) ReplicaStrategy {
	return func(in ReplicaConfig, err error) (ReplicaConfig, error) {
		// Do not process if there was a previous error
		if err != nil {
			return in, err
		}

		// Check to see if the file exists, if not; do nothing
		if _, err = os.Stat(path); os.IsNotExist(err) {
			return in, nil
		}

		// Attempt to read the file
		var data []byte
		if data, err = os.ReadFile(path); err != nil {
			return in, fmt.Errorf("could not read json config path %q: %s", path, err)
		}

		// Load the JSON
		cfgmap := make(map[uint64]ReplicaConfig)
		if err = json.Unmarshal(data, &cfgmap); err != nil {
			return in, fmt.Errorf("could not unmarshal JSON data: %s", err)
		}

		// If the PID of the replica is not in the map, continue
		conf, ok := cfgmap[in.PID]
		if !ok {
			return in, nil
		}

		// Otherwise, return the new configuration and StrategyBreak
		// NOTE: this requires a complete configuration, no default values!
		return conf, ErrStrategyBreak
	}
}
