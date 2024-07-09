git pull
docker build -t fun_banking .
docker kill fun_banking_container
docker container prune -f
docker run -d \
       -p 8082:8080 \
       -v $(pwd)/funbanking.sqlite:/app/funbanking.sqlite \
       -v $(pwd)/public/templates:/app/public/templates \
       -v $(pwd)/public/static:/app/public/static \
       -e DATABASE_URL=/app/funbanking.sqlite \
       -e TEMPLATES_PATH=/app/public/templates \
       -e STATIC_PATH=/app/public/static \
       --name fun_banking_container fun_banking
