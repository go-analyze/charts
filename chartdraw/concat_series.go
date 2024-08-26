package chartdraw

// ConcatSeries is a special type of series that concatenates its `InnerSeries`.
type ConcatSeries []Series

// Len returns the length of the concatenated set of series.
func (cs ConcatSeries) Len() int {
	total := 0
	for _, s := range cs {
		if typed, isValuesProvider := s.(ValuesProvider); isValuesProvider {
			total += typed.Len()
		}
	}

	return total
}

// GetValue returns the value at the (meta) index (i.e 0 => totalLen-1)
func (cs ConcatSeries) GetValue(index int) (x, y float64) {
	cursor := 0
	for _, s := range cs {
		if typed, isValuesProvider := s.(ValuesProvider); isValuesProvider {
			length := typed.Len()
			if index < cursor+length {
				x, y = typed.GetValues(index - cursor) //FENCEPOSTS.
				return
			}
			cursor += length
		}
	}
	return
}

// Validate validates the series.
func (cs ConcatSeries) Validate() error {
	for _, s := range cs {
		if err := s.Validate(); err != nil {
			return err
		}
	}
	return nil
}
