# lotterey

Чтобы протестировать работу сервиса необходимо:
  1. Запустить сервер через Makefile в термминале введите ```make start```. Убедитесь что сервер создался на порте 8080, потому что это значение вписано в `test.sh`
  2. Затем исполните `test.sh`. В терминал будет выведены результаты работы. Также изменения будут применены к lottery.db

Если запустить скрипт повторно, то он упадет с ошибкой (что тоже в своем роде является тестом того, что система работает правильно). Это происходит потому что после прохождения тестов не процесса уборки тестовых данных. Это было бы хорошо сделать, но я сконцетрировался на количестве и качестве тестов.

Ниже скриншоты выполнения `test.sh`.

<img width="513" height="1010" alt="Screenshot_2025-13-08_1755104503" src="https://github.com/user-attachments/assets/31957360-13a0-46e3-849d-6ccfebe0bf90" />
<img width="556" height="867" alt="Screenshot_2025-13-08_1755104725" src="https://github.com/user-attachments/assets/57aaea16-c533-428f-bb6d-cfa099cc055a" />
<img width="492" height="1042" alt="Screenshot_2025-13-08_1755104763" src="https://github.com/user-attachments/assets/de1ad9a5-d549-45a8-8d1d-abedeab4ff30" />
<img width="480" height="848" alt="Screenshot_2025-13-08_1755104787" src="https://github.com/user-attachments/assets/3c76448b-1fee-499e-9828-e36ac0c62dde" />
<img width="423" height="308" alt="Screenshot_2025-13-08_1755104811" src="https://github.com/user-attachments/assets/475cd451-a20b-42c3-b15f-1ba67fc63c31" />
<img width="527" height="932" alt="Screenshot_2025-13-08_1755104817" src="https://github.com/user-attachments/assets/119d7b95-14f3-4020-8678-549ccda80f36" />
<img width="590" height="922" alt="Screenshot_2025-13-08_1755104830" src="https://github.com/user-attachments/assets/bd831a5e-02ed-4de9-9c3f-311ff54ee46e" />
