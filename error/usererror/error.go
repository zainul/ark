package usererror

var OwnErrors map[string]string

const (
	ServiceName   = "01E"
	ReloadService = "RC"
)

const (
	NotFoundCode             = ServiceName + "0"
	InvalidToAccountNumber   = ServiceName + "1"
	InvalidFromAccountNumber = ServiceName + "2"
	InsufficientBalance      = ServiceName + "3"
	FailedToTransfer         = ServiceName + "4"
	MissingRequiredParam     = ServiceName + "5"
	FirstNameIsRequired      = ServiceName + "6"
	LastNameIsRequired       = ServiceName + "7"
	EmailIsRequired          = ServiceName + "8"
	EmailFormatNotCorrect    = ServiceName + "9"
	PasswordIsRequired       = ServiceName + "10"
	AccountTypeIsRequired    = ServiceName + "11"
	PhoneNumberIsRequired    = ServiceName + "12"
	UserAlreadyExist         = ServiceName + "13"
	InvalidAccountType       = ServiceName + "14"
	InternalServerError      = ServiceName + "15"
	ActivationCodeRequired   = ServiceName + "16"
	ActivationCodeNotFound   = ServiceName + "17"
	AccountAlreadyActivate   = ServiceName + "18"
	UserNotFound             = ServiceName + "19"
	AccountNotActivatedYet   = ServiceName + "20"

	ProductCodeRequired   = ReloadService + "1"
	AmountRequired        = ReloadService + "2"
	FulfillNumberRequired = ReloadService + "3"
	RefIDRequired         = ReloadService + "4"
	FulfillmentFailed     = ReloadService + "5"
)

func init() {
	errs := make(map[string]string)

	errs[NotFoundCode] = "Unexpected error"
	errs[MissingRequiredParam] = "Missing Some Required Param"
	errs[InvalidFromAccountNumber] = "Invalid from account number"
	errs[InvalidToAccountNumber] = "Invalid to account number"
	errs[InsufficientBalance] = "Insufficient balance"
	errs[FailedToTransfer] = "Failed to transfer from %s to %s"
	errs[FirstNameIsRequired] = "First name is required"
	errs[LastNameIsRequired] = "Last name is required"
	errs[EmailIsRequired] = "Email is required"
	errs[EmailFormatNotCorrect] = "Email format not valid"
	errs[PasswordIsRequired] = "Password is required"
	errs[AccountTypeIsRequired] = "Account type is required"
	errs[PhoneNumberIsRequired] = "Phone number is required"
	errs[UserAlreadyExist] = "User already exist"
	errs[InvalidAccountType] = "Invalid account type"
	errs[InternalServerError] = "Internal server error, please try again in few minute"
	errs[ActivationCodeRequired] = "Activation code required"
	errs[ActivationCodeNotFound] = "Activation code not found"
	errs[AccountAlreadyActivate] = "Account already activate"
	errs[UserNotFound] = "User Not found or not register in our system"
	errs[AccountNotActivatedYet] = "Account is not activated yet"

	errs[ProductCodeRequired] = "Product code is required"
	errs[AmountRequired] = "Amount is required"
	errs[FulfillNumberRequired] = "Fulfillment number is required"
	errs[RefIDRequired] = " Reference ID is required"
	errs[FulfillmentFailed] = "Fulfillment is failed"

	OwnErrors = errs
}

func GetErrors() map[string]string {
	return OwnErrors
}
func GetCode(err string) string {
	if val, ok := OwnErrors[err]; ok {
		return val
	}

	return OwnErrors[NotFoundCode]
}
