package psi_user_request

type PsiUserCreateRequest struct {
	Username string `json:"username" validate:"required"`

	// Identity Information
	FirstName      string `json:"first_name" validate:"required"`
	SecondName     string `json:"second_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	SecondLastName string `json:"second_last_name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	FPV            int    `json:"fpv" validate:"required"`
	CI             int    `json:"ci" validate:"required"`
	Nationality    string `json:"nationality" validate:"required"`
	BornDate       string `json:"born_date" validate:"required"`
	Genre          string `json:"genre" validate:"required"`

	// Contact Information
	ContactEmail   string `json:"contact_email" validate:"required,email"`
	PublicPhone    string `json:"public_phone" validate:"required"`
	ServiceAddress string `json:"service_address" validate:"required"`

	// Address Information
	// Carabobo Address
	MunicipalityCarabobo string `json:"municipality_carabobo"`
	PhoneCarabobo        string `json:"phone_carabobo"`
	CelPhoneCarabobo     string `json:"cel_phone_carabobo"`

	// Outside Carabobo Address
	StateOutside                string `json:"state_outside"`
	MunicipalityOutSideCarabobo string `json:"municipality_out_side_carabobo"`
	PhoneOutSideCarabobo        string `json:"phone_out_side_carabobo"`
	CelPhoneOutSideCarabobo     string `json:"cel_phone_out_side_carabobo"`

	// Education Information
	UniversityUndergraduate string `json:"university_undergraduate" validate:"required"`
	GraduateDate            string `json:"graduate_date" validate:"required"`
	MentionUndergraduate    string `json:"mention_undergraduate" validate:"required"`

	// Title Registration Information
	RegisterTitleState string `json:"register_title_state" validate:"required"`
	RegisterTitleDate  string `json:"register_title_date" validate:"required"`
	RegisterNumber     int    `json:"register_number" validate:"required"`
	RegisterFolio      string `json:"register_folio" validate:"required"`
	RegisterTome       string `json:"register_tome" validate:"required"`

	// Professional Information
	GuildDirector       bool `json:"guild_director" validate:"required"`
	SixtyFiveOrPlus     bool `json:"sixty_five_or_plus" validate:"required"`
	GuildCollaborator   bool `json:"guild_collaborator" validate:"required"`
	PublicEmployee      bool `json:"public_employee" validate:"required"`
	UniversityProfessor bool `json:"university_professor" validate:"required"`

	// Other Information
	DateOfLastSolvency string `json:"date_of_last_solvency" validate:"required"`
	DoubleGuild        bool   `json:"double_guild" validate:"required"`
	CPSM               bool   `json:"cpsm" validate:"required"`
}

type PsiUserUpdateRequest struct {
	ID       string  `json:"id" validate:"required"`
	Username *string `json:"username,omitempty"`

	// Identity Information
	FirstName      *string `json:"first_name,omitempty"`
	SecondName     *string `json:"second_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	SecondLastName *string `json:"second_last_name,omitempty"`
	Email          *string `json:"email,omitempty"`
	FPV            *int    `json:"fpv,omitempty"`
	CI             *int    `json:"ci,omitempty"`
	Nationality    *string `json:"nationality,omitempty"`
	BornDate       *string `json:"born_date,omitempty"`
	Genre          *string `json:"genre,omitempty"`

	// Contact Information
	ContactEmail   *string `json:"contact_email,omitempty"`
	PublicPhone    *string `json:"public_phone,omitempty"`
	ServiceAddress *string `json:"service_address,omitempty"`

	// Address Information
	// Carabobo Address
	MunicipalityCarabobo *string `json:"municipality_carabobo,omitempty"`
	PhoneCarabobo        *string `json:"phone_carabobo,omitempty"`
	CelPhoneCarabobo     *string `json:"cel_phone_carabobo,omitempty"`

	// Outside Carabobo Address
	StateOutside                *string `json:"state_outside,omitempty"`
	MunicipalityOutSideCarabobo *string `json:"municipality_out_side_carabobo,omitempty"`
	PhoneOutSideCarabobo        *string `json:"phone_out_side_carabobo,omitempty"`
	CelPhoneOutSideCarabobo     *string `json:"cel_phone_out_side_carabobo,omitempty"`

	// Education Information
	UniversityUndergraduate *string `json:"university_undergraduate,omitempty"`
	GraduateDate            *string `json:"graduate_date,omitempty"`
	MentionUndergraduate    *string `json:"mention_undergraduate,omitempty"`

	// Title Registration Information
	RegisterTitleState *string `json:"register_title_state,omitempty"`
	RegisterTitleDate  *string `json:"register_title_date,omitempty"`
	RegisterNumber     *int    `json:"register_number,omitempty"`
	RegisterFolio      *string `json:"register_folio,omitempty"`
	RegisterTome       *string `json:"register_tome,omitempty"`

	// Professional Information
	GuildDirector       *bool `json:"guild_director,omitempty"`
	SixtyFiveOrPlus     *bool `json:"sixty_five_or_plus,omitempty"`
	GuildCollaborator   *bool `json:"guild_collaborator,omitempty"`
	PublicEmployee      *bool `json:"public_employee,omitempty"`
	UniversityProfessor *bool `json:"university_professor,omitempty"`

	// Other Information
	DateOfLastSolvency *string `json:"date_of_last_solvency,omitempty"`
	DoubleGuild        *bool   `json:"double_guild,omitempty"`
	CPSM               *bool   `json:"cpsm,omitempty"`
}
