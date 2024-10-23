package main

import (
	"context"
	"fmt"
	"time"
)

// Struct untuk produk
type Product struct {
	ID    int
	Name  string
	Price int
}

// Struct untuk user session
type Session struct {
	Username string
	Active   bool
}

// Struct untuk keranjang belanja
type Cart struct {
	Items []Product
}

// Data produk tersedia
var products = []Product{
	{ID: 1, Name: "Laptop", Price: 15000000},
	{ID: 2, Name: "Smartphone", Price: 5000000},
	{ID: 3, Name: "Headset", Price: 500000},
}

// Fungsi login
func login(ctx context.Context, session *Session) {
	var username string
	fmt.Print("Masukkan username: ")
	fmt.Scan(&username)
	session.Username = username
	session.Active = true

	// Membuat timeout untuk sesi login (20 detik)
	go func() {
		select {
		case <-time.After(20 * time.Second):
			session.Active = false
			fmt.Println("\nSesi login habis. Silakan login kembali.")
		case <-ctx.Done():
			return
		}
	}()
}

// Fungsi menampilkan daftar produk
func displayProducts() {
	fmt.Println("Daftar Produk:")
	for _, p := range products {
		fmt.Printf("%d. %s - Rp%d\n", p.ID, p.Name, p.Price)
	}
}

// Fungsi menambahkan produk ke keranjang
func addToCart(cart *Cart) {
	var productID int
	displayProducts()
	fmt.Print("Pilih ID produk untuk ditambahkan ke keranjang: ")
	fmt.Scan(&productID)

	for _, product := range products {
		if product.ID == productID {
			cart.Items = append(cart.Items, product)
			fmt.Printf("Produk %s berhasil ditambahkan ke keranjang.\n", product.Name)
			return
		}
	}
	fmt.Println("ID produk tidak valid.")
}

// Fungsi melihat keranjang
func viewCart(cart *Cart) {
	if len(cart.Items) == 0 {
		fmt.Println("Keranjang kosong.")
		return
	}

	fmt.Println("Isi Keranjang:")
	for i, item := range cart.Items {
		fmt.Printf("%d. %s - Rp%d\n", i+1, item.Name, item.Price)
	}
}

// Fungsi checkout dan pembayaran
func checkout(cart *Cart) {
	if len(cart.Items) == 0 {
		fmt.Println("Keranjang kosong, tidak bisa checkout.")
		return
	}

	var total int
	fmt.Println("Checkout:")
	for _, item := range cart.Items {
		total += item.Price
		fmt.Printf("%s - Rp%d\n", item.Name, item.Price)
	}
	fmt.Printf("Total pembayaran: Rp%d\n", total)

	var payment int
	fmt.Print("Masukkan jumlah pembayaran: ")
	fmt.Scan(&payment)

	if payment < total {
		fmt.Println("Uang tidak cukup. Checkout dibatalkan.")
	} else {
		fmt.Println("Pembayaran berhasil! Barang segera dikirim.")
		cart.Items = nil
	}
}

// Fungsi utama
func main() {
	var session Session
	var cart Cart
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		// Mengecek apakah sesi login aktif
		if !session.Active {
			fmt.Println("Silakan login terlebih dahulu.")
			login(ctx, &session)
		}

		// Menampilkan menu
		fmt.Println("\nMenu:")
		fmt.Println("1. Lihat Produk")
		fmt.Println("2. Tambah Produk ke Keranjang")
		fmt.Println("3. Lihat Keranjang")
		fmt.Println("4. Checkout")
		fmt.Println("5. Keluar")
		fmt.Print("Pilih menu: ")

		var choice int
		fmt.Scan(&choice)

		// Mengecek apakah sesi login aktif sebelum memproses pilihan
		if !session.Active {
			fmt.Println("Sesi login habis. Silakan login kembali.")
			continue
		}

		switch choice {
		case 1:
			displayProducts()
		case 2:
			addToCart(&cart)
		case 3:
			viewCart(&cart)
		case 4:
			checkout(&cart)
		case 5:
			fmt.Println("Terima kasih telah berbelanja!")
			return
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}
