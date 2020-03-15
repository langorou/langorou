package model

func sortQuick(arr []uint32) {
	qsort(arr, 0, len(arr)-1)
}

func qsort(arr []uint32, low int, high int) {
	if low >= high {
		return
	}

	p := partition(arr, low, high)
	qsort(arr, low, p-1)
	qsort(arr, p+1, high)
}

func partition(arr []uint32, low int, high int) int {
	pivot := low

	for i := low + 1; i <= high; i++ {
		if arr[i] <= arr[low] {
			pivot++
			arr[pivot], arr[i] = arr[i], arr[pivot]
		}
	}

	arr[pivot], arr[low] = arr[low], arr[pivot]
	return pivot
}

type sortableU32 []uint32

func (u sortableU32) Len() int {
	return len(u)
}

func (u sortableU32) Less(i, j int) bool {
	return u[i] < u[j]
}

func (u sortableU32) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func classicSort(buf []uint32) {
	for j := range buf {
		for k := range buf {
			if buf[j] < buf[k] && j > k {
				buf[j], buf[k] = buf[k], buf[j]
			}
		}
	}
}
