package campaign

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

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
