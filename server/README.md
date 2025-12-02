# NanoKVM Server

This is the backend server implementation for NanoKVM.

For detailed documentation, please visit our [Wiki](https://wiki.sipeed.com/nanokvmpro).

## Structure

```shell
server
├── common       # Common utility components
├── config       # Server configuration
├── dl_lib       # Shared object libraries
├── include      # Header files for shared objects
├── logger       # Logging system
├── middleware   # Server middleware components
├── proto        # API request/response definitions
├── router       # API route handlers
├── service      # Core service implementations
├── utils        # Utility functions
└── main.go
```

## Configuration

The configuration file path is `/etc/kvm/server.yaml`.

```yaml
proto: http
port:
    http: 80
    https: 443
cert:
    crt: server.crt
    key: server.key

# Log level (debug/info/warn/error)
# Note: Use 'info' or 'error' in production, 'debug' only for development
logger:
    level: info
    file: stdout

# Authentication setting (enable/disable)
# Note: Only disable authentication in development environment
authentication: enable

jwt:
   # JWT secret key. If left empty, a random 64-byte key will be generated automatically.
   secretKey: ""
   # JWT token expiration time in seconds. Default: 2678400 (31 days)
   refreshTokenDuration: 2678400
   # Invalidate all JWT tokens when the user logs out. Default: true
   revokeTokensOnLogout: true

# Address for custom STUN server
# Note: You can disable the STUN service by setting it to 'disable' (e.g., in a LAN environment)
stun: stun.l.google.com:19302

# Address and authentication for custom TURN server
turn:
    turnAddr: example_addr
    turnUser: example_user
    turnCred: example_cred
```

## Build & Deploy

1. Install toolchains

```bash
cd ROOT/support/scripts
./toolchain_setup.sh
```

2. Build

```bash
cd ROOT/server
./build.sh
```

3. Deploy

```bash
scp NanoKVM-Server root@your-nanokvm-pro-ip:/kvmapp/server/
ssh root@your-nanokvm-pro-ip
systemctl restart nanokvm
```
