package dto

type SQLable interface {
	BriefBusiness | BriefProduct | BriefUser
}
