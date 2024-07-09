git pull
docker build -t fun_banking .
docker kill fun_banking_container
docker container prune -f
docker run -d \
       -p 8082:8080 \
       -v $(pwd)/funbanking.sqlite:/app/funbanking.sqlite \
       -v $(pwd)/templates:/app/templates \
       -e DATABASE_URL=/app/funbanking.sqlite \
       -e TEMPLATES_PATH=/app/templates \
       --name fun_banking_container fun_banking
