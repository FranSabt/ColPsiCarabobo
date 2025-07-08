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

type SpecialtyUpdate struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AdmindId    string `json:"admin_id"`
}
