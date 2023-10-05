package service_test

import (
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hulla-hoop/testSobes/internal/config"
	"github.com/hulla-hoop/testSobes/internal/modeldb"
	"github.com/hulla-hoop/testSobes/internal/service"
)

var s *service.Service

func TestMain(m *testing.M) {
	os.Setenv("AGEAPI", "https://api.agify.io/?name=%s")
	os.Setenv("NATIONAPI", "https://api.nationalize.io/?name=%s")
	os.Setenv("GENDERAPI", "https://api.genderize.io/?name=%s")
	e := log.New(os.Stdout, "ERROR:  ", log.Ldate|log.Lshortfile)
	cfgT := config.NewCfgApi()
	s = service.New(e, cfgT)
	exitVal := m.Run()
	os.Exit(exitVal)

}

func TestEncrement(t *testing.T) {

	userTest := modeldb.User{
		Name:       "Dmitriy",
		Surname:    "Olegovich",
		Patronymic: "Kirirlov",
	}

	userExpected := modeldb.User{
		Name:        "Dmitriy",
		Surname:     "Olegovich",
		Patronymic:  "Kirirlov",
		Age:         42,
		Gender:      "male",
		Nationality: "UA",
	}

	comparer := cmp.Comparer(func(x, y modeldb.User) bool {
		return x.Name == y.Name && x.Surname == y.Surname && x.Patronymic == y.Patronymic && x.Age == y.Age && x.Gender == y.Gender && x.Nationality == y.Nationality
	})

	result, err := s.Encriment(userTest)

	if diff := cmp.Diff(userExpected, result, comparer); diff != "" {
		t.Errorf(diff, err)
	}

}

func TestCheckErr(t *testing.T) {

	data := []struct {
		name     string
		user     modeldb.User
		want     string
		wantBool bool
	}{
		{
			name: "u1",
			user: modeldb.User{Name: "Shamil",
				Surname: "Suleimanov"},
			want:     "",
			wantBool: true,
		},
		{
			name: "u2",
			user: modeldb.User{Name: "Sha111mil",
				Surname: "Suleimanov"},
			want:     "Неверный формат",
			wantBool: false,
		},
		{
			name: "u3",
			user: modeldb.User{Name: "Shamil",
				Surname: "Suleima1111nov"},
			want:     "Неверный формат",
			wantBool: false,
		},
		{
			name: "u4",
			user: modeldb.User{Name: "",
				Surname: "Suleima1111nov"},
			want:     "Нет обязательного поля",
			wantBool: false,
		}, {
			name: "u5",
			user: modeldb.User{Name: "Shamil",
				Surname: ""},
			want:     "Нет обязательного поля",
			wantBool: false,
		},
		{
			name: "u6",
			user: modeldb.User{Name: "Sham%il",
				Surname: ""},
			want:     "Неверный формат",
			wantBool: false,
		},
		{
			name: "u7",
			user: modeldb.User{
				Surname: ""},
			want:     "Нет обязательного поля",
			wantBool: false,
		},
		{
			name:     "u8",
			user:     modeldb.User{Name: "Shamil"},
			want:     "Нет обязательного поля",
			wantBool: false,
		}, {
			name:     "u9",
			user:     modeldb.User{Name: "Sh am l", Surname: "Suleimanov"},
			want:     "Нет обязательного поля",
			wantBool: false,
		}, {
			name: "u10",
			user: modeldb.User{Name: "Shamil",
				Surname:    "Suleimanov",
				Patronymic: ""},
			want:     "",
			wantBool: true,
		}, {
			name: "u11",
			user: modeldb.User{Name: "Shamil",
				Surname:    "Suleimanov",
				Patronymic: "Alievich444"},
			want:     "",
			wantBool: false,
		},
		{
			name: "u12",
			user: modeldb.User{Name: " Shamil",
				Surname: "Suleimanov"},
			want:     "Неверный формат поля имя",
			wantBool: false,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result, ist := s.CheckErr(d.user)
			if ist != d.wantBool {
				t.Errorf("Expected %t | %v , got %t  | %v", d.wantBool, d.want, ist, result)
			}
		})
	}

}
