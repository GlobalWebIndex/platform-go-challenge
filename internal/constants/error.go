package constants

var (
	ErrInvalidUser                    = "[Err] user doesn't authorized:"
	ErrNotFoundBusiness               = "[Err] business doesn't exist:"
	ErrNotAllowCreateBusiness         = "[Err] can't create business:"
	ErrNotAllowCreateAPIKey           = "[Err] can't generate new api key:"
	ErrNotAllowCreateStripeCustomerId = "[Err] can't create stripe customerId:"
	ErrInvalidStripSessionID          = "[Err] failed to register stripe sessionID:"
	ErrInvalidRequest                 = "[Err] invalid request:"
	ErrInvalidBusinessInfo            = "[Err] invalid business information:"
)
