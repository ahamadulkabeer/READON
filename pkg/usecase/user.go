package usecase

import (
	"errors"
	"fmt"

	"readon/pkg/api/helpers"
	domain "readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"

	"github.com/jinzhu/copier"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(urepo interfaces.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: urepo,
	}
}

func (c *userUseCase) Save(user models.SignupData) (domain.User, error) {

	fmt.Println("user:", user)
	var User domain.User
	copier.Copy(&User, &user)

	err := helpers.ValidateUserData(&User)
	if err != nil {
		return User, err
	}
	err = c.userRepo.CheckForEmail(user.Email)
	if err == nil {
		return User, errors.New("Email already has an account ")
	}
	User, err = c.userRepo.Save(User)
	return User, err
}

func (c *userUseCase) UpdateUser(user models.UpdateData) (domain.User, error) {

	fmt.Println("user:", user)
	var User domain.User
	copier.Copy(&User, &user)

	err := helpers.ValidateUserUPdateData(&user)
	if err != nil {
		return User, err
	}
	User, err = c.userRepo.UpdateUser(User)
	return User, err
}

func (c *userUseCase) UserLogin(userinput models.Userlogindata) (int, bool, bool, error) {

	user, err := c.userRepo.FindByEmail(userinput.Email)
	if err != nil {
		return 0, false, false, errors.New("invalid username or password")
	}
	if user.Password != userinput.Password {
		//return 0, false, errors.New("password does not match")
		return 0, false, false, errors.New("invalid username or password")
	}
	if user.Permission == false {
		return 0, false, false, errors.New("user have benn blocked by the admin")
	}

	return int(user.ID), true, user.Premium, err
}

func (c userUseCase) GetUserProfile(id int) (domain.User, error) {
	var user domain.User
	user, err := c.userRepo.FindByID(uint(id))
	if err != nil {
		return user, err
	}
	return user, err
}

func (c userUseCase) DeleteUserAccount(id int) error {
	var user domain.User
	user, err := c.userRepo.FindByID(uint(id))
	if err != nil {
		return err
	}
	err = c.userRepo.DeleteUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (c userUseCase) VerifyAndSendOtp(email string) error {
	err := c.userRepo.CheckForEmail(email)
	if err != nil {
		return errors.New("Invalid email")
	}
	otp, err := helpers.GenerateAndSendOpt(email)
	fmt.Println("otp   : ", otp)
	if err != nil {
		return errors.New("could not send otp")
	}
	c.userRepo.SaveOtp(otp, email)
	return nil
}

func (c userUseCase) VerifyOtp(otp, email string) error {
	err := c.userRepo.VerifyOtp(otp, email)
	return err
}
