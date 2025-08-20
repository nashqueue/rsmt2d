# EDS Dataroot Generation Benchmark Analysis

**Generated:** 2025-06-17 14:10:25  
**Benchmark Runs:** 5 iterations per test  
**Source Data:** benchmark_results_20250617_135913.txt

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

| Approach | Time (ms) | Speedup vs NMT | Memory (MB) | Runs | Description |
|----------|-----------|-----------------|-------------|------|-------------|
| NMT | 3 ±0 | 1.00× | 22 | 5 | All trees use NMT |
| Stage1 | 2 ±0 | 1.33× | 13 | 5 | Merkle columns, NMT rows |
| Stage2 | 2 ±0 | 1.52× | 9 | 5 | + Merkle for Q2/Q3 rows |
| Stage3 | 2 ±0 | 1.69× | 7 | 5 | + Hybrid row trees (Q0:NMT, Q1:Merkle) |
| MerkleTree | 1 ±0 | 2.11× | 4 | 5 | All trees use Merkle |

### K=64 (EDS 128×128):

| Approach | Time (ms) | Speedup vs NMT | Memory (MB) | Runs | Description |
|----------|-----------|-----------------|-------------|------|-------------|
| NMT | 11 ±0 | 1.00× | 91 | 5 | All trees use NMT |
| Stage1 | 7 ±0 | 1.53× | 53 | 5 | Merkle columns, NMT rows |
| Stage2 | 6 ±0 | 1.88× | 35 | 5 | + Merkle for Q2/Q3 rows |
| Stage3 | 5 ±0 | 2.11× | 26 | 5 | + Hybrid row trees (Q0:NMT, Q1:Merkle) |
| MerkleTree | 4 ±0 | 2.50× | 16 | 5 | All trees use Merkle |

### K=128 (EDS 256×256):

| Approach | Time (ms) | Speedup vs NMT | Memory (MB) | Runs | Description |
|----------|-----------|-----------------|-------------|------|-------------|
| NMT | 48 ±1 | 1.00× | 397 | 5 | All trees use NMT |
| Stage1 | 29 ±1 | 1.65× | 218 | 5 | Merkle columns, NMT rows |
| Stage2 | 22 ±0 | 2.17× | 144 | 5 | + Merkle for Q2/Q3 rows |
| Stage3 | 20 ±1 | 2.44× | 109 | 5 | + Hybrid row trees (Q0:NMT, Q1:Merkle) |
| MerkleTree | 15 ±0 | 3.11× | 68 | 5 | All trees use Merkle |

### K=256 (EDS 512×512):

| Approach | Time (ms) | Speedup vs NMT | Memory (MB) | Runs | Description |
|----------|-----------|-----------------|-------------|------|-------------|
| NMT | 336 ±8 | 1.00× | 1347 | 5 | All trees use NMT |
| Stage1 | 208 ±7 | 1.61× | 786 | 5 | Merkle columns, NMT rows |
| Stage2 | 149 ±3 | 2.25× | 509 | 5 | + Merkle for Q2/Q3 rows |
| Stage3 | 128 ±3 | 2.63× | 376 | 5 | + Hybrid row trees (Q0:NMT, Q1:Merkle) |
| MerkleTree | 65 ±18 | 5.19× | 225 | 5 | All trees use Merkle |

### K=512 (EDS 1024×1024):

| Approach | Time (ms) | Speedup vs NMT | Memory (MB) | Runs | Description |
|----------|-----------|-----------------|-------------|------|-------------|
| NMT | 1093 ±335 | 1.00× | 5506 | 5 | All trees use NMT |
| Stage1 | 575 ±14 | 1.90× | 3181 | 5 | Merkle columns, NMT rows |
| Stage2 | 405 ±9 | 2.70× | 2095 | 5 | + Merkle for Q2/Q3 rows |
| Stage3 | 451 ±131 | 2.43× | 1561 | 5 | + Hybrid row trees (Q0:NMT, Q1:Merkle) |
| MerkleTree | 403 ±13 | 2.72× | 968 | 5 | All trees use Merkle |

