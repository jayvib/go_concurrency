package ordone

func OrDone(dones ...<-chan interface{}) <-chan interface{} {
	switch len(dones) {
	case 0:
		return nil
	case 1:
		return dones[0]
	}

	ordone := make(chan interface{})

	go func() {
		defer close(ordone)

		switch len(dones) {
		case 2:
			select {
			case <-dones[0]:
			case <-dones[1]:
			}
		default: // more than 2 dones to be managed.
			select {
			case <-dones[0]:
			case <-dones[1]:
			case <-dones[2]:
			case <-OrDone(append(dones[3:], ordone)...):
			}
		}
	}()
	return ordone
}
