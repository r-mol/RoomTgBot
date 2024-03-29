# 📚 Project Description
У нас были разные бытовые проблемы (такие как: кто должен принести воду, отчет по тому кто что купил и тд) с соседями по комнате. Из -за этого мы решили телеграмм бота на Golang 🙃.  Теперь мы оптимизировали данные процессы в нашей комнате. А также мы начали соревноваться, кто принес больше воды!!!

## 🎬 Demo
![Screenshot](Menu.png)


## 📝 Usage
    

## ✏️ Features

| Feature                                          | Supported |
|--------------------------------------------------|:---------:|
| Отмечать кто принес воду                         |     ✅     |
| Отмечать кто прибрался в комнате                 |     ✅     |
| Отправлять/получать отчет по покупкам в магазине |     ✅     |
| Отправлять/получать новости для комнаты          |     ✅     |
| Получать статистику по приносу воды              |     ✅     |
| Настройка получения уведомлений от бота          |     ✅     |


## 📝 Getting started
Мы рекомендуем использовать docker, как наиболее простое и независимое
от платформы решение.

### 🎚️ Setup credits
- Откройте папку `env`.
- `*.env.tpl` файлы это файлы шаблоны с необходимыми параметрами.
скопируйте все `<file>.env.tpl` в `<file>.env` и введите секреты для приложения
(телеграм токен, пароль для баз данных). В сумме всего 6 строк. Поделены на 
файлы для вашей безопасности :)

### 🐳 Using Docker
- `docker-compose up` для запуска основного приложения.
- Ну и все :^)

### ⚙️ Configure manually
1. `go build -o bot-exe ./cmd` - компилирует бота
3. `source ./env/*` добавляет переменные из env в шелл.
5. `./bot-exe` - бот подключится к базам данных, которые
Вы указали в файле и начнет работу.

## 🔧 Technologies & Tools
![Code Go](https://img.shields.io/badge/Code-Go-informational?style=flat&logo=Go&logoColor=white&color=2bbc8a)
![DB Mongo](https://img.shields.io/badge/DB-Mongo-informational?style=flat&logo=Mongo&logoColor=white&color=2bbc8a)
![DB Redis](https://img.shields.io/badge/DB-Redis-informational?style=flat&logo=Redis&logoColor=white&color=2bbc8a)
![Linter GoLinter](https://img.shields.io/badge/Linter-GoLinter-informational?style=flat&logo=GoLinter&logoColor=white&color=2bbc8a)

## 🎛️ Hints for developers
- MongoDB `localhost:27017`
- Redis `localhost:6379`
- MongoExpress `localhost:8081`
  - `docker-compose --profile=debug up` для запуска
- После измменений в коде, пересоберите бота
  - `docker-compose build app`

