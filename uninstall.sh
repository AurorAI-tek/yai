#!/bin/bash

BINNAME="${BINNAME:-yai}"
BINDIR="${BINDIR:-/usr/local/bin}"

echo "Uninstallation of Yai ..."
echo

sudo rm $BINDIR/$BINNAME
sudo rm $HOME/.config/yai.json

echo
echo "Uninstallation of Yai complete!"