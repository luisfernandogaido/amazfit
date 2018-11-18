package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	urlbase     = "https://amazfitwatchfaces.com"
	folderFiles = "./files/"
)

var (
	reIds  *regexp.Regexp
	reFile *regexp.Regexp
	reImg  *regexp.Regexp
)

func init() {
	reIds = regexp.MustCompile(`"\/pace\/view\/\?id=([0-9]+)"`)
	reFile = regexp.MustCompile(`"\/pace\/download\?file=([^"]+)"`)
	reImg = regexp.MustCompile(`\/pace\/resource\/img\/([^"]+)`)
}

func main() {
	t0 := time.Now()
	var (
		err     error
		paginas int
	)
	if len(os.Args) == 1 {
		paginas = 118
	} else {
		paginas, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal("número de páginas inválido")
		}
	}
	if err := os.MkdirAll(folderFiles, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Procurando ids em %v páginas.\n", paginas)
	ids, err := getAllIds(1, paginas, 16)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v ids encontrados. Baixando arquivos e imagens...\n", len(ids))
	if err := getAllFiles(ids, 16); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Operação concluída em %v.\n", time.Since(t0))
	fmt.Println(time.Since(t0))
}

func getIds(p int) ([]string, error) {
	res, err := http.Get(urlbase + "/pace/p/" + strconv.Itoa(p))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	matches := reIds.FindAllStringSubmatch(string(b), -1)
	ids := make([]string, 0, 16)
	for _, m := range matches {
		ids = append(ids, m[1])
	}
	return ids, nil
}

func getAllIds(ini, fim, c int) ([]string, error) {
	n := fim - ini + 1
	chIds := make(chan []string, n)
	chErr := make(chan error, n)
	sem := make(chan struct{}, c)
	for i := ini; i <= fim; i++ {
		sem <- struct{}{}
		go func(i int) {
			defer func() { <-sem }()
			ids, err := getIds(i)
			if err != nil {
				chErr <- err
				return
			}
			chIds <- ids
		}(i)
	}
	con := 0
	allIds := make([]string, 0, 1000)
	for {
		select {
		case err := <-chErr:
			return nil, err
		case ids := <-chIds:
			allIds = append(allIds, ids...)
			con++
			if con == n {
				return allIds, nil
			}
		}
	}
}

func getFile(id string) error {
	res, err := http.Get(urlbase + "/pace/view/?id=" + id)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	conteudo := string(b)
	download := reFile.FindString(conteudo)
	img := reImg.FindString(conteudo)
	endereco := urlbase + strings.Replace(download, `"`, "", 2)
	enderecoImg := urlbase + img
	p := strings.Index(endereco, "file")
	q := strings.Index(endereco, ".wfz")
	if q == -1 {
		q = strings.Index(endereco, ".apk")
	}
	fileName := id + "-" + endereco[p+5:q+4]
	fileNameImg := id + "-" + filepath.Base(enderecoImg)
	res2, err := http.Get(endereco)
	if err != nil {
		return err
	}
	defer res2.Body.Close()
	b, err = ioutil.ReadAll(res2.Body)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(folderFiles+fileName, b, 0644); err != nil {
		return err
	}
	res3, err := http.Get(enderecoImg)
	if err != nil {
		return err
	}
	defer res3.Body.Close()
	b, err = ioutil.ReadAll(res3.Body)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(folderFiles+fileNameImg, b, 0644); err != nil {
		return err
	}
	return nil
}

func getAllFiles(ids []string, c int) error {
	sem := make(chan struct{}, c)
	chErr := make(chan error)
	for i := 0; i < len(ids); i++ {
		select {
		case err := <-chErr:
			return err
		default:

		}
		sem <- struct{}{}
		go func(i int) {
			defer func() { <-sem }()
			if err := getFile(ids[i]); err != nil {
				chErr <- err
			}
			fmt.Printf("%v/%v (%v%%)\n", i+1, len(ids), math.Round(100*float64(i+1)/float64(len(ids))))
		}(i)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
	return nil
}
