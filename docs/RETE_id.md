# Grule's RETE Algorithm

[![Eng](https://github.com/gosquared/flags/blob/master/flags/flags/flat/24/United-Kingdom.png?raw=true)](RETE_en.md)
[![Ind](https://github.com/gosquared/flags/blob/master/flags/flags/flat/24/Indonesia.png?raw=true)](RETE_id.md)

[Tutorial](Tutorial_id.md) | [Rule Engine](RuleEngine_id.md) | [GRL](GRL_id.md) | [GRL JSON](GRL_JSON_id.md) | [RETE Algorithm](RETE_id.md) | [Functions](Function_id.md) | [FAQ](FAQ_id.md) | [Benchmark](Benchmarking_id.md)

---

Dari Wikipedia : Algoritma Rette (/ˈriːtiː/ REE-tee, /ˈreɪtiː/ RAY-tee, /ˈriːt/ REET, /rɛˈteɪ/ reh-TAY) adalah algorima pencian pola yang dipergunakan dalam implementasi sistem berasarkan __peraturan__ (__rule__)
Algoritma ini dibangun untuk dapat secara efektif menjalankan banyak __rule__ atau pola-pola kepada banyak data atau fakta dalam sebuah basis pengetahuan (__knowledgebase__).
Dengan demikan sebuah implementasi Rete dapat menentukan mana dari sekian banyak rule harus dijalankan bedasarkan
data dan fakta yang ada.

Beberapa bentuk dari algoritma RETE ini juga diimplementasikan dalam `grule-rule-engine` semenjak versi `1.1.0`
yang mengantikan pendekatan sebelumnya yang `Naive` saat melakukan evaluasi __rule__ untuk menambahkan rule tersebut kdalam
`ConflictSet`

`ExpressionAtom` dalam DRL dikompilasi dan dipastikan tidak terjadi duplikasi di dalam __working memory__ di dalam Grule. 
Cara ini mengikatkan kinerja mesin secara berarti terlebih jika ada banyak __rule__ yang tersimpan dan tiap-tiap rule tersebut
memiliki banyak duplikasi __expression__ atau banyak memiliki fungsi-fungsi yang __berat__.

Implementasi RETE pada Grule tidak mengenal __selector__ `Class` karena dalam Grule, sebuah __expression__ bisa 
berisi definisi dari banyak class. Contohnya, sebuah __expression__ seperti:

```.go
when
    ClassA.attr == ClassB.attr + ClassC.AFunc()
then
    ...
```

Pada __expression__ diatas, melibatkan atribut dan fungsi yang berasal dari 3 __class__ berbeda.
Fitur ini menjadikan pemilahan __class__ dalam RETE menjadi rumit.

Anda dapat membaca lebih jauh mengenai algoritma RETE disini:

* https://en.wikipedia.org/wiki/Rete_algorithm
* https://www.drdobbs.com/architecture-and-design/the-rete-matching-algorithm/184405218
* https://www.sparklinglogic.com/rete-algorithm-demystified-part-2/ 

### Mengapa algoritma RETE perlu dipergunakan

Asumsikan kita memiliki data fakta.

```go
type Fact struct {
    StringValue string
}

func (f *Fact) VeryHeavyAndLongFunction() bool {
    ...
}
```

Kemudian  fakta tersebut dimasukan kedalam konteks data.

```go
f := &Fact{}
dctx := context.NewDataContext()
err := dctx.Add("Fact", f)
```

Kemudian kita juga punya DRL seperti ...

```go
rule ... {
    when
        Fact.VeryHeavyAndLongFunction() && Fact.StringValue == "Fish"
    then
        ...
}
rule ... {
    when
        Fact.VeryHeavyAndLongFunction() && Fact.StringValue == "Bird"
    then
        ...
}
rule ... {
    when
        Fact.VeryHeavyAndLongFunction() && Fact.StringValue == "Mammal"
    then
        ...
}
...
// and alot more of simillar rule
...
rule ... {
    when
        Fact.VeryHeavyAndLongFunction() && Fact.StringValue == "Insect"
    then
        ...
}
```

Menjalankan DRL diatas bisa saja "membunuh" mesin __rule__ karena ia akan mencoba 
mengevaluasi dan memanggil setiap `VeryHeavyAndLongFunction` yang ada di dalam skop `When`
untuk menentukan rule mana yang menjadi kandidat untuk dieksekusi. 

Karenanya, mesin __rule__ tidak menjalankan setiap fungsi `Fact.VeryHeavyAndLongFunction` yang ada
dalam skrip. Algoritma Rete hanya melakukan evaluasi terhadap satu saja pemanggilan fungsi ini dan mengingat
hasil panggilan fungsi ini. Jadi saat fungsi ini seharusnya dieksekusi kembali, Rete cukup __mengembalikan__ hasil
eksekusi yang iya ingat.

Hal yang sama juga untuk `Fact.StringValue`. Algoritma Rete akan memuat nilai dari sebuah variabel dari dalam konteks
data dan mengingatnya. Hingga variabel tersebut berubah dalam skop `Then`, seperti ..

```go
rule ... {
    when
        ...
    then
        Fact.StringValue = "something else";
}
```

### Apa isi dari Working-Memory dalam Grule

Grule akan berusaha mengingat seluruh `Expression` yang di definiskan didalam skop `when` di setiap __rule__ 

Pertama, ia akan berusaha agar tidak ada satu pun node AST (Abstract Syntax Tree) yang terduplikasi. 

Kedua, untuk setiap node AST tersebut hanya boleh di-evaluasi satu kali, hingga elemen variabel didalam node tersebut
berubah. Seperti :

Boolean Expression :

```text
    when
    Fact.A == Fact.B + Fact.Func(Fact.C) - 20
```

__Expression__ ini bisa dipecah menjadi beberapa __expression__.

```text
Expression "Fact.A" --> Sebuah variabel
Expression "Fact.B" --> Sebuah variabel
Expression "Fact.C" --> Sebuah variabel
Expression "Fact.Func(Fact.C)"" --> Sebuah fungsi yang memiliki argumen Fact.C
Expression "20" --> Sebuah konstanta
Expression "Fact.B + Fact.Func(Fact.C)" --> Sebuah operasi matematika yange berisi 2 variabel; Fact.B dan Fact.C
Expression "(Fact.B + Fact.Func(Fact.C))" - 20 -- Sebuah operasi matematika yang juga berisi 2 variabel.
```

Setiap dari __expression__ diatas akan diingat nilai yang dimiliki/dihasilkan setiap 
kali nilai tersebut dimuat pertamakali. Jadi evaluasi berikutnya, fungsi atau variabel tidak akan
dipanggil atau dimuat ulang karena nilai nya itu sendiri otomatis di kembalikan.

Jika satu dari variabel tersebut berubah nilainya dalam skop `then`, sebagai contoh

```text
    then
        Fact.B = Fact.A * 20
```

Kita lihat `Fact.B` nilainya berubah, maka semua __Ekspression__ yang berisi `Fact.B` akan dihilangkan
dari __Working Memory__:

```text
Expression "Fact.B" --> Sebuah variabel
Expression "Fact.B + Fact.Func(Fact.C)" --> Sebuah operasi matematika yange berisi 2 variabel; Fact.B dan Fact.C
Expression "(Fact.B + Fact.Func(Fact.C))" - 20 -- Sebuah operasi matematika yang juga berisi 2 variabel.
```

Untuk __expression__ yang dihilangkan dari __working memory__, nilai mereka akan di re-evaluasi pada siklus berikutnya.

### Masalah yang dihadapi oleh RETE berkenaan dengan fungsi dan __method__

Saat Grule mencoba mengingat nilai variabel yang ia evaluasi dalam skop `when` dan `then`, jika anda mengubah
nilai mereka dari luar __rule engine__, contohnya merubah nilai-nilai mereka dari dalam pemanggilan fungsi,
maka Grule tidak dapat "melihat" perubahan nilai ini, karenanya Grule dapat keliru saat mengevaluasi nilai-nilai variabel
dan fungsi.

Anggaplah ada __fact__ berikut:

```go
type Fact struct {
    StringValue string
}

func (f *Fact) SetStringValue(newValue string) {
    f.StringValue = newValue
}
```

Kemudian anda masukan __fact__ ini kedalam konteks data

```go
f := &Fact{
    StringValue: "One",
}
dctx := context.NewDataContext()
err := dctx.Add("Fact", f)
```

Dalam GRL, anda melakukan seperti ini

```go
rule one "One" {
    when
        Fact.StringValue == "One"
        // here grule remembers that Fact.StringValue value is "One"
    then
        Fact.SetStringValue("Two");
        // here grule does not know if Fact.StringValue has changed inside the function.
        // What grule know is Fact.StringValue is still "One"
}

rule two "Two" {
    when
        Fact.StringValue == "Two"
        // Because of that, this will never evaluated true.
    then
        Fact.SetStringValue("Three");
}
```

Maka __rule engine__ akan selesai eksekusi tanpa ada kesalahan, tapi dari hasil yang ada terjadi kesalahan
dimana `Fact.StringValue` seharusnya `Two` tidak terapai.

Untuk mengatasi masalah ini, anda harus memberi "petunjuk" kepada grule bahwa variabel tersebut telah
berubah nilainya menggunakan fungsi `Changed`.  

```go
rule one "One" {
    when 
        Fact.StringValue == "One"
        // here grule remember that Fact.StringValue value is "One"
    then
        Fact.SetStringValue("Two");
        // here grule does not know if Fact.StringValue has changed inside the function.
        // What grule know is Fact.StringValue is still "One"

        // We should tell Grule that the variable changed within the Fact
        Changed("Fact.StringValue")
}
```
