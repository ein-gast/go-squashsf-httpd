# squashsf-httpd

Специализированный HTTP-сервер, который умеет отдавать файлы из образа SquashFS, не распаковывая и не монтируя его.

Сервер написан на golang. Для работы со SquashFS используется пакет [github.com/diskfs/go-diskfs](https://github.com/diskfs/go-diskfs).

## Как собрать из исходников

Нужен golang>=1.24:
```bash
git clone https://github.com/ein-gast/go-squashsf-httpd.git
make all
stat squashfs-httpd
```

Можно собрать через докер (не требуется golang на хост-машине):
```bash
git clone https://github.com/ein-gast/go-squashsf-httpd.git
make dockerbuild
stat squashfs-httpd.bin
```

## Как установить и настроить

Если вы использовали `make all` или скачали сборку, то испольниый файл `squashfs-httpd` можно положить куда вам угодно.

Если использовать пакетный менеджер golang, то постаить сервер можно так:
```bash
go install github.com/ein-gast/go-squashsf-httpd/cmd/squashsf-httpd@latest
```

Простейший способ запуска для раздачи фйлов из одного образа squashfs:

```bash
./squashfs-httpd -host 127.0.0.1 -port 8080 -squash ./examples/data/potree-lion.sq
```

Такая команда выдаст по адресу `http://127.0.0.1:8080/index.html` облако точек, взятое из примеров проекта [Potree](https://github.com/potree/potree). 

Подробнее о параметрах: `./squashfs-httpd --help`

Чтобы раздавать контент из нескольких файлов нужно написать конфигурационный файл и запустить сервер с ним:

```bash
./squashfs-httpd -config squashfs-httpd.yaml
```

Пример конфигурации:

```yaml
#  -- squashfs-httpd.yaml --
# tcp адрес и порт на котором сервер приниает соединнеия
bind_addr: 127.0.0.1
bind_port: 8080
# кодировка текстовых файлов, будет добавлена к "content-type: text/...; charset=..."
charset: utf-8
# размер буфера чтения из SquashFS-файла, обработка каждого запроса с 200-м ответом создаст такой буфер
buffer: 10240
# пути до log-файлов, можно указыввать /dev/stdout, /dev/stderr, /dev/null
# относительные пути логгов будут построены от каталога конфигурационного фала
error_log: "./var/logs/error.log"
access_log: "./var/logs/access.log"
# отключает запись в access_log
access_log_off: false
# если pid-файл не указан, то путь будет выбран автоматически
#pid_file: "/run/squashfs-httpd.pid"
pid_file_off: false
# через столько секунд соединение будет разорвано, если не получилось считать или записать в него данные
client_timeout: 5.0
# способы подключения файлов:
#  squash    - один файл с образом SquashFS
#  squashdir - каталог, в котором лежат образы SquashFS
# относительные пути будут построены от каталога конфигурационного фала
routes:
  # файл index.html из архива potree-lion.sq
  # будет доступен как http://127.0.0.1:8080/one/index.html
  - prefix: /one/
    squash: ./examples/data/potree-lion.sq
  # файл index.html из архива potree-lion.sq
  # будет доступен как http://127.0.0.1:8080/two/potree-lion.sq/index.html
  - prefix: /two/
    squashdir: ./examples/data/
```

По сигналу **USR1** сервер переоткрывает логи. По сигналу **USR2** сервер "отпускает" файлы, открытые из роутов `squashdir`.

## Как использовать

**squashsf-httpd** создан чтобы работать в контейнерах совместно с nginx и раздавать запакованные в SquashFS папки с большим количеством мелких файлов.

Главная область применения — тайловые кэши ортофотопланов и оптимизированные облака точек. Для многих применений данные «тайлятся» один раз и просто хранятся в виде каталога с сотнями тысяч или миллионами подкаталогов и файлов. Такие каталоги неудобны в эксплуатации: они расходуют inode'ы, их неудобно бэкапить или переносить в другое место — много времени уходит на построение списка файлов. Удобнее упаковать тайловый кэш в архив и примонтировать, для этого идеален SquashFS (можно взять zip + fuse-zip, но SquashFS позволяет тоньше настроить сжатие).

Смонтированные архивы начинают доставлять проблемы, если идёт постоянный приток данных (некоторые промышленные объекты снимают раз в день) — приходится постоянно подмонтировать новые файлы. Задачу монтирования можно возложить на приложение, но тогда нужно следить за его правами, а если приложение запущено в непривилегированном контейнере, то даже через FUSE ничего не смонтируешь. Из этой цепочки затруднений (много мелких файлов, squashfs, постоянное пополнение данных, приложения в непривилегированных контейнерах) появился запрос на простой HTTP-сервер, который выдаёт файлы сразу из SquashFS, не монтируя образ. Это основное назначение **squashfs-httpd**.

Простейший пример связки nginx + squashfs-httpd в докере есть в папке[./examples/docker/](./examples/docker/)

## Ссылки
- https://github.com/plougher/squashfs-tools
- https://github.com/diskfs/go-diskfs
- https://github.com/CalebQ42/squashfs
- https://github.com/h2non/filetype
- https://github.com/potree/potree
