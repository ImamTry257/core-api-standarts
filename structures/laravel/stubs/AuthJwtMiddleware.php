<?php

// STUB STRUKTUR — bukan kode aktif.
// Core 1 Auth — rujuk: core_auth_api.md §1.3 (sesi FULL_AUTH) & §3 (error middleware).
//
// Dipasang di routes/api.php sebagai group middleware:
//   Route::middleware(['auth.jwt'])->group(...)
//
// Alur:
//   1. Baca header Authorization: Bearer <access-token>.
//   2. Verifikasi token (JwtService::parse) → user_id, role, exp.
//   3. WAJIB menolak token Partial/OTP — endpoint terproteksi hanya menerima FULL_AUTH.
//   4. Gagal → middlewareError() bentuk {"error": "...", "code": "AUTH_*"} (401):
//        AUTH_MISSING_HEADER | AUTH_INVALID_HEADER_FORMAT |
//        AUTH_INVALID_OR_EXPIRED_TOKEN | AUTH_USER_NOT_FOUND | AUTH_USER_NOT_VERIFIED
//   5. Sukses → simpan user_id/role ke $request->attributes, lanjut $next.

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;

class AuthJwtMiddleware
{
    public function handle(Request $request, Closure $next)
    {
        // TODO: implementasi 5 langkah di atas.

        return $next($request);
    }
}
