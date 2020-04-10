package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/operator-framework/operator-sdk/hack/generate/changelog/util"

	"github.com/blang/semver"
	log "github.com/sirupsen/logrus"
)

const repo = "github.com/operator-framework/operator-sdk"

func main() {
	var (
		tag           string
		fragmentsDir  string
		changelogFile string
		migrationDir  string
		validateOnly  bool
	)

	flag.StringVar(&tag, "tag", "",
		"Title for generated CHANGELOG and migration guide sections")
	flag.StringVar(&fragmentsDir, "fragments-dir", filepath.Join("changelog", "fragments"),
		"Path to changelog fragments directory")
	flag.StringVar(&changelogFile, "changelog", "CHANGELOG.md",
		"Path to CHANGELOG")
	flag.StringVar(&migrationDir, "migration-guide-dir",
		filepath.Join("website", "content", "en", "docs", "migration"),
		"Path to migration guide directory")
	flag.BoolVar(&validateOnly, "validate-only", false,
		"Only validate fragments")
	flag.Parse()

	if tag == "" && !validateOnly {
		log.Fatalf("flag '-tag' is required without '-validate-only'")
	}

	entries, err := util.LoadEntries(fragmentsDir, repo)
	if err != nil {
		log.Fatalf("failed to load fragments: %v", err)
	}
	if len(entries) == 0 {
		log.Fatalf("no entries found")
	}

	if validateOnly {
		return
	}

	version, err := semver.Parse(strings.TrimPrefix(tag, "v"))
	if err != nil {
		log.Fatalf("flag '-tag' is not a valid semantic version: %v", err)
	}
	if len(version.Pre) > 0 || len(version.Build) > 0 {
		log.Fatalf("flag '-tag' must not include a build number or pre-release identifiers")
	}

	cl := util.ChangelogFromEntries(version, entries)
	if err := cl.WriteFile(changelogFile); err != nil {
		log.Fatalf("failed to update CHANGELOG: %v", err)
	}

	mg := util.MigrationGuideFromEntries(version, entries)
	mgFile := filepath.Join(migrationDir, fmt.Sprintf("v%s.md", version))
	if err := mg.WriteFile(mgFile); err != nil {
		log.Fatalf("failed to create migration guide: %v", err)
	}
}
