package pkg

import "time"

type App struct {
	ClientId                 string       `json:"clientId,omitempty" bson:"_id,omitempty"`
	LastUpdated              *time.Time   `json:"lastUpdated,omitempty" bson:"lastUpdated,omitempty"`
	LastCrawl                *time.Time   `json:"lastCrawl,omitempty" bson:"lastCrawl,omitempty"`
	Name                     string       `json:"name,omitempty" bson:"name,omitempty"`
	Url                      string       `json:"url,omitempty" bson:"url,omitempty"`
	IconUrl                  string       `json:"iconUrl,omitempty" bson:"iconUrl,omitempty"`
	PrimaryCategory          string       `json:"primaryCategory,omitempty" bson:"primaryCategory,omitempty"`
	Languages                []string     `json:"languages,omitempty" bson:"languages,omitempty"`
	Introduction             string       `json:"introduction,omitempty" bson:"introduction,omitempty"`
	Details                  string       `json:"details,omitempty" bson:"details,omitempty"`
	Features                 string       `json:"features,omitempty" bson:"features,omitempty"`
	DemoStoreUrl             string       `json:"demoStoreUrl,omitempty" bson:"demoStoreUrl,omitempty"`
	FeatureMediaUrl          string       `json:"featureMediaUrl,omitempty" bson:"featureMediaUrl,omitempty"`
	Screenshots              []Screenshot `json:"screenshots,omitempty" bson:"screenshots,omitempty"`
	Integrations             []string     `json:"integrations,omitempty" bson:"integrations,omitempty"`
	SupportEmail             string       `json:"supportEmail,omitempty" bson:"supportEmail,omitempty"`
	SupportPhone             string       `json:"supportPhone,omitempty" bson:"supportPhone,omitempty"`
	PrivacyPolicyUrl         string       `json:"privacyPolicyUrl,omitempty" bson:"privacyPolicyUrl,omitempty"`
	DeveloperWebsiteUrl      string       `json:"developerWebsite,omitempty" bson:"developerWebsite,omitempty"`
	FaqUrl                   string       `json:"faqUrl,omitempty" bson:"faqUrl,omitempty"`
	ChangelogUrl             string       `json:"changelogUrl,omitempty" bson:"changelogUrl,omitempty"`
	SupportPortalUrl         string       `json:"supportPortalUrl,omitempty" bson:"supportPortalUrl,omitempty"`
	TutorialUrl              string       `json:"tutorialUrl,omitempty" bson:"tutorialUrl,omitempty"`
	AdditionalAppDocumentUrl string       `json:"additionalAppDocumentUrl,omitempty" bson:"additionalAppDocumentUrl,omitempty"`
	Pricing                  []Pricing    `json:"pricing,omitempty" bson:"pricing,omitempty"`
	Subtitle                 string       `json:"subtitle,omitempty" bson:"subtitle,omitempty"`
}

type SitemapEntry struct {
	Location           string
	ParsedLastModified *time.Time
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
