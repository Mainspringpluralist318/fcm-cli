# 🚀 FCM CLI

⚡ Send Firebase push notifications from terminal in one command

🧩 With YAML config, profiles, `.env` support, `fcm init`, `--tokens-file`, and JSON output

> A lightweight, production-ready alternative to using Firebase Admin SDK directly when you want to send push notifications from terminal, scripts, CI pipelines, or backend jobs.

```bash
fcm -t TOKEN -n '{"title":"Hi","body":"Hello"}'
```

<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./assets/logo-light.png">
    <source media="(prefers-color-scheme: light)" srcset="./assets/logo-dark.png">
    <img src="./assets/logo-light.png" width="180" alt="FCM CLI logo" style="border-radius:50%;">
  </picture>
  <h3>Production-ready Firebase Cloud Messaging CLI for Go</h3>
</div>

<p align="center">
  <a href="https://github.com/interdev7/fcm-cli/releases">
    <img src="https://img.shields.io/github/v/release/interdev7/fcm-cli" />
  </a>
  <img src="https://img.shields.io/github/actions/workflow/status/interdev7/fcm-cli/release.yml" />
  <img src="https://img.shields.io/badge/Go-1.25-blue" />
  <img src="https://img.shields.io/badge/license-MPL%202.0-blue" />
</p>

---

# ✨ Overview

**FCM CLI** is a lightweight, production-ready command-line tool for sending Firebase Cloud Messaging notifications without having to wire Firebase SDK code into every project.

It is designed for:

- 🧑‍💻 Developers who want fast local testing
- ⚙️ DevOps / CI pipelines that need scriptable notification delivery
- 🧪 QA teams testing push behavior on multiple devices
- 🚀 Deployment workflows sending release, incident, or internal alerts
- 🌐 Backend services in any language that can execute shell commands

It works especially well when you want one consistent tool across:

- local development
- shell scripts
- GitHub Actions / CI
- cron jobs
- Node.js / Python / PHP / Java / C# backends

---

# ⚡ Why FCM CLI?

Instead of:

- opening Firebase Console
- pasting JSON manually
- writing temporary scripts again and again
- re-implementing the same notification logic in multiple languages

You can use one reusable command:

```bash
fcm -t TOKEN -n '{"title":"Hi","body":"Hello"}'
```

Or define repeatable workflows through YAML profiles:

```bash
fcm --config fcm.yaml --profile prod
```

---

# ✅ Features

- 🚀 Send FCM messages in seconds
- ⚡ Parallel multicast with progress bar
- 🔁 Built-in retry with exponential backoff
- 🎯 Supports token, topic, condition, and token files
- 📦 Zero dependencies at runtime (single binary)
- 🧠 Production-ready logging (color + JSON logs)
- 🧩 YAML config + profiles support
- 🌱 Easy onboarding with `fcm init`
- 🔐 `.env` support for local development
- 🤖 Structured JSON output for automation
- 🌐 Works from any language through CLI execution

---

# 🚀 Quick Start

If you just want to test one push as fast as possible:

```bash
export FCM_KEY=service-account.json

fcm -t TOKEN \
  -n '{"title":"Hello","body":"World"}'
```

That is enough to get started.

---

# 🔐 How to get service-account.json

To send messages, `fcm-cli` needs a Firebase service account key.

