package dto

// APIResponse — bentuk response sukses & gagal bisnis {status, message, data}
// (core_auth_api.md §1.1). Error middleware memakai bentuk {error, code} — lihat middleware.
type APIResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// APIListResponse — envelope pagination seragam {contents, page, size, totalData}
// (core_crud_api.md §4.1). FE menghitung total halaman: ceil(totalData/size).
type APIListResponse struct {
	Contents  interface{} `json:"contents"`
	Page      int         `json:"page"`
	Size      int         `json:"size"`
	TotalData int         `json:"totalData"`
}

func Success(message string, data interface{}) APIResponse {
	return APIResponse{Status: true, Message: message, Data: data}
}

func Fail(message string) APIResponse {
	return APIResponse{Status: false, Message: message, Data: ""}
}
