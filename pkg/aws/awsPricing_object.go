package aws

import "time"

type TEST struct {
	FormatVersion   string             `json:"formatVersion,omitempty"`
	Disclaimer      string             `json:"disclaimer,omitempty"`
	OfferCode       string             `json:"offerCode,omitempty"`
	Version         string             `json:"version,omitempty"`
	PublicationDate time.Time          `json:"publicationDate,omitempty"`
	Products        map[string]Product `json:"products,omitempty"`
}

type PriceList struct {
	Product         Product   `json:"product,omitempty"`
	PublicationDate time.Time `json:"publicationDate,omitempty"`
	ServiceCode     string    `json:"serviceCode,omitempty"`
	Terms           struct {
		OnDemand map[string]TermsAttributes `json:"OnDemand,omitempty"`
	} `json:"terms,omitempty"`
	Version string `json:"version,omitempty"`
}

type Product struct {
	Attributes struct {
		Servicecode  string `json:"servicecode,omitempty"`
		Location     string `json:"location,omitempty"`
		LocationType string `json:"locationType,omitempty"`
		Availability string `json:"availability,omitempty"`
		StorageClass string `json:"storageClass,omitempty"`
		VolumeType   string `json:"volumeType,omitempty"`
		Usagetype    string `json:"usagetype,omitempty"`
		Operation    string `json:"operation,omitempty"`
		Durability   string `json:"durability,omitempty"`
		RegionCode   string `json:"regionCode,omitempty"`
		Servicename  string `json:"servicename,omitempty"`
	} `json:"attributes,omitempty"`
	ProductFamily string `json:"productFamily,omitempty"`
	Sku           string `json:"sku,omitempty"`
}

type TermsAttributes struct {
	EffectiveDate   time.Time                 `json:"effectiveDate,omitempty"`
	OfferTermCode   string                    `json:"offerTermCode,omitempty"`
	PriceDimensions map[string]PriceDimension `json:"priceDimensions,omitempty"`
	Sku             string                    `json:"sku,omitempty"`
	TermAttributes  struct {
	} `json:"termAttributes,omitempty"`
}

type PriceDimension struct {
	AppliesTo    []any  `json:"appliesTo,omitempty"`
	BeginRange   string `json:"beginRange,omitempty"`
	Description  string `json:"description,omitempty"`
	EndRange     string `json:"endRange,omitempty"`
	PricePerUnit struct {
		Usd string `json:"USD,omitempty"`
	} `json:"pricePerUnit,omitempty"`
	RateCode string `json:"rateCode,omitempty"`
	Unit     string `json:"unit,omitempty"`
}
