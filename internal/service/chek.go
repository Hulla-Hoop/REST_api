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
			User := <-u
			chekErr, chek := s.CheckErr(User)

			if chek {
				User, err := s.EncrimentAge(User)
				if err != nil {
					fmt.Println(err)
				}

				User, err = s.EncrimentGender(User)
				if err != nil {
					fmt.Println(err)
				}

				User, err = s.EncrimentCountry(User)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("age encriment", User)
			} else {
				UserFail := UserFailed{
					Name:       User.Name,
					Surname:    User.Surname,
					Patronymic: User.Patronymic,
					Failed:     chekErr,
				}
				uFailed <- UserFail
			}

		}

	}

}
