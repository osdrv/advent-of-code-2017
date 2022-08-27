package main

const (
	FACT_A = uint64(16807)
	FACT_B = uint64(48271)
	MOD_F  = uint64(2147483647)
)

func part1() {
	println("===== part 1 =====")
	startA, startB := uint64(699), uint64(124)

	matched := 0
	valA, valB := startA, startB
	for i := 0; i < 40_000_000; i++ {
		valA = (valA * FACT_A) % MOD_F
		valB = (valB * FACT_B) % MOD_F
		//printf("valA: %d, valB: %d", valA, valB)
		if (valA & 0xFFFF) == (valB & 0xFFFF) {
			matched++
		}
	}

	printf("matched: %d", matched)
}

func part2() {
	println("===== part 2 =====")
	startA, startB := uint64(699), uint64(124)

	matched := 0
	valA, valB := startA, startB
	for i := 0; i < 5_000_000; i++ {
		for {
			valA = (valA * FACT_A) % MOD_F
			if (valA & 3) == 0 {
				break
			}
		}
		for {
			valB = (valB * FACT_B) % MOD_F
			if (valB & 7) == 0 {
				break
			}
		}
		//printf("valA: %d, valB: %d", valA, valB)
		if (valA & 0xFFFF) == (valB & 0xFFFF) {
			matched++
		}
	}

	printf("matched: %d", matched)
}

func main() {
	part1()
	part2()
}
