package pkg

import (
	"os"
)

type Scanner struct {
	config     Config
	gitManager GitManager
}

func NewScanner(config *Config, gitManager GitManager) *Scanner {
	return &Scanner{
		config:     *config,
		gitManager: gitManager,
	}
}

type ScannerStatus struct {
	IsCloned      bool
	Uncommitted   bool
	CurrentBranch string
}

func NewScannerStatus() *ScannerStatus {
	return &ScannerStatus{IsCloned: false, Uncommitted: false, CurrentBranch: ""}
}

func (m *Scanner) Scan(projectName string) (*ScannerStatus, error) {
	var result = NewScannerStatus()
	var err error
	result.IsCloned, err = m.IsCloned(projectName)

	if err != nil {
		return nil, err
	}

	if result.IsCloned == false {
		return result, nil
	}

	result.CurrentBranch, err = m.gitManager.GetCurrentBranch(projectName)

	if err != nil {
		return nil, err
	}

	result.Uncommitted, err = m.gitManager.HasUncommitedChanges(projectName)

	if err != nil {
		return nil, err
	}

	if result.Uncommitted == false {
		return result, nil
	}

	return result, nil
}

func (m *Scanner) IsCloned(projectName string) (bool, error) {
	cfg, err := m.config.GetProjectByName(projectName)

	if err != nil {
		return false, err
	}

	if _, err := os.Stat(cfg.Repository.Target); os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}
