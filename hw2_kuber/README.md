Часть1.

В лабе 1 поднимал Postgres через docker и docker compose, поэтому решил поднять теперь Postgres через kuber.

1. Запуск Kubernetes кластер через Minikube 
<img width="792" height="160" alt="start_minicube" src="https://github.com/user-attachments/assets/67312767-98e8-4e3d-bc6a-e3a982658966" />

Проверил состояние через терминал

kubectl get nodes

<img width="412" height="32" alt="image" src="https://github.com/user-attachments/assets/59edc81b-f058-4a3b-8158-b606da06b890" />

kubectl get all

<img width="596" height="31" alt="kuber_status" src="https://github.com/user-attachments/assets/028bbcaf-1d1b-4f26-bbfb-e38d503c343f" />

и vscode

<img width="231" height="42" alt="start_minicube_vscode" src="https://github.com/user-attachments/assets/77f3fee5-25c3-46ea-9554-350791e671ac" />

Все ок!

P.s. Про Lens помню, еще не дошел

2. Создание YAML-конфигурации

Для Postgres нужен Volume для хранения данных где-то снаружи, нужно как-то задать пароль для БД.
C Volume понятно, нужен был просто тип ресурса PersistentVolumeClaim, в нем особо выставлять ничего не надо.
С паролем не так было очевидно, разбираося как использовать секрет.

Yaml готов, можно запускать

kubectl apply -f postgres.yaml

<img width="339" height="50" alt="image" src="https://github.com/user-attachments/assets/05b3050e-a52d-4e42-b917-cb7651335f41" />

3. Проверка состояния

kubectl get pods

И все плохо

<img width="638" height="39" alt="image" src="https://github.com/user-attachments/assets/70c1ea89-d079-4e80-a850-6e079ee1f16f" />

И стало еще хуже

<img width="597" height="47" alt="image" src="https://github.com/user-attachments/assets/87c61d84-dad4-4442-a4a2-b54a737ae7d9" />

Проблемы было 2:

  - Путь к данным Postgres указал типовой, а он изменился в версиях 18+ (В докере также напоролся)

  kubectl logs <pod_name>
  
  - А секрета то нет, надо добавить

  kubectl create secret generic postgres-secret --from-literal=password=my_password

  <img width="300" height="37" alt="image" src="https://github.com/user-attachments/assets/a0375d24-1567-40b6-a358-01c0b99e72ec" />

4. Рестарт

Если была бы только ошибка с секретом, то можно было бы просто перезапустить

kubectl rollout restart deployment postgres

Из-за директории пришлось все удалить и запустить заново

kubectl delete deployment postgres

kubectl delete service postgres-service

kubectl delete pvc postgres-pvc

kubectl delete -f postgres.yaml

kubectl apply -f postgres.yaml

Проверил состояние сервисов

kubectl get nodes

<img width="484" height="35" alt="image" src="https://github.com/user-attachments/assets/407f2dc3-785c-48b0-9cef-8f80b585a455" />

kubectl get all

<img width="701" height="194" alt="image" src="https://github.com/user-attachments/assets/6fb9e828-73a3-446b-b9e8-01d81cad82ce" />

Все ок!

5. Проверка БД

Проверил подключение к БД

kubectl exec -it deploy/postgres -- psql -U admin -d mydb -c "SELECT version();"

<img width="724" height="73" alt="image" src="https://github.com/user-attachments/assets/994c4f2b-50cd-492a-adab-3925b70cad9d" />

Все ок!

Сервис имеет тип NodePort, Minikube пробросил порт, можно подключиться через клиента

minikube service postgres-service --url

<img width="734" height="35" alt="image" src="https://github.com/user-attachments/assets/c9ac29f6-3131-4cd5-a35d-eb0cd6d8bc17" />

Через указанный порт подключился через pgAdmin

<img width="571" height="175" alt="image" src="https://github.com/user-attachments/assets/121e5fc3-d352-437c-8614-3e0da477f064" />

Удалил под

kubectl delete pod -l app=postgres

Новый под должен был создаться, а данные в БД сохраниться.

Провека

kubectl exec -it deploy/postgres -- psql -U admin -d mydb -c "SELECT * FROM test_table;"

<img width="237" height="67" alt="image" src="https://github.com/user-attachments/assets/1ad946e1-de5c-437b-b346-07af83fa7837" />

Все ок!


 
