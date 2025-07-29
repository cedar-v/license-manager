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
- **Backend**: Go 1.23+ with Gorilla Mux router and Logrus logging
- **Database**: PostgreSQL 12+ / MySQL 12+
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

- **Security**: HMAC-SHA256 signing, hardware fingerprint binding, HTTPS encryption
- **Performance**: Supports 100+ concurrent users with <2s API response time
- **Reliability**: Comprehensive error handling and logging

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd license-manager

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
# Build Docker image
docker build -t license-manager .

# Run with Docker
docker run -p 8080:8080 -v $(pwd)/config.yaml:/app/config.yaml license-manager
```

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

---

## Contributing

We welcome contributions! Please feel free to submit a Pull Request.

## Support

If you have any questions or need support, please open an issue.