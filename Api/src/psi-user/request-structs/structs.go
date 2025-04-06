package psi_user_request

type PsiUserCreateRequest struct {
	Username string `json:"username"`
	// Identity
	FirstName      string `json:"firstName"`
	SecondName     string `json:"secondName"`
	LastName       string `json:"lastName"`
	SecondLastName string `json:"secondLastName"`
	Email          string `json:"email"`
	FPV            int    `json:"fpv"`
	CI             int    `json:"ci"`
	Nationality    string `json:"nationality"`
	BornDate       string `json:"bornDate"`
	Genre          string `json:"genre"`

	// Contact
	ContactEmail   string `json:"contactEmail"`
	PublicPhone    string `json:"publicPhone"`
	ServiceAddress string `json:"serviceAddress"`

	// Carabobo Direction
	MunicipalityCarabobo string `json:"municipalityCarabobo"`
	PhoneCarabobo        string `json:"phoneCarabobo"`
	CelPhoneCarabobo     string `json:"celPhoneCarabobo"`

	// Outside Carabobo Direction
	StateOutside                string `json:"stateOutside"`
	MunicipalityOutSideCarabobo string `json:"municipalityOutSideCarabobo"`
	PhoneOutSideCarabobo        string `json:"phoneOutSideCarabobo"`
	CelPhoneOutSideCarabobo     string `json:"celPhoneOutSideCarabobo"`

	// ------ PsiUSerColData ------ //
	UniversityUndergraduate string `json:"universityUndergraduate"`
	GraduateDate            string `json:"graduateDate"`
	MentionUndergraduate    string `json:"mentionUndergraduate"`

	// Undergraduate Data Title Register
	RegisterTitleState string `json:"registerTitleState"`
	RegisterTitleDate  string `json:"registerTitleDate"`
	RegisterNumber     int    `json:"registerNumber"`
	RegisterFolio      string `json:"registerFolio"`
	RegisterTome       string `json:"registerTome"`

	// Professional Data
	GuildDirector       bool `json:"guildDirector"`
	SixtyFiveOrPlus     bool `json:"sixtyFiveOrPlus"`
	GuildCollaborator   bool `json:"guildCollaborator"`
	PublicEmployee      bool `json:"publicEmployee"`
	UniversityProfessor bool `json:"universityProfessor"`

	// Otros campos
	DateOfLastSolvency string `json:"dateOfLastSolvency"`
	DoubleGuild        bool   `json:"doubleGuild"`
	CPSM               bool   `json:"cpsm"`
}
