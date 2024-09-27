## Порядок команд:

1. `eval $(minikube docker-env)` - чтобы minikube видел локальные образы
2. `make build` - создал docker образы контейнеров
3. `kubectl apply -f charts/templates/calendar-deployment.yaml `
4. `kubectl logs calendar-deployment-5489c97d56-6ffx6` - посмотреть логи.
5. `minikube addons enable ingress`
6. `kubectl apply -f charts/templates/service.yaml `
7. `kubectl apply -f charts/templates/ingress.yaml `
8. `echo -e "127.0.0.1 calendar.localhost" | sudo tee -a /etc/hosts`
9. `minikube tunnel`

>статья: https://habr.com/ru/articles/752586/



