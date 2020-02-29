package campaign

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"git.ctrlz.es/mgdelacroix/campaigner/model"
)

func Save(campaign *model.Campaign) error {
	marshaledCampaign, err := json.MarshalIndent(campaign, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("./campaign.json", marshaledCampaign, 0600); err != nil {
		return fmt.Errorf("cannot save campaign: %w", err)
	}
	return nil
}

func Read() (*model.Campaign, error) {
	if _, err := os.Stat("."); err != nil {
		return nil, fmt.Errorf("cannot read campaign: %w", err)
	}

	fileContents, err := ioutil.ReadFile("./campaign.json")
	if err != nil {
		return nil, fmt.Errorf("there was a problem reading the campaign file: %w", err)
	}

	var campaign model.Campaign
	if err := json.Unmarshal(fileContents, &campaign); err != nil {
		return nil, fmt.Errorf("there was a problem parsing the campaign file: %w", err)
	}

	return &campaign, nil
}
