<?php

// STUB STRUKTUR — bukan kode aktif.
// Core 2 RBAC — rujuk: core_role_permission_api.md §5 (error middleware).
//
// Pemakaian di routes/api.php:
//   Route::get('/users', [...])->middleware('permission:users,read');
//   Route::post('/users', [...])->middleware('permission:users,create');
//
// Aturan (prinsip: core_principles.md §3):
//   1. Role SUPERADMIN (atau setara) SELALU bypass.
//   2. Selain itu: cocokkan "<resource>.<action>" (huruf kecil) dengan permission role aktif.
//   3. Default DENY — tidak cocok → middlewareError('permission denied', 403, 'AUTHZ_PERMISSION_DENIED').
//   4. Semua error memakai bentuk {"error", "code"} — BUKAN {"status","message"}.

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;

class PermissionMiddleware
{
    public function handle(Request $request, Closure $next, string $resource, string $action)
    {
        // TODO: bypass SUPERADMIN → load role.permissions → cek "<resource>.<action>" → deny 403.

        return $next($request);
    }
}
