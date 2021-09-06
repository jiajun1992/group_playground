package main

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
)

var Max2000 = big.NewInt(2000)

func main() {

	var seed int64
	seed = 123
	//Generate a random prime order group
	rng := rand.New(rand.NewSource(seed))

	var p, q big.Int

	p.Rand(rng, Max2000)
	fmt.Println("p = ", p, "q = ", q)
	flag := p.ProbablyPrime(50)
	fmt.Println("flag = ", flag)
	q, err := crand.Prime(seed, 100)

}
