package main

// SPDX-License-Identifier: GPL-3.0-or-later
import (
	"context"
	"time"

	"github.com/lopster568/phantomDNS/internal/blocklist"
	"github.com/lopster568/phantomDNS/internal/config"
	"github.com/lopster568/phantomDNS/internal/dnsengine"
	"github.com/lopster568/phantomDNS/internal/logger"
	"github.com/lopster568/phantomDNS/internal/policy"
	"github.com/lopster568/phantomDNS/internal/storage/db"
	"github.com/lopster568/phantomDNS/internal/storage/models"
	"github.com/lopster568/phantomDNS/internal/storage/repositories"
)

func main() {
	logger.Log.Info("Starting PhantomDNS Data Plane...")
	// 1. Initialize DB
	db.InitDB("/app/data/phantomdns.db")
	// 2. Initialize Repositories (store)
	repos := repositories.NewStore(db.DB)
	// 2.1 Blocklist Engine
	blEngine := blocklist.NewEngine(repos.Blocklist)
	// 4. Create a fake blocklist source (can be a small plain-text domain list)
	src := models.BlocklistSource{
		ID:        "test-source",
		Name:      "StevenBlack Blocklist",
		URL:       "https://raw.githubusercontent.com/StevenBlack/hosts/master/data/StevenBlack/hosts", // or your own small text list
		Format:    "hosts",                                                                             // must match a parser in your system
		Category:  "test",
		Enabled:   true,
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 5. Run single-source update
	if err := blEngine.UpdateSource(ctx, src, ""); err != nil {
		logger.Log.Fatalf("Smoke test failed: %v", err)
	}

	logger.Log.Info("Blocklist fetch/parse/store successful, plus stored in DB.")

	logger.Log.Info("Dumping first few blocklisted hosts for verification...")
	hosts, err := blEngine.List()
	if err != nil {
		logger.Log.Fatalf("failed to list blocklist entries: %v", err)
	}
	for i, h := range hosts {
		if i >= 10 {
			break
		}
		logger.Log.Infof("Blocklisted: %s", h)
	}

	// 3. Initialize Policy Engine
	policyEngine := policy.NewPolicyEngine()
	policies, err := policy.LoadPoliciesFromFile("/app/configs/policies.json")
	if err != nil {
		logger.Log.Fatalf("failed to load policies from file: %v", err)
	}
	if err := policyEngine.LoadPolicies(policies); err != nil {
		logger.Log.Fatalf("failed to load snapshot: %v", err)
	}
	// 4. Initialize DNS Engine with default config and repos
	engine, err := dnsengine.NewDNSEngine(config.DefaultConfig.DataPlane, repos, policyEngine)
	if err != nil {
		logger.Log.Fatal("Failed to create DNS engine: " + err.Error())
	}

	// 4.1 Attach blocklist checker to DNS engine
	engine.AttachBlocklistChecker(repos.Blocklist)
	// 5. Initialize and Run Server with the engine
	srv, err := dnsengine.NewServer(config.DefaultConfig.DataPlane, engine)
	if err != nil {
		logger.Log.Fatal("Failed to create server: " + err.Error())
	}
	srv.Run()
}
