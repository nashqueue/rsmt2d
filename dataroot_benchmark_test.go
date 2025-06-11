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
// Tests both regular Merkle trees and Namespaced Merkle Trees (NMT).
func BenchmarkDatarootGeneration(b *testing.B) {
	// Test k values: 32, 64, 128, 256, 512, 1024 (original data square size)
	// EDS will be 2k x 2k: 64, 128, 256, 512, 1024, 2048
	kValues := []int{32, 64, 128, 256, 512, 1024}
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
			// Generate random data for k x k original data square
			data := genRandomDataSquare(k, shareSize)
			
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				// Create EDS from original data
				eds, err := ComputeExtendedDataSquare(data, codec, NewDefaultTree)
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
		
		// Test with NMT (Namespaced Merkle Tree)
		b.Run(fmt.Sprintf("NMT_k=%d_EDS=%dx%d", k, k*2, k*2), func(b *testing.B) {
			// For NMT, we need to generate the full EDS (already extended) with sorted data
			// This is different from the erasure coding approach above
			edsSize := k * 2
			
			// Generate sorted namespaced data for the full EDS
			fullEDSData := genRandomNamespacedDataSquare(edsSize, shareSize, namespaceSize)
			
			// Create NMT tree constructor
			treeConstructor := newErasuredNamespacedMerkleTreeConstructor(
				uint64(edsSize),
				nmt.NamespaceIDSize(namespaceSize),
				nmt.IgnoreMaxNamespace(true),
				nmt.InitialCapacity(edsSize),
			)
			
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				// Import the pre-generated EDS (no erasure coding, just tree construction)
				eds, err := ImportExtendedDataSquare(fullEDSData, codec, treeConstructor)
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
		
		// Test with Hybrid: Merkle columns + NMT rows
		b.Run(fmt.Sprintf("Hybrid_k=%d_EDS=%dx%d", k, k*2, k*2), func(b *testing.B) {
			// For hybrid, we need to generate the full EDS with sorted data for rows
			edsSize := k * 2
			
			// Generate sorted namespaced data for the full EDS
			fullEDSData := genRandomNamespacedDataSquare(edsSize, shareSize, namespaceSize)
			
			// Create hybrid tree constructor (NMT for rows, Merkle for columns)
			treeConstructor := newHybridTreeConstructor(
				uint64(edsSize),
				nmt.NamespaceIDSize(namespaceSize),
				nmt.IgnoreMaxNamespace(true),
				nmt.InitialCapacity(edsSize),
			)
			
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				// Import the pre-generated EDS
				eds, err := ImportExtendedDataSquare(fullEDSData, codec, treeConstructor)
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

// genRandomNamespacedDataSquare generates a k x k square of namespaced data
// Each share has a namespace prefix followed by random data
// Data is sorted by namespace ID as required by NMT
func genRandomNamespacedDataSquare(k, shareSize, namespaceSize int) [][]byte {
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