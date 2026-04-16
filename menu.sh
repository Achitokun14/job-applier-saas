#!/bin/bash
# ============================================================
#  Job Applier SaaS - Main Menu
#  Compatible with Bash (Linux/macOS/Git Bash on Windows)
# ============================================================

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

clear_screen() {
    clear 2>/dev/null || echo -e "\033[2J\033[H"
}

show_header() {
    echo -e "${PURPLE}"
    echo "============================================================"
    echo "          ____  _             _    _                         "
    echo "         |  _ \\| |           | |  | |                        "
    echo "         | |_) | |_   _  __ _| | _| |     ___   __ _        "
    echo "         |  _ <| | | | |/ _\` | |/ /\ |    / _ \\ / _\` |      "
    echo "         | |_) | | |_| | (_| |   <| |___| (_) | (_| |      "
    echo "         |____/|_|\\__,_|\\__,_|_|\\_\\______\\___/ \\__, |      "
    echo "                                                __/ |      "
    echo "                                               |___/       "
    echo "              Automated Job Application SaaS              "
    echo "============================================================"
    echo -e "${NC}"
}

show_menu() {
    echo -e "${CYAN}Main Menu:${NC}"
    echo ""
    echo -e "  ${GREEN}=== Run All Services ===${NC}"
    echo "  1) Start All (Local)"
    echo "  2) Start All (Docker)"
    echo "  3) Stop All Services"
    echo ""
    echo -e "  ${BLUE}=== Individual Services ===${NC}"
    echo "  4) Start Backend Only"
    echo "  5) Start Frontend Only"
    echo "  6) Start Python Service Only"
    echo "  7) Start TUI Only"
    echo ""
    echo -e "  ${YELLOW}=== Build & Install ===${NC}"
    echo "  8) Install All Dependencies"
    echo "  9) Build All"
    echo " 10) Build Backend"
    echo "  11) Build Frontend"
    echo "  12) Build TUI"
    echo ""
    echo -e "  ${RED}=== Docker ===${NC}"
    echo " 13) Docker Build"
    echo " 14) Docker Logs"
    echo " 15) Docker Stop"
    echo ""
    echo -e "  ${PURPLE}=== Utilities ===${NC}"
    echo "  0) Exit"
    echo ""
    echo -n "Select option: "
}

run_script() {
    local script="$1"
    if [ -f "$script" ]; then
        bash "$script"
    elif [ -f "${script}.sh" ]; then
        bash "${script}.sh"
    else
        echo -e "${RED}Script not found: $script${NC}"
    fi
}

while true; do
    clear_screen
    show_header
    show_menu
    read -r choice

    case $choice in
        1)
            run_script "$SCRIPT_DIR/scripts/start-all.sh"
            ;;
        2)
            run_script "$SCRIPT_DIR/scripts/start-docker.sh"
            ;;
        3)
            run_script "$SCRIPT_DIR/scripts/stop-all.sh"
            ;;
        4)
            run_script "$SCRIPT_DIR/scripts/backend/start.sh"
            ;;
        5)
            run_script "$SCRIPT_DIR/scripts/frontend/start.sh"
            ;;
        6)
            run_script "$SCRIPT_DIR/scripts/python-service/start.sh"
            ;;
        7)
            run_script "$SCRIPT_DIR/scripts/tui/start.sh"
            ;;
        8)
            run_script "$SCRIPT_DIR/scripts/install-all.sh"
            ;;
        9)
            echo "Building all services..."
            run_script "$SCRIPT_DIR/scripts/backend/build.sh"
            run_script "$SCRIPT_DIR/scripts/frontend/build.sh"
            run_script "$SCRIPT_DIR/scripts/tui/build.sh"
            echo "All builds complete!"
            ;;
        10)
            run_script "$SCRIPT_DIR/scripts/backend/build.sh"
            ;;
        11)
            run_script "$SCRIPT_DIR/scripts/frontend/build.sh"
            ;;
        12)
            run_script "$SCRIPT_DIR/scripts/tui/build.sh"
            ;;
        13)
            run_script "$SCRIPT_DIR/scripts/docker/build.sh"
            ;;
        14)
            run_script "$SCRIPT_DIR/scripts/docker/logs.sh"
            ;;
        15)
            run_script "$SCRIPT_DIR/scripts/docker/stop.sh"
            ;;
        0)
            echo -e "${GREEN}Goodbye!${NC}"
            exit 0
            ;;
        *)
            echo -e "${RED}Invalid option. Press Enter to continue...${NC}"
            read -r
            ;;
    esac

    echo ""
    echo -e "${CYAN}Press Enter to return to menu...${NC}"
    read -r
done
