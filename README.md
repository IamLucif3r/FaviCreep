
![Favicreep](https://github.com/user-attachments/assets/37e7b040-93fc-41a6-86fc-6d55a48fcf49)


# FaviCreep

> "Every forgotten panel tells a story... and leaves a door wide open."  
> — Unknown Bug Bounty Hunter

---

## 🎯 What is FaviCreep?

**FaviCreep** is a CLI tool built for **bug bounty hunters**, **red teamers**, and **offensive security analysts**. It helps you **discover shadow assets**, **exposed portals**, and **forgotten infrastructure** by analyzing something often overlooked:

> 🔥 **Favicons** — those tiny images your browser loads, even on staging servers.

---

## 💡 Core Idea

Many companies reuse the **same favicon** across environments: dev, staging, prod, internal tools, CI/CD dashboards, forgotten test portals.

FaviCreep leverages this quirk:

1. 🔍 **Enumerate subdomains** of a target domain
2. 🎨 **Fetch favicons** from each subdomain
3. 🧠 **Hash them** using mmh3 (like Shodan does)
4. 📦 **Cluster subdomains** sharing the same hash (same branding = same app)
5. 🌐 **Search Shodan** for **external systems** using that same hash

Result?  
A goldmine of attack surface you were never supposed to see.

---

## 🚀 Features

- 🔄 Concurrency for fast favicon fetching
- 🧮 mmh3 hashing like Shodan
- 🧠 Clustering by favicon hash
- 🌍 Shodan integration (via API) to find global assets using the same favicon
- 🧾 JSON export of clustered data
- 🔧 Uses [Subfinder](https://github.com/projectdiscovery/subfinder) for subdomain enumeration

---

## 🧱 Installation
### Use the tool directly

```bash
go install github.com/iamlucif3r/favicreep@latest
```
**Requirements:**

- Go 1.20+
- Shodan API key (for Shodan module)
- [Subfinder](https://github.com/projectdiscovery/subfinder) installed and in your `$PATH`

```bash
git clone https://github.com/iamlucif3r/favicreep.git
cd favicreep
go build -o favicreep main.go
```

## 🧪 Usage

#### 1. Scan a domain and cluster favicon hashes

```bash
./favicreep scan --domain example.com

Options:
  -c, --concurrency int   Max concurrent favicon fetches (default: 10)
  -o, --output string     Save clustered result to a JSON file
Example:

./favicreep scan --domain hackerone.com -c 20 -o hackerone_clusters.json
```
##### NOTE: You need to set your Shodan API key in environment variable :

```bash
export SHODAN_API_KEY="YOUR_API_KEY"
```

#### 2. Hunt Internet-wide via Shodan
```bash
export SHODAN_API_KEY="your_api_key"
./favicreep shodan --hash 12345678
```
This will query Shodan for http.favicon.hash:12345678 and return matching public IPs, ports, and hostnames.

#### 3. Example Flow

```bash
./favicreep scan --domain target.com -o output.json
# Pick a hash from the output, then:
./favicreep shodan --hash 873172492
```

Boom 💥 — You've just pivoted from target.com to dozens of exposed servers worldwide.

## 🧠 Why It Works (The Hacker’s Intuition)

Developers are lazy efficient — they reuse favicons across environments.

Staging/dev/test instances often:

- Lack authentication
- Run older, vulnerable versions
- Lie outside normal scopes
- Favicon hash is deterministic and indexed by Shodan
- Use it to correlate infra your target team may have forgotten.

## ⚠️ Legal

This tool is for educational and authorized security testing only.
Do not scan or interact with systems you don't have permission to test.

