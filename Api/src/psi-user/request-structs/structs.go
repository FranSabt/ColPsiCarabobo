package psi_user_request

type PsiUserCreateRequest struct {
	Username string `json:"username" validate:"required"`

	// Identity Information
	FirstName      string `json:"firstName" validate:"required"`
	SecondName     string `json:"secondName" validate:"required"`
	LastName       string `json:"lastName" validate:"required"`
	SecondLastName string `json:"secondLastName" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	FPV            int    `json:"fpv" validate:"required"`
	CI             int    `json:"ci" validate:"required"`
	Nationality    string `json:"nationality" validate:"required"`
	BornDate       string `json:"bornDate" validate:"required"`
	Genre          string `json:"genre" validate:"required"`

	// Contact Information
	ContactEmail   string `json:"contactEmail" validate:"required,email"`
	PublicPhone    string `json:"publicPhone" validate:"required"`
	ServiceAddress string `json:"serviceAddress" validate:"required"`

	// Address Information
	// Carabobo Address
	MunicipalityCarabobo string `json:"municipalityCarabobo"`
	PhoneCarabobo        string `json:"phoneCarabobo"`
	CelPhoneCarabobo     string `json:"celPhoneCarabobo"`

	// Outside Carabobo Address
	StateOutside                string `json:"stateOutside"`
	MunicipalityOutSideCarabobo string `json:"municipalityOutSideCarabobo"`
	PhoneOutSideCarabobo        string `json:"phoneOutSideCarabobo"`
	CelPhoneOutSideCarabobo     string `json:"celPhoneOutSideCarabobo"`

	// Education Information
	UniversityUndergraduate string `json:"universityUndergraduate" validate:"required"`
	GraduateDate            string `json:"graduateDate" validate:"required"`
	MentionUndergraduate    string `json:"mentionUndergraduate" validate:"required"`

	// Title Registration Information
	RegisterTitleState string `json:"registerTitleState" validate:"required"`
	RegisterTitleDate  string `json:"registerTitleDate" validate:"required"`
	RegisterNumber     int    `json:"registerNumber" validate:"required"`
	RegisterFolio      string `json:"registerFolio" validate:"required"`
	RegisterTome       string `json:"registerTome" validate:"required"`

	// Professional Information
	GuildDirector       bool `json:"guildDirector" validate:"required"`
	SixtyFiveOrPlus     bool `json:"sixtyFiveOrPlus" validate:"required"`
	GuildCollaborator   bool `json:"guildCollaborator" validate:"required"`
	PublicEmployee      bool `json:"publicEmployee" validate:"required"`
	UniversityProfessor bool `json:"universityProfessor" validate:"required"`

	// Other Information
	DateOfLastSolvency string `json:"dateOfLastSolvency" validate:"required"`
	DoubleGuild        bool   `json:"doubleGuild" validate:"required"`
	CPSM               bool   `json:"cpsm" validate:"required"`
}

type PsiUserUpdateRequest struct {
	ID       string  `json:"id" validate:"required"`
	Username *string `json:"username,omitempty"`

	// Identity Information
	FirstName      *string `json:"firstName,omitempty"`
	SecondName     *string `json:"secondName,omitempty"`
	LastName       *string `json:"lastName,omitempty"`
	SecondLastName *string `json:"secondLastName,omitempty"`
	Email          *string `json:"email,omitempty"`
	FPV            *int    `json:"fpv,omitempty"`
	CI             *int    `json:"ci,omitempty"`
	Nationality    *string `json:"nationality,omitempty"`
	BornDate       *string `json:"bornDate,omitempty"`
	Genre          *string `json:"genre,omitempty"`

	// Contact Information
	ContactEmail   *string `json:"contactEmail,omitempty"`
	PublicPhone    *string `json:"publicPhone,omitempty"`
	ServiceAddress *string `json:"serviceAddress,omitempty"`

	// Address Information
	// Carabobo Address
	MunicipalityCarabobo *string `json:"municipalityCarabobo,omitempty"`
	PhoneCarabobo        *string `json:"phoneCarabobo,omitempty"`
	CelPhoneCarabobo     *string `json:"celPhoneCarabobo,omitempty"`

	// Outside Carabobo Address
	StateOutside                *string `json:"stateOutside,omitempty"`
	MunicipalityOutSideCarabobo *string `json:"municipalityOutSideCarabobo,omitempty"`
	PhoneOutSideCarabobo        *string `json:"phoneOutSideCarabobo,omitempty"`
	CelPhoneOutSideCarabobo     *string `json:"celPhoneOutSideCarabobo,omitempty"`

	// Education Information
	UniversityUndergraduate *string `json:"universityUndergraduate,omitempty"`
	GraduateDate            *string `json:"graduateDate,omitempty"`
	MentionUndergraduate    *string `json:"mentionUndergraduate,omitempty"`

	// Title Registration Information
	RegisterTitleState *string `json:"registerTitleState,omitempty"`
	RegisterTitleDate  *string `json:"registerTitleDate,omitempty"`
	RegisterNumber     *int    `json:"registerNumber,omitempty"`
	RegisterFolio      *string `json:"registerFolio,omitempty"`
	RegisterTome       *string `json:"registerTome,omitempty"`

	// Professional Information
	GuildDirector       *bool `json:"guildDirector,omitempty"`
	SixtyFiveOrPlus     *bool `json:"sixtyFiveOrPlus,omitempty"`
	GuildCollaborator   *bool `json:"guildCollaborator,omitempty"`
	PublicEmployee      *bool `json:"publicEmployee,omitempty"`
	UniversityProfessor *bool `json:"universityProfessor,omitempty"`

	// Other Information
	DateOfLastSolvency *string `json:"dateOfLastSolvency,omitempty"`
	DoubleGuild        *bool   `json:"doubleGuild,omitempty"`
	CPSM               *bool   `json:"cpsm,omitempty"`
}
