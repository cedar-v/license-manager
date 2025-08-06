# License Manager

[中文](README.md) | [English](README_EN.md) 

---

## Overview

License Manager is an independent software licensing platform that provides license code generation, distribution, validation, and management services for IoT platforms and other software systems. It offers both online and offline licensing modes with hardware-based binding for enhanced security.

## Key Features

- 🔧 **Customer Management**: Complete customer information management with status control
- 🔐 **License Generation**: Online/offline license modes with hardware fingerprinting
- 📊 **License Management**: Real-time status monitoring and license lifecycle management
- 📦 **Deployment Packages**: Automatic generation of deployment packages with configurations
- 🌐 **API Services**: RESTful APIs for validation, activation, and heartbeat monitoring
- ⚙️ **System Management**: Admin authentication and monitoring dashboard
- 🛠️ **Cross-platform Tools**: Hardware information extraction tools for multiple platforms

## Technical Stack

- **Frontend**: Vue.js 3+ with modern UI components
- **Backend**: Go 1.23+ with Gin framework, GORM, Viper configuration and Logrus logging
- **Database**: PostgreSQL 12+ / MySQL 8+
- **Configuration**: YAML format configuration files
- **Deployment**: Docker, single machine, or system service

## API Endpoints

```
POST /api/validate      - License validation
POST /api/activate      - License activation
POST /api/heartbeat     - Heartbeat reporting
GET  /api/license/{code} - License information query
GET  /api/customers     - Customer list API
GET  /tools/{tool}      - Tool download
```

## Security & Performance

- **Security Features**:
  - JWT authentication and authorization
  - HMAC-SHA256 signature verification
  - Hardware fingerprint binding for anti-piracy
  - HTTPS transport encryption
  - Multi-language error message support (Chinese/English/Japanese)
  
- **Performance**:
  - **High Concurrency**: Go native goroutines support, theoretically supports 10,000+ concurrent connections
  - **Low Latency**: Average API response time < 50ms
  - **Memory Optimized**: Go GC optimization, memory usage < 100MB
  - **Database Connection Pool**: Connection reuse for maximum database performance
  
- **Reliability**:
  - Comprehensive error handling with multi-language error messages
  - Structured logging and monitoring
  - Automatic database migrations
  - Graceful shutdown and resource cleanup

## Installation

```bash
# Clone the repository
git clone https://github.com/cedar-v/license-manager.git
cd license-manager/backend/cmd

# Build the application
go build -o license-manager

# Configure the application
cp config.example.yaml config.yaml
# Edit config.yaml with your settings

# Run the application
./license-manager
```

## Docker Deployment

```bash
# Build Docker image (multi-stage build with frontend and backend)
docker build -t license-manager .

# Create configuration file
cp backend/configs/config.example.yaml config.yaml
# Edit config.yaml to configure database and other settings

# Run with Docker
docker run -d \
  --name license-manager \
  -p 18888:18888 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  license-manager

# Check running status
docker logs license-manager

# Health check
curl http://localhost:18888/health
```

### Docker Compose Deployment (Recommended)

```yaml
# docker-compose.yml
version: '3.8'
services:
  license-manager:
    build: .
    ports:
      - "18888:18888"
    volumes:
      - ./config.yaml:/app/config.yaml
    environment:
      - GIN_MODE=release
    restart: unless-stopped
    
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root@123
      MYSQL_DATABASE: license_manager
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

volumes:
  mysql_data:
```

## Open Source License

This project is licensed under the **GNU General Public License v3.0 (GPL-3.0)**.

### License Description

- **Free Use**: You are free to use, study, modify, and distribute this software
- **Open Source Requirement**: If you distribute modified versions, they must also be open source under GPL-3.0 license
- **Copyright Protection**: Derivative works using this software must retain the original copyright notice
- **No Warranty**: The software is provided "as is" without any express or implied warranty

### Commercial Use

- Commercial use is permitted, but derivative works must also be open source
- For proprietary licensing or commercial support, please contact the project maintainers

### Full License Text

For detailed license terms, please see the [LICENSE](LICENSE) file in the project root directory, or visit:
https://www.gnu.org/licenses/gpl-3.0.html

---

## Contributing

We welcome contributions! Please feel free to submit a Pull Request.

## Support

If you have any questions or need support, please open an issue.