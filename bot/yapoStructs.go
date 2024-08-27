package bot

type FeaturesResponse struct {
	Features []Features `json:"features,omitempty"`
}
type Features struct {
	Label string `json:"label,omitempty"`
	Value string `json:"value,omitempty"`
}

type YapoResponse struct {
	Pagination Pagination   `json:"pagination,omitempty"`
	Ads        []Ads        `json:"ads,omitempty"`
	Categories []Categories `json:"categories,omitempty"`
}
type Pagination struct {
	TotalRecords int `json:"totalRecords,omitempty"`
	Page         int `json:"page,omitempty"`
	Limit        int `json:"limit,omitempty"`
}
type Category struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	ParentID   int    `json:"parentId,omitempty"`
	ParentName string `json:"parentName,omitempty"`
}
type Location struct {
	RegionID    int    `json:"regionId,omitempty"`
	RegionName  string `json:"regionName,omitempty"`
	CommuneID   int    `json:"communeId,omitempty"`
	CommuneName string `json:"communeName,omitempty"`
}
type Media struct {
	URLThumbnails string `json:"urlThumbnails,omitempty"`
	URLImage      string `json:"urlImage,omitempty"`
	SeqNo         int    `json:"seqNo,omitempty"`
}
type Ads struct {
	ListID        int      `json:"listId,omitempty"`
	ListTime      string   `json:"listTime,omitempty"`
	Type          string   `json:"type,omitempty"`
	Name          string   `json:"name,omitempty"`
	Subject       string   `json:"subject,omitempty"`
	Body          string   `json:"body,omitempty"`
	PublisherType string   `json:"publisherType,omitempty"`
	Price         int      `json:"price,omitempty"`
	OldPrice      int      `json:"oldPrice,omitempty"`
	Category      Category `json:"category,omitempty"`
	Location      Location `json:"location,omitempty"`
	Media         []Media  `json:"media,omitempty"`
	EstateType    string   `json:"estateType,omitempty"`
	Currency      string   `json:"currency,omitempty"`
	Rooms         int      `json:"rooms,omitempty"`
	Bathrooms     int      `json:"bathrooms,omitempty"`
	Size          int      `json:"size,omitempty"`
	UtilSize      int      `json:"utilSize,omitempty"`
	IsBuilder     bool     `json:"isBuilder,omitempty"`
	SellerName    string   `json:"sellerName,omitempty"`
	IsActioned    bool     `json:"isActioned,omitempty"`
	IsEcommerce   bool     `json:"isEcommerce,omitempty"`
	Email         string   `json:"email,omitempty"`
	GarageSpaces  int      `json:"garageSpaces,omitempty"`
	Condominio    int      `json:"condominio,omitempty"`
}
type Categories struct {
	Category Category `json:"category,omitempty"`
	DocCount int      `json:"docCount,omitempty"`
}
