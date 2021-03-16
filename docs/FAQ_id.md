# FAQ

[![Eng](https://github.com/gosquared/flags/blob/master/flags/flags/flat/24/United-Kingdom.png?raw=true)](FAQ_en.md)
[![Ind](https://github.com/gosquared/flags/blob/master/flags/flags/flat/24/Indonesia.png?raw=true)](FAQ_id.md)

[Tutorial](Tutorial_id.md) | [Rule Engine](RuleEngine_id.md) | [GRL](GRL_id.md) | [GRL JSON](GRL_JSON_id.md) | [RETE Algorithm](RETE_id.md) | [Functions](Function_id.md) | [FAQ](FAQ_id.md) | [Benchmark](Benchmarking_id.md)

---

## 1. Grule Panik pada Siklus Maksimum

**Pertanyaan**: Saya mendapat pesan panik ini saat Grule engine dijalankan.

```Shell
panic: GruleEngine successfully selected rule candidate for execution after 5000 cycles, this could possibly caused by rule entry(s) that keep added into execution pool but when executed it does not change any data in context. Please evaluate your rule entries "When" and "Then" scope. You can adjust the maximum cycle using GruleEngine.MaxCycle variable.
```

**Jawaban**: Error ini mengindikasikan masalah yang ada pada __rule__ yang anda buat
dan dievaluasi oleh engine. Grule akan terus menjalankan __jaringan RETE__ didalam
__working memory__ hingga tidak ada lagi tindakan yang bisa dilakukan dalam __conflict set__,
yang mana disebut sebagai kondisi terminasi yang natural/normal. Jika dalam kumpulan rule tidak pernah 
mengizinkan __jaringan RETE__ untuk mencapai kondisi akhir ini, maka eksekusi akan terus terjadi selamanya.
Secara __default__ konfigurasi untuk `GruleEngine.MaxCycle` adalah `5000`, dimana nilai ini untuk
melindungi eksekusi tidak berujung karena tidak pernah mencapai kondisi terminasi.

Anda dapat meningkatkan nilai ini jika menurut anda sistem __rule__ anda membutuhkan siklus lebih
banyak untuk bisa mencapai terminasi, tapi jika anda merasa ragu jika menambah nilai ini 
akan menghentikan pesan panik, maka kemungkinan anda memiliki kumpulan __rule__ yang tidak 
punya kondisi akhir.

Asumsikan __fact__ berikut ini:

```go
type Fact struct {
   Payment int
   Cashback int
}
```

Dan __rule-rule__ seperti berikut:

```Shell
rule GiveCashback "Give cashback if payment is above 100" {
    When 
         F.Payment > 100
    Then
         F.Cashback = 10;
}

rule LogCashback "Emit log if cashback is given" {
    When 
         F.Cashback > 5
    Then
         Log("Cashback given :" + F.Cashback);
}
```

Kita akan menjalankan __rule__ tadi pada sebuah turunan fakta ...

```go
&Fact {
     Payment: 500,
}
```

... eksekusi ini tidak akan mencapai terminasi. 

```
Siklus 1: Menjalankan "GiveCashback" .... karena F.Payment > 100 adalah kondisi yang valid
Siklus 2: Menjalankan "GiveCashback" .... karena F.Payment > 100 adalah kondisi yang valid
Siklus 3: Menjalankan "GiveCashback" .... karena F.Payment > 100 adalah kondisi yang valid
...
Siklus 5000: Menjalankan "GiveCashback" .... karena F.Payment > 100 adalah tetap kondisi yang valid
panik
```

Grule menjalankan __rule__ yang sama lagi dan lagi karena kondisi pada **WHEN**
terus menerus memberikan hasil yang valid.

Satu cara untuk memecahkan masalah ini adalah merubah __rule__ "GiveCashback" menjadi seperti:

```Shell
rule GiveCashback "Give cashback if payment is above 100" {
    When 
         F.Payment > 100 &&
         F.Cashback == 0
    Then
         F.Cashback = 10;
}
```

Dengan demikian, __rule__ `GiveCashback` turut memperhitungkan perubahan nilai yang terjadi.
Yang tadinya nilai variabel `Cashback` adalah 0, dikarenakan perubahan yang terjadi membuat
evaluasi ini menjadi tidak valid lagi pada siklus berikutnya, hingga menyebabkan evaluasi
pindah pada __rule__ yang lain hingga selesai.

Cara diatas adalah cara untuk mengotrol eksekusi __rule__ secara "natural" hingga setelah
serangkaian siklus engine akan berhenti secara normal karena tidak ada lagi __rule__ yang bisa 
di eksekusi. Namun, ada kalanya anda tidak bisa menghetikan eksekusi dengan cara seperti ini.
Alternatif lain adalah untuk mengubah rule menjadi seperti berikut:

```Shell
rule GiveCashback "Give cashback if payment is above 100" {
    When 
         F.Payment > 100
    Then
         F.Cashback = 10;
         Retract("GiveCashback");
}
```

Fungsi `Retract` akan menghilangkan sementara __rule__ "GiveCashback" dari dalam __knowledge base__
hingga siklus berakhir. Karena __rule__ ini tidak lagi tersedia dalam siklus ini, maka __rule__
tersebut tidak dapat lagi dievaluasi hingga akhir. Perlu anda ketahui, bahwa __rule__ akan hilang
sementara saja setelah `Retract` dipanggil. Pada siklus-siklus setelahnya, rule tersebut akan tersedia
kembali.

---

## 2. Menyimpan Rule kedalam Database

**Pertanyaan**: Apakah ada rencana untuk mengintegrasikan Grule dengan penyimpanan di Database?

**Jawaban**: Tidak. Walaupun ini adalah ide yang baik untuk menyimpan __rule__ kedalam
database, Grule tidak akan membuat sebuah koneksi kepada sebuah database untuk secara otomatis menyimpan
dan mengambil __rule__. Anda dapat dengan mudah membuat mekanisme ini sendiri menggunakan
cara-cara yang sudah ada: menggunakan *Reader*, *File*, *Byte Array*, *String* dan *Git*.
Sebuah string dapat dengan mudah dimasukan dan baca dari database, untuk menyimpan/mengambil
__rule__ dan memasukannya kedalam __Knowledgebase__ dalam Grule.

Kami tidak ingin membuat keterikatan pada database apapun.

---

## 3. Jumlah maksimal rule dalam satu Knowledgebase

**Pertanyaan**: Berapa banyak __rule entry__ uang bisa dimasukan kedalam __knowledgebase__?

**Jawaban**: Anda dapat menambahkan berapapun __rule__ yang anda perlukan, selama minimal ada 1 rule dalam 
sebuah __knowledgebase__

---

## 4. Mengetahui __rule__ apa saja yang valid untuk sebuah __fact__

**Pertanyaan**: Bagaimana saya mengetahui __rule__ - __rule__ mana saja yang valid terhadap sebuah __fact__?

**Jawaban**: Anda dapat menggunakan fungsi `engine.FetchMatchingRule`. Silahkan merujuk pada
[Matching Rules Doc](MatchingRules_id.md) untuk informasi lebih lengkap.

---

## 5. Use-Case untuk Rule Engine

**Pertanyaan**: Saya sudah membaca-baca tentang __rule engine__, tapi apa sebenarnya keuntungan yang didapat? Berikan kami contoh Use-Case.

**Jawaban**: Berikut ini adalah contoh situasi yang sebaiknya diselesaikan menggunakan solusi __rule-engine__ menurut hemat kami.

1. Sebuah sistem pakar yang harus mengevaluasi fakta-fakta guna memberikan sebuah kesimpulan yang nyata.
   Jika tidak menggunakan model RETE dan __rule-engine__, seorang developer akan membuat kode program
   yang berisi `if`/`else` yang beranak pinak dan permutasi terhadap kombinasi kondisi-kondisi yang ada
   membuat manajemen kode menjadi mustahil. Pendekatan __rule engine__ menggunakan tabel mungkin bisa
   memecahkan masalah, namun pendekatan ini menjadikan solusinya kaku dan tidak begitu mudah di
   buat kode program nya. Sistem seperti Grule ini memudahkan anda untuk mendeskripsikan peraturan terhadap
   data yang dipergunakan dalam sistem, dan melepaskan anda dari kebutuhan untuk mengimplementasi bagaimana
   sebenarnya evaluasi logika peraturan itu terlaksana, menyebunyikan banyak kompleksitas dari anda.

2. Sistem pemberian Rating atau Skor. For example, a bank system may want to create a "score" for
   each customer based on the customer's transaction records (facts).  We could
   see their score change based on how often they interact with the bank, how
   much money they transfer in and out, how quickly they pay their bills, how
   much interest they accrue earn for themselves or for the bank, and so on. A
   rule engine is provided by the developer and the specification of the facts
   and rules can then be supplied by subject matter experts in the bank's
   customer business. Decoupling these different teams puts the responsbilities
   where they should be.

3. Computer games. Player status, rewards, penalties, damage, scores and
   probability systems are many different examples of where rule play a
   significant part of nearly all computer games. These rules can interact in
   very complex ways, often times in ways that the developer didn't imagine.
   Coding these dynamic situations in a scripting language (e.g. LUA) can get
   quite complex, and a rule engine can help simplify the work tremendously.

4. Classification systems. This is actually a generalization of the rating
   system described above.  Using a rule engine, we can classify things such as
   credit eligibility, bio chemical identification, risk assessment for
   insurance products, potential security threats, and many more.

5. Advice/Suggestion system. A "rule" is simply another kind of data, which
   makes it a prime candidate for definition by another program.  This program
   can be another expert system or artificial intelligence.  Rules can be
   manipulated by another system in order to deal with new types of facts or
   newly discovered information about the domain which the rule set is intending
   to model.

There are so many other use-cases that would benefit from the use of
Rule-Engine. The above cases represent only a small number of the potential. 

However it is important to state that a Rule-Engine not a silver bullet, of
course.  Many alternatives exist to solve "knowledge" problems in software and
those should be employed when they are most appropriate. One would not employ a
rule engine where a simple `if` / `else` branch would suffice, for instance.

Theres's someting else to note: some rule engine implementations are extremely
expensive yet many businesses gain so much value from them that the cost of
running them is easily offset by that value.  For even moderately complex use
cases, the benefit of a strong rule engine that can decouple teams and tame
business complexity seems to be quite clear.

---

## 6. Logging

**Pertanyaan**: Grule's logs are extremely verbose.  Can I turn off Grule's logger?

**Jawaban**: Yes. You can reduce (or completely stop) Grule's logging by increasing it's log level.

```go
import (
    "github.com/hyperjumptech/grule-rule-engine/logger"
    "github.com/sirupsen/logrus"
)
...
...
logger.SetLogLevel(logrus.PanicLevel)
```

This will set Grule's log to `Panic` level, where it will only emits log when it panicked.

Of course, modifying the log level reduces your ability to debug the system so
we suggest that a higher log level setting only be instituted in production
environments.
