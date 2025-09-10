#!/usr/bin/env python3

import csv
import subprocess
import sys
from collections import defaultdict
import statistics
from datetime import datetime
import os

def get_adaptive_iterations(k):
    """Get adaptive iteration count based on k value for better statistical confidence"""
    # Adaptive iterations: more for small k, fewer for large k
    if k <= 32:
        return 50
    elif k <= 64:
        return 30
    elif k <= 128:
        return 20
    elif k <= 256:
        return 15
    elif k <= 512:
        return 10
    elif k <= 1024:
        return 5
    else:
        return 3

def run_benchmarks(timeout_minutes=360, max_k=1024):
    """Run the benchmark suite and return the results file path"""
    # Create k-value specific benchmark regex
    k_values = [32, 64, 128, 256, 512, 1024, 2048]
    selected_k = [k for k in k_values if k <= max_k]
    
    if not selected_k:
        print(f"No valid k values <= {max_k}")
        return None
    
    # Calculate adaptive iterations for each k
    k_iterations = {k: get_adaptive_iterations(k) for k in selected_k}
    
    # Better time estimation based on k-value complexity and iterations
    time_estimates = {32: 0.02, 64: 0.17, 128: 1, 256: 5, 512: 20, 1024: 60, 2048: 240}  # minutes per run
    estimated_minutes = sum(time_estimates.get(k, k/4) * k_iterations[k] for k in selected_k)
    
    print(f"Starting benchmark run with adaptive iterations")
    print(f"Testing k values with iterations:")
    for k in selected_k:
        print(f"    k={k}: {k_iterations[k]} iterations")
    print(f"Estimated time: {estimated_minutes:.1f} minutes ({estimated_minutes/60:.1f} hours)")
    print()
    
    results_file = f"results/benchmark_results_{datetime.now().strftime('%Y%m%d_%H%M%S')}.txt"
    
    # Run benchmarks for each k value and tree type separately to show detailed progress
    tree_types = ["MerkleTree", "NMT", "Stage1", "Stage2", "Stage3"]
    total_combinations = len(selected_k) * len(tree_types)
    current_combo = 0
    
    try:
        with open(results_file, 'w') as f:
            for k in selected_k:
                print(f"Starting k={k} benchmarks...")
                
                for tree_type in tree_types:
                    current_combo += 1
                    iterations_for_k = k_iterations[k]
                    print(f"   {tree_type} ({current_combo}/{total_combinations}) [{iterations_for_k} runs]...", end=" ", flush=True)
                    
                    # Run specific tree type for this k value with adaptive iterations
                    bench_pattern = f"BenchmarkDatarootGeneration/{tree_type}_k={k}_"
                    cmd = [
                        "go", "test", 
                        f"-bench={bench_pattern}", 
                        "-benchmem", 
                        f"-count={iterations_for_k}",
                        f"-timeout={timeout_minutes}m"
                    ]
                    
                    result = subprocess.run(cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, 
                                          text=True, timeout=timeout_minutes*60, cwd="..")
                    
                    if result.returncode != 0:
                        print(f"FAILED")
                        print(f"Error details for {tree_type} k={k}:")
                        print("=" * 60)
                        print(result.stdout)
                        print("=" * 60)
                        f.write(f"# Failed for {tree_type} k={k}:\n{result.stdout}\n")
                        return None
                    
                    # Write results to file
                    f.write(result.stdout)
                    f.flush()  # Ensure data is written immediately
                    
                    # Count completed benchmarks for this tree type
                    benchmark_lines = [line for line in result.stdout.split('\n') if 'BenchmarkDatarootGeneration' in line and 'ns/op' in line]
                    
                    print(f"OK ({len(benchmark_lines)} runs)")
                
                print()  # Add spacing between k values
        
        print(f"All benchmarks completed. Results saved to {results_file}")
        return results_file
        
    except subprocess.TimeoutExpired:
        print(f"Benchmark timed out after {timeout_minutes} minutes")
        return None
    except Exception as e:
        print(f"Error running benchmarks: {e}")
        return None

