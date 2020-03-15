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
	pivot := arr[low]
	i := low

	for j := low; j <= high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}

	arr[i], arr[high] = arr[high], arr[i]
	return i
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
