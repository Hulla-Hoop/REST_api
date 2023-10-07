package service

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hulla-hoop/testSobes/internal/modeldb"
	"github.com/hulla-hoop/testSobes/internal/psql"
)

type UserFailed struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Failed     string `json:"failed"`
}

type Servicer interface {
	Encriment(u modeldb.User) (modeldb.User, error)
	EncrimentAge(uName string) (int, error)
	EncrimentGender(uName string) (string, error)
	EncrimentCountry(uName string) (string, error)
	CheckErr(U modeldb.User) (string, bool)
}

// Функция обогащает верные сообщения и ложит в БД, невернные сообщения отправляются в очередь FIO_FAILED

func Distribution(s Servicer, u chan modeldb.User, uFailed chan UserFailed, infoLogger *log.Logger, errLogger *log.Logger, wg *sync.WaitGroup, db psql.DB) {

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			infoLogger.Println("Выход из горутины Service прекращено сигналом - ", sig)
			wg.Done()
			close(uFailed)
			run = false
		default:
			User := <-u
			chekErr, chek := s.CheckErr(User)
			if chek {
				User, err := s.Encriment(User)
				infoLogger.Println("Сообщение готово к хранению в БД", User)
				err = db.Create(&User)
				if err != nil {
					errLogger.Println(err)
				}

			} else {
				UserFail := UserFailed{
					Name:       User.Name,
					Surname:    User.Surname,
					Patronymic: User.Patronymic,
					Failed:     chekErr,
				}

				infoLogger.Println("Сообщение не прошло проверку и отправлено в очередь FIO_FAIL")
				uFailed <- UserFail
			}

		}

	}

}