### K=1024 (EDS 2048×2048):

| Approach | Time (ms) | Speedup vs NMT | Memory (MB) | Runs | Description |
|----------|-----------|-----------------|-------------|------|-------------|
| NMT | 4743 ±892 | 1.00× | 21871 | 5 | All trees use NMT |
| Stage1 | 2535 ±112 | 1.87× | 13031 | 5 | Merkle columns, NMT rows |
| Stage2 | 2201 ±483 | 2.16× | 8698 | 5 | + Merkle for Q2/Q3 rows |
| Stage3 | 2231 ±94 | 2.13× | 6589 | 5 | + Hybrid row trees (Q0:NMT, Q1:Merkle) |
| MerkleTree | 1751 ±136 | 2.71× | 4207 | 5 | All trees use Merkle |

## Performance Visualization

```

Performance Chart - K=32 (EDS=64x64):
==================================================
NMT       : ████████████████████████████████████████      3ms
Stage1    : ██████████████████████████████                2ms (1.33x faster)
Stage2    : ██████████████████████████                    2ms (1.52x faster)
Stage3    : ███████████████████████                       2ms (1.69x faster)
MerkleTree: ██████████████████                            1ms (2.11x faster)

Performance Chart - K=64 (EDS=128x128):
==================================================
NMT       : ████████████████████████████████████████     11ms
Stage1    : ██████████████████████████                    7ms (1.53x faster)
Stage2    : █████████████████████                         6ms (1.88x faster)
Stage3    : ██████████████████                            5ms (2.11x faster)
MerkleTree: ███████████████                               4ms (2.50x faster)

Performance Chart - K=128 (EDS=256x256):
==================================================
NMT       : ████████████████████████████████████████     48ms
Stage1    : ████████████████████████                     29ms (1.65x faster)
Stage2    : ██████████████████                           22ms (2.17x faster)
Stage3    : ████████████████                             20ms (2.44x faster)
MerkleTree: ████████████                                 15ms (3.11x faster)

Performance Chart - K=256 (EDS=512x512):
==================================================
NMT       : ████████████████████████████████████████    336ms
Stage1    : ████████████████████████                    208ms (1.61x faster)
Stage2    : █████████████████                           149ms (2.25x faster)
Stage3    : ███████████████                             128ms (2.63x faster)
MerkleTree: ███████                                      65ms (5.19x faster)

Performance Chart - K=512 (EDS=1024x1024):
==================================================
NMT       : ████████████████████████████████████████   1093ms
Stage1    : █████████████████████                       575ms (1.90x faster)
Stage2    : ██████████████                              405ms (2.70x faster)
Stage3    : ████████████████                            451ms (2.43x faster)
MerkleTree: ██████████████                              403ms (2.72x faster)

Performance Chart - K=1024 (EDS=2048x2048):
==================================================
NMT       : ████████████████████████████████████████   4743ms
Stage1    : █████████████████████                      2535ms (1.87x faster)
Stage2    : ██████████████████                         2201ms (2.16x faster)
Stage3    : ██████████████████                         2231ms (2.13x faster)
MerkleTree: ██████████████                             1751ms (2.71x faster)
```

## Analysis

### Stage 3 Scaling Analysis

Stage 3 performance vs NMT baseline:
- K=32: 1.69× faster than NMT
- K=64: 2.11× faster than NMT
- K=128: 2.44× faster than NMT
- K=256: 2.63× faster than NMT
- K=512: 2.43× faster than NMT
- K=1024: 2.13× faster than NMT

**Scaling Trend**: Stage 3 efficiency is improving with larger EDS sizes.

### Methodology

- **Environment**: Go benchmark framework with `-benchmem` flag
- **Iterations**: 5 runs per test for statistical reliability
- **Timeout**: Extended timeout for large EDS sizes
- **Validation**: Multiple measurement points for consistency

---
*Generated by automated benchmark analysis pipeline*
