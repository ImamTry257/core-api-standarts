<?php

// STUB STRUKTUR — bukan kode aktif.
// Core 3+4 CRUD — rujuk: core_crud_api.md §1.2 (pola 6 operasi) & §4 (pagination).
//
// | Operasi | Method & Path              | Status | message baku                              |
// |---------|----------------------------|--------|-------------------------------------------|
// | Create  | POST   /example            | 201    | "Example created successfully"             |
// | List    | GET    /example            | 200    | "Successfully retrieved all examples"      |
// | Detail  | GET    /example/{id}       | 200    | "Successfully retrieved example"           |
// | Update  | PUT    /example/{id}       | 200    | "Example updated successfully"             |
// | Delete  | DELETE /example/{id}       | 200    | "Example deleted successfully", data null  |
// | Lookup  | GET    /example/names      | 200    | "Successfully retrieved all examples name" |
//
// Kode error: 400 "Invalid ID" / "<field> is required" / duplikat,
//             404 "example tidak ditemukan", 403 permission, 500 hanya gangguan server.
// List: query search,page,size (default 1/10) + filter domain → envelope {contents,page,size,totalData}.
// Soft delete bila data master; validasi field wajib & relasi SEBELUM insert (bukan try/catch 500).

namespace App\Http\Controllers\Api\V1\Dashboard;

use App\Http\Controllers\Controller;
use App\Support\ApiResponse;
use Illuminate\Http\Request;

class ExampleCrudController extends Controller
{
    public function index(Request $request)
    {
        // TODO: parse search/page/size (max(1,...)), count total, ORDER BY deterministik, envelope.
    }

    public function show(string $id)
    {
        // TODO: ctype_digit → 400 "Invalid ID"; tak ada → 404.
    }

    public function store(Request $request)
    {
        // TODO: validasi wajib/relasi/duplikat → 400; sukses 201.
    }

    public function update(Request $request, string $id)
    {
        // TODO: id → 400/404; validasi → 400; sukses 200.
    }

    public function destroy(string $id)
    {
        // TODO: id → 400/404; soft delete; sukses 200 data null.
    }

    public function getAllNames(Request $request)
    {
        // TODO: lookup terpaginasi ringan {id,name} — core_crud_api.md §3.
    }
}
