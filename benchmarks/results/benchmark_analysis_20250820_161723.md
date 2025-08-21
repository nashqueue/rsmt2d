# EDS Dataroot Generation Benchmark Analysis

**Generated:** 2025-08-20 16:17:23  
**Benchmark Runs:** 1 iterations per test  
**Source Data:** benchmark_results_20250820_153611.txt

## Executive Summary

This report analyzes the performance of different tree construction approaches for Extended Data Square (EDS) dataroot generation across various sizes.

## EDS Quadrant Layout Reference

```
Q0 | Q1    (Original data | Row parity)
---+---
Q2 | Q3    (Column parity | Intersection parity)  
```

**Stage Optimizations:**
- **Stage 1**: Merkle columns, NMT rows
- **Stage 2**: + Merkle for Q2/Q3 row roots (bottom half)
- **Stage 3**: + Hybrid row trees (Merkle for Q1, NMT for Q0)

## Performance Results

### K=32 (EDS 64×64):

| Approach | Time (ms) | Speedup vs NMT | Memory (MB) | Runs | CV% | Description |
|----------|-----------|-----------------|-------------|------|-----|-------------|
| NMT | 6 ±0 | 1.00× | 23 | 50 | 3.0% | All trees use NMT |
| Stage1 | 5 ±0 | 1.31× | 14 | 50 | 3.2% | Merkle columns, NMT rows |
| Stage2 | 4 ±0 | 1.50× | 9 | 50 | 3.5% | + Merkle for Q2/Q3 rows |
| Stage3 | 4 ±0 | 1.63× | 7 | 50 | 2.8% | + Hybrid row trees (Q0:NMT, Q1:Merkle) |
| MerkleTree | 3 ±0 | 1.89× | 4 | 50 | 1.4% | All trees use Merkle |

### K=64 (EDS 128×128):

| Approach | Time (ms) | Speedup vs NMT | Memory (MB) | Runs | CV% | Description |
|----------|-----------|-----------------|-------------|------|-----|-------------|
| NMT | 19 ±1 | 1.00× | 96 | 30 | 5.7% | All trees use NMT |
| Stage1 | 14 ±1 | 1.35× | 54 | 30 | 6.3% | Merkle columns, NMT rows |
| Stage2 | 12 ±0 | 1.52× | 37 | 30 | 2.8% | + Merkle for Q2/Q3 rows |
| Stage3 | 11 ±0 | 1.68× | 28 | 30 | 2.5% | + Hybrid row trees (Q0:NMT, Q1:Merkle) |
| MerkleTree | 9 ±0 | 2.02× | 18 | 30 | 2.2% | All trees use Merkle |

### K=128 (EDS 256×256):

| Approach | Time (ms) | Speedup vs NMT | Memory (MB) | Runs | CV% | Description |
|----------|-----------|-----------------|-------------|------|-----|-------------|
| NMT | 76 ±8 | 1.00× | 410 | 20 | 10.9% | All trees use NMT |
| Stage1 | 46 ±3 | 1.66× | 227 | 20 | 5.9% | Merkle columns, NMT rows |
| Stage2 | 35 ±4 | 2.17× | 149 | 20 | 10.6% | + Merkle for Q2/Q3 rows |
| Stage3 | 29 ±1 | 2.61× | 111 | 20 | 4.4% | + Hybrid row trees (Q0:NMT, Q1:Merkle) |
| MerkleTree | 25 ±1 | 3.08× | 72 | 20 | 3.0% | All trees use Merkle |

### K=256 (EDS 512×512):

| Approach | Time (ms) | Speedup vs NMT | Memory (MB) | Runs | CV% | Description |
|----------|-----------|-----------------|-------------|------|-----|-------------|
| NMT | 332 ±27 | 1.00× | 1348 | 15 | 8.1% | All trees use NMT |
| Stage1 | 190 ±7 | 1.75× | 787 | 15 | 3.6% | Merkle columns, NMT rows |
| Stage2 | 132 ±3 | 2.51× | 509 | 15 | 2.6% | + Merkle for Q2/Q3 rows |
| Stage3 | 109 ±2 | 3.05× | 377 | 15 | 2.0% | + Hybrid row trees (Q0:NMT, Q1:Merkle) |
| MerkleTree | 84 ±1 | 3.95× | 226 | 15 | 1.7% | All trees use Merkle |

### K=512 (EDS 1024×1024):

