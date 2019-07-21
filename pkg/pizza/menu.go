package pizza

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Menu struct {
	Meta     MenuMeta               `json:"Misc"`
	Flavors  map[string]*ObjectMap  `json:"Flavors"`
	Products map[string]*Product    `json:"Products"`
	Sizes    map[string]*ObjectMap  `json:"Sizes"`
	Toppings map[string]*ObjectMap  `json:"Toppings"`
	Variants map[string]*Variant    `json:"Variants"`
	Coupons  map[string]*ObjectInfo `json:"Coupons"`
}

type MenuMeta struct {
	Status        int    `json:"Status"`
	StoreID       string `json:"StoreID"`
	BusinessDate  string `json:"BusinessDate"`
	StoreAsOfTime string `json:"StoreAsOfTime"`
	LanguageCode  string `json:"LanguageCode"`
	Version       string `json:"Version"`
	ExpiresOn     string `json:"ExpiresOn"`
}

type ObjectMap map[string]ObjectInfo

type ObjectInfo struct {
	Code        string                 `json:"Code,omitempty"`
	ImageCode   string                 `json:"ImageCode,omitempty"`
	Description string                 `json:"Description,omitempty"`
	Local       bool                   `json:"Local,omitempty"`
	Name        string                 `json:"Name,omitempty"`
	SortSeq     string                 `json:"SortSeq,omitempty"`
	Tags        map[string]interface{} `json:"Tags,omitempty"`
}

type Product struct {
	AvailableToppings string                 `json:"AvailableToppings,omitempty"`
	AvailableSides    string                 `json:"AvailableSides,omitempty"`
	Code              string                 `json:"Code,omitempty"`
	DefaultToppings   string                 `json:"DefaultToppings,omitempty"`
	DefaultSides      string                 `json:"DefaultSides,omitempty"`
	Description       string                 `json:"Description,omitempty"`
	ImageCode         string                 `json:"ImageCode,omitempty"`
	Local             bool                   `json:"Local,omitempty"`
	Name              string                 `json:"Name,omitempty"`
	ProductType       string                 `json:"ProductType,omitempty"`
	Tags              map[string]interface{} `json:"Tags,omitempty"`
	Variants          []string               `json:"Variants,omitempty"`
}

type Variant struct {
	Code                       string                 `json:"Code,omitempty"`
	FlavorCode                 string                 `json:"FlavorCode,omitempty"`
	ImageCode                  string                 `json:"ImageCode,omitempty"`
	Local                      bool                   `json:"Local,omitempty"`
	Name                       string                 `json:"Name,omitempty"`
	Price                      string                 `json:"Price,omitempty"`
	ProductCode                string                 `json:"ProductCode,omitempty"`
	SizeCode                   string                 `json:"SizeCode,omitempty"`
	Tags                       map[string]interface{} `json:"Tags,omitempty"`
	AllowedCookingInstructions string                 `json:"AllowedCookingInstructions,omitempty"`
	DefaultCookingInstructions string                 `json:"DefaultCookingInstructions,omitempty"`
	Prepared                   bool                   `json:"Prepared,omitempty"`
}

// not sure if needed yet
// type Categorization struct{}
// type Sides map[string]Side

func (c *Client) GetStoreMenu(storeID string) (*Menu, error) {
	url, err := url.Parse(fmt.Sprintf(menuURL, storeID))
	if err != nil {
		return nil, err
	}

	q := url.Query()
	q.Add("lang", "en")
	q.Add("structured", "true")
	url.RawQuery = q.Encode()

	resp, err := c.Get(url.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	menu := &Menu{}
	if err := json.NewDecoder(resp.Body).Decode(menu); err != nil {
		return nil, err
	}

	return menu, nil
}

// Rather than include a ton of lookup logic for coupons,
// just look up the ridiculous 50% off coupon that sometimes
// exists.
// TODO: real coupon lookups
func (m *Menu) GetFiftyPercentCouponCode() string {
	for _, coupon := range m.Coupons {
		if coupon.ImageCode == "OLO50" {
			return coupon.Code
		}
	}

	return ""
}
