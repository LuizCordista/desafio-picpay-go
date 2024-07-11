package services

import (
	"desafio-picpay/models"
	"desafio-picpay/repositories"
	"desafio-picpay/utils"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
)

type UserService interface {
	CreateUser(user *models.User) error
	Transfer(payerCPF, payeeCPF string, value float64) error
}

type userService struct {
	repo repositories.GormUserRepository
}

func NewUserService(repo repositories.GormUserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(user *models.User) error {
	user.Currency = 0

	user.CPF = utils.CleanCPF(user.CPF)

	if !utils.ValidateCPF(user.CPF) {
		return errors.New("CPF inválido")
	}

	if !utils.ValidateEmail(user.Email) {
		return errors.New("Email inválido")
	}

	if _, err := s.repo.FindByCPF(user.CPF); err == nil {
		return errors.New("CPF já existe no banco de dados")
	}

	if _, err := s.repo.FindByEmail(user.Email); err == nil {
		return errors.New("Email já existe no banco de dados")
	}

	if err := s.repo.Create(user); err != nil {
		return err
	}

	return nil
}

func (s *userService) Transfer(payerCPF, payeeCPF string, value float64) error {
	err := s.repo.GetInstance().Transaction(func(tx *gorm.DB) error {

		var err error

		payer, err := s.repo.FindByCPF(payerCPF)
		if err != nil {
			return errors.New("pagador não encontrado")
		}

		payee, err := s.repo.FindByCPF(payeeCPF)
		if err != nil {
			return errors.New("recebedor não encontrado")
		}

		if payer.Currency < value {
			return fmt.Errorf("saldo insuficiente")
		}

		if payer.IsShop {
			return fmt.Errorf("lojistas não podem transferir dinheiro")
		}

		resp, err := http.Get("https://util.devi.tools/api/v2/authorize")
		if err != nil {
			return fmt.Errorf("falha ao contatar serviço de autorização")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("não autorizado")
		}

		payer.Currency -= value
		payee.Currency += value

		if err := tx.Save(&payer).Error; err != nil {
			return err
		}

		if err := tx.Save(&payee).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
