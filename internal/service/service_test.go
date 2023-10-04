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

func Test_Encrement(t *testing.T) {

	e := log.New(os.Stdout, "ERROR:  ", log.Ldate|log.Lshortfile)
	cfgT := config.NewCfgApi()
	s := service.New(e, cfgT)

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
