package ubuntu

import (
	"fmt"
	"net"
	"strings"
)

// Upgrade the package cache and update the installed packages (using apt).
func UpdatePackages() *ShellCommand {
	return And("apt-get update", "DEBIAN_FRONTEND=noninteractive apt-get upgrade -y")
}

// Upgrade the package cache and update the installed packages (using apt).
func UpdatePackagesOmitError() *ShellCommand {
	return Or("apt-get update", "DEBIAN_FRONTEND=noninteractive apt-get upgrade -y")
}

// Update the package cache for a given repository only. Repo selection is done
// via the name of apt's configuration file taken from /etc/apt/sources.list.d.
// This is much faster if you just added a repo and want to install software as
// you need not update all other packages too (which most probably happened
// just recently during provisioning).
func UpdateSelectedRepoPackages(repoConfigPath string) *ShellCommand {
	return &ShellCommand{
		Command: fmt.Sprintf(
			`apt-get update -o Dir::Etc::sourcelist="sources.list.d/%s" -o Dir::Etc::sourceparts="-" -o APT::Get::List-Cleanup="0"`,
			repoConfigPath,
		),
	}
}

// Install the given packages using apt-get. At least one package must be given (pkgs can be left empty).
func InstallPackages(pkg string, pkgs ...string) *ShellCommand {
	return &ShellCommand{
		Command: fmt.Sprintf("DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends --force-yes %s %s", pkg, strings.Join(pkgs, " ")),
	}
}

func InstallPackagesWithoutForce(pkg string, pkgs ...string) *ShellCommand {
	return &ShellCommand{
		Command: fmt.Sprintf("DEBIAN_FRONTEND=noninteractive apt-get install %s %s", pkg, strings.Join(pkgs, " ")),
	}
}

// Remove the given packages using apt-get. At least one package must be given (pkgs can be left empty).
func RemovePackages(pkg string, pkgs ...string) *ShellCommand {
	return &ShellCommand{
		Command: fmt.Sprintf("DEBIAN_FRONTEND=noninteractive apt-get autoremove -y --no-install-recommends "),
	}
}

func RemovePackage(pkg string, pkgs ...string) *ShellCommand {
	return &ShellCommand{
		Command: fmt.Sprintf("DEBIAN_FRONTEND=noninteractive apt-get remove -y --no-install-recommends %s %s", pkg, strings.Join(pkgs, " ")),
	}
}

func PurgePackages(pkg string, pkgs ...string) *ShellCommand {
	return &ShellCommand{
		Command: fmt.Sprintf("DEBIAN_FRONTEND=noninteractive apt-get purge -y --no-install-recommends %s %s", pkg, strings.Join(pkgs, " ")),
	}
}

// PinPackage pins package via dpkg --set-selections
func PinPackage(name string) *ShellCommand {
	return Shell(fmt.Sprintf(`echo "%s hold" | dpkg --set-selections`, name))
}

// StartOrRestart starts or restarts a service configured with upstart
func StartOrRestart(service string) *ShellCommand {
	return Shell(fmt.Sprintf("if status %s | grep running; then { stop %s && start %s; }; else start %s; fi", service, service, service, service))
}

// EnsureRunning will start the service if not yet running. This should be used whenever a restart
// might break stuff (think ElasticSearch cluster instances in an ES update).
func EnsureRunning(service string) *ShellCommand {
	return Shell(fmt.Sprintf("status %s | grep running || start %s", service, service))
}

// IPString returns the non loopback local IP of the host
func IPNet(Netif string) *net.IPNet {
	var ipnet_ptr *net.IPNet
	//addrs, err := net.InterfaceAddrs()
	interfaces, err :=  net.Interfaces()
	if err != nil {
		return nil
	}

  for _,inter := range interfaces {
		if addrs,err := inter.Addrs(); err == nil {
			for _,addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if Netif != "" {
						if ipnet.IP.To4() != nil && inter.Name == Netif {
					    ipnet_ptr = ipnet
				    }
					} else {
						if ipnet.IP.To4() != nil {
					    ipnet_ptr = ipnet
			  	  }
					}
			 }
		 }
	 }
 }
	return ipnet_ptr
}

// IPString returns the non loopback local IP of the host
func IP(netif string) string {
	ipnet := IPNet(netif)
	if ipnet != nil {
		return ipnet.IP.String()
	}
	return ""
}

func ArraytoString(prefix,suffix string,value []string) string {
	str := ""
	for _,i := range value {
     str = str + " " + prefix + i + suffix
	}
	return str
}
