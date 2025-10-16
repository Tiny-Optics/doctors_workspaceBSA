# Doctor's Workspace - BLOODSA

A secure, custom-built web platform for haematologists, physicians, and data capturers to access Standard Operating Procedures (SOPs), manage referrals, and interact with the African HOPeR Registry.

## 🎯 Project Overview

The Doctor's Workspace is a dedicated sub-domain platform for BLOODSA that provides authenticated healthcare professionals with:

- **Secure Access** to clinical documentation and SOPs
- **Referral Integration** for transplant cases via REDCap
- **Registry Management** for the African HOPeR Registry
- **Document Management** for ethics approvals and clinical guidelines

**Project Duration:** 9 weeks (October 10 - December 15, 2025)  
**Budget:** R40,000.00  
**Website:** [tbc].bloodsa.org.za

## ✨ Key Features

### 1. Secure Access & User Management
- Robust authentication system for healthcare professionals
- Role-based access control (Admin, Doctor, Data Capturer)
- Admin panel for user creation and permission management
- Public gateway from main BLOODSA website

### 2. Standard Operating Procedures (SOPs)
- Categorized SOP library:
  - Anemia
  - Lymphoma
  - Myeloma
  - General Business
- PDF viewer and download functionality
- De-identified documents stored on Dropbox, served via AWS
- Searchable document repository

### 3. Referral Integration
- One-click "Refer for Transplant" access from dashboard
- Integration with REDCap referral forms
- Automated email notifications to 10 pre-defined specialists
- Secure submission tracking

### 4. African HOPeR Registry Documents
- Structured documentation repository
- Training video library (YouTube private/unlisted)
- Downloadable guidelines and example documents
- Ethics approval upload system:
  - Authenticated user uploads
  - Automatic storage in `Dropbox/BLDS_approvals/{username}/`
  - Admin email notifications on successful upload

## 🛠️ Technology Stack

### Core Technologies
- **Backend:** Go 1.24.4 with Gin framework
- **Frontend:** Vue.js 3.5 + TypeScript + Vite
- **Database:** MongoDB
- **State Management:** Pinia
- **Routing:** Vue Router
- **Authentication:** Custom secure authentication system
- **Containerization:** Docker

### Third-Party Services
- **REDCap** - Research form management (existing client account)
- **Dropbox** - Document storage (existing client account)
- **Email API** - Free tier (up to 300 emails/day)

### Hosting
- **Provider:** [To be determined]
- **Cost:** R285/month or R3,420/year
- **Resources:**
  - vCPU: 1 vCore
  - Memory: 2 GB
  - Storage: 50 GB

## 📅 Project Timeline

### Phase 1: Discovery & Core MVP (6 weeks)

| Sprint | Timeframe | Focus | Payment |
|--------|-----------|-------|---------|
| **Sprint 1** | Oct 17-23 | Planning & Architecture | - |
| **Sprint 2** | Oct 24-30 | Authentication & SOP Framework | R3,750 |
| **Sprint 3** | Oct 31 - Nov 6 | Referral & Registry MVP | R3,750 |
| **Sprint 4** | Nov 7-13 | Document Management & Upload | R3,750 |
| **Sprint 5** | Nov 14-20 | Development Continuation | R3,750 |
| **Sprint 6** | Nov 21-27 | QA & UAT | R3,750 |

### Phase 2: Enhancements, UX & Launch (3 weeks)

| Sprint | Timeframe | Focus | Payment |
|--------|-----------|-------|---------|
| **Sprint 7** | Nov 28 - Dec 4 | UX & UI Enhancements | R3,750 |
| **Sprint 8** | Dec 5-11 | Performance & Security | R3,750 |
| **Sprint 9** | Dec 12-15 | Final QA & Launch Prep | R3,750 |

**Total Sprints:** 9  
**Launch Target:** Mid-December 2025

## 💰 Payment Structure

| Milestone | Amount | Percentage |
|-----------|--------|------------|
| Project Kick-off (Deposit) | R10,000 | 25% |
| Phase 1 Complete | R10,000 | 25% |
| Phase 2 Complete | R10,000 | 25% |
| Sprint Payments (9 × R3,750) | R10,000 | 25% |
| **Total** | **R40,000** | **100%** |

## 🚀 Getting Started

### Prerequisites
- Access credentials for AWS
- Access credentials for Dropbox
- REDCap account and form URLs
- Admin access to BLOODSA.org.za

### Environment Setup

