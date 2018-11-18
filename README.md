# amazfit
Baixa arquivos do site https://amazfitwatchfaces.com.

## Uso

```
~$ amazfit <paginas>
```

O programa cria uma pasta chamada `file` na mesma pasta que o programa (se não existir) e
baixa os arquivos .wfz e .apk do site da primeira página até o número de páginas informado.
Para cada arquivo .wtz ou .apk, baixa a imagem correspondente. Exemplos de um par de arquivos baixado
na pasta:

```
1950-5be8924a93b06_11112018.wfz
1950-5bef2d9528d49_16112018.png
```
 
O que associa os dois arquivos é o número antes do hífen. Neste caso, 1950. Este número é o id do watch
face cadastrado no site.

Digitando o nome do programa sem a quantidade de páginas:

```
~$ amazfit
```

o programa baixa as 118 primeiras páginas, que é a quantidade de páginas que existia no site em 17/11/2018.