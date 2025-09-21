# StremThru Proxy

A secure HTTP(S) proxy with authentication, JWT encryption and advanced tunneling support.

## ✨ Features

- 🔐 **Secure HTTP(S) proxy** with Basic Auth authentication
- 👥 **Multi-user support** with fine-grained permission management
- 🔑 **JWT tokens** with automatic expiration and AES-256-GCM encryption
- 🌐 **Advanced tunneling** via SOCKS5, HTTP proxy, Cloudflare WARP
- ⚡ **Byte serving** for optimal streaming
- 📊 **Custom headers** and request management
- 🛡️ **Enhanced security**: configurable secrets, structured logs
- 🐳 **Docker ready** with simplified configuration

## 🚀 Installation

### Manual Installation

```bash
# 1. Download and compile
git clone https://github.com/Dydhzo/stremthru-proxy
cd stremthru-proxy
go build -o stremthru-proxy

# 2. Set required authentication (MANDATORY)
export STREMTHRU_PROXY_AUTH="admin:your-password"
export STREMTHRU_JWT_SECRET="your-secret-key"  # Generate with: openssl rand -base64 32

# 3. Run the proxy
./stremthru-proxy
```

### Docker Compose Installation

```yaml
# docker-compose.yml
services:
  stremthru-proxy:
    image: ghcr.io/dydhzo/stremthru-proxy:latest
    ports:
      - "8080:8080"
    environment:
      - STREMTHRU_BASE_URL=http://localhost:8080
      - STREMTHRU_PROXY_AUTH=admin:your-password  # REQUIRED
      - STREMTHRU_JWT_SECRET=your-secret-key  # Generate: openssl rand -base64 32
    restart: unless-stopped
```

```bash
docker-compose up -d
```

## ⚙️ Configuration

### Environment variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `STREMTHRU_PORT` | Listening port | `8080` | No |
| `STREMTHRU_BASE_URL` | Proxy base URL | `http://localhost:8080` | No |
| `STREMTHRU_JWT_SECRET` | JWT secret key (IMPORTANT!) | *random* | **Recommended** |
| `STREMTHRU_PROXY_AUTH` | User authentication | - | **REQUIRED** |
| `STREMTHRU_HTTP_PROXY` | External proxy for tunneling | - | No |
| `STREMTHRU_TUNNEL` | Tunneling configuration by hostname | - | No |
| `STREMTHRU_LOG_LEVEL` | Log level (DEBUG/INFO/WARN/ERROR) | `INFO` | No |
| `STREMTHRU_LOG_FORMAT` | Log format (json/text) | `json` | No |

## 🛠️ Available endpoints

| Endpoint | Method | Description | Auth Required |
|----------|---------|-------------|---------------|
| `/` | GET | Landing page with server information | No |
| `/v0/health` | GET | Service health check | No |
| `/v0/health/__debug__` | GET | Debug health check (detailed info) | No |
| `/v0/stats` | GET | Real-time statistics (bandwidth, connections) | **Yes** |
| `/v0/proxy` | GET | Create proxy links (simple mode) | **Yes** |
| `/v0/proxy` | POST | Create proxy links (advanced mode) | **Yes** |
| `/v0/proxy/{token}` | GET | Access proxied content via JWT token | No |
| `/v0/proxy/{token}` | HEAD | Headers only (without downloading) | No |
| `/v0/proxy/{token}/{filename}` | GET | Access with custom filename | No |

## 📄 License

MIT License - see LICENSE for details.

## 🙏 Credits

Based on the original [StremThru](https://github.com/MunifTanjim/stremthru) project by MunifTanjim.
This version focuses exclusively on proxy functionality with security enhancements.