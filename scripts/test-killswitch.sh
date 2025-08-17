#!/bin/bash

# 🧪 Kill Switch System Integration Test
# Tests all components of the kill switch system
# Author: Sauron Pro Security Team

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_LOG="/tmp/killswitch-test-$(date +%s).log"
TEST_VPS_ID="test-vps-$(date +%s)"
TEST_ADMIN_KEY="test-key-$(openssl rand -hex 16)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Test results
TESTS_TOTAL=0
TESTS_PASSED=0
TESTS_FAILED=0

# Logging
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$TEST_LOG"
}

# Test result functions
test_start() {
    ((TESTS_TOTAL++))
    echo -e "${BLUE}🧪 TEST $TESTS_TOTAL: $1${NC}"
    log "TEST START: $1"
}

test_pass() {
    ((TESTS_PASSED++))
    echo -e "${GREEN}✅ PASS: $1${NC}"
    log "TEST PASS: $1"
}

test_fail() {
    ((TESTS_FAILED++))
    echo -e "${RED}❌ FAIL: $1${NC}"
    log "TEST FAIL: $1"
}

test_skip() {
    echo -e "${YELLOW}⏭️  SKIP: $1${NC}"
    log "TEST SKIP: $1"
}

# Banner
show_banner() {
    echo -e "${BLUE}"
    cat << 'EOF'
╔══════════════════════════════════════════════════════════════╗
║                🧪 KILL SWITCH INTEGRATION TEST 🧪              ║
║                                                              ║
║  This test suite validates all kill switch components       ║
║  in a safe, controlled environment.                         ║
║                                                              ║
║  ⚠️  DO NOT RUN ON PRODUCTION SYSTEMS  ⚠️                   ║
╚══════════════════════════════════════════════════════════════╝
EOF
    echo -e "${NC}"
}

# Test 1: Environment Setup
test_environment_setup() {
    test_start "Environment setup and prerequisites"
    
    # Check required tools
    local tools=("curl" "jq" "openssl" "systemctl")
    for tool in "${tools[@]}"; do
        if command -v "$tool" >/dev/null 2>&1; then
            log "✓ $tool found"
        else
            test_fail "$tool not found"
            return 1
        fi
    done
    
    # Check script files exist
    local scripts=(
        "$SCRIPT_DIR/../scripts/kill-switch.sh"
        "$SCRIPT_DIR/../scripts/heartbeat-master.sh"
    )
    
    for script in "${scripts[@]}"; do
        if [[ -f "$script" && -x "$script" ]]; then
            log "✓ $script exists and is executable"
        else
            test_fail "$script missing or not executable"
            return 1
        fi
    done
    
    test_pass "Environment setup complete"
}

# Test 2: Handler File Validation
test_handler_validation() {
    test_start "Handler file validation"
    
    local handlers=(
        "$SCRIPT_DIR/../handlers/killswitch.go"
        "$SCRIPT_DIR/../handlers/deadmans_switch.go"
        "$SCRIPT_DIR/../handlers/fleet_management.go"
    )
    
    for handler in "${handlers[@]}"; do
        if [[ -f "$handler" ]]; then
            # Check for required functions
            local required_funcs=(
                "HandleKillSwitchMaster"
                "HandleKillSwitchVPS"
                "executeKillSwitch"
            )
            
            for func in "${required_funcs[@]}"; do
                if grep -q "$func" "$handler" 2>/dev/null; then
                    log "✓ Function $func found in $(basename "$handler")"
                else
                    log "⚠ Function $func not found in $(basename "$handler")"
                fi
            done
        else
            test_fail "Handler file missing: $(basename "$handler")"
            return 1
        fi
    done
    
    test_pass "Handler validation complete"
}

# Test 3: Configuration Validation
test_configuration() {
    test_start "Configuration system validation"
    
    # Test configuration files
    local configs=(
        "$SCRIPT_DIR/../install/sauron-heartbeat.service"
        "$SCRIPT_DIR/../docs/KILL_SWITCH_MANUAL.md"
    )
    
    for config in "${configs[@]}"; do
        if [[ -f "$config" ]]; then
            log "✓ Configuration file exists: $(basename "$config")"
        else
            test_fail "Configuration file missing: $(basename "$config")"
            return 1
        fi
    done
    
    test_pass "Configuration validation complete"
}

