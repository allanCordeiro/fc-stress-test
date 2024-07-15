# fc-stress-test

Sistema CLI em go para realizar testes de carga em um serviço web.

A imagem se encontra no docker hub. 

`docker pull allancordeiros/fc-stress-test:latest`


Para executar disparar o comando docker run, como abaixo:

`docker run allancordeiros/fc-stress-test --url=[url desejada] --requests=[quantidade de requests] --concurrency=[chamadas simultaneas]`

Sendo os parametros:

- URL (Obrigatório): a url desejada
- Requets: quantidade de chamadas a URL. Caso seja omitido, será por padrão 100.
- Concurrency: chamadas simultâneas. Ao ser omitido assume por padrão 10.
