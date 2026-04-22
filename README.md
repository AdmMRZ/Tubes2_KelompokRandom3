# DOM Traversal & Search Engine

Tugas Besar 2 IF2211 Strategi Algoritma - Aplikasi Web Pencarian Elemen pada Struktur Pohon DOM menggunakan algoritma **Breadth-First Search (BFS)** dan **Depth-First Search (DFS)**. Aplikasi ini dibangun secara _Fullstack_ dengan Golang murni sebagai _Backend_ (Parser & Traversal Engine) dan Vanilla Web (HTML/JS/CSS) sebagai _Frontend_ UI/UX Premium.

## Penjelasan Algoritma yang Diimplementasikan

Pemrosesan dan pencarian elemen HTML (DOM *subtree*) dalam aplikasi ini mengimplementasikan dua metode menelusuri Graf/Pohon demi menentukan jalur dan kecocokan terhadap sebuah CSS Selector. Keduanya mendukung struktur *Combinator* CSS kompleks seperti `Child (>)`, `Descendant ( )`, `Adjacent Sibling (+)`, dan `General Sibling (~)`.

1. **Breadth-First Search (BFS)**
   Algoritma BFS bekerja secara **melebar**. Pencarian dimulai dari akar (*root*) dokumen HTML (seperti tag `<html>`), kemudian bergerak memeriksa semua _node/tag_ yang berada pada level kedalaman yang sama sebelum turun ke level/tingkat bersarang berikutnya. 
   - **Implementasi:** Pada _backend_ (Go), BFS diatur menggunakan struktur data tipe **Antrean (Queue)**. Elemen yang pertama kali masuk ke antrean akan diproses terlebih dulu (FIFO).
   - **Karakteristik:** Sangat optimal untuk menemukan hasil (khususnya jika dibatasi `Top N`) yang terdekat dengan akar dokumen. Jalur traversal akan terlihat memeriksa *siblings* secara horizontal.

2. **Depth-First Search (DFS)**
   Algoritma DFS menelusuri pohon secara **mendalam**. Algoritma akan memilih satu dahan turunan dan terus menyelam hingga mencapai *leaf* (anak terdalam/terakhir) sebelum melakukan langkah mundur (*backtracking*) untuk melahap *node* anak lainnya pada saudara tetangga.
   - **Implementasi:** Pada _backend_ (Go), algoritma ini diimplementasikan menggunakan **Stack Eksplisit (LIFO)**. Anak-anak dari suatu *node* ditekan (*push*) ke dalam tumpukan dalam susunan memutar (*reverse order* dari kanan ke kiri) demi menjamin penelusuran selalu terjun memprioritaskan "Anak Cabang Kiri" terlebih dulu.
   - **Karakteristik:** Sangat ampuh menembus komponen struktur bersarang dalam.

---

## Requirement Program & Instalasi 

**Kebutuhan Eksekusi:**
- **Docker Desktop / Docker Engine** (disarankan)
- **Docker Compose**

---

## Langkah-langkah Kompilasi & *Build*

### Menggunakan Docker Compose 
Metode ini adalah cara paling instan dan minim kendala, berkat arsitektur _Multi-stage build_ Golang dan _Reverse Proxy_ Nginx di sisi Frontend yang telah dirakit sempurna di konfigurasi kami.

1. Buka Terminal / Command Prompt pada *root directory*.
2. Nyalakan mesin Docker dengan eksekusi perintah:
   ```bash
   docker-compose up -d --build
   ```
3. Tunggu hingga proses kompilasi _images_ usai (*Backend* di- _build_ dalam OS _Alpine_ kecil).
4. Selesai! Buka browser Anda dan kunjungi halaman utama pada alamat:
   **`http://localhost:3000`**
---

## _Authors_

1. **Muhammad Adam Mirza** – 18223015
2. **Harfhan Ikhtiar Ahmad Ridzky** – 18223123
