package bitmap

type Bitmap struct {
	bits     []byte
	capacity int
}

func NewBitmap(cap int) *Bitmap {
	return &Bitmap{
		capacity: cap,
		bits:     make([]byte, (cap>>3)+1),
	}
}
