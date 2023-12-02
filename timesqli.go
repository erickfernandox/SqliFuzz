package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	red   = "\033[31m"
	gray  = "\033[90m"
	reset = "\033[0m"
)

func medirTempoRequisicao(url string) float64 {
	inicio := time.Now()

	_, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erro na requisição para %s: %s\n", url, err)
	}

	fim := time.Now()
	tempoTotal := fim.Sub(inicio).Seconds()
	return tempoTotal
}

func testarURLs(tempoSQLi float64) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := scanner.Text()
		url = url // Remover espaços em branco e quebras de linha
		if medirTempoRequisicao(url) >= tempoSQLi {
			fmt.Printf("%sVulnerable: %s%s - {%f}\n", red, url, reset, medirTempoRequisicao(url))
		} else {
			fmt.Printf("%sNot Vulnerable: %s%s\n", gray, url, reset)
		}
	}
}

func main() {
	var tempoSQLi float64

	flag.Float64Var(&tempoSQLi, "t", 0, "Tempo a ser testado.")
	flag.Parse()

	if tempoSQLi == 0 {
		fmt.Println("Uso: TimeSQLi -t <tempo>")
		os.Exit(1)
	}

	testarURLs(tempoSQLi)
}
