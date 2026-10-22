<?php

// STUB STRUKTUR — bukan kode aktif.
// Core 3+4 response — rujuk: core_auth_api.md §1.1 (3 bentuk) & core_crud_api.md §4.1 (envelope list).
//
// TIGA bentuk WAJIB, jangan dicampur:
//   success()          → 200/201  {"status": true,  "message": "...", "data": ...}
//   error()            → 400/401/404/409/500  {"status": false, "message": "...", "data": ""}
//   middlewareError()  → 401/403 middleware  {"error": "...", "code": "AUTH_*|AUTHZ_*"}
//
// Envelope list (SEMUA endpoint GET daftar, termasuk /names):
//   {"contents": [...], "page": 1, "size": 10, "totalData": 123}

namespace App\Support;

use Illuminate\Http\JsonResponse;

class ApiResponse
{
    public static function success(string $message, mixed $data = null, int $status = 200): JsonResponse
    {
        // TODO: response()->json(['status' => true, 'message' => $message, 'data' => $data], $status);
    }

    public static function error(string $message, int $status, ?string $code = null): JsonResponse
    {
        // TODO: {"status": false, "message": $message, "data": ""} + "code" bila $code diisi.
    }

    public static function middlewareError(string $message, int $status, ?string $code = null): JsonResponse
    {
        // TODO: {"error": $message, "code": $code}.
    }
}
