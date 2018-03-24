package character

type Character struct {
	*Actor
}

func GetCharacterInstance(a *Actor) *Character {
	return &Character{a}
}
