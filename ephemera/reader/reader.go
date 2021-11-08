package reader

const (
	ItemId    = "id"
	ItemType  = "type"
	ItemValue = "value"
)

func At(m Map) string {
	return m.StrOf(ItemId)
}
