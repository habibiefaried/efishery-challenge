# Deskripsi
App simpel untuk config/env management

# Penggunaan
Aplikasi ini sudah dihubungkan secara live di host `kampustsl.id` port `1337`. Segala perubahan pada branch `main` akan otomatis di update ke server tersebut. Berikut adalah contoh penggunaan

Operator: 
* SET <key> <value> : Set key dengan value
* GET <key> : Mengambil value dari sebuah key
* LIST: Mengambil semua key yang ada pada sistem
* UNSET <key> : Menghapus key
* IMPORT <tipe> <url> : Import key value file yang tersimpan pada format .env/yaml/json di URL

```
$ nc -vvv kampustsl.id 1337
Connection to kampustsl.id 1337 port [tcp/*] succeeded!
SET kunci jawaban
Penulisan key berhasil

GET kunci
jawaban

LIST
APP_ID
APP_SECRET
kunci

GET APP_ID
1234567

UNSET APP_ID
key APP_ID berhasil dihapus

UNSET APP_SECRET
key APP_SECRET berhasil dihapus

LIST
kunci

IMPORT YAML https://raw.githubusercontent.com/habibiefaried/efishery-challenge/main/yaml.sample
import berhasil

LIST
APP_ID
APP_SECRET
kunci

GET APP_SECRET
abcdef
```
