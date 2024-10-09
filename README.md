## targets-filter
Converting trickest and chaos bbp targets in json, updates every 12 hour

## Usage
Get `amazon.com` domains from trickest
```
curl -s "https://raw.githubusercontent.com/rix4uni/targets-filter/refs/heads/main/trickest-targets.json" | jq -r '.[] | select(.domain == "amazon.com") | .hostnames' | xargs -I{} curl -s "{}"
```

Get `amazon.com` domains `zip` from chaos
```
curl -s "https://raw.githubusercontent.com/rix4uni/targets-filter/refs/heads/main/chaos-targets.json" | jq -r '.[] | select(.domain == "amazon.com") | .ZIP' | xargs -I{} wget -q {}
```