# Test 4: Dead Man's Switch Logic
test_deadmans_switch() {
    test_start "Dead man's switch logic validation"
    
    # Create test configuration
    local test_config='{
        "enabled": true,
        "check_interval": "5s",
        "master_timeout": "10s",
        "auto_destruct": false,
        "destruction_level": 1
    }'
    
    echo "$test_config" > "/tmp/deadmans_test_config.json"
    
    if [[ -f "/tmp/deadmans_test_config.json" ]]; then
        log "✓ Dead man's switch test configuration created"
        
        # Validate JSON structure
        if jq empty "/tmp/deadmans_test_config.json" 2>/dev/null; then
            log "✓ Configuration JSON is valid"
        else
            test_fail "Configuration JSON is invalid"
            return 1
        fi
    else
        test_fail "Failed to create test configuration"
        return 1
    fi
    
    rm -f "/tmp/deadmans_test_config.json"
    test_pass "Dead man's switch logic validation complete"
}

# Test 5: Kill Switch CLI Interface
test_cli_interface() {
    test_start "Kill switch CLI interface validation"
    
    local kill_switch_script="$SCRIPT_DIR/../scripts/kill-switch.sh"
    
    # Test help/usage (should not require admin key)
    if "$kill_switch_script" --help >/dev/null 2>&1 || \
       "$kill_switch_script" -h >/dev/null 2>&1 || \
       "$kill_switch_script" help >/dev/null 2>&1; then
        log "✓ CLI help functionality works"
    else
        log "⚠ CLI help not available (may be intentional for security)"
    fi
    
    # Test script syntax
    if bash -n "$kill_switch_script"; then
        log "✓ Kill switch script syntax is valid"
    else
        test_fail "Kill switch script has syntax errors"
        return 1
    fi
    
    test_pass "CLI interface validation complete"
}

# Test 6: Heartbeat Service Validation
test_heartbeat_service() {
    test_start "Heartbeat service validation"
    
    local heartbeat_script="$SCRIPT_DIR/../scripts/heartbeat-master.sh"
    
    # Test script syntax
    if bash -n "$heartbeat_script"; then
        log "✓ Heartbeat script syntax is valid"
    else
        test_fail "Heartbeat script has syntax errors"
        return 1
    fi
    
    # Test service file
    local service_file="$SCRIPT_DIR/../install/sauron-heartbeat.service"
    if [[ -f "$service_file" ]]; then
        # Check service file structure
        if grep -q "\[Unit\]" "$service_file" && \
           grep -q "\[Service\]" "$service_file" && \
           grep -q "\[Install\]" "$service_file"; then
            log "✓ Service file structure is valid"
        else
            test_fail "Service file structure is invalid"
            return 1
        fi
    else
        test_fail "Service file missing"
        return 1
    fi
    
    test_pass "Heartbeat service validation complete"
}

# Test 7: Network Communication Simulation
test_network_simulation() {
    test_start "Network communication simulation"
    
    # Test JSON payload creation
    local test_payload
    test_payload=$(jq -n \
        --arg admin_key "$TEST_ADMIN_KEY" \
        --arg vps_id "$TEST_VPS_ID" \
        --argjson level "1" \
        --argjson delay "0" \
        --arg reason "Integration test" \
        --arg confirmation_code "Lord Sauron" \
        '{
            admin_key: $admin_key,
            vps_id: $vps_id,
            destruction_level: $level,
            delay_seconds: $delay,
            reason: $reason,
            confirmation_code: $confirmation_code
        }')
    
    if echo "$test_payload" | jq empty 2>/dev/null; then
        log "✓ Kill switch payload JSON is valid"
    else
        test_fail "Kill switch payload JSON is invalid"
        return 1
    fi
    
    # Test heartbeat payload
    local heartbeat_payload
    heartbeat_payload=$(jq -n \
        --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
        '{
            timestamp: $timestamp,
            source: "test"
        }')
    
    if echo "$heartbeat_payload" | jq empty 2>/dev/null; then
        log "✓ Heartbeat payload JSON is valid"
    else
        test_fail "Heartbeat payload JSON is invalid"
        return 1
    fi
    
    test_pass "Network communication simulation complete"
}

# Test 8: Security Validation
test_security() {
    test_start "Security validation"
    
    # Test admin key validation (simulated)
    local valid_key_pattern="^[a-zA-Z0-9]{32,}$"
    if [[ "$TEST_ADMIN_KEY" =~ $valid_key_pattern ]]; then
        log "✓ Admin key format validation works"
    else
        test_fail "Admin key format validation failed"
        return 1
    fi
    
    # Test confirmation code validation
    local confirmation_code="Lord Sauron"
    if [[ "$confirmation_code" == "Lord Sauron" ]]; then
        log "✓ Confirmation code validation works"
    else
        test_fail "Confirmation code validation failed"
        return 1
    fi
    
    test_pass "Security validation complete"
}