1.  Go to the [Firebase Console](https://console.firebase.google.com/).
2.  Select your project.
3.  Click the **Gear icon** (Project Settings) > **Service accounts**.
4.  Click **Generate new private key**.
5.  Save the JSON file as `service-account.json` in your project directory.

> [!IMPORTANT]
> Never commit this file to version control. Add it to your `.gitignore`.

---

# 📦 Installation

## ⚡ Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/interdev7/fcm-cli/main/install.sh | bash
```

This installs the latest release binary for your platform.

---

## 🍺 Homebrew

```bash
brew tap interdev7/fcm
brew install fcm
```

This is the easiest option on macOS if you prefer Homebrew-managed installs.

---

## 🛠 Manual Install from Source

```bash
git clone https://github.com/interdev7/fcm-cli.git
cd fcm-cli
go mod tidy
go build -o fcm ./cmd/fcm/main.go
```

Then run it locally:

```bash
./fcm -v
```

Optional system-wide install:

```bash
sudo mv fcm /usr/local/bin/fcm
```

---

## 🔍 Verify Installation

```bash
fcm -v
```

You should see the current version.

You can also check help:

```bash
fcm -h
```

---

# ⚙️ Configuration

FCM CLI supports configuration through:

- CLI flags
- YAML config file
- profiles inside YAML
- environment variables
- `.env` file

This makes it flexible for both local development and production automation.

---

## Environment Variables

Supported environment variables:

```bash
export FCM_KEY=service-account.json
export FCM_CONFIG=fcm.yaml
export FCM_LOG=info
```

### What they mean

- `FCM_KEY` — path to your Firebase service account JSON
- `FCM_CONFIG` — default YAML config file path
- `FCM_LOG` — default log mode (`info`, `debug`, or `json`)

---

## `.env` File Support

You can create a `.env` file in the project root:

```env
FCM_KEY=service-account.json
FCM_CONFIG=fcm.yaml
FCM_LOG=debug
```

FCM CLI loads `.env` automatically if it exists.

You can also load a custom env file:

```bash
fcm --env-file .env.local --profile prod
```

This is useful for:

- local overrides
- per-environment settings
- keeping shell setup minimal

---

# 🧩 YAML Configuration

FCM CLI supports reusable configuration through `fcm.yaml`.

This is ideal for:

- repeatable commands
- team workflows
- CI/CD
- multi-environment setups
- named profiles like `prod`, `staging`, `smoke`

---

## Generate Starter Config

```bash
fcm init
```

Overwrite an existing config:

```bash
fcm init --force
```

Generate to a custom path:

```bash
fcm init --file config/fcm.yaml
```

---

## Example `fcm.yaml`

```yaml
notification:
  title: Hello
  body: World

data:
  env: dev

log: info

profiles:
  prod:
    topic: production
    notification:
      title: Deploy
      body: New version released
    data:
      version: "1.2.0"

  smoke:
    tokens:
      - token1
      - token2
      - token3
    notification:
      title: Smoke test
      body: Batch notification test
```

---

## Using a Profile

```bash
fcm --config fcm.yaml --profile prod
```

This loads the `prod` profile from the config file and sends using its values.

---

## Overriding YAML with CLI Flags

CLI flags always win over YAML values.

Example:

```bash
fcm --config fcm.yaml --profile prod \
  -n '{"title":"Hotfix","body":"Immediate update"}'
```

That means:

- profile values are loaded
- then CLI overrides are applied on top

---

## Priority Order

Values are resolved in this order:

1. CLI flags
2. Profile (`--profile`)
3. Root config (`fcm.yaml`)
4. Environment variables / `.env`

This makes behavior predictable and easy to debug.

---

## Secrets

Do not store sensitive values in YAML if you can avoid it.

Recommended approach:

- keep `service-account.json` outside the repo
- use `FCM_KEY` via env or `.env`
- commit only non-secret config

Example:

```bash
export FCM_KEY=service-account.json
fcm --config fcm.yaml --profile prod
```

---

# 🚀 Usage

Basic form:

```bash
fcm [options]
```

Config generation:

```bash
fcm init [options]
```

---

# 🧩 Flags

| Flag                   | Description                         |
| ---------------------- | ----------------------------------- |
| `-k`, `--key`          | Firebase key file                   |
| `-t`, `--token`        | Single token                        |
| `--tokens`             | Multiple tokens (comma-separated)   |
| `--tokens-file`        | File with one token per line        |
| `-topic`               | Topic                               |
| `-c`, `--condition`    | Condition                           |
| `-n`, `--notification` | Notification JSON                   |
| `-d`, `--data`         | Data JSON                           |
| `-f`, `--config`       | YAML config file                    |
| `--profile`            | Profile from config                 |
| `--env-file`           | Load additional `.env` file         |
| `-l`, `--log`          | Log level (`info`, `debug`, `json`) |
| `--json`               | JSON output for automation          |
| `-v`, `--version`      | Version                             |
| `-h`, `--help`         | Help                                |

---

# 📱 Examples

## Send to One Device

```bash
fcm -t TOKEN -n '{"title":"Hi","body":"Hello"}'
```

---

## Specify Key File Directly

If you don't want to use environment variables, you can pass the key file using `-k` or `--key`:

```bash
fcm -k service-account.json \
  -t TOKEN \
  -n '{"title":"Test","body":"Using manual key file"}'
```

---

## Send to Multiple Tokens

```bash
fcm --tokens "t1,t2,t3" -n '{"title":"Batch","body":"Hello"}'
```

---

## Batch from File

Create `tokens.txt`:

```text
token1
token2
token3
```

Send:

```bash
fcm --tokens-file tokens.txt -n '{"title":"Batch","body":"Hello from file"}'
```

This is useful when:

- QA has a device list
- you export tokens from another system
- you want repeatable batch tests

---

## Send to Topic

```bash
fcm -topic news -n '{"title":"News","body":"Update"}'
```

---

## Send with Condition

```bash
fcm -c "'news' in topics" -n '{"title":"Cond","body":"Test"}'
```

---

## Send with Data Payload

```bash
fcm -t TOKEN \
  -n '{"title":"Chat","body":"New message"}' \
  -d '{"chatId":"123"}'
```

---

## Use Config File

```bash
fcm --config fcm.yaml --profile prod
```

---

## Use `.env` Defaults

If `.env` already contains `FCM_KEY` and `FCM_CONFIG`:

```bash
fcm --profile prod
```

---

## Debug Mode

```bash
fcm -l debug
```

Use this when you want more visibility into what the CLI is doing.

---

# 🤖 JSON Output (Important)

FCM CLI can return structured machine-readable output using `--json`.

This is critical for:

- CI/CD pipelines
- backend integrations
- automation scripts

---

## Single Message Success

```bash
fcm --profile prod --json
```

Example output:

```json
{
  "success": true,
  "message_id": "projects/.../messages/123"
}
```

This `message_id` is useful for:

- logging
- traceability
- downstream processing
- verifying successful delivery request submission to Firebase

---

## Error Output

Example:

```json
{
  "success": false,
  "error": "Invalid token"
}
```

This makes it easy to handle failures programmatically.

---

## Batch / Multicast JSON Output

When sending to many tokens, `--json` can return structured results for automation.

Example shape:

```json
{
  "success": false,
  "success_count": 2,
  "failure_count": 1,
  "results": [
    {
      "token": "token1",
      "success": true,
      "message_id": "projects/.../messages/111"
    },
    {
      "token": "token2",
      "success": false,
      "error": "Invalid registration token"
    }
  ]
}
```

This is especially useful if you want to:

- retry only failed tokens
- log exact failures
- remove invalid tokens from your system
- process results from Node.js / Python / CI

---

## Why JSON Output Matters

Without `--json`, CLI output is optimized for humans.

With `--json`, output is optimized for code.

That means you can safely parse it in backend services and automation pipelines.

---

# 📊 Multicast Example

Human-readable progress output:

```text
Progress: 3/5 (60%)
```

This is shown during batch sending when not using JSON automation mode.

---

# 🔁 Retry Strategy

FCM CLI retries failed requests automatically with exponential backoff:

| Attempt | Delay |
| ------- | ----- |
| 1       | 1s    |
| 2       | 2s    |
| 3       | 4s    |

This helps smooth over transient network or API issues.

---

# 🌍 Works with Any Language

Use `fcm` from any backend via CLI.
No Firebase SDK required.

This is one of the biggest strengths of the project:
you can keep a single notification workflow regardless of whether your backend is written in Go, Node.js, Python, PHP, C#, or Java.

---

## 🟦 Node.js

```js
const { execFile } = require("node:child_process");

execFile("fcm", ["--profile", "prod", "--json"], (err, stdout) => {
  const res = JSON.parse(stdout);

  if (!res.success) {
    console.error(res.error);
    return;
  }

  console.log("Message ID:", res.message_id);
});
```

---

## 🟢 Go

```go
package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("fcm", "--profile", "prod", "--json")
	out, _ := cmd.Output()

	var res map[string]interface{}
	json.Unmarshal(out, &res)

	fmt.Println(res)
}
```

---

## 🟣 Python

```python
import subprocess
import json

