package pkg

import "time"

type App struct {
	ClientId                 string       `json:"clientId"`
	LastUpdated              *time.Time   `json:"lastUpdated"`
	Name                     string       `json:"name"`
	Url                      string       `json:"url"`
	IconUrl                  string       `json:"iconUrl"`
	PrimaryCategory          string       `json:"primaryCategory"`
	Languages                []string     `json:"languages"`
	Introduction             string       `json:"introduction"`
	Details                  string       `json:"details"`
	Features                 string       `json:"features"`
	DemoStoreUrl             string       `json:"demoStoreUrl"`
	FeatureMediaUrl          string       `json:"featureMediaUrl"`
	Screenshots              []Screenshot `json:"screenshots"`
	Integrations             []string     `json:"integrations"`
	SupportEmail             string       `json:"supportEmail"`
	SupportPhone             string       `json:"supportPhone"`
	PrivacyPolicyUrl         string       `json:"privacyPolicyUrl"`
	DeveloperWebsiteUrl      string       `json:"developerWebsite"`
	FaqUrl                   string       `json:"faqUrl"`
	ChangelogUrl             string       `json:"changelogUrl"`
	SupportPortalUrl         string       `json:"supportPortalUrl"`
	TutorialUrl              string       `json:"tutorialUrl"`
	AdditionalAppDocumentUrl string       `json:"additionalAppDocumentUrl"`
	Pricing                  []Pricing    `json:"pricing"`
	Subtitle                 string       `json:"subtitle"`
}

type PricingType string

const (
	Free      PricingType = "Free"
	Recurring PricingType = "Recurring"
	OneTime   PricingType = "OneTime"
)

type Pricing struct {
	Name              string      `json:"name"`
	Features          []string    `json:"features"`
	AdditionalCharges string      `json:"additionalCharges"`
	Type              PricingType `json:"type"`
}

type Screenshot struct {
	ImageUrl    string `json:"imageUrl"`
	Description string `json:"description"`
}
