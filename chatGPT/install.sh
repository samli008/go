mv gpt /opt/
cat > /usr/lib/systemd/system/gpt.service << EOF
[Unit]
Description=chatGPT
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt
ExecStart=/opt/gpt
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now gpt
systemctl status doc
ss -nlpt |grep 8181
