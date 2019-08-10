package main

// MinOf arbitary no. of bytes
func MinOf(vars ...byte) byte {
	mini := vars[0]
	for _, i := range vars {
		if mini > i {
			mini = i
		}
	}

	return mini
}

// MaxOf arbitary no. of bytes
func MaxOf(vars ...byte) byte {
	maxi := vars[0]
	for _, i := range vars {
		if i > maxi {
			maxi = i
		}
	}

	return maxi
}
