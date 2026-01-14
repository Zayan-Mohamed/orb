# Production Deployment

Guidelines for deploying Orb in production environments.

## Architecture

### Recommended Setup

```
Internet
   ↓
Load Balancer (TLS termination)
   ↓
Reverse Proxy (Nginx/Caddy)
   ↓
Orb Relay Servers (multiple instances)
   ↓
Session Database (Redis/PostgreSQL)
```

## TLS Configuration

### Generate Certificates

```bash
# Let's Encrypt
sudo certbot certonly --standalone -d relay.example.com

# Self-signed (development only)
openssl req -x509 -newkey rsa:4096 -nodes -out cert.pem -keyout key.pem -days 365
```

### Configure TLS Relay

Use reverse proxy (Nginx/Caddy) for TLS termination.

## High Availability

### Load Balancing

```nginx
upstream orb_relays {
    least_conn;
    server relay1.internal:8080;
    server relay2.internal:8080;
    server relay3.internal:8080;
}
```

### Session Persistence

Implement sticky sessions or shared session store.

## Monitoring and Logging

### Prometheus Metrics

- Connection count
- Active sessions
- Bandwidth usage
- Error rates

### Logging

```bash
# Structured logging
orb relay 2>&1 | tee -a /var/log/orb/relay.log
```

## Backup and Recovery

- Configuration backup
- Log retention
- Disaster recovery plan

## Security Hardening

- Minimal permissions
- SELinux/AppArmor
- Firewall rules
- DDoS protection
- Rate limiting
- Regular audits

## Scaling

### Horizontal Scaling

Add more relay servers behind load balancer.

### Vertical Scaling

Increase CPU/RAM for existing servers.

## Next Steps

- [Relay Server Setup](relay-server.md)
- [Docker Deployment](docker.md)
