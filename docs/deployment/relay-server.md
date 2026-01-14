# Relay Server Deployment

Guide for deploying and managing Orb relay servers.

## Overview

The relay server is a stateless WebSocket server that forwards encrypted messages between peers. It cannot decrypt traffic and requires minimal resources.

## Quick Start

### Local Development

```bash
# Start relay
orb relay

# Or with custom port
orb relay --port 9090
```

### Production Deployment

```bash
# Bind to all interfaces
orb relay --host 0.0.0.0 --port 8080
```

## System Requirements

### Minimum

- CPU: 1 core
- RAM: 256 MB
- Disk: 100 MB
- Network: 10 Mbps

### Recommended

- CPU: 2 cores
- RAM: 512 MB
- Disk: 1 GB
- Network: 100 Mbps

## systemd Service

Create `/etc/systemd/system/orb-relay.service`:

```ini
[Unit]
Description=Orb Relay Server
After=network.target

[Service]
Type=simple
User=orb
ExecStart=/usr/local/bin/orb relay --host 0.0.0.0 --port 8080
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable orb-relay
sudo systemctl start orb-relay
sudo systemctl status orb-relay
```

## Reverse Proxy (Nginx)

```nginx
upstream orb {
    server localhost:8080;
}

server {
    listen 443 ssl http2;
    server_name relay.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location /ws {
        proxy_pass http://orb;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }

    location / {
        proxy_pass http://orb;
    }
}
```

## Monitoring

```bash
# Check status
curl http://localhost:8080/health

# View logs
journalctl -u orb-relay -f

# Monitor connections
netstat -an | grep :8080
```

## Security

- Use TLS (wss://)
- Firewall rules
- Rate limiting
- Log monitoring
- Regular updates

## Next Steps

- [Production Deployment](production.md)
- [Docker Deployment](docker.md)
