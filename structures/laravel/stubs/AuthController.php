<?php

// STUB STRUKTUR — bukan kode aktif.
// Core 1 Auth — rujuk: core_auth_api.md §2.1–§2.9 (request, pesan, kode status persis kontrak).
//
// Pemetaan route (routes/api.php):
//   Group /auth TERBUKA      : login, verify-code, resend-code,
//                              activation/check, activation/confirm, forgot/password, reset/password
//   Group /auth + auth.jwt   : refresh-token  (hanya sesi FULL_AUTH — core_auth_api.md §1.5)
//   Group users + auth.jwt   : users/{id}/resend-activation → permission users,resend_activation
//                              (§2.6 — path blueprint /auth/activation/resend TIDAK ADA di backend)
//
// Layer: AuthController → AuthService (state sesi/OTP/activation) → JwtService/OtpService → Jobs (email).
// Semua method di bawah = penanda "isi di sini" — isi logika saat implementasi, jangan dijalankan apa adanya.

namespace App\Http\Controllers\Api\V1\Dashboard;

use App\Http\Controllers\Controller;
use App\Support\ApiResponse;

class AuthController extends Controller
{
    // §2.1 POST /auth/login — {email, password(sha256 hex 64)}; sukses → sesi PARTIAL_AUTH;
    // gagal → pesan generik "email or password is wrong".
    public function login()
    {
        // TODO: isi sesuai core_auth_api.md §2.1 → ApiResponse::success/error.
    }

    // §2.2 POST /auth/verify-code — kode OTP → sesi FULL_AUTH (core_auth_api.md §1.5).
    public function verifyCode()
    {
        // TODO: isi sesuai core_auth_api.md §2.2.
    }

    // §2.3 POST /auth/resend-code — kirim ulang OTP login (batas & cooldown di kontrak).
    public function resendCode()
    {
        // TODO: isi sesuai core_auth_api.md §2.3.
    }

    // §2.4 POST /auth/activation/check — cek status aktivasi akun.
    public function activationCheck()
    {
        // TODO: isi sesuai core_auth_api.md §2.4.
    }

    // §2.5 POST /auth/activation/confirm — konfirmasi kode aktivasi.
    public function activationConfirm()
    {
        // TODO: isi sesuai core_auth_api.md §2.5.
    }

    // §2.6 POST /users/{id}/resend-activation — di group USERS (bukan /auth):
    // FULL_AUTH + permission users,resend_activation; activation JWT baru 30 menit.
    // Letak method riil: UserController (stub ini hanya penunjuk urutan Core 1).
    public function resendActivation(string $id)
    {
        // TODO: isi sesuai core_auth_api.md §2.6 (validasi id digit, 400/401/403/404).
    }

    // §2.7 POST /auth/forgot/password — mulai alur reset via email.
    public function forgotPassword()
    {
        // TODO: isi sesuai core_auth_api.md §2.7.
    }

    // §2.8 POST /auth/reset/password — tukar token reset dengan password baru.
    public function resetPassword()
    {
        // TODO: isi sesuai core_auth_api.md §2.8.
    }

    // §2.9 POST /auth/refresh-token — refresh token sekali pakai (rotasi) → access token baru;
    // terproteksi auth.jwt (hanya FULL_AUTH).
    public function refreshToken()
    {
        // TODO: isi sesuai core_auth_api.md §2.9.
    }
}
