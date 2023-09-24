package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type UserFailed struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Failed     string `json:"failed"`
}

func (s *Service) Distribution(u chan User, uFailed chan UserFailed) {

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Println("Ybito", sig)
			s.wg.Done()
			close(uFailed)
			run = false
		default:
			User1 := <-u
			chekErr, chek := s.CheckErr(User1)

			if chek {
				User2, err := s.EncrimentAge(User1)
				if err != nil {
					fmt.Println(err)
				}

				User2, err = s.EncrimentGender(User2)
				if err != nil {
					fmt.Println(err)
				}

				User2, err = s.EncrimentCountry(User2)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("age encriment", User2)
			} else {
				UserFail := UserFailed{
					Name:       User1.Name,
					Surname:    User1.Surname,
					Patronymic: User1.Patronymic,
					Failed:     chekErr,
				}
				uFailed <- UserFail
			}

		}

	}

}
