package difference

func NewDifference(ID, oldModel, newModel interface{}) *Difference {
	return &Difference{
		OldData: oldModel,
		NewData: newModel,
		ID:      ID,
		data:    make(map[string]interface{}),
	}
}
