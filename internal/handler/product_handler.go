package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"product-api/internal/model"
	"product-api/internal/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// Create godoc
// @Summary      Buat produk baru
// @Description  Menambahkan produk baru ke sistem
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      model.ProductCreateRequest  true  "Data produk"
// @Success      201      {object}  model.Product
// @Failure      400      {object}  map[string]string
// @Router       /products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.ProductCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "data tidak valid")
		return
	}

	created, err := h.service.CreateProduct(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// GetList godoc
// @Summary      Ambil daftar produk
// @Description  Menampilkan semua produk yang ada
// @Tags         products
// @Produce      json
// @Success      200  {array}  model.Product
// @Router       /products [get]
func (h *ProductHandler) GetList(w http.ResponseWriter, r *http.Request) {
	products := h.service.GetAllProducts()
	writeJSON(w, http.StatusOK, products)
}

// GetDetail godoc
// @Summary      Ambil detail produk
// @Description  Menampilkan detail produk berdasarkan ID
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  model.Product
// @Failure      404  {object}  map[string]string
// @Router       /products/{id} [get]
func (h *ProductHandler) GetDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}

	product, err := h.service.GetProductByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, product)
}

// Update godoc
// @Summary      Update produk
// @Description  Mengubah data produk berdasarkan ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      int            true  "Product ID"
// @Param        product  body      model.ProductUpdateRequest  true  "Data produk baru"
// @Success      200      {object}  model.Product
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}

	var req model.ProductUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "data tidak valid")
		return
	}

	updated, err := h.service.UpdateProduct(id, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// Delete godoc
// @Summary      Hapus produk
// @Description  Menghapus produk berdasarkan ID
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}

	if err := h.service.DeleteProduct(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "produk berhasil dihapus"})
}