def parse_benchmark_results(results_file):
    """Parse benchmark results into structured data"""
    data = defaultdict(list)
    
    with open(results_file, 'r') as f:
        for line in f:
            if 'BenchmarkDatarootGeneration' in line and 'ns/op' in line:
                parts = line.split()
                if len(parts) >= 7:
                    benchmark_name = parts[0]
                    ns_per_op = int(parts[2])
                    bytes_per_op = int(parts[4])
                    allocs_per_op = int(parts[6])
                    
                    # Extract k value and approach
                    if 'k=' in benchmark_name:
                        k_value = int(benchmark_name.split('k=')[1].split('_')[0])
                        
                        if 'MerkleTree_' in benchmark_name:
                            approach = 'MerkleTree'
                        elif 'NMT_' in benchmark_name:
                            approach = 'NMT'
                        elif 'Stage1_' in benchmark_name:
                            approach = 'Stage1'
                        elif 'Stage2_' in benchmark_name:
                            approach = 'Stage2'
                        elif 'Stage3_' in benchmark_name:
                            approach = 'Stage3'
                        else:
                            continue
                        
                        data[(k_value, approach)].append({
                            'time_ns': ns_per_op,
                            'memory_bytes': bytes_per_op,
                            'allocs': allocs_per_op
                        })
    
    return data

def calculate_averages(data):
    """Calculate averages and statistics from benchmark data"""
    averages = {}
    
    for (k_value, approach), measurements in data.items():
        if measurements:
            times = [m['time_ns'] for m in measurements]
            memories = [m['memory_bytes'] for m in measurements]
            allocs = [m['allocs'] for m in measurements]
            
            # Calculate both mean and median for better statistical insight
            time_mean = statistics.mean(times) / 1_000_000
            time_median = statistics.median(times) / 1_000_000
            time_std = statistics.stdev(times) / 1_000_000 if len(times) > 1 else 0
            
            # Calculate coefficient of variation (relative standard deviation)
            cv = (time_std / time_mean * 100) if time_mean > 0 else 0
            
            averages[(k_value, approach)] = {
                'time_ms': time_mean,
                'time_median': time_median,
                'time_std': time_std,
                'cv': cv,  # Coefficient of variation as percentage
                'memory_mb': statistics.mean(memories) / (1024 * 1024),
                'allocs': statistics.mean(allocs),
                'count': len(measurements)
            }
    
    return averages

def generate_ascii_charts(averages):
    """Generate ASCII bar charts for the results"""
    k_values = sorted(set(k for k, _ in averages.keys()))
    approaches = ['NMT', 'Stage1', 'Stage2', 'Stage3', 'MerkleTree']
    
    charts = []
    
    for k in k_values:
        chart_lines = [f"\nPerformance Chart - K={k} (EDS={k*2}x{k*2}):", "=" * 50]
        
        # Get data for this k value
        k_data = {}
        for approach in approaches:
            if (k, approach) in averages:
                k_data[approach] = averages[(k, approach)]
        
        if not k_data:
            continue
            
        # Find max time for scaling
        max_time = max(data['time_ms'] for data in k_data.values())
        scale = 40 / max_time if max_time > 0 else 1
        
        for approach in approaches:
            if approach in k_data:
                time_ms = k_data[approach]['time_ms']
                bar_length = int(time_ms * scale)
                speedup = ""
                
                if approach != 'NMT' and 'NMT' in k_data:
                    speedup_val = k_data['NMT']['time_ms'] / time_ms
                    speedup = f" ({speedup_val:.2f}x faster)"
                
                bar = "█" * bar_length
                chart_lines.append(f"{approach:<10}: {bar:<40} {time_ms:6.0f}ms{speedup}")
        
        charts.extend(chart_lines)
    
    return "\n".join(charts)

