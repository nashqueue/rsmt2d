# Dataroot Generation Benchmark

This directory contains benchmarks for measuring the performance of generating Dataroots over Extended Data Squares (EDS) with varying original data square sizes (k values). The benchmark compares **regular Merkle trees** vs **Namespaced Merkle Trees (NMT)**.

## Overview

The benchmark tests Dataroot generation for k values that are powers of 2:
- **K = 32** → EDS = 64×64 (4,096 total shares)
- **K = 64** → EDS = 128×128 (16,384 total shares)  
- **K = 128** → EDS = 256×256 (65,536 total shares)
- **K = 256** → EDS = 512×512 (262,144 total shares)
- **K = 512** → EDS = 1024×1024 (1,048,576 total shares)
- **K = 1024** → EDS = 2048×2048 (4,194,304 total shares)

Where **k** is the width of the original data square, and the EDS is **2k × 2k** (doubled in each dimension).

### Tree Types Compared

1. **Merkle Tree**: Standard SHA-256 based Merkle tree with erasure coding
2. **NMT (Namespaced Merkle Tree)**: Namespace-aware trees used in Celestia for data availability sampling
3. **Hybrid**: NMT for rows + Merkle trees for columns (mixed approach)

## Files

### Core Benchmark
- `dataroot_benchmark_test.go` - The main benchmark function
- `run_dataroot_benchmark.sh` - Script to run benchmarks and generate reports

### Analysis & Visualization
- `analyze_benchmark_simple.py` - Comparison analysis for Merkle Tree vs NMT (no dependencies)
- `analyze_benchmark.py` - Original analysis script (no dependencies)
- `simple_benchmark_plot.py` - Matplotlib-based visualization (requires matplotlib)
- `create_benchmark_plot.py` - Advanced visualization (requires pandas, matplotlib, seaborn)

### Generated Files
- `benchmark_results.txt` - Raw benchmark output from Go
- `benchmark_data.csv` - Parsed benchmark data
- `benchmark_summary.csv` - Processed summary statistics
- `dataroot_benchmark_results.png` - Visualization (if Python packages available)

## Usage

### Quick Start
```bash
# Run the complete benchmark suite
./run_dataroot_benchmark.sh
```

### Manual Steps
```bash
# 1. Run just the benchmark
go test -bench=BenchmarkDatarootGeneration -benchmem -count=5

# 2. Run with analysis
go test -bench=BenchmarkDatarootGeneration -benchmem -count=5 > benchmark_results.txt
python3 analyze_benchmark.py
```

### Custom Benchmark
```bash
# Run with specific parameters
go test -bench=BenchmarkDatarootGeneration -benchtime=10s -count=3
```

## Results Summary

Based on the benchmark results comparing Merkle Tree vs NMT:

### Performance Comparison

| K Value | EDS Size   | Merkle (ms) | Hybrid (ms) | NMT (ms) | Merkle vs NMT | Hybrid vs NMT | Memory Hybrid | Memory NMT |
|---------|------------|-------------|-------------|----------|---------------|---------------|---------------|------------|
| 32      | 64×64      | 1.5         | 5.6         | 8.4      | **5.6x faster** | **1.5x faster** | 11.0 MB       | 20.1 MB    |
| 64      | 128×128    | 4.5         | 11.6        | 18.4     | **4.1x faster** | **1.6x faster** | 43.3 MB       | 80.1 MB    |
| 128     | 256×256    | 16.8        | 34.5        | 57.2     | **3.4x faster** | **1.7x faster** | 173.9 MB      | 320.8 MB   |
| 256     | 512×512    | 60.9        | 182.2       | 213.3    | **3.5x faster** | **1.2x faster** | 698.7 MB      | 1,286.8 MB |
| 512     | 1024×1024  | 442.7       | 708.5       | 1,354.2  | **3.1x faster** | **1.9x faster** | 2,824.0 MB    | 5,147.1 MB |
| 1024    | 2048×2048  | 1,173.2     | 2,893.5     | 3,497.3  | **3.0x faster** | **1.2x faster** | 11,189.0 MB   | 20,501.4 MB|

### Performance Characteristics

1. **NMT (Namespaced Merkle Tree) - Baseline**:
   - Full namespace functionality for both rows and columns
   - Highest memory usage due to namespace handling
   - No erasure coding (uses pre-generated EDS)
   - Most feature-complete implementation

2. **Hybrid (NMT rows + Merkle columns)**:
   - **1.2-1.9x faster than NMT** (best at k=512: 1.9x faster)
   - Retains namespace functionality for rows
   - **~45% memory savings** compared to full NMT
   - Good balance of performance and namespace features

3. **Merkle Tree (Maximum Performance)**:
   - **3.0-5.6x faster than NMT** (consistent 3-4x at large sizes)
   - Includes full erasure coding + root calculation
   - **~80% memory savings** compared to NMT
   - No namespace functionality but maximum efficiency

