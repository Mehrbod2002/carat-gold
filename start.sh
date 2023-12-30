systemctl start nginx
chmod -x start.sh

cp meta.timer /etc/systemd/system
cp meta.service /etc/systemd/system
# systemctl daemon-reload
# sudo systemctl enable meta.timer
# sudo systemctl restart meta.timer
# sudo systemctl enable meta.service
# sudo systemctl restart meta.service

mkdir -p /var/www/html
rm -rf /var/www/html/*
cp -r /root/server/meta/chart/* /var/www/html/
cp cert.pem /etc/nginx/ssl
cp key.pem /etc/nginx/ssl
cp nginx.conf /etc/nginx/sites-available/default
systemctl restart nginx