#### Backend Setup
```bash
# Navigate to backend directory
cd backend

# Install Go dependencies
go mod download

# Configure environment variables
cp .env.example .env
# Edit .env with your credentials (MongoDB, AWS, Dropbox, Email API)

# Start MongoDB with Docker
docker-compose up -d

# Run backend development server
make run
# or
go run cmd/main.go
```

#### Frontend Setup
```bash
# Navigate to frontend directory
cd frontend

# Install Node.js dependencies (requires Node.js 20.19.0+ or 22.12.0+)
npm install

# Configure environment variables
cp .env.example .env
# Edit .env with API endpoints

# Run frontend development server
npm run dev

# Build for production
npm run build
```

#### Docker Setup (Recommended)

For a simpler setup with all services containerized:

```bash
# Copy environment variables
cp .env.example .env
# Edit .env with your credentials

# Start all services (frontend, backend, and MongoDB)
make dev
# or
docker compose up

# Start in background
make up
# or
docker compose up -d
```

**Available Docker Commands:**
- `make dev` - Start development environment
- `make prod` - Start production environment
- `make logs` - View all logs
- `make down` - Stop all services
- `make clean` - Remove containers and volumes

See [DOCKER.md](./DOCKER.md) for detailed Docker documentation.

### Configuration Required
1. **AWS Integration** - Document storage and serving
2. **Dropbox Integration** - File storage and retrieval
3. **REDCap Integration** - Referral form linking
4. **Email API** - Notification system
5. **User Roles** - Admin, Doctor, Data Capturer permissions

## 📂 Project Structure

```
doctors_workspaceBSA/
├── backend/
│   ├── cmd/                 # Application entry points
│   │   └── api/            # API server
│   ├── internal/            # Private application code
│   ├── .air.toml           # Air hot-reload configuration
│   ├── .dockerignore       # Docker ignore file
│   ├── docker-compose.yml   # MongoDB standalone config
│   ├── Dockerfile          # Production Dockerfile
│   ├── Dockerfile.dev      # Development Dockerfile with Air
│   ├── go.mod              # Go module dependencies
│   ├── Makefile            # Build and run commands
│   └── README.md
├── frontend/
│   ├── src/
│   │   ├── components/     # Vue components
│   │   ├── router/         # Vue Router configuration
│   │   ├── stores/         # Pinia state management
│   │   ├── App.vue         # Root component
│   │   └── main.ts         # Application entry point
│   ├── public/             # Static assets
│   ├── .dockerignore       # Docker ignore file
│   ├── Dockerfile          # Production Dockerfile with Nginx
│   ├── Dockerfile.dev      # Development Dockerfile with Vite
│   ├── nginx.conf          # Nginx configuration for production
│   ├── package.json        # Node.js dependencies
│   ├── vite.config.ts      # Vite configuration
│   └── tsconfig.json       # TypeScript configuration
├── .env.example            # Environment variables template
├── docker-compose.yml      # Development orchestration
├── docker-compose.prod.yml # Production orchestration
├── DOCKER.md               # Docker documentation
├── Makefile                # Root-level Docker commands
└── README.md               # This file
```

## 🔒 Security Features

- SSL/TLS encryption
- Secure authentication and session management
- Role-based access control (RBAC)
- Secure file upload validation
- Regular security audits and updates
- Automated backups
- Admin permission hardening

## 🧪 Testing

### QA Testing
- Authentication flow testing
- SOP access and download testing
- Referral form integration testing
- Registry upload functionality
- Responsive design testing (mobile, tablet, desktop)
- Cross-browser compatibility testing

### User Acceptance Testing (UAT)
- Client-led UAT during Sprint 6
- Feedback consolidation within 3 working days
- Bug triaging and prioritization

## 📋 Client Responsibilities

- Provide accurate data before sprint deadlines
- Consolidate feedback within 3 working days
- Provide branding assets (logo, fonts, colors)
- Supply sample SOPs and registry documents
- REDCap form URLs and doctor email lists
- UAT participation and sign-off
- Final content upload before launch

## 🔄 Change Management

- Scope limited to features described in project documentation
- Change requests must be documented and approved
- Changes may affect timeline and cost
- 10% delay fee per week for client-caused delays

## 📞 Support & Contact

**Project Team:** Tiny Optics  
**Client:** BLOODSA

### Hosting Support
Post-launch hosting maintenance will be managed under BLOODSA's existing website retainer agreement at R285/month.

## 📄 License

[To be determined - specify license type]

## 🙏 Acknowledgments

- BLOODSA team for project collaboration
- Healthcare professionals participating in UAT
- REDCap platform for research data capture

---

**Last Updated:** October 2025  
**Status:** In Development - Sprint [Current Sprint Number]  
**Next Milestone:** [Next deliverable/sprint]
