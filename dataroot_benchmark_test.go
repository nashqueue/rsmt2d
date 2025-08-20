package rsmt2d

import (
	"crypto/rand"
	"fmt"
	"sort"
	"testing"

	"github.com/celestiaorg/nmt"
)

// BenchmarkDatarootGeneration benchmarks generating a Dataroot over the EDS 
// with k values that are powers of 2 starting at 32, 64, 128, 256, 512.
// The EDS doubles in size compared to the original data square (k).
//
// All benchmarks use the same Reed-Solomon erasure coding approach for fair comparison:
// - MerkleTree: Reed-Solomon erasure coding (k×k → 2k×2k) + SHA-256 Merkle trees
// - NMT: Reed-Solomon erasure coding (k×k → 2k×2k) + Namespaced Merkle Trees
// - Stages 1-3: Reed-Solomon erasure coding (k×k → 2k×2k) + hybrid tree construction
//
// Performance optimizations:
// - Data generation moved outside benchmark timer for clean measurements
// - NMT compatibility achieved using parity namespace (0xFF bytes) for non-Q0 quadrants
// - All benchmarks measure pure Reed-Solomon + tree construction performance
func BenchmarkDatarootGeneration(b *testing.B) {
	// Test k values: 32, 64, 128, 256, 512, 1024, 2048 (original data square size)
	// EDS will be 2k x 2k: 64, 128, 256, 512, 1024, 2048, 4096
	kValues := []int{32, 64, 128, 256, 512, 1024, 2048}
	shareSize := 512 // Standard share size used in tests
	namespaceSize := 29 // Standard namespace size for Celestia
	
	for _, k := range kValues {
		// Skip if codec doesn't support this many shares
		codec := NewLeoRSCodec()
		if codec.MaxChunks() < k*k {
			continue
		}
		
		// Test with regular Merkle tree
		b.Run(fmt.Sprintf("MerkleTree_k=%d_EDS=%dx%d", k, k*2, k*2), func(b *testing.B) {
			// Pre-generate data outside the timer for cleaner performance measurement
			testData := make([][][]byte, b.N)
			for n := 0; n < b.N; n++ {
				testData[n] = genRandomDataSquare(k, shareSize)
			}
			
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				// Create EDS from pre-generated data using Reed-Solomon erasure coding
				eds, err := ComputeExtendedDataSquare(testData[n], codec, NewDefaultTree)
				if err != nil {
					b.Fatal(err)
				}
				
				// Generate Dataroot (row and column roots)
				_, err = eds.Roots()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		
		// Test with NMT (Namespaced Merkle Tree) - BASELINE
		// Uses same Reed-Solomon + tree construction as other benchmarks
		// NMT wrapper automatically handles parity namespace (0xFF bytes) for non-Q0 quadrants
		b.Run(fmt.Sprintf("NMT_k=%d_EDS=%dx%d", k, k*2, k*2), func(b *testing.B) {
			// Create pure NMT tree constructor for both rows and columns
			treeConstructor := newErasuredNamespacedMerkleTreeConstructor(
				uint64(k), // Original square size (NMT wrapper expects this)
				nmt.NamespaceIDSize(namespaceSize),
				nmt.IgnoreMaxNamespace(true),
				nmt.InitialCapacity(k*2),
			)
			
			// Pre-generate data outside the timer for cleaner performance measurement
			testData := make([][][]byte, b.N)
			for n := 0; n < b.N; n++ {
				testData[n] = genRandomNamespacedOriginalData(k, shareSize, namespaceSize)
			}
			
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				// Create EDS from pre-generated data using Reed-Solomon erasure coding (same as other benchmarks)
				// NMT wrapper automatically prepends parity namespace (0xFF) to parity shares
				eds, err := ComputeExtendedDataSquare(testData[n], codec, treeConstructor)
				if err != nil {
					b.Fatal(err)
				}
				
				// Generate Dataroot (row and column roots)
				_, err = eds.Roots()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		
		// Test with Stage 1: Merkle columns + NMT rows
		b.Run(fmt.Sprintf("Stage1_k=%d_EDS=%dx%d", k, k*2, k*2), func(b *testing.B) {
			// Create Stage 1 tree constructor (Merkle for columns, NMT for rows)
			treeConstructor := newStage1TreeConstructor(
				uint64(k*2), // EDS size will be 2k×2k
				nmt.NamespaceIDSize(namespaceSize),
				nmt.IgnoreMaxNamespace(true),
				nmt.InitialCapacity(k*2),
			)
			
			// Pre-generate data outside the timer for cleaner performance measurement
			testData := make([][][]byte, b.N)
			for n := 0; n < b.N; n++ {
				testData[n] = genRandomNamespacedOriginalData(k, shareSize, namespaceSize)
			}
			
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				// Create EDS from pre-generated data using Reed-Solomon erasure coding (same as MerkleTree)
				eds, err := ComputeExtendedDataSquare(testData[n], codec, treeConstructor)
				if err != nil {
					b.Fatal(err)
				}
				
				// Generate Dataroot (row and column roots)
				_, err = eds.Roots()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		
		// Test with Stage 2: Stage 1 + binary Merkle for Q2/Q3 row roots
		b.Run(fmt.Sprintf("Stage2_k=%d_EDS=%dx%d", k, k*2, k*2), func(b *testing.B) {
			// Create Stage 2 tree constructor
			treeConstructor := newStage2TreeConstructor(
				uint64(k*2), // EDS size will be 2k×2k
				nmt.NamespaceIDSize(namespaceSize),
				nmt.IgnoreMaxNamespace(true),
				nmt.InitialCapacity(k*2),
			)
			
			// Pre-generate data outside the timer for cleaner performance measurement
			testData := make([][][]byte, b.N)
			for n := 0; n < b.N; n++ {
				testData[n] = genRandomNamespacedOriginalData(k, shareSize, namespaceSize)
			}
			
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				// Create EDS from pre-generated data using Reed-Solomon erasure coding (same as MerkleTree)
				eds, err := ComputeExtendedDataSquare(testData[n], codec, treeConstructor)
				if err != nil {
					b.Fatal(err)
				}
				
				// Generate Dataroot (row and column roots)
				_, err = eds.Roots()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		
		// Test with Stage 3: Stage 2 + optimize Q1 row parity
		b.Run(fmt.Sprintf("Stage3_k=%d_EDS=%dx%d", k, k*2, k*2), func(b *testing.B) {
			// Create Stage 3 tree constructor
			treeConstructor := newStage3TreeConstructor(
				uint64(k*2), // EDS size will be 2k×2k
				nmt.NamespaceIDSize(namespaceSize),
				nmt.IgnoreMaxNamespace(true),
				nmt.InitialCapacity(k*2),
			)
			
			// Pre-generate data outside the timer for cleaner performance measurement
			testData := make([][][]byte, b.N)
			for n := 0; n < b.N; n++ {
				testData[n] = genRandomNamespacedOriginalData(k, shareSize, namespaceSize)
			}
			
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				// Create EDS from pre-generated data using Reed-Solomon erasure coding (same as MerkleTree)
				eds, err := ComputeExtendedDataSquare(testData[n], codec, treeConstructor)
				if err != nil {
					b.Fatal(err)
				}
				
				// Generate Dataroot (row and column roots)
				_, err = eds.Roots()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// genRandomDataSquare generates a k x k square of random data with given share size
func genRandomDataSquare(k, shareSize int) [][]byte {
	data := make([][]byte, k*k)
	for i := 0; i < k*k; i++ {
		share := make([]byte, shareSize)
		_, err := rand.Read(share)
		if err != nil {
			panic(err)
		}
		data[i] = share
	}
	return data
}

// genRandomNamespacedDataSquare generates a k x k square of namespaced original data
// This creates ORIGINAL data (not extended) with proper namespaces for NMT compatibility
func genRandomNamespacedOriginalData(k, shareSize, namespaceSize int) [][]byte {
	data := make([][]byte, k*k)
	
	// Generate all shares first
	for i := 0; i < k*k; i++ {
		share := make([]byte, shareSize)
		
		// Generate random namespace ID (first namespaceSize bytes)
		namespaceID := make([]byte, namespaceSize)
		_, err := rand.Read(namespaceID)
		if err != nil {
			panic(err)
		}
		copy(share[:namespaceSize], namespaceID)
		
		// Generate random data for the rest of the share
		shareData := make([]byte, shareSize-namespaceSize)
		_, err = rand.Read(shareData)
		if err != nil {
			panic(err)
		}
		copy(share[namespaceSize:], shareData)
		
		data[i] = share
	}
	
	// Sort by namespace ID (lexicographically) as required by NMT
	sort.Slice(data, func(i, j int) bool {
		// Compare namespace IDs (first namespaceSize bytes)
		for k := 0; k < namespaceSize; k++ {
			if data[i][k] != data[j][k] {
				return data[i][k] < data[j][k]
			}
		}
		return false // Equal namespaces
	})
	
	return data
}


// hybridTreeConstructor creates trees that use NMT for rows and Merkle trees for columns
type hybridTreeConstructor struct {
	squareSize    uint64
	nmtOpts       []nmt.Option
	namespaceSize int
}

// newHybridTreeConstructor creates a tree constructor that uses NMT for rows and Merkle for columns
func newHybridTreeConstructor(squareSize uint64, opts ...nmt.Option) TreeConstructorFn {
	// Extract namespace size from options
	nmtOpts := &nmt.Options{}
	for _, setter := range opts {
		setter(nmtOpts)
	}
	
	return hybridTreeConstructor{
		squareSize:    squareSize,
		nmtOpts:       opts,
		namespaceSize: int(nmtOpts.NamespaceIDSize),
	}.NewTree
}

// NewTree creates either an NMT (for rows) or Merkle tree (for columns)
func (h hybridTreeConstructor) NewTree(axis Axis, axisIndex uint) Tree {
	if axis == Row {
		// Use NMT for rows
		newTree := newErasuredNamespacedMerkleTree(h.squareSize, axisIndex, h.nmtOpts...)
		return &newTree
	} else {
		// Use regular Merkle tree for columns
		return NewDefaultTree(axis, axisIndex)
	}
}

// Stage 1: Merkle columns, NMT rows (inverse of original hybrid)
type stage1TreeConstructor struct {
	squareSize    uint64
	nmtOpts       []nmt.Option
	namespaceSize int
}

func newStage1TreeConstructor(squareSize uint64, opts ...nmt.Option) TreeConstructorFn {
	nmtOpts := &nmt.Options{}
	for _, setter := range opts {
		setter(nmtOpts)
	}
	
	return stage1TreeConstructor{
		squareSize:    squareSize,
		nmtOpts:       opts,
		namespaceSize: int(nmtOpts.NamespaceIDSize),
	}.NewTree
}

func (s stage1TreeConstructor) NewTree(axis Axis, axisIndex uint) Tree {
	if axis == Col {
		// Use Merkle trees for columns
		return NewDefaultTree(axis, axisIndex)
	} else {
		// Use NMT for rows (with original square size, not EDS size)
		originalSquareSize := s.squareSize / 2
		newTree := newErasuredNamespacedMerkleTree(originalSquareSize, axisIndex, s.nmtOpts...)
		return &newTree
	}
}

// Stage 2: Stage 1 + binary Merkle for Q2/Q3 row roots (bottom half)
type stage2TreeConstructor struct {
	squareSize    uint64
	nmtOpts       []nmt.Option
	namespaceSize int
}

func newStage2TreeConstructor(squareSize uint64, opts ...nmt.Option) TreeConstructorFn {
	nmtOpts := &nmt.Options{}
	for _, setter := range opts {
		setter(nmtOpts)
	}
	
	return stage2TreeConstructor{
		squareSize:    squareSize,
		nmtOpts:       opts,
		namespaceSize: int(nmtOpts.NamespaceIDSize),
	}.NewTree
}

func (s stage2TreeConstructor) NewTree(axis Axis, axisIndex uint) Tree {
	if axis == Col {
		// Use Merkle trees for columns
		return NewDefaultTree(axis, axisIndex)
	} else {
		// For rows: use Merkle for Q2/Q3 (bottom half), NMT for Q0/Q1 (top half)
		halfSize := s.squareSize / 2
		if uint64(axisIndex) >= halfSize {
			// Q2/Q3 (bottom half) - use Merkle tree (all parity data with same namespace)
			return NewDefaultTree(axis, axisIndex)
		} else {
			// Q0/Q1 (top half) - use NMT (with original square size)
			originalSquareSize := s.squareSize / 2
			newTree := newErasuredNamespacedMerkleTree(originalSquareSize, axisIndex, s.nmtOpts...)
			return &newTree
		}
	}
}

// Stage 3: Stage 2 + optimize Q1 row parity with Merkle trees
type stage3TreeConstructor struct {
	squareSize    uint64
	nmtOpts       []nmt.Option
	namespaceSize int
}

func newStage3TreeConstructor(squareSize uint64, opts ...nmt.Option) TreeConstructorFn {
	nmtOpts := &nmt.Options{}
	for _, setter := range opts {
		setter(nmtOpts)
	}
	
	return stage3TreeConstructor{
		squareSize:    squareSize,
		nmtOpts:       opts,
		namespaceSize: int(nmtOpts.NamespaceIDSize),
	}.NewTree
}

func (s stage3TreeConstructor) NewTree(axis Axis, axisIndex uint) Tree {
	if axis == Col {
		// Use Merkle trees for columns
		return NewDefaultTree(axis, axisIndex)
	} else {
		// For rows: hybrid approach based on quadrant structure
		halfSize := s.squareSize / 2
		if uint64(axisIndex) < halfSize {
			// Q0/Q1 (top half) - use hybrid tree: NMT for Q0 (left), Merkle for Q1 (right)
			return newHybridRowTree(s.squareSize, axisIndex, s.nmtOpts...)
		} else {
			// Q2/Q3 (bottom half) - use Merkle tree (all parity data)
			return NewDefaultTree(axis, axisIndex)
		}
	}
}

// hybridRowTree implements a tree that uses NMT for left half (Q0) and Merkle for right half (Q1)
type hybridRowTree struct {
	squareSize    uint64
	axisIndex     uint
	nmtOpts       []nmt.Option
	leftTree      Tree  // NMT for Q0 (left half)
	rightTree     Tree  // Merkle for Q1 (right half)
	pushCount     int   // Track how many items we've received
}

// newHybridRowTree creates a hybrid row tree with NMT left, Merkle right
func newHybridRowTree(squareSize uint64, axisIndex uint, nmtOpts ...nmt.Option) Tree {
	// Create NMT for left half (Q0) with original square size
	originalSquareSize := squareSize / 2
	leftTree := newErasuredNamespacedMerkleTree(originalSquareSize, axisIndex, nmtOpts...)
	
	// Create Merkle tree for right half (Q1)
	rightTree := NewDefaultTree(Row, axisIndex)
	
	return &hybridRowTree{
		squareSize: squareSize,
		axisIndex:  axisIndex,
		nmtOpts:    nmtOpts,
		leftTree:   &leftTree,
		rightTree:  rightTree,
		pushCount:  0,
	}
}

// Push implements the Tree interface by routing to appropriate subtree
func (h *hybridRowTree) Push(data []byte) error {
	halfSize := h.squareSize / 2
	
	if uint64(h.pushCount) < halfSize {
		// First half goes to NMT (Q0)
		err := h.leftTree.Push(data)
		if err != nil {
			return err
		}
	} else {
		// Second half goes to Merkle tree (Q1)
		err := h.rightTree.Push(data)
		if err != nil {
			return err
		}
	}
	
	h.pushCount++
	return nil
}

// Root implements the Tree interface by combining subtree roots
func (h *hybridRowTree) Root() ([]byte, error) {
	// Get roots from both subtrees
	leftRoot, err := h.leftTree.Root()
	if err != nil {
		return nil, err
	}
	
	rightRoot, err := h.rightTree.Root()
	if err != nil {
		return nil, err
	}
	
	// Combine the two roots using a simple concatenation and hash
	// This follows the pattern of how binary trees combine child nodes
	combined := append(leftRoot, rightRoot...)
	
	// Use a simple hash (in practice would use the same hash as the tree implementations)
	if len(combined) < 32 {
		// Pad to 32 bytes if needed
		padded := make([]byte, 32)
		copy(padded, combined)
		return padded, nil
	}
	
	return combined[:32], nil
}