result = subprocess.run(
    ["fcm", "--profile", "prod", "--json"],
    capture_output=True,
    text=True
)

res = json.loads(result.stdout)

if not res["success"]:
    print("Error:", res["error"])
else:
    print("Message ID:", res.get("message_id"))
```

---

## 🟡 PHP

```php
<?php
$output = shell_exec("fcm --profile prod --json");
$res = json_decode($output, true);

if (!$res["success"]) {
    echo "Error: " . $res["error"];
} else {
    echo "Message ID: " . $res["message_id"];
}
```

---

## 🔵 C#

```csharp
using System.Diagnostics;
using System.Text.Json;

var process = new Process();
process.StartInfo.FileName = "fcm";
process.StartInfo.Arguments = "--profile prod --json";
process.StartInfo.RedirectStandardOutput = true;
process.Start();

string output = process.StandardOutput.ReadToEnd();

var res = JsonSerializer.Deserialize<dynamic>(output);

Console.WriteLine(res);
```

---

## ☕ Java

```java
import java.io.*;

ProcessBuilder pb = new ProcessBuilder("fcm", "--profile", "prod", "--json");
Process process = pb.start();

BufferedReader reader = new BufferedReader(
    new InputStreamReader(process.getInputStream())
);

String output = reader.readLine();
System.out.println(output);
```

---

# 📦 Releases

Download prebuilt binaries from:

👉 https://github.com/interdev7/fcm-cli/releases

Use release binaries if you want:

- fast setup
- reproducible installs
- easy CI / server provisioning
- no local Go toolchain dependency

---

# 🧠 Architecture

FCM CLI is built around a simple but practical model:

- Go + goroutines for concurrency
- OAuth2 service account authentication
- FCM HTTP v1 API
- Stateless CLI design
- YAML + `.env` driven local workflows
- Structured JSON result output for automation

This keeps the tool:

- small
- portable
- easy to integrate
- language-agnostic

---

# 🔐 Security

Recommended practices:

- never commit `service-account.json`
- use environment variables in production
- keep secrets outside version control
- rotate credentials when needed
- prefer per-environment `.env` or CI secrets

---

# 🤝 Contributing

PRs are welcome.

If you have ideas, bug reports, or feature requests:

- open an issue
- submit a PR
- share feedback from real-world FCM usage

Areas especially worth improving:

- templates
- rate limiting
- metrics
- dry-run mode
- richer batch tooling

---

# ⭐ Support

If you find this project useful, give it a star ⭐

That helps the project grow and makes it easier for other developers to discover it.

---

# ❤️ Sponsor

If FCM CLI saves you time in development, testing, or CI — consider supporting the project.

Your support helps with:

- maintenance
- new features
- documentation
- long-term stability

👉 https://github.com/sponsors/interdev7

---

# 📄 License

Licensed under the Mozilla Public License 2.0
See the [LICENSE](LICENSE) file for details.

---

# 📝 Changelog

See [CHANGELOG.md](CHANGELOG.md)
