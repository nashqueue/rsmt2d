# Automated Benchmark Analysis Pipeline

## Quick Start

```bash
# Test pipeline with k=64 max (fast)
python3 generate_analysis_report.py 1 10 64

# Full pipeline up to k=1024 (60 minutes)  
python3 generate_analysis_report.py 1 60 1024

# Full pipeline up to k=2048 (4+ hours)
python3 generate_analysis_report.py 1 300 2048
```

## Usage

```bash
python3 generate_analysis_report.py <runs> [timeout_minutes] [max_k]
```

**Parameters:**
- `runs`: Number of benchmark iterations (1-10 recommended)
- `timeout_minutes`: Timeout in minutes (default: 180)
- `max_k`: Maximum k value to test (default: 2048)
  - k=32: ~1 second
  - k=64: ~10 seconds  
  - k=128: ~1 minute
  - k=256: ~5 minutes
  - k=512: ~20 minutes
  - k=1024: ~60 minutes
  - k=2048: ~4+ hours

## What It Does

1. **Runs Benchmarks**: Executes Go benchmarks for all stages up to max_k
2. **Parses Results**: Extracts performance data from benchmark output
3. **Generates Report**: Creates comprehensive Markdown analysis with:
   - Performance tables for each k value
   - ASCII bar charts
   - Scaling analysis
   - Key findings
   - Implementation details

## Output Files

- `benchmark_results_YYYYMMDD_HHMMSS.txt` - Raw benchmark data
- `benchmark_analysis_YYYYMMDD_HHMMSS.md` - Complete analysis report

## Examples

**Quick test (2 minutes):**
```bash
python3 generate_analysis_report.py 1 10 64
```

**Full single run up to k=1024 (60 minutes):**
```bash
python3 generate_analysis_report.py 1 60 1024  
```

**Full single run up to k=2048 (4+ hours):**
```bash
python3 generate_analysis_report.py 1 300 2048
```

**Production analysis up to k=1024 (6+ hours):**
```bash
python3 generate_analysis_report.py 10 360 1024
```

**Production analysis up to k=2048 (40+ hours):**
```bash
python3 generate_analysis_report.py 10 2400 2048
```

The pipeline automatically handles benchmark execution, data parsing, statistical analysis, and report generation - no manual intervention required!