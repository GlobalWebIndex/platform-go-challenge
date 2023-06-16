package id

import "github.com/google/uuid"

// UUID creates a UUID Version 4 based on random (crypto/rand)
func UUID() uuid.UUID {
	return uuid.New()
}

// UUID1 creates a UUID Version 1 based on the current NodeID and clock
func UUID1() uuid.UUID {
	u, _ := uuid.NewUUID()

	return u
}

// UUID creates a UUID Version 4 based on random (crypto/rand)
func UUID4() uuid.UUID {
	return uuid.New()
}
