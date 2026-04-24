# stackwatch

> A lightweight CLI tool for monitoring CloudFormation and Terraform stack drift in real time.

---

## Installation

```bash
go install github.com/yourusername/stackwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/stackwatch.git
cd stackwatch
go build -o stackwatch .
```

---

## Usage

Monitor a CloudFormation stack for drift:

```bash
stackwatch watch --provider cloudformation --stack my-production-stack
```

Monitor a Terraform workspace:

```bash
stackwatch watch --provider terraform --dir ./infra --interval 60s
```

Check drift status once without continuous monitoring:

```bash
stackwatch check --provider cloudformation --stack my-production-stack
```

### Flags

| Flag | Description | Default |
|------------|--------------------------------------|---------|
| `--provider` | Cloud provider (`cloudformation`, `terraform`) | required |
| `--stack` | Stack name (CloudFormation) | — |
| `--dir` | Terraform working directory | `.` |
| `--interval` | Polling interval | `5m` |
| `--output` | Output format (`text`, `json`) | `text` |

---

## Requirements

- Go 1.21+
- AWS credentials configured (for CloudFormation)
- Terraform CLI installed (for Terraform stacks)

---

## License

[MIT](LICENSE) © 2024 yourusername