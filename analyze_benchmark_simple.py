#!/usr/bin/env python3
"""
Simple benchmark analysis script for both Merkle Tree and NMT comparisons.
"""

import csv
import statistics

def analyze_benchmarks():
    print("="*80)
    print("DATAROOT GENERATION BENCHMARK ANALYSIS")
    print("NMT Baseline: Comparing efficiency vs full Namespaced Merkle Tree")
    print("="*80)
    
    # Read and organize data
    merkle_data = {}
    nmt_data = {}
    hybrid_data = {}
    
    try:
        with open('benchmark_data.csv', 'r') as f:
            reader = csv.DictReader(f)
            for row in reader:
                benchmark_name = row.get('benchmark_name', '')
                k = int(row['k_value'])
                
                data_point = {
                    'ns_per_op': float(row['ns_per_op']),
                    'bytes_per_op': float(row['bytes_per_op']) if row['bytes_per_op'] else 0,
                    'allocs_per_op': float(row['allocs_per_op']) if row['allocs_per_op'] else 0,
                    'eds_size': row['eds_size']
                }
                
                if 'MerkleTree' in benchmark_name:
                    if k not in merkle_data:
                        merkle_data[k] = []
                    merkle_data[k].append(data_point)
                elif 'NMT' in benchmark_name:
                    if k not in nmt_data:
                        nmt_data[k] = []
                    nmt_data[k].append(data_point)
                elif 'Hybrid' in benchmark_name:
                    if k not in hybrid_data:
                        hybrid_data[k] = []
                    hybrid_data[k].append(data_point)
                    
    except FileNotFoundError:
        print("Error: benchmark_data.csv not found")
        return
    
    # Calculate averages for each k value and tree type
    merkle_summary = {}
    nmt_summary = {}
    hybrid_summary = {}
    
    for k in sorted(set(merkle_data.keys()) | set(nmt_data.keys()) | set(hybrid_data.keys())):
        if k in merkle_data:
            measurements = merkle_data[k]
            merkle_summary[k] = {
                'avg_ms': statistics.mean([m['ns_per_op'] for m in measurements]) / 1_000_000,
                'avg_mb': statistics.mean([m['bytes_per_op'] for m in measurements]) / (1024 * 1024),
                'avg_allocs': int(statistics.mean([m['allocs_per_op'] for m in measurements])),
                'eds_size': measurements[0]['eds_size']
            }
        
        if k in nmt_data:
            measurements = nmt_data[k]
            nmt_summary[k] = {
                'avg_ms': statistics.mean([m['ns_per_op'] for m in measurements]) / 1_000_000,
                'avg_mb': statistics.mean([m['bytes_per_op'] for m in measurements]) / (1024 * 1024),
                'avg_allocs': int(statistics.mean([m['allocs_per_op'] for m in measurements])),
                'eds_size': measurements[0]['eds_size']
            }
        
        if k in hybrid_data:
            measurements = hybrid_data[k]
            hybrid_summary[k] = {
                'avg_ms': statistics.mean([m['ns_per_op'] for m in measurements]) / 1_000_000,
                'avg_mb': statistics.mean([m['bytes_per_op'] for m in measurements]) / (1024 * 1024),
                'avg_allocs': int(statistics.mean([m['allocs_per_op'] for m in measurements])),
                'eds_size': measurements[0]['eds_size']
            }
    
    # Print comparison table
    print("\nPERFORMANCE COMPARISON:")
    print("-" * 120)
    print(f"{'K Value':<8} {'EDS Size':<10} {'Merkle (ms)':<12} {'Hybrid (ms)':<12} {'NMT (ms)':<12} {'Merkle vs NMT':<15} {'Hybrid vs NMT':<15}")
    print("-" * 120)
    
    for k in sorted(set(merkle_summary.keys()) | set(nmt_summary.keys()) | set(hybrid_summary.keys())):
        merkle_time = merkle_summary.get(k, {}).get('avg_ms', 0)
        nmt_time = nmt_summary.get(k, {}).get('avg_ms', 0)
        hybrid_time = hybrid_summary.get(k, {}).get('avg_ms', 0)
        eds_size = merkle_summary.get(k, nmt_summary.get(k, hybrid_summary.get(k, {}))).get('eds_size', 'N/A')
        
        merkle_str = f"{merkle_time:.1f}" if merkle_time > 0 else "N/A"
        hybrid_str = f"{hybrid_time:.1f}" if hybrid_time > 0 else "N/A"
        nmt_str = f"{nmt_time:.1f}" if nmt_time > 0 else "N/A"
        
        merkle_vs_nmt_str = f"{nmt_time/merkle_time:.1f}x faster" if merkle_time > 0 and nmt_time > 0 else "N/A"
        hybrid_vs_nmt_str = f"{nmt_time/hybrid_time:.1f}x faster" if hybrid_time > 0 and nmt_time > 0 else "N/A"
        
        print(f"{k:<8} {eds_size:<10} {merkle_str:<12} {hybrid_str:<12} {nmt_str:<12} {merkle_vs_nmt_str:<15} {hybrid_vs_nmt_str:<15}")
    
    # Detailed breakdown
    print("\nDETAILED BREAKDOWN:")
    print("-" * 100)
    
    print("\nMERKLE TREE RESULTS:")
    print(f"{'K Value':<8} {'EDS Size':<10} {'Time (ms)':<12} {'Memory (MB)':<12} {'Allocations':<12}")
    print("-" * 60)
    for k in sorted(merkle_summary.keys()):
        data = merkle_summary[k]
        print(f"{k:<8} {data['eds_size']:<10} {data['avg_ms']:<12.1f} {data['avg_mb']:<12.1f} {data['avg_allocs']:<12}")
    
    print("\nNMT (NAMESPACED MERKLE TREE) RESULTS:")
    print(f"{'K Value':<8} {'EDS Size':<10} {'Time (ms)':<12} {'Memory (MB)':<12} {'Allocations':<12}")
    print("-" * 60)
    for k in sorted(nmt_summary.keys()):
        data = nmt_summary[k]
        print(f"{k:<8} {data['eds_size']:<10} {data['avg_ms']:<12.1f} {data['avg_mb']:<12.1f} {data['avg_allocs']:<12}")
    
    print("\nHYBRID (NMT ROWS + MERKLE COLUMNS) RESULTS:")
    print(f"{'K Value':<8} {'EDS Size':<10} {'Time (ms)':<12} {'Memory (MB)':<12} {'Allocations':<12}")
    print("-" * 60)
    for k in sorted(hybrid_summary.keys()):
        data = hybrid_summary[k]
        print(f"{k:<8} {data['eds_size']:<10} {data['avg_ms']:<12.1f} {data['avg_mb']:<12.1f} {data['avg_allocs']:<12}")
    
    # Performance overhead analysis
    print("\nEFFICIENCY ANALYSIS (vs NMT Baseline):")
    print("-" * 80)
    print("Comparing performance gains relative to the NMT baseline implementation.")
    print()
    
    merkle_keys = sorted(merkle_summary.keys())
    nmt_keys = sorted(nmt_summary.keys())
    hybrid_keys = sorted(hybrid_summary.keys())
    common_keys = sorted(set(merkle_keys) & set(nmt_keys) & set(hybrid_keys))
    
    if len(common_keys) >= 2:
        print("Performance improvements vs NMT:")
        for i, k in enumerate(common_keys):
            merkle_time = merkle_summary[k]['avg_ms']
            nmt_time = nmt_summary[k]['avg_ms']
            hybrid_time = hybrid_summary[k]['avg_ms']
            
            merkle_improvement = nmt_time / merkle_time
            hybrid_improvement = nmt_time / hybrid_time
            
            print(f"K={k:3d}: Merkle {merkle_improvement:.1f}x faster ({merkle_time:.1f}ms), Hybrid {hybrid_improvement:.1f}x faster ({hybrid_time:.1f}ms) vs NMT ({nmt_time:.1f}ms)")
    
    # ASCII comparison chart
    print("\nASCII PERFORMANCE COMPARISON:")
    print("-" * 80)
    
    max_time = 0
    for k in common_keys:
        max_time = max(max_time, merkle_summary[k]['avg_ms'], nmt_summary[k]['avg_ms'], hybrid_summary[k]['avg_ms'])
    
    chart_width = 50
    for k in common_keys:
        merkle_time = merkle_summary[k]['avg_ms']
        nmt_time = nmt_summary[k]['avg_ms']
        hybrid_time = hybrid_summary[k]['avg_ms']
        eds_size = merkle_summary[k]['eds_size']
        
        merkle_bar_length = int((merkle_time / max_time) * chart_width)
        nmt_bar_length = int((nmt_time / max_time) * chart_width)
        hybrid_bar_length = int((hybrid_time / max_time) * chart_width)
        
        merkle_bar = '█' * merkle_bar_length
        nmt_bar = '▓' * nmt_bar_length
        hybrid_bar = '▒' * hybrid_bar_length
        
        print(f"K={k} ({eds_size})")
        print(f"  Merkle |{merkle_bar:<{chart_width}} {merkle_time:.1f}ms")
        print(f"  Hybrid |{hybrid_bar:<{chart_width}} {hybrid_time:.1f}ms")
        print(f"  NMT    |{nmt_bar:<{chart_width}} {nmt_time:.1f}ms")
        print()
    
    print("Legend: █ = Merkle Tree, ▒ = Hybrid (NMT rows + Merkle cols), ▓ = NMT")
    print("="*80)

if __name__ == "__main__":
    analyze_benchmarks()