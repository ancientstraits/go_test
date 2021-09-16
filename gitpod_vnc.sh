#!/bin/sh

# This script easily gives gitpod graphical 
# capabilities by using tightvncserver and noVNC.
# To be lightweight, it also uses awesomeWM as the
# graphical environment, but you can change it
# by editing this script.

# Package management
sudo apt update

sudo apt install awesome tightvncserver xterm

# awesomeWM
mkdir ~/.config/awesome
cp /etc/xdg/awesome/rc.lua ~/.config/awesome/rc.lua

# Run vncserver once to setup username and password
vncserver
vncserver -kill :1

# Change the startup script to start awesomeWM
rm -rf ~/.vnc/xstartup
cat << EOF > ~/.vnc/xstartup
#!/bin/sh
awesome &
EOF
# You can configure vncserver with the -geometry option,
# which will change the screen size.
# This gives a display with a width of 600
# and a height of 400.
vncserver -geometry 600x400

# novnc
git clone https://github.com/novnc/novnc ~/.novnc
~/.novnc/utils/novnc_proxy --vnc localhost:5901

# Gitpod will notify you that 3 ports are open.
# Make the 6080 port public, then open it in a new tab or
# the preview by clicking the "Remote Explorer" sidebar
# button on the side of Gitpod Code.