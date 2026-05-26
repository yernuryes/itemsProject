package entities

type Item struct {
	ID             int64
	Name           string
	Amount         int
	Price          int
	Promocode      string
	CategoryID     int
	Category       *Category       `gorm:"foreignKey:CategoryID"`
	ClientSegments []ClientSegment `gorm:"many2many:item_client_segments"`
}

type Category struct {
	ID   int64
	Name string
}

type ClientSegment struct {
	ID    int64
	Name  string
	Items []Item `gorm:"many2many:item_client_segments"`
}

type ItemClientSegments struct {
	ItemID          int64
	ClientSegmentID int64
}

type Car struct {
	ID        int64
	Name      string
	Model     string
	Year      int
	Promocode string
}
