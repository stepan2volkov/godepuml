# GoDepUML

Небольшой скрипт для создания простой диаграммы зависимостей на PlantUML.

## Использование

```bash
godepuml -p <path-to-project> -o <output-puml-file> [exlude-1 ... exclud-N]
```

## Пример

Для проекта [csvdb](https://github.com/stepan2volkov/csvdb) будет сгенерирован файл со следующим содержанием:
```puml
@startuml 'github.com/stepan2volkov/csvdb'

[cmd.csvdb] --> [internal.app]
[cmd.csvdb] --> [internal.app.table]
[cmd.csvdb] --> [internal.app.table.formatter]
[cmd.csvdb] --> [internal.app.table.loader]
[internal.app.parser] --> [internal.app.table]
[internal.app.parser] --> [internal.app.table.operation]
[internal.app.parser] --> [internal.app.scanner]
[internal.app] --> [internal.app.parser]
[internal.app] --> [internal.app.scanner]
[internal.app] --> [internal.app.table]
[internal.app.table.formatter] --> [internal.app.table]
[internal.app.table.value] --> [internal.app.table]
[internal.app.table.loader] --> [internal.app.table]
[internal.app.table.loader] --> [internal.app.table.value]
[internal.app.table.operation] --> [internal.app.table]

@enduml
```

Для генерации изображения можно установить puml локально, или воспользоваться [сервером](http://www.plantuml.com/plantuml/uml).

![csvdb logo](/images/csvdb.png)