git pull
docker build -t fun_banking .
docker kill fun_banking_container
docker container prune -f
docker run -d \
       -p 8082:8080 \
       -v $(pwd)/fun_banking.db:/app/fun_banking.db \
       -e DATABASE_URL=/app/fun_banking.db \
       --name fun_banking_container fun_banking
