package specialties_structs

type SpecialtiesRequest struct {
	Name        string
	Description string
	AdmindId    string
}

type SpecialtyName struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
