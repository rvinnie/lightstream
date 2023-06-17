if [ ! -f /etc/ssl/certs/nginx.crt ]; then
  echo "Nginx: setting up ssl ...";
  openssl req -x509 -nodes -days 365 -newkey rsa:4096 \
      -keyout /etc/ssl/private/nginx-selfsigned.key \
      -out /etc/ssl/certs/nginx-selfsigned.crt \
      -subj "/C=RU/ST=Moscow/L=Moscow/O=lightstream/CN=lightstream.ru";
  echo "Nginx: ssl is set up!";
fi
