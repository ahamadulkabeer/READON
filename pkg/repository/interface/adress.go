package interfaces

import "readon/pkg/domain"

type AdddressRepository interface {
	AddAdress(adress domain.Adress) error
	EditAddress(newadress domain.Adress) error
	ListAdresses(userid int) ([]domain.Adress, error)
	GetAdress(userid, adresid int) (domain.Adress, error)
	DeleteAdress(userid, adresid int) error
}