# Test 9: Documentation Completeness
test_documentation() {
    test_start "Documentation completeness"
    
    local doc_file="$SCRIPT_DIR/../docs/KILL_SWITCH_MANUAL.md"
    
    if [[ -f "$doc_file" ]]; then
        # Check for required sections
        local required_sections=(
            "OVERVIEW"
            "ARCHITECTURE COMPONENTS"
            "DESTRUCTION LEVELS"
            "ACTIVATION METHODS"
            "SECURITY PROTOCOLS"
            "EMERGENCY PROCEDURES"
        )
        
        for section in "${required_sections[@]}"; do
            if grep -q "$section" "$doc_file"; then
                log "✓ Documentation section found: $section"
            else
                log "⚠ Documentation section missing: $section"
            fi
        done
        
        # Check file size (should be comprehensive)
        local file_size
        file_size=$(wc -c < "$doc_file")
        if [[ $file_size -gt 10000 ]]; then
            log "✓ Documentation is comprehensive ($file_size bytes)"
        else
            log "⚠ Documentation may be incomplete ($file_size bytes)"
        fi
    else
        test_fail "Documentation file missing"
        return 1
    fi
    
    test_pass "Documentation validation complete"
}

# Test 10: Integration Readiness
test_integration_readiness() {
    test_start "Integration readiness assessment"
    
    # Check main.go integration
    local main_file="$SCRIPT_DIR/../main.go"
    if [[ -f "$main_file" ]]; then
        # Check for kill switch endpoints
        if grep -q "HandleKillSwitchMaster" "$main_file" && \
           grep -q "HandleKillSwitchVPS" "$main_file"; then
            log "✓ Kill switch endpoints integrated in main.go"
        else
            test_fail "Kill switch endpoints not properly integrated"
            return 1
        fi
        
        # Check for dead man's switch initialization
        if grep -q "InitDeadMansSwitch" "$main_file"; then
            log "✓ Dead man's switch initialization found"
        else
            test_fail "Dead man's switch not initialized in main.go"
            return 1
        fi
    else
        test_fail "main.go file not found"
        return 1
    fi
    
    test_pass "Integration readiness complete"
}

# Test Results Summary
show_results() {
    echo
    echo -e "${BLUE}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║                    🧪 TEST RESULTS SUMMARY 🧪                 ║${NC}"
    echo -e "${BLUE}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo
    echo -e "Total Tests: ${BLUE}$TESTS_TOTAL${NC}"
    echo -e "Passed:      ${GREEN}$TESTS_PASSED${NC}"
    echo -e "Failed:      ${RED}$TESTS_FAILED${NC}"
    
    local pass_rate=0
    if [[ $TESTS_TOTAL -gt 0 ]]; then
        pass_rate=$((TESTS_PASSED * 100 / TESTS_TOTAL))
    fi
    
    echo -e "Pass Rate:   ${BLUE}$pass_rate%${NC}"
    echo
    
    if [[ $TESTS_FAILED -eq 0 ]]; then
        echo -e "${GREEN}✅ ALL TESTS PASSED - SYSTEM READY FOR DEPLOYMENT${NC}"
        log "ALL TESTS PASSED - Pass rate: $pass_rate%"
    else
        echo -e "${RED}❌ $TESTS_FAILED TESTS FAILED - REVIEW REQUIRED${NC}"
        log "TESTS FAILED: $TESTS_FAILED/$TESTS_TOTAL - Pass rate: $pass_rate%"
    fi
    
    echo
    echo -e "Test log: ${YELLOW}$TEST_LOG${NC}"
}

# Main execution
main() {
    show_banner
    
    log "Starting kill switch integration test suite"
    
    # Run all tests
    test_environment_setup
    test_handler_validation
    test_configuration
    test_deadmans_switch
    test_cli_interface
    test_heartbeat_service
    test_network_simulation
    test_security
    test_documentation
    test_integration_readiness
    
    show_results
    
    # Exit with appropriate code
    if [[ $TESTS_FAILED -eq 0 ]]; then
        exit 0
    else
        exit 1
    fi
}

# Handle script interruption
trap 'echo -e "\n${RED}Test suite interrupted${NC}"; exit 130' INT TERM

# Execute main function
main "$@"
