package chart

type EnumMarshaler interface {
	GetEnum() (key string, value string)
	SetEnum(keyOrValue string) bool
}
