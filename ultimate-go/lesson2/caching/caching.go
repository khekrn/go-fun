package caching

// Create a square matrix of 16,777,216 bytes.
const (
	rows = 4 * 1024
	cols = 4 * 1024
)

// matrix represents a matrix with a large number of
// columns per row.
var matrix [rows][cols]byte

// data represents data node for the linkedlist
type data struct {
	v byte
	p *data
}

var list *data

func init() {
	var last *data

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			var d data
			if last == nil {
				list = &d
			}
			if last != nil {
				last.p = &d
			}
			last = &d

			// Add a value to all even elements.
			if row%2 == 0 {
				matrix[row][col] = 0xFF
				d.v = 0xFF
			}
		}

		// Count the number of elements in the link list.
		var ctr int
		d := list
		for d != nil {
			ctr++
			d = d.p
		}

		//fmt.Println("Elements in the link list", ctr)
		//fmt.Println("Elements in the matrix", rows*cols)
	}
}

// LinkedListTraverse traverses the linked list linearly.
// The linked list is twice as slow as the row traversal mainly because
// there are cache line misses but fewer TLB (Translation Lookaside Buffer)
// misses. A bulk of the nodes connected in the list exist inside the same
// OS pages.

// Each running program is given a full memory map of virtual memory by the
// OS and that running program thinks they have all of the physical memory
// on the machine. However, physical memory needs to be shared with all the
// running programs. The operating system shares physical memory by breaking
// the physical memory into pages and mapping pages to virtual memory for
// any given running program. Each OS can decide the size of a page, but 4k,
// 8k, 16k are reasonable and common sizes.

// The TLB is a small cache inside the processor that helps to reduce latency
// on translating a virtual address to a physical address within the scope of
// an OS page and offset inside the page. A miss against the TLB cache can cause
// large latencies because now the hardware has to wait for the OS to scan its
// page table to locate the right page for the virtual address in question.
// If the program is running on a virtual machine (like in the cloud) then the
// virtual machine paging table needs to be scanned first.
func LinkedListTraversal() int {
	var ctr int

	d := list

	for d != nil {
		if d.v == 0xFF {
			ctr++
		}
		d = d.p
	}

	return ctr
}

// ColumnTraverse traverses the matrix linearly down each column.
// Column Traverse is the worst by an order of magnitude because this access
// pattern crosses over OS page boundaries on each memory access. This
// causes no predictability for cache line prefetching and becomes essentially
// random access memory.
func ColumnTraverse() int {
	var ctr int

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			if matrix[row][col] == 0xFF {
				ctr++
			}
		}
	}

	return ctr
}

// RowTraverse traverses the matrix linearly down each row.
// Row traverse will have the best performance because it walks through
// memory, cache line by connected cache line, which creates a predictable
// access pattern. Cache lines can be prefetched and copied into the L1 or
// L2 cache before the data is needed.
func RowTraverse() int {
	var ctr int

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if matrix[row][col] == 0xFF {
				ctr++
			}
		}
	}
	return ctr
}
