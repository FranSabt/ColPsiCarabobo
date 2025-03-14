package psi_user_controller

import (
	"encoding/csv"
	"fmt"
	"io"
)

type PsiUserCsv struct {
	UserName                    string
	Email                       string
	Password                    string
	FirstName                   string
	SecondName                  string
	LastName                    string
	SecondLastName              string
	FPV                         string
	CI                          string
	Letter                      string
	Nationality                 string
	BornDate                    string
	Genre                       string
	ContactEmail                string
	ShowContactEmail            string
	PublicPhone                 string
	ShowPublicPhone             string
	ServiceAddress              string
	ShowPublicServiceAddress    string
	Solvent                     string
	MunicipalityCarabobo        string
	PhoneCarabobo               string
	CelPhoneCarabobo            string
	State                       string
	MunicipalityOutSideCarabobo string
	PhoneOutSideCarabobo        string
	CelPhoneOutSideCarabobo     string
	UniversityUndergraduate     string
	GraduateDate                string
	MentionUndergraduate        string
	RegisterTitleState          string
	RegisterTitleDate           string
	RegisterNumber              string
	RegisterFolio               string
	RegisterTome                string
	GuildDirector               string
	SixtyFiveOrPlus             string
	GuildCollaborator           string
	PublicEmployee              string
	UniversityProfessor         string
	DateOfLastSolvency          string
	DoubleGuild                 string
	CPSM                        string
}

func ProcessCsv(src io.Reader) (*[]PsiUserCsv, error) {
	// Crea un nuevo lector CSV
	reader := csv.NewReader(src)

	// Lee y omite la primera línea (encabezados)
	_, err := reader.Read()
	if err != nil {
		fmt.Println("Error al leer el encabezado:", err)
		return nil, err
	}

	var psi_users []PsiUserCsv

	// Lee el resto de las líneas
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error al leer el archivo CSV:", err)
			return nil, err
		}

		// Crea una estructura PsiUserCsv con los datos del registro
		user := PsiUserCsv{
			UserName:                    record[0],
			Email:                       record[1],
			Password:                    RandomPass(),
			FirstName:                   record[3],
			SecondName:                  record[4],
			LastName:                    record[5],
			SecondLastName:              record[6],
			FPV:                         record[7],
			CI:                          record[8],
			Letter:                      record[9],
			Nationality:                 record[10],
			BornDate:                    record[11],
			Genre:                       record[12],
			ContactEmail:                record[13],
			ShowContactEmail:            record[14],
			PublicPhone:                 record[15],
			ShowPublicPhone:             record[16],
			ServiceAddress:              record[17],
			ShowPublicServiceAddress:    record[18],
			Solvent:                     record[19],
			MunicipalityCarabobo:        record[20],
			PhoneCarabobo:               record[21],
			CelPhoneCarabobo:            record[22],
			State:                       record[23],
			MunicipalityOutSideCarabobo: record[24],
			PhoneOutSideCarabobo:        record[25],
			CelPhoneOutSideCarabobo:     record[26],
			UniversityUndergraduate:     record[27],
			GraduateDate:                record[28],
			MentionUndergraduate:        record[29],
			RegisterTitleState:          record[30],
			RegisterTitleDate:           record[31],
			RegisterNumber:              record[32],
			RegisterFolio:               record[33],
			RegisterTome:                record[34],
			GuildDirector:               record[35],
			SixtyFiveOrPlus:             record[36],
			GuildCollaborator:           record[37],
			PublicEmployee:              record[38],
			UniversityProfessor:         record[39],
			DateOfLastSolvency:          record[40],
			DoubleGuild:                 record[41],
			CPSM:                        record[42],
		}

		// Imprime la estructura PsiUserCsv
		fmt.Printf("%+v\n", user)
		psi_users = append(psi_users, user)
	}
	return &psi_users, nil
}
