package images

// Transpose performs a 90 degree counter-clockwise rotation of the byte slice
// representation of a grayscale image
func Transpose(dst, src []byte, long, short int) {
	for srcy := 0; srcy < short; srcy++ {
		for srcx := 0; srcx < long; srcx++ {
			dst[(long-srcx-1)*short+srcy] = src[srcy*long+srcx]
		}
	}
}
