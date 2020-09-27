package seaweedfs

import (
	"fmt"
	"os"

	"github.com/chrislusf/seaweedfs/weed/command"
)

//StartSeaweedServer starts a seaweed server with one master, one volume and one filer
func StartSeaweedServer(dir string) {
	commands := command.Commands
	for _, cmd := range commands {
		if cmd.Name() == "server" && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			cmd.Flag.Parse([]string{"-dir=" + dir, "-filer=true", "-volume.max=50"})
			args := cmd.Flag.Args()
			if !cmd.Run(cmd, args) {
				fmt.Fprintf(os.Stderr, "\n")
				cmd.Flag.Usage()
				fmt.Fprintf(os.Stderr, "Default Parameters:\n")
				cmd.Flag.PrintDefaults()
			}
		}
	}
	return
}

//StartSeaweedVolume starts an extra volume server for seaweedfs
func StartSeaweedVolume(dir string, port string) {
	commands := command.Commands
	for _, cmd := range commands {
		if cmd.Name() == "volume" && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			cmd.Flag.Parse([]string{"-max=100", "-mserver=localhost:9333", "-dir=" + dir, "-port=" + port})
			args := cmd.Flag.Args()
			if !cmd.Run(cmd, args) {
				fmt.Fprintf(os.Stderr, "\n")
				cmd.Flag.Usage()
				fmt.Fprintf(os.Stderr, "Default Parameters:\n")
				cmd.Flag.PrintDefaults()
			}
		}
	}
	return
}

//ShowSeaweedVersion shows seaweed version
func ShowSeaweedVersion() {
	commands := command.Commands
	for _, cmd := range commands {
		if cmd.Name() == "version" && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			cmd.Flag.Parse([]string{})
			args := cmd.Flag.Args()
			if !cmd.Run(cmd, args) {
				fmt.Fprintf(os.Stderr, "\n")
				cmd.Flag.Usage()
				fmt.Fprintf(os.Stderr, "Default Parameters:\n")
				cmd.Flag.PrintDefaults()
			}
		}
	}
	return
}

//MountDirectory mounts the given directory onto the SeaweedFS server
func MountDirectory(dir string, dirname string) {
	commands := command.Commands
	for _, cmd := range commands {
		if cmd.Name() == "mount" && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			cmd.Flag.Parse([]string{"-filer=localhost:8888", "-dir=" + dir, "-filer.path=/" + dirname})
			args := cmd.Flag.Args()
			if !cmd.Run(cmd, args) {
				fmt.Fprintf(os.Stderr, "\n")
				cmd.Flag.Usage()
				fmt.Fprintf(os.Stderr, "Default Parameters:\n")
				cmd.Flag.PrintDefaults()
			}
		}
	}
	return
}
