package service

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hulla-hoop/testSobes/internal/modeldb"
)

type UserFailed struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Failed     string `json:"failed"`
}

func (s *Service) Distribution(u chan modeldb.User, uFailed chan UserFailed) {

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			s.inflogger.Println("Выход из горутины Service прекращено сигналом - ", sig)
			s.wg.Done()
			close(uFailed)
			run = false
		default:
			User := <-u
			chekErr, chek := s.CheckErr(User)

			if chek {
				User, err := s.EncrimentAge(User)
				if err != nil {
					s.errLogger.Panicln(err)
				}

				User, err = s.EncrimentGender(User)
				if err != nil {
					s.errLogger.Println(err)
				}

				User, err = s.EncrimentCountry(User)
				if err != nil {
					s.errLogger.Panicln(err)
				}
				User.CreatedAt = time.Now()
				User.UpdatedAt = time.Now()
				s.inflogger.Println("Сообщение готово к хранению в БД", User)
				s.db.Create(User)

			} else {
				UserFail := UserFailed{
					Name:       User.Name,
					Surname:    User.Surname,
					Patronymic: User.Patronymic,
					Failed:     chekErr,
				}
				s.inflogger.Println("Сообщение не прошло проверку и отправлено в очередь FIO_FAIL")
				uFailed <- UserFail
			}

		}

	}

}
