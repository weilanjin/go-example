package channel

import "unsafe"

/*
	select {
		case ch <- 1:
		default:
	}
*/

// nb 代表 non-blocking
func selectnbsend(c *hchan, elem unsafe.Pointer) (selected bool) {
	return chansend(c, elem, false, getcallerpc())
}

/*
	select {
		case <-ch:
		default:
	}
*/

func selectnbrecv(elem unsafe.Pointer, c *hchan) (selected, received bool) {
	return chanrecv(c, elem, false)
}