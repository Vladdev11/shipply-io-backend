#!/bin/bash
systemctl stop shipply-client-backend
rm /opt/shipply-client-backend/shipply-client-backend
cp ../bin_release/shipply-client-backend /opt/shipply-client-backend/shipply-client-backend
chmod +x /opt/shipply-client-backend/shipply-client-backend
systemctl start shipply-client-backend