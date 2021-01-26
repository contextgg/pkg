package types

var local Register

func init() {
	local = NewRegistry()
}

func GetAll() map[string]*TypeData {
	return local.GetAll()
}
func GetTypeData(name string) (*TypeData, bool) {
	return local.GetTypeData(name)
}
func SetTypeDataByFactory(fn TypeFunc, internalType bool) *TypeData {
	return local.SetTypeDataByFactory(fn, internalType)
}
func SetTypeData(t interface{}, internalType bool) *TypeData {
	return local.SetTypeData(t, internalType)
}
func UnmarshalByName(name string, raw []byte, legacy bool) (interface{}, error) {
	return local.UnmarshalByName(name, raw, legacy)
}
