package tseg

import (
	"os"
	"strings"
	"unicode"
)

type freqDict map[string]float64

type dict map[string]int

// Структура для хранения всех разбиений строки вместо глобальной переменной
type segsAccum struct {
	segs [][]string
}

// Добавляет разбиение в slice разбиений
func (sa *segsAccum) addSeg(seg []string) {
	if sa.segs == nil {
		sa.segs = make([][]string, 0)
	}
	sa.segs = append(sa.segs, seg)
}

// Функция, читающая файл и возвращающая его содержимое в string
func readFile(path string, bufSize int) (fileStr string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	buf := make([]byte, bufSize)
	n, err := file.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

// Читает словарь и возвращает slice из слов в нем
func parseDict(path string) (words []string, err error) {
	dict, err := readFile(path, 1024*1024)
	if err != nil {
		return nil, err
	}
	return strings.Fields(dict), nil
}

// Читает обучающий текст, переводит в нижний регистр,
// удаляет все символы, кроме букв и возвращает
// slice из слов
func parseText(path string) (words []string, err error) {
	text, err := readFile(path, 3*1024*1024)
	if err != nil {
		return nil, err
	}
	text = strings.ToLower(text)
	remover := func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		return ' '
	}
	text = strings.Map(remover, text)
	return strings.Fields(text), nil
}

// Добавляет все слова из обучающего текста в slice словаря
func addWordsToDictFromText(dictSlice []string, textSlice []string) (newDict []string) {
	return append(dictSlice, textSlice...)
}

// Создает словарь map[string]int, содержащий все слова
// из словаря и обучающего текста
func createDict(dictSlice []string) dict {
	d := make(dict)
	for _, val := range dictSlice {
		d[val] = 0
	}
	return d
}

// Создает частотный словарь всех биграм обучающего текста
func createFreqDict(textSlice []string) freqDict {
	fd := make(freqDict)
	numberTextWords := len(textSlice)
	for i := 1; i < numberTextWords; i++ {
		fd[textSlice[i-1] + " " + textSlice[i]]++
	}
	for key := range fd {
		fd[key] /= float64(numberTextWords)
	}
	return fd
}

// Создает все возможные сегментации текста согласно словарю и обучающему тексту и
// записывает их в глобальную переменную segs
func getTextSegs(str string, d dict, fd freqDict, seg []string, sa *segsAccum) {
	if len(str) == 0 {
		sa.addSeg(seg)
	}
	var word string
	for _, r := range str {
		word = word + string(r)
		if _, ok := d[word]; ok {
			getTextSegs(str[len(word):], d, fd, append(seg, word), sa)
		}
	}
}

// Выбирает наиболее вероятную сегментацию согласно частотному словарю биграм обучающего текста
func chooseBest(fd freqDict, sa *segsAccum) (bestSeg []string) {
	var bestFreq float64
	var bestNumber int
	for i := range sa.segs {
		freq := 1.0
		lenSegsI := len(sa.segs[i])
		for j := 1; j < lenSegsI; j++ {
			bigram := sa.segs[i][j-1] + " " + sa.segs[i][j]
			if val, ok := fd[bigram]; ok {
				freq *= val
			} else {
				freq = 0
			}
		}
		if freq >= bestFreq {
			bestFreq = freq
			bestNumber = i
		}
	}
	return sa.segs[bestNumber]
}

//Основная функция, возвращающая сегментацию текста
func GetTextSegmentation(str string, dictPath string, textPath string) ([]string, error) {
	dictSlice, err := parseDict(dictPath)
	if err != nil {
		return nil, err
	}
	textSlice, err := parseText(textPath)
	if err != nil {
		return nil, err
	}
	dictSlice = addWordsToDictFromText(dictSlice, textSlice)
	d := createDict(dictSlice)
	fd := createFreqDict(textSlice)
	sa := &segsAccum{segs: make([][]string, 0)}
	getTextSegs(str, d, fd, make([]string, 0), sa)
	return chooseBest(fd, sa), nil
}