| Approach | Time (ms) | Speedup vs NMT | Memory (MB) | Runs | CV% | Description |
|----------|-----------|-----------------|-------------|------|-----|-------------|
| NMT | 1265 ±96 | 1.00× | 5453 | 10 | 7.5% | All trees use NMT |
| Stage1 | 795 ±24 | 1.59× | 3181 | 10 | 3.0% | Merkle columns, NMT rows |
| Stage2 | 599 ±31 | 2.11× | 2104 | 10 | 5.1% | + Merkle for Q2/Q3 rows |
| Stage3 | 504 ±18 | 2.51× | 1564 | 10 | 3.5% | + Hybrid row trees (Q0:NMT, Q1:Merkle) |
| MerkleTree | 365 ±11 | 3.47× | 970 | 10 | 3.0% | All trees use Merkle |

### K=1024 (EDS 2048×2048):

| Approach | Time (ms) | Speedup vs NMT | Memory (MB) | Runs | CV% | Description |
|----------|-----------|-----------------|-------------|------|-----|-------------|
| NMT | 5152 ±265 | 1.00× | 21878 | 5 | 5.1% | All trees use NMT |
| Stage1 | 3298 ±192 | 1.56× | 13039 | 5 | 5.8% | Merkle columns, NMT rows |
| Stage2 | 2462 ±95 | 2.09× | 8699 | 5 | 3.9% | + Merkle for Q2/Q3 rows |
| Stage3 | 2113 ±98 | 2.44× | 6593 | 5 | 4.7% | + Hybrid row trees (Q0:NMT, Q1:Merkle) |
| MerkleTree | 1682 ±161 | 3.06× | 4210 | 5 | 9.6% | All trees use Merkle |

## Performance Visualization

```

Performance Chart - K=32 (EDS=64x64):
==================================================
NMT       : ████████████████████████████████████████      6ms
Stage1    : ██████████████████████████████                5ms (1.31x faster)
Stage2    : ██████████████████████████                    4ms (1.50x faster)
Stage3    : ████████████████████████                      4ms (1.63x faster)
MerkleTree: █████████████████████                         3ms (1.89x faster)

Performance Chart - K=64 (EDS=128x128):
==================================================
NMT       : ████████████████████████████████████████     19ms
Stage1    : █████████████████████████████                14ms (1.35x faster)
Stage2    : ██████████████████████████                   12ms (1.52x faster)
Stage3    : ███████████████████████                      11ms (1.68x faster)
MerkleTree: ███████████████████                           9ms (2.02x faster)

Performance Chart - K=128 (EDS=256x256):
==================================================
NMT       : ████████████████████████████████████████     76ms
Stage1    : ████████████████████████                     46ms (1.66x faster)
Stage2    : ██████████████████                           35ms (2.17x faster)
Stage3    : ███████████████                              29ms (2.61x faster)
MerkleTree: ████████████                                 25ms (3.08x faster)

Performance Chart - K=256 (EDS=512x512):
==================================================
NMT       : ████████████████████████████████████████    332ms
Stage1    : ██████████████████████                      190ms (1.75x faster)
Stage2    : ███████████████                             132ms (2.51x faster)
Stage3    : █████████████                               109ms (3.05x faster)
MerkleTree: ██████████                                   84ms (3.95x faster)

Performance Chart - K=512 (EDS=1024x1024):
==================================================
NMT       : ████████████████████████████████████████   1265ms
Stage1    : █████████████████████████                   795ms (1.59x faster)
Stage2    : ██████████████████                          599ms (2.11x faster)
Stage3    : ███████████████                             504ms (2.51x faster)
MerkleTree: ███████████                                 365ms (3.47x faster)

Performance Chart - K=1024 (EDS=2048x2048):
==================================================
NMT       : ████████████████████████████████████████   5152ms
Stage1    : █████████████████████████                  3298ms (1.56x faster)
Stage2    : ███████████████████                        2462ms (2.09x faster)
Stage3    : ████████████████                           2113ms (2.44x faster)
MerkleTree: █████████████                              1682ms (3.06x faster)
```

## Analysis

### Stage 3 Scaling Analysis

Stage 3 performance vs NMT baseline:
- K=32: 1.63× faster than NMT
- K=64: 1.68× faster than NMT
- K=128: 2.61× faster than NMT
- K=256: 3.05× faster than NMT
- K=512: 2.51× faster than NMT
- K=1024: 2.44× faster than NMT

**Scaling Trend**: Stage 3 efficiency is improving with larger EDS sizes.

### Methodology

- **Environment**: Go benchmark framework with `-benchmem` flag
- **Iterations**: Adaptive per k-value (50 for k=32, down to 3-5 for k=2048)
- **Statistics**: Mean (μ), Median (M), Standard Deviation, CV% (Coefficient of Variation)
- **High Variance**: When CV% > 20%, both mean and median shown as "mean(μ)/median(M)"
- **Timeout**: Extended timeout for large EDS sizes
- **Validation**: Multiple measurement points for consistency

---
*Generated by automated benchmark analysis pipeline*
