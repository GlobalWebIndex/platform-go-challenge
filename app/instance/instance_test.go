package instance_test

import (
	"fmt"

	"x-gwi/app/instance"
)

func ExampleInstance_PassVerifyHash() {
	i1 := instance.NewInstance()
	i1h1 := i1.PassHash()
	i1h2 := i1.PassHash()
	i2 := instance.NewInstance()
	i2h1 := i2.PassHash()
	i2h2 := i2.PassHash()

	fmt.Printf("- i1\n")
	fmt.Printf("i1.Valid(): bool: %v\n", i1.Valid())
	fmt.Printf("i1.PassVerifyHash(i1h1): bool: %v\n", i1.PassVerifyHash(i1h1))
	fmt.Printf("i1.PassVerifyHash(i1h2): bool: %v\n", i1.PassVerifyHash(i1h2))
	fmt.Printf("i1.PassVerifyHash(i2h1): bool: %v\n", i1.PassVerifyHash(i2h1))
	fmt.Printf("i1.PassVerifyHash(i2h2): bool: %v\n", i1.PassVerifyHash(i2h2))

	fmt.Printf("- i2\n")
	fmt.Printf("i2.Valid(): bool: %v\n", i2.Valid())
	fmt.Printf("i2.PassVerifyHash(i1h1): bool: %v\n", i2.PassVerifyHash(i1h1))
	fmt.Printf("i2.PassVerifyHash(i1h2): bool: %v\n", i2.PassVerifyHash(i1h2))
	fmt.Printf("i2.PassVerifyHash(i2h1): bool: %v\n", i2.PassVerifyHash(i2h1))
	fmt.Printf("i2.PassVerifyHash(i2h2): bool: %v\n", i2.PassVerifyHash(i2h2))

	fmt.Printf("- hashes\n")
	fmt.Printf("i1h1 == i1h1: bool: %v\n", string(i1h1) == string(i1h1))
	fmt.Printf("i1h1 == i1h2: bool: %v\n", string(i1h1) == string(i1h2))
	fmt.Printf("i1h1 == i2h1: bool: %v\n", string(i1h1) == string(i2h1))
	fmt.Printf("i1h1 == i2h2: bool: %v\n", string(i1h1) == string(i2h2))

	// fmt.Println(" dev view only")
	// fmt.Println(i1)
	// fmt.Println(string(i1h1))
	// fmt.Println(string(i1h2))
	// fmt.Println(i2)
	// fmt.Println(string(i2h1))
	// fmt.Println(string(i2h2))

	// passHash, _ := bcrypt.GenerateFromPassword([]byte(""), bcrypt.MinCost)
	// fmt.Printf("passHash: %s\n", string(passHash))

	// Output:
	//
	// - i1
	// i1.Valid(): bool: true
	// i1.PassVerifyHash(i1h1): bool: true
	// i1.PassVerifyHash(i1h2): bool: true
	// i1.PassVerifyHash(i2h1): bool: false
	// i1.PassVerifyHash(i2h2): bool: false
	// - i2
	// i2.Valid(): bool: true
	// i2.PassVerifyHash(i1h1): bool: false
	// i2.PassVerifyHash(i1h2): bool: false
	// i2.PassVerifyHash(i2h1): bool: true
	// i2.PassVerifyHash(i2h2): bool: true
	// - hashes
	// i1h1 == i1h1: bool: true
	// i1h1 == i1h2: bool: false
	// i1h1 == i2h1: bool: false
	// i1h1 == i2h2: bool: false
	//
}
