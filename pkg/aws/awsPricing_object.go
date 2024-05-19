package aws

import "time"

type PriceList struct {
	Product struct {
		Attributes struct {
			Availability string `json:"availability,omitempty"`
			Durability   string `json:"durability,omitempty"`
			Location     string `json:"location,omitempty"`
			LocationType string `json:"locationType,omitempty"`
			Operation    string `json:"operation,omitempty"`
			RegionCode   string `json:"regionCode,omitempty"`
			Servicecode  string `json:"servicecode,omitempty"`
			Servicename  string `json:"servicename,omitempty"`
			StorageClass string `json:"storageClass,omitempty"`
			Usagetype    string `json:"usagetype,omitempty"`
			VolumeType   string `json:"volumeType,omitempty"`
		} `json:"attributes,omitempty"`
		ProductFamily string `json:"productFamily,omitempty"`
		Sku           string `json:"sku,omitempty"`
	} `json:"product,omitempty"`
	PublicationDate time.Time `json:"publicationDate,omitempty"`
	ServiceCode     string    `json:"serviceCode,omitempty"`
	Terms           struct {
		OnDemand map[string]TermsAttributes `json:"OnDemand,omitempty"`
	} `json:"terms,omitempty"`
	Version string `json:"version,omitempty"`
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
