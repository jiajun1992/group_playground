package main

import (
	crand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"math/rand"
)

var one = big.NewInt(1)
var two = big.NewInt(2)
var Max2000 = big.NewInt(2000)

// GenSafePrime generates a pair of safe prime
func GenSafePrime(bits int) (*big.Int, *big.Int) {
	var p, q big.Int

	flag := false
	for !flag {
		temp, err := crand.Prime(crand.Reader, bits)
		if err != nil {
			println("Err = ", err)
		}
		p.Mul(two, temp)
		p.Add(one, &p)
		flag = p.ProbablyPrime(40)
		if flag {
			q = *temp
		}
	}
	return &p, &q
}

func FindGenerator(p, q *big.Int, seed int64) *big.Int {
	rng := rand.New(rand.NewSource(seed))
	var g big.Int

	for {
		g.Rand(rng, p)
		var temp big.Int
		temp.Exp(&g, q, p)
		if temp.Cmp(one) == 0 {
			return &g
		}
	}

}

func oneRound(seed int64) (int, int) {
	bitLength := 17
	//Generate a random prime order group
	p, q := GenSafePrime(bitLength)
	g := FindGenerator(p, q, seed)

	fmt.Println("p = ", *p)
	fmt.Println("q = ", *q)
	fmt.Println("g = ", *g)

	var sqrtP big.Int
	sqrtP.Sqrt(p)
	fmt.Println("sqrtP = ", sqrtP)

	// find a random x in Z_p, and g^x mod p is not a prime
	rng := rand.New(rand.NewSource(seed))
	var x big.Int

	for {
		x.Rand(rng, q)
		var tempGx big.Int
		tempGx.Exp(g, &x, p)
		flag := tempGx.ProbablyPrime(40)
		if !flag {
			break
		}
	}

	// find a pair of y and a, s.t. x = y-a mod q
	var y, a big.Int
	// We denote A = g^y mod p, B = g^a mod p
	var A, B big.Int

	iCounter := 0
	maxRound := math.Pow(2, float64(bitLength))
	for {
		y.Rand(rng, q)
		a.Sub(&y, &x)
		a.Mod(&a, q)
		A.Exp(g, &y, p)
		B.Exp(g, &a, p)

		if A.Cmp(&sqrtP) == -1 && B.Cmp(&sqrtP) == -1 {
			fmt.Println("A = ", A)
			fmt.Println("B = ", B)
			fmt.Println("Rounds to break DI = ", iCounter)
			break
		}
		if iCounter >= int(maxRound) {
			return iCounter, 1
		}
		iCounter++
	}
	return iCounter, 0
}

func main() {
	var iCounter int64 = 777
	total := 0
	totalNotFound := 0
	maxRound := 1000
	for i := 0; i < maxRound; i++ {
		temp1, temp2 := oneRound(iCounter + int64(i))
		if temp2 == 0 {
			total += temp1
		}
		totalNotFound += temp2
	}
	fmt.Println("Total = ", total)
	fmt.Println("TotalNotFound = ", totalNotFound)
	fmt.Println("Average = ", total/(maxRound-totalNotFound))
}
