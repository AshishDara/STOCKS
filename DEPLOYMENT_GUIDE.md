# Complete Deployment Guide: GrowTrade Application

## Table of Contents
1. [Accessing Your Code](#accessing-your-code)
2. [Setting Up GitHub Repository](#setting-up-github-repository)
3. [Preparing for Deployment](#preparing-for-deployment)
4. [AWS Deployment Options](#aws-deployment-options)
5. [Option 1: AWS EC2 Deployment (Recommended)](#option-1-aws-ec2-deployment-recommended)
6. [Option 2: AWS Elastic Beanstalk](#option-2-aws-elastic-beanstalk)
7. [Option 3: AWS App Runner](#option-3-aws-app-runner)
8. [Frontend Deployment (Vercel/Netlify)](#frontend-deployment-vercelnetlify)
9. [Post-Deployment Configuration](#post-deployment-configuration)
10. [Troubleshooting](#troubleshooting)

---

## Accessing Your Code

### Your Code Location
All your application code is already in your workspace:
```
/Users/ashishdara634/Desktop/STOCKS/
├── backend/          # Golang backend
├── frontend/         # React frontend
├── README.md
├── INTERVIEW_GUIDE.md
└── Other files...
```

### Verify Your Code
```bash
# Navigate to your project
cd ~/Desktop/STOCKS

# Check backend files
ls -la backend/

# Check frontend files
ls -la frontend/
```

**Your code is ready!** All files are in the STOCKS directory.

---

## Setting Up GitHub Repository

### Step 1: Create GitHub Account (if you don't have one)
1. Go to https://github.com
2. Sign up for a free account
3. Verify your email

### Step 2: Create New Repository on GitHub
1. Log in to GitHub
2. Click the **"+"** icon in top right → **"New repository"**
3. Repository name: `growtrade`
4. Description: "Full-stack trading dashboard with Golang backend and React frontend"
5. Visibility: **Public** (or Private if you prefer)
6. **DO NOT** initialize with README, .gitignore, or license (we already have files)
7. Click **"Create repository"**

### Step 3: Initialize Git in Your Project
```bash
# Navigate to your project directory
cd ~/Desktop/STOCKS

# Initialize git repository
git init

# Add all files
git add .

# Create initial commit
git commit -m "Initial commit: Trading dashboard application"

# Add remote repository (replace YOUR_USERNAME with your GitHub username)
git remote add origin https://github.com/YOUR_USERNAME/growtrade.git

# Push to GitHub
git branch -M main
git push -u origin main
```

**Example:**
```bash
# If your GitHub username is "ashishdara634"
git remote add origin https://github.com/ashishdara634/growtrade.git
```

### Step 4: Verify Repository
1. Go to https://github.com/YOUR_USERNAME/growtrade
2. You should see all your files

### Step 5: Create .gitignore (if not exists)
Create `.gitignore` in root directory:

```bash
# Backend
backend/*.db
backend/*.sqlite
backend/*.sqlite3
backend/trading-server

# Frontend
frontend/node_modules/
frontend/dist/
frontend/.vite/

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Environment
.env
.env.local
```

---

## Preparing for Deployment

### Step 1: Update Backend Configuration

#### Create Environment File
1. Copy the template file:
   ```bash
   cd backend
   cp env.example .env
   ```
2. Update `.env` (for local development or deployment):
   ```env
   JWT_SECRET=your-super-secret-key-change-in-production
   PORT=8080
   DB_PATH=trading.db
   ALLOWED_ORIGINS=http://localhost:3000
   ```
   - `JWT_SECRET`: required in production
   - `PORT`: API port (defaults to 8080)
   - `DB_PATH`: SQLite file path
   - `ALLOWED_ORIGINS`: comma-separated list of allowed origins for CORS (leave empty to allow all during local development)

#### Update main.go to Use Environment Variables
We need to modify the backend to read from environment variables:

```go
// In main.go, update the JWT secret section:
func main() {
    // Get JWT secret from environment or use default
    if secret := os.Getenv("JWT_SECRET"); secret != "" {
        jwtSecret = []byte(secret)
    } else {
        jwtSecret = []byte("your-secret-key-change-in-production")
    }
    
    // Get port from environment
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    // ... rest of code ...
    
    log.Printf("Server starting on :%s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

### Step 2: Update Frontend Configuration

#### Update API URLs for Production
1. Copy the template file:
   ```bash
   cd frontend
   cp env.example .env
   ```
2. Update `.env` for local development:
   ```env
   VITE_API_URL=http://localhost:8080
   VITE_WS_URL=ws://localhost:8080/ws
   ```
3. For production (optionally create `.env.production`):
   ```env
   VITE_API_URL=https://your-backend-domain.com
   VITE_WS_URL=wss://your-backend-domain.com/ws
   ```
4. Environment variables are consumed through `src/config.js`, so no component-level changes are required when switching environments.

### Step 3: Build Frontend for Production
```bash
cd frontend
npm run build
```

This creates a `dist/` folder with production-ready files.

### Step 4: Create Deployment Scripts

#### Backend Build Script (`backend/build.sh`)
```bash
#!/bin/bash
echo "Building GrowTrade backend..."
go mod download
go build -o growtrade-server main.go
echo "Build complete! Binary: growtrade-server"
```

Make it executable:
```bash
chmod +x backend/build.sh
```

#### Frontend Build Script (`frontend/build.sh`)
```bash
#!/bin/bash
echo "Building GrowTrade frontend..."
npm install
npm run build
echo "Build complete! Files in dist/ folder"
```

Make it executable:
```bash
chmod +x frontend/build.sh
```

---

## AWS Deployment Options

### Comparison

| Option | Pros | Cons | Best For |
|--------|------|------|----------|
| **EC2** | Full control, flexible | More setup required | Learning, custom configs |
| **Elastic Beanstalk** | Easy deployment, auto-scaling | Less control | Quick deployment |
| **App Runner** | Serverless, auto-scaling | Limited customization | Simple apps |
| **ECS/Fargate** | Container-based, scalable | Complex setup | Production apps |

**Recommendation**: Start with **EC2** for learning, then consider **Elastic Beanstalk** for easier management.

---

## Option 1: AWS EC2 Deployment (Recommended)

### Step 1: Create AWS Account
1. Go to https://aws.amazon.com
2. Click "Create an AWS Account"
3. Follow the signup process
4. **Note**: You'll need a credit card (free tier available)

### Step 2: Launch EC2 Instance

1. **Log in to AWS Console**
   - Go to https://console.aws.amazon.com
   - Search for "EC2" in services

2. **Launch Instance**
   - Click "Launch Instance"
   - Name: `growtrade-backend`

3. **Choose AMI (Amazon Machine Image)**
   - Select **Ubuntu Server 22.04 LTS** (Free tier eligible)
   - Architecture: 64-bit (x86)

4. **Instance Type**
   - Select **t2.micro** (Free tier eligible)
   - 1 vCPU, 1 GB RAM

5. **Key Pair**
   - Click "Create new key pair"
   - Name: `growtrade-key`
   - Type: RSA
   - Format: .pem
   - Click "Create key pair"
   - **IMPORTANT**: Download and save the .pem file securely!

6. **Network Settings**
   - Allow SSH traffic from: My IP
   - Click "Add security group rule"
     - Type: Custom TCP
     - Port: 8080
     - Source: 0.0.0.0/0 (allow from anywhere)
   - Click "Add security group rule"
     - Type: Custom TCP
     - Port: 80
     - Source: 0.0.0.0/0

7. **Configure Storage**
   - 8 GB (Free tier)

8. **Launch Instance**
   - Click "Launch Instance"
   - Wait for instance to be "Running"

### Step 3: Connect to EC2 Instance

#### On Mac/Linux:
```bash
# Navigate to where you saved the .pem file
cd ~/Downloads  # or wherever you saved it

# Set correct permissions
chmod 400 growtrade-key.pem

# Connect to instance
ssh -i growtrade-key.pem ubuntu@YOUR_EC2_PUBLIC_IP
```

**Find your EC2 Public IP:**
- Go to EC2 Dashboard
- Click on your instance
- Copy "Public IPv4 address"

#### On Windows:
Use **PuTTY** or **Windows Terminal**:
1. Download PuTTY: https://www.putty.org/
2. Convert .pem to .ppk using PuTTYgen
3. Connect using PuTTY

### Step 4: Install Dependencies on EC2

Once connected via SSH:

```bash
# Update system
sudo apt update
sudo apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Verify Go installation
go version

# Install Node.js (for building frontend if needed)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# Install nginx (for serving frontend)
sudo apt install nginx -y

# Install SQLite (if not already installed)
sudo apt install sqlite3 -y
```

### Step 5: Deploy Backend to EC2

#### Option A: Using Git (Recommended)

```bash
# On your local machine
cd ~/Desktop/STOCKS
git add .
git commit -m "Prepare for deployment"
git push origin main

# On EC2 instance (via SSH)
cd ~
git clone https://github.com/YOUR_USERNAME/growtrade.git
cd growtrade/backend

# Set environment variables
export JWT_SECRET="your-super-secret-production-key"
export PORT=8080

# Build and run
go mod download
go build -o growtrade-server main.go

# Run in background using screen or systemd
screen -S growtrade
./growtrade-server
# Press Ctrl+A then D to detach
```

#### Option B: Using SCP (Direct Copy)

```bash
# On your local machine
cd ~/Desktop/STOCKS/backend
go build -o growtrade-server main.go

# Copy to EC2
scp -i ~/Downloads/growtrade-key.pem growtrade-server ubuntu@YOUR_EC2_IP:~/
scp -i ~/Downloads/growtrade-key.pem -r . ubuntu@YOUR_EC2_IP:~/growtrade-backend/

# On EC2
chmod +x growtrade-server
export JWT_SECRET="your-secret-key"
./growtrade-server
```

### Step 6: Create Systemd Service (Keep Backend Running)

```bash
# On EC2, create service file
sudo nano /etc/systemd/system/growtrade.service
```

Add this content:
```ini
[Unit]
Description=GrowTrade Backend Server
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/growtrade/backend
Environment="JWT_SECRET=your-super-secret-production-key"
Environment="PORT=8080"
ExecStart=/home/ubuntu/growtrade/backend/growtrade-server
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable growtrade
sudo systemctl start growtrade
sudo systemctl status growtrade
```

### Step 7: Deploy Frontend to EC2

```bash
# On your local machine, build frontend
cd ~/Desktop/STOCKS/frontend
npm run build

# Copy dist folder to EC2
scp -i ~/Downloads/growtrade-key.pem -r dist ubuntu@YOUR_EC2_IP:~/

# On EC2, configure nginx
sudo nano /etc/nginx/sites-available/growtrade
```

Add nginx configuration:
```nginx
server {
    listen 80;
    server_name YOUR_EC2_PUBLIC_IP;  # or your domain

    # Serve frontend
    location / {
        root /home/ubuntu/dist;
        try_files $uri $uri/ /index.html;
    }

    # Proxy API requests to backend
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # Proxy WebSocket
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
}
```

Enable site:
```bash
sudo ln -s /etc/nginx/sites-available/growtrade /etc/nginx/sites-enabled/
sudo nginx -t  # Test configuration
sudo systemctl restart nginx
```

### Step 8: Get Your Application URL

Your application will be available at:
- **Frontend**: `http://YOUR_EC2_PUBLIC_IP`
- **Backend API**: `http://YOUR_EC2_PUBLIC_IP/api`

---

## Option 2: AWS Elastic Beanstalk

### Step 1: Install EB CLI

```bash
# On Mac
brew install awsebcli

# On Linux
pip install awsebcli --upgrade --user
```

### Step 2: Configure Backend for Elastic Beanstalk

Create `backend/.ebextensions/01-go.config`:
```yaml
option_settings:
  aws:elasticbeanstalk:container:golang:
    GoVersion: 1.21
  aws:elasticbeanstalk:application:environment:
    JWT_SECRET: your-secret-key
    PORT: 8080
```

Create `backend/Procfile`:
```
web: ./growtrade-server
```

### Step 3: Initialize Elastic Beanstalk

```bash
cd backend
eb init -p go -r us-east-1 growtrade-backend
eb create growtrade-env
eb deploy
```

### Step 4: Get Application URL

```bash
eb status
eb open
```

---

## Option 3: AWS App Runner

### Step 1: Create Dockerfile for Backend

Create `backend/Dockerfile`:
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o growtrade-server main.go

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/growtrade-server .
EXPOSE 8080
CMD ["./growtrade-server"]
```

### Step 2: Push to Container Registry

```bash
# Build and push to ECR (Elastic Container Registry)
aws ecr create-repository --repository-name growtrade-backend
# Follow AWS CLI instructions to push image
```

### Step 3: Create App Runner Service

1. Go to AWS Console → App Runner
2. Create service
3. Select container image
4. Configure service
5. Deploy

---

## Frontend Deployment (Vercel/Netlify)

### Option A: Vercel (Recommended for React)

#### Step 1: Install Vercel CLI
```bash
npm install -g vercel
```

#### Step 2: Deploy
```bash
cd frontend
vercel
```

Follow prompts:
- Link to existing project? No
- Project name: growtrade
- Directory: ./
- Override settings? No

#### Step 3: Configure Environment Variables
1. Go to Vercel Dashboard
2. Select your project
3. Settings → Environment Variables
4. Add:
   - `VITE_API_URL`: `https://your-backend-url.com`
   - `VITE_WS_URL`: `wss://your-backend-url.com/ws`

#### Step 4: Redeploy
```bash
vercel --prod
```

### Option B: Netlify

#### Step 1: Install Netlify CLI
```bash
npm install -g netlify-cli
```

#### Step 2: Deploy
```bash
cd frontend
netlify deploy --prod
```

#### Step 3: Configure Environment Variables
1. Netlify Dashboard → Site settings → Environment variables
2. Add `VITE_API_URL` and `VITE_WS_URL` (include `/ws` in the WebSocket URL)

### Option C: AWS S3 + CloudFront

#### Step 1: Create S3 Bucket
```bash
aws s3 mb s3://growtrade-frontend
```

#### Step 2: Upload Frontend Build
```bash
cd frontend
npm run build
aws s3 sync dist/ s3://growtrade-frontend --delete
```

#### Step 3: Enable Static Website Hosting
```bash
aws s3 website s3://growtrade-frontend --index-document index.html
```

#### Step 4: Create CloudFront Distribution
1. AWS Console → CloudFront
2. Create distribution
3. Origin: S3 bucket
4. Deploy

---

## Post-Deployment Configuration

### Step 1: Update Frontend Environment Variables

Update your frontend deployment with:
- `VITE_API_URL`: Your backend URL
- `VITE_WS_URL`: Your WebSocket URL (e.g., `wss://api.yourdomain.com/ws`)

### Step 2: Configure CORS

Update backend `main.go` CORS settings:
```go
config := cors.DefaultConfig()
config.AllowOrigins = []string{
    "https://your-frontend-domain.com",
    "http://localhost:3000",  // For local development
}
config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
r.Use(cors.New(config))
```

### Step 3: Set Up Domain (Optional)

#### Using Route 53
1. AWS Console → Route 53
2. Register domain or use existing
3. Create A record pointing to EC2 IP

#### Using Other DNS Providers
1. Add A record: `@` → Your EC2 IP
2. Add CNAME: `www` → Your domain

### Step 4: Enable HTTPS (SSL Certificate)

#### Using Let's Encrypt (Free)
```bash
# On EC2
sudo apt install certbot python3-certbot-nginx -y
sudo certbot --nginx -d yourdomain.com
```

#### Using AWS Certificate Manager
1. Request certificate in ACM
2. Validate domain
3. Configure in CloudFront or ALB

### Step 5: Set Up Monitoring

#### CloudWatch (AWS)
- Monitor EC2 instance
- Set up alarms
- View logs

#### Application Monitoring
- Add logging to backend
- Monitor error rates
- Track performance

---

## Troubleshooting

### Backend Not Starting

**Check logs:**
```bash
# If using systemd
sudo journalctl -u growtrade -f

# If running manually
./growtrade-server
```

**Common issues:**
- Port 8080 already in use: `sudo lsof -i :8080`
- Database permissions: `chmod 666 trading.db`
- Environment variables not set

### Frontend Not Connecting to Backend

**Check:**
1. Environment variables set correctly
2. CORS configured properly
3. Backend URL is correct
4. Security groups allow traffic

### WebSocket Not Working

**Check:**
1. WebSocket URL uses `wss://` (not `ws://`) for HTTPS
2. Nginx proxy configuration correct
3. Security groups allow WebSocket upgrade

### Database Issues

**Backup database:**
```bash
sqlite3 trading.db .dump > backup.sql
```

**Restore database:**
```bash
sqlite3 trading.db < backup.sql
```

### Performance Issues

**Optimize:**
1. Use larger EC2 instance
2. Enable database indexing
3. Add caching (Redis)
4. Use CDN for frontend

---

## Cost Estimation (AWS)

### Free Tier (First Year)
- EC2 t2.micro: **Free** (750 hours/month)
- S3: **Free** (5 GB storage)
- Data Transfer: **Free** (15 GB out)

### After Free Tier
- EC2 t2.micro: ~$8-10/month
- S3: ~$0.023/GB/month
- Data Transfer: ~$0.09/GB

**Estimated Monthly Cost**: $10-15/month (small traffic)

---

## Security Checklist

- [ ] Change default JWT secret
- [ ] Use HTTPS (SSL certificate)
- [ ] Restrict SSH access (only your IP)
- [ ] Enable firewall (ufw)
- [ ] Regular security updates
- [ ] Database backups
- [ ] Environment variables secured
- [ ] CORS properly configured
- [ ] Rate limiting implemented
- [ ] Input validation on all endpoints

---

## Maintenance

### Regular Tasks

**Weekly:**
- Check application logs
- Monitor server resources
- Review security groups

**Monthly:**
- Update dependencies
- Backup database
- Review costs

**As Needed:**
- Scale resources
- Update application
- Fix security issues

---

## Quick Deployment Commands

### Backend Deployment
```bash
# Build
cd backend
go build -o growtrade-server main.go

# Deploy to EC2
scp -i key.pem growtrade-server ubuntu@EC2_IP:~/
ssh -i key.pem ubuntu@EC2_IP
sudo systemctl restart growtrade
```

### Frontend Deployment
```bash
# Build
cd frontend
npm run build

# Deploy to Vercel
vercel --prod

# Or deploy to S3
aws s3 sync dist/ s3://growtrade-frontend --delete
```

---

## Next Steps After Deployment

1. **Test all features** on production
2. **Set up monitoring** and alerts
3. **Configure backups** for database
4. **Set up CI/CD** pipeline (GitHub Actions)
5. **Document** your deployment process
6. **Monitor** costs and optimize

---

## Support Resources

- **AWS Documentation**: https://docs.aws.amazon.com
- **Go Documentation**: https://go.dev/doc
- **React Documentation**: https://react.dev
- **Vercel Documentation**: https://vercel.com/docs

---

**Good luck with your deployment! 🚀**

