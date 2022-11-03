package main

type op int

const (
	createOp op = iota
	updateOp
	deleteOp
)

func reconcileLists(plan, fact []string) (update, create, delete []string) {
	work := make(map[string]op)
	for _, f := range fact {
		work[f] = deleteOp
	}
	for _, p := range plan {
		if _, ok := work[p]; ok {
			work[p] = updateOp
		} else {
			work[p] = createOp
		}
	}
	for k, v := range work {
		switch v {
		case createOp:
			create = append(create, k)
		case updateOp:
			update = append(update, k)
		case deleteOp:
			delete = append(delete, k)
		}
	}
	return create, update, delete
}
