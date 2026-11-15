FROM nginx:1.27-alpine

# Кодировка для корректного отображения русских символов в HTML
ENV LANG=C.UTF-8

# Копируем файлы сайта в веб-папку nginx
COPY index.html /usr/share/nginx/html/index.html
COPY filmotheka.html /usr/share/nginx/html/filmotheka.html

# Копируем конфиг nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# healthcheck: проверяем, что nginx отвечает
HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
  CMD wget -qO- http://localhost/ || exit 1

EXPOSE 80