### ASCII Performance Comparison
```
K=32 (64x64)
  Merkle |                                                   1.4ms
  Hybrid |                                                   5.7ms
  NMT    |                                                   8.5ms

K=64 (128x128)
  Merkle |                                                   4.7ms
  Hybrid |                                                   12.1ms
  NMT    |                                                   17.6ms

K=128 (256x256)
  Merkle |                                                   15.9ms
  Hybrid |▒                                                  35.7ms
  NMT    |▓▓▓                                                60.1ms

K=256 (512x512)
  Merkle |███                                                62.1ms
  Hybrid |▒▒▒▒▒▒                                             135.4ms
  NMT    |▓▓▓▓▓▓▓▓▓▓                                         209.8ms

K=512 (1024x1024)
  Merkle |██████                                             442.7ms
  Hybrid |▒▒▒▒▒▒▒▒▒▒                                         708.5ms
  NMT    |▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓                                1354.2ms

K=1024 (2048x2048)
  Merkle |████████████████                                   1173.2ms
  Hybrid |▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒          2893.5ms
  NMT    |▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 3497.3ms
```
Legend: █ = Merkle Tree, ▒ = Hybrid (NMT rows + Merkle cols), ▓ = NMT

## Key Insights

### Hybrid Performance Sweet Spot
The **hybrid approach shows peak efficiency at k=512**, achieving **1.9x better performance** than full NMT while retaining namespace functionality for rows. This represents the optimal performance-to-namespace-functionality ratio.

### Scaling Characteristics (vs NMT Baseline)
- **Merkle Tree**: **3-5.6x faster** than NMT, consistent performance advantage across all sizes
- **Hybrid**: **1.2-1.9x faster** than NMT, best efficiency at k=512
- **NMT**: Full namespace functionality comes at significant performance cost but provides complete feature set

### Memory Implications
At **k=1024 (2048×2048 EDS)**:
- **NMT**: 20.5 GB memory usage (baseline)
- **Hybrid**: 11.2 GB memory usage (**45% memory savings**)
- **Merkle**: 3.8 GB memory usage (**81% memory savings**)

### Practical Recommendations (Performance vs Features)
- **Use NMT** when you need complete namespace functionality for both rows and columns
- **Use Hybrid** when namespace-aware row operations are sufficient and you want **~2x performance boost** 
- **Use Merkle** when namespace features aren't required and you need **~3-5x performance boost**

## Technical Details

### What is Measured

**Merkle Tree Benchmark**:
1. **EDS Construction**: Creating an Extended Data Square from random original data
2. **Erasure Coding**: Reed-Solomon encoding to fill the extended quadrants  
3. **Root Calculation**: Computing Merkle roots for all rows and columns (Dataroot)

**NMT Benchmark**:
1. **EDS Import**: Import pre-generated, namespace-sorted EDS data
2. **Root Calculation**: Computing namespaced Merkle roots for all rows and columns

**Hybrid Benchmark**:
1. **EDS Import**: Import pre-generated, namespace-sorted EDS data
2. **Root Calculation**: Computing NMT roots for rows and Merkle roots for columns

**Note**: The NMT and Hybrid benchmarks use `ImportExtendedDataSquare` with pre-sorted data rather than `ComputeExtendedDataSquare` because NMT requires lexicographically sorted namespace data, which is incompatible with the erasure coding process that reorganizes data.

### Implementation Notes
- Uses **LeoRS codec** (Leopard Reed-Solomon) for erasure coding
- **Share size**: 512 bytes (standard test size)
- **Tree construction**: Default Merkle tree implementation
- **Random data**: Cryptographically secure random data for each benchmark run

### Memory Allocation Patterns
The high allocation counts reflect:
- Share creation and copying during EDS construction
- Merkle tree node allocations for root calculation
- Intermediate buffers for Reed-Solomon encoding

## Reproducing Results

### Prerequisites
- Go 1.19+ with the rsmt2d module
- Python 3.x (optional, for visualization)

### Dependencies
The benchmark runs without external dependencies. For visualization:
```bash
# For basic plotting
pip install matplotlib

# For advanced analysis
pip install pandas matplotlib seaborn
```

### Environment
Results may vary based on:
- CPU architecture and performance
- Available memory
- Go version and optimization settings
- System load during benchmark execution

## Interpreting Results

### Performance Scaling
- **Good scaling** (close to theoretical O(n²)): Indicates efficient implementation
- **Worse than O(n²)**: May suggest memory pressure or algorithmic inefficiencies
- **Better than O(n²)**: Could indicate CPU cache effects at smaller sizes

### Memory Considerations
- Memory usage grows quadratically with EDS size
- For large k values (512+), ensure sufficient RAM is available
- Memory allocations scale with the total number of operations

This benchmark provides insights into the practical performance characteristics of Dataroot generation for different data square sizes, useful for capacity planning and performance optimization.