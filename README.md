# Сегментация текста на Golang
Использованный алгоритм найден на [хабре](https://habr.com/post/141228/).

`tseg.go` основной файл с реализацией.

`tseg_test.go` представляет собой тест. Чтобы начать тестирование, запустите `go test`.

`dict.txt` хранит словарь.

`text.txt` хранит обучающий текст, на основе которого и будет происходить сегментация. 
Нужен для того, чтобы определить вероятность расположения вместе двух слов.
Сегментация сильно зависит от обучающего текста.

### Установка

`go get github.com/mov-inc-xor/tseg`

### Использование

```go
package main

import (
	"fmt"
	"log"
	"tseg"
)

func main() {
	str := "ilovetee"
	sr, err := tseg.NewSegmentator("dict.txt", "text.txt")
	if err != nil {
		log.Fatalln(err)
		return
	}
	seg, err := sr.GetSegmentation(str)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(seg)
}
```

`-> [i love tee]`
