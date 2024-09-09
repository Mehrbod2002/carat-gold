pm2 kill
pm2 start ~/carat-gold/service/carat-gold --name app
pm2 start ~/carat-gold/service/history/wsgi.py --name history
pm2 start ~/carat-gold/service/socket/socket --name socket
pm2 start ~/carat-gold/service/metatrader/metatrader --name metatrader
cd ~/carat-user && pm2 start npm --name "ui" -- run start && cd
cd ~/carat-admin && pm2 start npm --name "admin" -- run start && cd
pm2 save