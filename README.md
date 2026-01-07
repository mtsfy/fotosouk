<h1 align="center">fotosouk</h1>
<p align="center" style="font-style: italic;">- ፎቶ ሱቅ - </p>

> [!NOTE]
> Work in progress.

## :page_facing_up: Description

Fast and scalable image processing backend service. Upload images, apply transformations (resize, crop, rotate, filters), retrieve in different formats, and manage your image library with user authentication. Built with Go, PostgreSQL, and AWS S3.

## :hammer_and_wrench: Development

- Go 1.21+
- PostgreSQL 14+
- Docker & Docker Compose
- AWS S3 (or compatible)

## :sparkles: Features

- **User Authentication**

  - Sign-up and login with JWT
  - Refresh token support

- **Image Management**

  - Upload images (JPEG, PNG)
  - List all uploaded images with metadata
  - Retrieve image by ID
  - Delete images with storage cleanup

- **Image Transformations**

  - Resize with high-quality resampling
  - Crop (center-based)
  - Rotate (90, 180, 270 degrees - lossless)
  - Format conversion (JPEG ↔ PNG)
  - Filters: Grayscale, Sepia
  - Chainable transformations

- **Storage**
  - AWS S3 integration
  - Efficient file uploads/downloads
  - Public image URLs

## :rocket: Quick Start

### Using Docker Compose

```bash
# Clone and setup
git clone git@github.com:mtsfy/fotosouk.git
cd fotosouk

# Create .env file
cp .env.example .env
# Edit .env with your AWS credentials and database config

# Build and run
docker-compose up --build
```

### Environment Variables

```env
# Server
PORT=3000

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=fotosouk
DATABASE_URL=postgresql://...

# JWT
JWT_ACCESS_SECRET=your_secret_key_here
JWT_REFRESH_SECRET=your_refresh_secret_key

# AWS S3
AWS_S3_BUCKET=your-bucket-name
AWS_S3_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret
```

### API Endpoints

#### Authentication

```bash
# Sign up
POST /register
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "securepass123"
}

# Login
POST /login
{
  "username": "johndoe",
  "password": "securepass123"
}

# Refresh token
POST /refresh
```

#### Images

```bash
# Upload image
POST /images
Content-Type: multipart/form-data
image: <file>

# Get all images
GET /images

# Get image by ID
GET /images/:id

# Transform image
POST /images/:id/transform
{
  "transformations": {
    "resize": { "width": 800, "height": 600 },
    "rotate": 90,
    "filters": { "grayscale": true, "sepia": true },
    "format": "jpeg"
  }
}

# Delete image
DELETE /images/:id
```
