package jsn

func BoxBool(v *bool) (ret string) {
	if *v {
		ret = "$TRUE"
	} else {
		ret = "$FALSE"
	}
	return
}
func BoxFloat64(v *float64) float64 {
	return *v
}
func BoxString(v *string) string {
	return *v
}
