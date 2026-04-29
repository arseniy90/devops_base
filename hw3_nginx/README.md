Часть 1.

Пет-проекта с сервером пока нет, поэтому использовал из [дз](https://dev.to/maxcore/basic-ultimate-guide-nginx-1ngd)

/app/server.py

Для запуска серверов и nginx написан Docker Compose file.

В список хостов добавлены домены:

- 127.0.0.1 app1.local

- 127.0.0.1 app2.local

Настройки для перенаправления HTTP в HTTPS +
виртуальные хосты для обслуживания нескольких доменных имен на одном сервере.

<img width="407" height="117" alt="image" src="https://github.com/user-attachments/assets/70f41f00-a077-4db1-b6ae-0ad90464f3fa" />

<br>

Настройка первого домена (для второго аналогично)

- Сервер слушает порт 443(HTTPS), установлены SSL-сертификаты
- Сгенерировал ключ и сертификат с командой:
  
  sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  
  -keyout app1_key.pem -out app1_cert.pem \
  
  -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

  Решил ими не делиться, поэтому здесь их нет

<img width="543" height="147" alt="image" src="https://github.com/user-attachments/assets/7ade1017-d332-409c-9745-e0b6bef0bfc6" />

<br>

- Использование alias для статических файлов

<img width="322" height="77" alt="image" src="https://github.com/user-attachments/assets/2d8c5c30-7ca3-4170-be30-d8310c6d64d9" />

<br>

- Настройка Nginx как Reverse Proxy
- Добавление файла с настройками proxy_params (на основе [дз](https://dev.to/maxcore/basic-ultimate-guide-nginx-1ngd), одинаковый в обоих доменах)

<img width="403" height="96" alt="image" src="https://github.com/user-attachments/assets/4184b70c-61f6-4321-b128-5ae08ea820d1" />

<br>

Запускаем docker compose.... Успех! ( с n-ой попытки...)

<img width="452" height="110" alt="image" src="https://github.com/user-attachments/assets/1b2e5a65-5b31-436e-8934-2e987b413a37" />

<br>

Проверяем домен 1... ок

<img width="222" height="111" alt="image" src="https://github.com/user-attachments/assets/c957d73a-e625-45e6-abb8-d7b7810a2002" />

<br>

Проверяем домен 2... ок

<img width="242" height="172" alt="image" src="https://github.com/user-attachments/assets/d9f41be1-27ec-4139-8c52-bf56954cb865" />

<br>

Часть 2

Проверяем сайт https://ororo.tv/ - сериалы/фильиы в оригинале с субтитрами на разных языках

0. Что вообще покажет простой GET запрос?

<img width="572" height="52" alt="image" src="https://github.com/user-attachments/assets/6ddab6ab-2e2e-4419-be1b-72f862985166" />

<br>

В хедере много всего, включая content-security-policy, permissions-policy и т.д. И это:

<img width="661" height="42" alt="image" src="https://github.com/user-attachments/assets/6913001d-7b2c-405b-a2cc-50d601ff39b0" />

<br>

Cервис для защитить есть.

1. Перебор страниц Fuzzing

ffuf -u https://ororo.tv/FUZZ -w wordlist.txt -mc 200,301

wordlist.txt - здесь типовые имена для подстановки. Сгенерировал список наиболее частых имен для проверки.

В нем уже есть расширения, поэтому опция -e уже была излишней

<img width="880" height="447" alt="image" src="https://github.com/user-attachments/assets/ed5bb2b7-2c34-449b-a5d5-adbc3b05a919" />

<br>

Что-то нашлось..."robots.txt"

Этот файл публичный, используется для web search engine crawlers

<img width="677" height="242" alt="image" src="https://github.com/user-attachments/assets/ece9e3a1-b9fa-4742-ab54-d6af1cdcf0b8" />

<br>

disallow стоит еще разок прогнать с ffuf

ffuf -u https://ororo.tv/FUZZ/users -w wordlist.txt -mc 200,301

<img width="881" height="435" alt="image" src="https://github.com/user-attachments/assets/f81eb154-a60d-48e7-8a86-ffaf9bc1854a" />

<br>

Чисто!

2. path traversal

ffuf -u http://ororo.tv/index.php?page=FUZZ -w traversal.txt -mr "root:x:0:0:"

проверка /etc/passwd с различной глубиной

<img width="906" height="538" alt="image" src="https://github.com/user-attachments/assets/b1e81de9-82dc-4a68-8e0b-13060096d0cc" />

<br>

Все ок, что неудивительно, Cloudfare защищает от таких запросов.

3. Insecure Direct Object Reference (IDOR)

Аккаунт у меня есть, но в URL id не указывается. Возможно id в cookies...

Пробовал определить скрытый URL адрес через dev tools, но не нашел.

Защита IDOR обеспечена на уровне проектирования

