package sql

type Rows struct {
	Count  int
	cursor int
}

func (r *Rows) Next() bool { return r.cursor <= r.Count }
func (r *Rows) Scan(dest ...any) error {
	for _, d := range dest {
		switch v := d.(type) {
		case *int:
			*v = 21
		case *string:
			*v = "Alice"
		}
	}
	r.cursor++
	return nil
}
