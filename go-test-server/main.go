package main

import (
	"fmt"
	"net/http"
)

type String string

// 素因数分解する関数
func factorize(numbers []int, c chan<- []int) {
	for _, number := range numbers {
		var a []int
		for i := 1; i < number+1; i++ {
			if number%i == 0 {
				a = append(a, i)
			}
		}
		c <- a
	}
	// 送信側がチャネルをクローズする
	close(c)
}

func (s String) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	numbers := []int{53541233, 21235343, 11421443, 5423123}
	c := make(chan []int)
	go factorize(numbers, c)
	fmt.Fprint(w, s)
}

func main() {
	http.Handle("/", String("Success!"))
	http.ListenAndServe(":8000", nil)
}
