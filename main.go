package ciproject

import "fmt"

func main() {
	a := 10
	b := 9
	fmt.Println(Result(a,b))
}

func Result(a int, b int)(g int){
	g = a + b
	fmt.Println("a + b =", g)
	return
}
