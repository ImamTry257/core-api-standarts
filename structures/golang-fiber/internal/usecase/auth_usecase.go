package usecase

// AuthUseCase — Core 1: lapisan usecase auth (STUB penunjuk, bukan implementasi).
// Pola interface + impl mengikuti learning-auth-services/pkg/auth/usecase
// (dipakai bersama repository gorm: pkg/auth/repository).
// Sembilan metode = core_auth_api.md §2.1–§2.9; isi logikanya di impl saat implementasi.
type AuthUseCase interface {
	// Login — §2.1: verifikasi kredensial → sesi PARTIAL_AUTH.
	Login(email, password string) (interface{}, error)
	// VerifyOtp — §2.2: kode OTP → sesi FULL_AUTH.
	VerifyOtp(email, otpCode string) (interface{}, error)
	// ResendOtp — §2.3: kirim ulang OTP login.
	ResendOtp(email string) (interface{}, error)
	// CheckActivation — §2.4: cek status aktivasi.
	CheckActivation(email string) (interface{}, error)
	// ConfirmActivation — §2.5: konfirmasi kode aktivasi.
	ConfirmActivation(activationCode string) (interface{}, error)
	// ResendActivation — §2.6: kirim ulang kode aktivasi (panggilan dari domain users).
	ResendActivation(userID int) (interface{}, error)
	// ForgotPassword — §2.7: mulai alur reset password via email.
	ForgotPassword(email string) (interface{}, error)
	// ResetPassword — §2.8: tukar token reset dengan password baru.
	ResetPassword(resetToken, newPassword string) (interface{}, error)
	// Refresh — §2.9: rotasi refresh token sekali pakai → access token baru.
	Refresh(refreshToken string) (interface{}, error)
}

type authUseCase struct {
	// repo repository.Repository — aktifkan (pola learning-auth-services).
}

func NewAuthUseCase() AuthUseCase {
	return &authUseCase{}
}

func stub() (interface{}, error) {
	// TODO: implementasi per metode mengikuti core_auth_api.md §2.1–§2.9 —
	// repository → verifikasi → terbitkan token (pola utils.SetJWT) → audit.
	return nil, nil
}

func (u *authUseCase) Login(email, password string) (interface{}, error) { return stub() }

func (u *authUseCase) VerifyOtp(email, otpCode string) (interface{}, error) { return stub() }

func (u *authUseCase) ResendOtp(email string) (interface{}, error) { return stub() }

func (u *authUseCase) CheckActivation(email string) (interface{}, error) { return stub() }

func (u *authUseCase) ConfirmActivation(activationCode string) (interface{}, error) {
	return stub()
}

func (u *authUseCase) ResendActivation(userID int) (interface{}, error) { return stub() }

func (u *authUseCase) ForgotPassword(email string) (interface{}, error) { return stub() }

func (u *authUseCase) ResetPassword(resetToken, newPassword string) (interface{}, error) {
	return stub()
}

func (u *authUseCase) Refresh(refreshToken string) (interface{}, error) { return stub() }
