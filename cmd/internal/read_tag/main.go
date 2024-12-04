package main

import (
	"flag"
	"os"

	"github.com/sagernet/sing-box/cmd/internal/build_shared"
	"github.com/sagernet/sing-box/log"
	F "github.com/sagernet/sing/common/format"
)

var nightly bool

func init() {
	flag.BoolVar(&nightly, "nightly", false, "Print nightly tag")
}

func main() {
	flag.Parse()
	if nightly {
		version, err := build_shared.ReadTagVersionRev()
		if err != nil {
			log.Fatal(err)
		}
		var (
			versionStr   string
			isPrerelease bool
		)
		if version.PreReleaseIdentifier != "" {
			isPrerelease = true
			versionStr = version.VersionString() + "-nightly"
		} else {
			version.Patch++
			versionStr = version.VersionString() + "-nightly"
		}
		if build_shared.IsDevBranch() {
			isPrerelease = true
		}
		err = setGitHubOutput("version", versionStr)
		if err != nil {
			log.Fatal(err)
		}
		err = setGitHubOutput("prerelease", F.ToString(isPrerelease))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		tag, err := build_shared.ReadTag()
		if err != nil {
			log.Error(err)
			os.Stdout.WriteString("unknown\n")
		} else {
			os.Stdout.WriteString(tag + "\n")
		}
	}
}

func setGitHubOutput(name string, value string) error {
	outputFile, err := os.OpenFile(os.Getenv("GITHUB_ENV"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	_, err = outputFile.WriteString(name + "=" + value + "\n")
	if err != nil {
		outputFile.Close()
		return err
	}
	err = outputFile.Close()
	if err != nil {
		return err
	}
	os.Stderr.WriteString(name + "=" + value + "\n")
	return nil
}
