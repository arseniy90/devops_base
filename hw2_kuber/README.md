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

Проверил подключение к БД (лучше через имя пода)

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

Часть2.

1. Создание директории Helm-чарта

helm create pg-hc

А также отмечаем, что не надо использовать '_' в именах.

2. Создание YAML-конфигураций с шаблонами

Использовал разные сущности в шаблонах:

- .Values - внешние настройки
  
- .Release - данные об установке
  
- include вместе с _helpers.tpl - повторное использование

Опробовал if условие для корректной установки порта

3. Создание релиза

helm install pg-dev ./pg-hc

<img width="329" height="111" alt="image" src="https://github.com/user-attachments/assets/9db85f0d-e5cd-4daf-9145-abd845901d86" />

4. Проверка БД

minikube service pg-dev-pg-hc

<img width="739" height="226" alt="image" src="https://github.com/user-attachments/assets/f48e3214-476f-46c4-8f99-3b3292340d1a" />

Порт 30001 - значение из values.yaml корректно подставилось в шаблоне

Проверил таблицу в БД, я ведь ее ранее создавал

kubectl exec -it <pod_name> -- psql -U admin -d mydb -c "SELECT * FROM test_table;"

<img width="360" height="64" alt="image" src="https://github.com/user-attachments/assets/e033e08f-8344-4fa3-af19-e2d056a74c5a" />

Подключение к БД есть, а таблицы нет(((

5. Фикс и апгрейд

В volume исправил имя в metadata и слелал upgrade

helm upgrade pg-dev ./pg-hc

и получил это.....

<img width="1107" height="61" alt="image" src="https://github.com/user-attachments/assets/e41c2fa3-af09-432c-b876-014270be38f4" />

Все очень плохо....за что мне это?

В ошибке есть подсказка с причиной, и что нужно исправить, нашел такой вариант:

 - Добавить аннотацию с именем релиза
    
 - Добавить аннотацию с пространством имен
    
 - Добавить метку, что ресурсом управляет Helm

kubectl annotate pvc postgres-pvc meta.helm.sh/release-name=pg-dev --overwrite

kubectl annotate pvc postgres-pvc meta.helm.sh/release-namespace=default --overwrite

kubectl label pvc postgres-pvc app.kubernetes.io/managed-by=Helm --overwritepersistentvolumeclaim/postgres-pvc annotated

Снова upgrade

helm upgrade pg-dev ./pg-hc

<img width="408" height="129" alt="image" src="https://github.com/user-attachments/assets/d403efe9-287f-4ada-8ebd-00d322859140" />

Проверил таблицу в БД, яведь ее ранее создавал

kubectl exec -it <pod_name> -- psql -U admin -d mydb -c "SELECT * FROM test_table;"

<img width="233" height="64" alt="image" src="https://github.com/user-attachments/assets/fc9ad550-1973-489c-a3f7-689cf036c8b9" />

Все ок!

Часть 3

1. Шаблонизация в Helm:  меняем нужные переменные в values.yaml

2. Версионирование и история релизов

helm history pg-dev

3. Удобный откат

Да rollback тоже опробовал 

helm rollback pg-dev 1

<img width="914" height="79" alt="image" src="https://github.com/user-attachments/assets/2677c4f8-4db9-4e66-b7a4-fbbfb9b8ddd8" />