def generate_analysis_report(averages, results_file):
    """Generate the complete analysis report"""
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    report_file = f"results/benchmark_analysis_{datetime.now().strftime('%Y%m%d_%H%M%S')}.md"
    
    with open(report_file, 'w') as f:
        f.write(f"""# Benchmark Results

**Generated:** {timestamp}  
**Source:** {results_file.split('/')[-1]}

""")
        
        # Generate performance tables
        k_values = sorted(set(k for k, _ in averages.keys()))
        approaches = ['NMT', 'Stage1', 'Stage2', 'Stage3', 'MerkleTree']
        
        for k in k_values:
            f.write(f"### K={k} (EDS {k*2}×{k*2}):\n\n")
            f.write("| Approach | Time (ms) | Speedup | Memory (MB) | Runs | CV% |\n")
            f.write("|----------|-----------|---------|-------------|------|-----|\n")
            
            nmt_time = None
            for approach in approaches:
                if (k, approach) in averages:
                    data = averages[(k, approach)]
                    time_ms = data['time_ms']
                    memory_mb = data['memory_mb']
                    count = data['count']
                    std = data['time_std']
                    
                    if approach == 'NMT':
                        nmt_time = time_ms
                        speedup = "1.00×"
                    elif nmt_time:
                        speedup = f"{nmt_time / time_ms:.2f}×"
                    else:
                        speedup = "N/A"
                    
                    # Show median if significantly different from mean (high variance)
                    median = data.get('time_median', time_ms)
                    cv = data.get('cv', 0)
                    
                    if cv > 20:  # High variance threshold
                        time_str = f"{time_ms:.0f}(μ)/{median:.0f}(M) ±{std:.0f}"
                    else:
                        time_str = f"{time_ms:.0f}" + (f" ±{std:.0f}" if count > 1 and std > 0 else "")
                    
                    cv_str = f"{data.get('cv', 0):.1f}%" if data.get('cv', 0) > 0 else "N/A"
                    f.write(f"| {approach} | {time_str} | {speedup} | {memory_mb:.0f} | {count} | {cv_str} |\n")
            
            f.write("\n")
        
        # Add ASCII charts
        f.write("## Performance Visualization\n\n```\n")
        f.write(generate_ascii_charts(averages))
        f.write("\n```\n\n")
        
        f.write("\n")
    
    print(f"Report generated: {report_file}")
    return report_file

def main():
    if len(sys.argv) > 1 and sys.argv[1] in ['-h', '--help', 'help']:
        print("Usage: python3 benchmark.py [timeout_minutes] [max_k]")
        print("Example: python3 benchmark.py 360 1024")
        print("  timeout_minutes: timeout in minutes (default: 360)")  
        print("  max_k: maximum k value to test (default: 1024)")
        print("\nAdaptive iterations are used automatically:")
        print("  k=32: 50 iterations")
        print("  k=64: 30 iterations")
        print("  k=128: 20 iterations")
        print("  k=256: 15 iterations")
        print("  k=512: 10 iterations")
        print("  k=1024: 5 iterations")
        sys.exit(0)
    
    timeout_minutes = int(sys.argv[1]) if len(sys.argv) > 1 else 360
    max_k = int(sys.argv[2]) if len(sys.argv) > 2 else 1024
    
    print(f"Starting automated benchmark pipeline...")
    print(f"Configuration: {timeout_minutes}m timeout, max_k={max_k}")
    
    # Step 1: Run benchmarks
    results_file = run_benchmarks(timeout_minutes, max_k)
    if not results_file:
        print("Benchmark failed. Exiting.")
        sys.exit(1)
    
    # Step 2: Parse results
    print("Parsing benchmark results...")
    data = parse_benchmark_results(results_file)
    if not data:
        print("No valid benchmark data found. Exiting.")
        sys.exit(1)
    
    # Step 3: Calculate averages
    print("Calculating statistics...")
    averages = calculate_averages(data)
    
    # Step 4: Generate report
    print("Generating report...")
    report_file = generate_analysis_report(averages, results_file)
    
    print(f"\nPipeline completed successfully.")
    print(f"Report: {report_file}")
    print(f"Raw data: {results_file}")

if __name__ == "__main__":
    main()