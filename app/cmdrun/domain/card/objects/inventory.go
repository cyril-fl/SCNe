package objects

// Inventory holds one object per unit; identical objects share a name.
type Inventory []Object

// Stack is a run of objects with the same name.
type Stack struct {
	Object
	Qty int
}

// Getters & Setters

// Count returns how many objects are named name.
func (inv Inventory) Count(name string) int {
	count := 0
	for _, object := range inv {
		if object.GetName() == name {
			count++
		}
	}
	return count
}

// Stacks groups the objects by name, in order of first appearance.
func (inv Inventory) Stacks() []Stack {
	stacks := []Stack{}
	index := map[string]int{}

	for _, object := range inv {
		if i, ok := index[object.GetName()]; ok {
			stacks[i].Qty++
			continue
		}
		index[object.GetName()] = len(stacks)
		stacks = append(stacks, Stack{Object: object, Qty: 1})
	}
	return stacks
}

// Methods

// Add appends qty copies of object.
func (inv *Inventory) Add(object Object, qty int) {
	for range qty {
		*inv = append(*inv, object)
	}
}

// Remove drops qty objects named name, the last ones first; it removes nothing
// and reports false when fewer are held.
func (inv *Inventory) Remove(name string, qty int) bool {
	if qty <= 0 || inv.Count(name) < qty {
		return false
	}

	for i := len(*inv) - 1; i >= 0 && qty > 0; i-- {
		if (*inv)[i].GetName() == name {
			*inv = append((*inv)[:i], (*inv)[i+1:]...)
			qty--
		}
	}
	return true
}
