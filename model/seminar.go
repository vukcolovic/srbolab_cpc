package model

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

const (
	SEMINAR_STATUS_OPENED      = 1
	SEMINAR_STATUS_FILLED      = 2
	SEMINAR_STATUS_IN_PROGRESS = 3
	SEMINAR_STATUS_CLOSED      = 4
)

type Seminar struct {
	gorm.Model
	Start                  time.Time       `json:"start_date"`
	ClassRoomID            uint            `json:"class_room_id"`
	ClassRoom              ClassRoom       `json:"class_room"`
	Trainees               []ClientSeminar `json:"trainees"`
	SeminarThemeID         uint            `json:"seminar_theme_id"`
	SeminarTheme           SeminarTheme    `json:"seminar_theme"`
	SeminarStatusID        uint            `json:"seminar_status_id"`
	SeminarStatus          SeminarStatus   `json:"seminar_status"`
	Days                   []SeminarDay    `json:"days"`
	Documents              []*File         `json:"documents" gorm:"many2many:seminar_file;"`
	SerialNumberByLocation int             `json:"serial_number_by_location"`
}

// Šifra obuke treba da se generiše automatski po principu: 7.3-226-LA-23-10-06
// gde je 7 - oznaka za periodičnu obuku (kada je dodatna obuka umesto
// 7.3 biće 35, kada je osnovna obuka biće 140),
// 3 je redni broj teme,
// 226 redni broj seminara u Lazarevcu (svaka lokacija kreće sa
// računanjem od 1 i broji za sebe),
// LA oznaka za lokaciju (LA-Lazarevac, BG-Beograd Batajnica,
// SRB-Srbobran, PO-Požarevac),
// 23-10-06 datum seminara,
func (s Seminar) GetCode() string {
	code := ""
	if s.SeminarTheme.BaseSeminarType.Code == "BASIC" {
		code = code + "140"
	}
	if s.SeminarTheme.BaseSeminarType.Code == "ADDITIONAL" {
		code = code + "35"
	}
	if s.SeminarTheme.BaseSeminarType.Code == "CYCLE" {
		code = code + "7" + s.SeminarTheme.Code
	}

	code = code + "-" + strconv.Itoa(s.SerialNumberByLocation)
	code = code + "-" + s.ClassRoom.Location.Code
	code = code + "-" + s.Start.Format("2006-02-01")
	return code
}

type BaseSeminarType struct {
	gorm.Model
	Code string `json:"code"`
	Name string `json:"name"`
}

func (b BaseSeminarType) GetSeminarTypeForSentence() string {
	if b.Code == "ADDITIONAL" {
		return "додатној"
	}
	if b.Code == "BASE" {
		return "основној"
	}
	return "периодичној"
}

type SeminarTheme struct {
	gorm.Model
	BaseSeminarTypeID uint            `json:"base_seminar_type_id"`
	BaseSeminarType   BaseSeminarType `json:"base_seminar_type"`
	Code              string          `json:"code"`
	Name              string          `json:"name"`
	NumberOfDays      int             `json:"number_of_days"`
}

func (s SeminarTheme) GetSeminarThemeWithBaseTheme() string {
	base := ""
	if s.BaseSeminarType.Code != "BASIC" && s.BaseSeminarType.Code != "ADDITIONAL" {
		base = s.BaseSeminarType.Name
	}
	return base + " " + s.Name
}

type SeminarStatus struct {
	gorm.Model
	Code string `json:"code"`
	Name string `json:"name"`
}

type SeminarDay struct {
	gorm.Model
	Date      time.Time        `json:"date"`
	Number    int              `json:"number"`
	Name      string           `json:"name"`
	Classes   []SeminarClass   `json:"classes"`
	SeminarID uint             `json:"seminar_id"`
	Seminar   Seminar          `json:"seminar"`
	TestID    *uint            `json:"test_id"`
	Test      Test             `json:"test"`
	Presence  []ClientPresence `json:"presence"`
	Documents []*File          `json:"documents" gorm:"many2many:seminarday_file;"`
}

type SeminarClass struct {
	gorm.Model
	Number int    `json:"number"`
	Name   string `json:"name"`
	//Teachers     User       `json:"teachers" gorm:"many2many:class_teacher;"`
	TeacherID    *int       `json:"teacher_id" gorm:"default:null"`
	Teacher      *User      `json:"teacher"`
	SeminarDayID uint       `json:"seminar_day_id"`
	SeminarDay   SeminarDay `json:"seminar_day"`
}

func GetSeminarClassByNumber(classes []SeminarClass, number int) *SeminarClass {
	for _, class := range classes {
		if class.Number == number {
			return &class
		}
	}

	return nil
}

type SeminarClassName struct {
	gorm.Model
	SeminarThemeID uint         `json:"seminar_theme_id"`
	SeminarTheme   SeminarTheme `json:"seminar_theme"`
	DayNumber      int          `json:"day_number"`
	ClassNumber    int          `json:"class_number"`
	ClassName      string       `json:"class_name"`
}

type SeminarDayThemeName struct {
	gorm.Model
	SeminarThemeID uint         `json:"seminar_theme_id"`
	SeminarTheme   SeminarTheme `json:"seminar_theme"`
	DayNumber      int          `json:"day_number"`
	ThemeName      string       `json:"theme_name"`
}
