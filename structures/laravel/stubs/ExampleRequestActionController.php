<?php

// STUB STRUKTUR — bukan kode aktif.
// Core 5 Step-up Security — rujuk: core_stepup_security_api.md §1.4 (lifecycle).
//
// Urutan yang harus didukung controller ini (+SecurityPolicyController & challenge di Core 2):
//   1. POST /security-policies          daftar target (module+action → riskLevel+validationType)
//   2. POST /request-actions            buat intent:
//        - wajib ada SecurityPolicyTarget aktif untuk module+action (404 bila tidak)
//        - status: validationType CONFIRM → PENDING, selain itu → CHALLENGE_REQUIRED
//        - payload_hash = sha256(json), expires_at = +30 menit, 201
//   3. (challenge OTP via Core 2, bila terkait requestActionId → status CHALLENGE_VERIFIED)
//   4. POST /request-actions/{id}/consume → EXECUTED + consumed_at (409 bila sudah EXECUTED)
//   5. CRUD list/detail/update/delete — lihat core_stepup_security_api.md §3 (409 finalized, 404 not found)
//
// Proteksi route: auth.jwt + role.strict:SUPERADMIN (TANPA step-up — request-action adalah wadah intent).

namespace App\Http\Controllers\Api\V1\Dashboard;

use App\Http\Controllers\Controller;
use App\Support\ApiResponse;
use Illuminate\Http\Request;

class ExampleRequestActionController extends Controller
{
    public function store(Request $request)
    {
        // TODO: validasi module/action/payload → cari target aktif → buat row intent → 201.
    }

    public function markConsumed(string $id)
    {
        // TODO: 400 Invalid ID / 404 not found / 409 already executed → set EXECUTED, 200.
    }
}
