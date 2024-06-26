# rate-limiter
O rate limiter é utilizado para limitar o número máximo de requisições por intervalo de tempo com base em um endereço IP e/ou em um token de acesso. 
### Configuraçao
#### No arquivo .env que está na pasta cmd, são feitas as configurações da rate limiter. Nele, temos as seguintes variáveis:
##### IP_LIMIT: número máximo de requisições por IP dentro de um intervalo de tempo.
##### IP_EXPIRATION_INTERVAL: tempo de bloqueio do endereço IP que atingiu o limite de requisições (variável IP_LIMIT). Durante esse período o limiter bloqueará qualquer requisição que tenha o IP bloqueado como origem. Após esse intervalo, novas requisições são permitidas.
##### IP_WINDOW_INTERVAL: intervalo de tempo no qual são permitidas novas requisições por IP, até o limite definido na variável IP_LIMIT. Após esse intervalo, novas requisições são permitidas.
##### TOKENS: lista de tokens (separados por vírgula) válidos que têm permissão para fazer as requisições.
##### TOKENS_LIMIT: lista com o número máximo de requisições por token dentro de um intervalo de tempo. Deve seguir a mesma ordem da lista de tokens atribuída à variável TOKENS. Valores separados pro vírgula.
##### TOKENS_EXPIRATION_INTERVAL: lista com o tempo de bloqueio do token que atingiu o limite de requisições (variável TOKENS_LIMIT). Durante esse período o limiter bloqueará qualquer requisição ASSINADA COM O TOKEN BLOQUEADO. Após esse intervalo, novas requisições são permitidas.
##### TOKEN_WINDOW_INTERVAL: intervalo de tempo no qual são permitidas novas requisições por token, até o limite definido na variável TOKENS_LIMIT. Após esse intervalo, novas requisições são permitidas. Deve seguir a mesma ordem da lista de tokens atribuída à variável TOKENS. Valores separados pro vírgula.
##### LIMIT_BY: variável que permite configurar o rate limiter para limitar as requisições por IP ou por token ou por IP e token (nesse caso, o token tem prioridade sobre o IP). Valores possíveis: IP (limitar por IP); API_KEY (limitar apenas por token); IP,API_KEY (limitar por token e IP).

#### Exemplo dde configuração do .env
IP_LIMIT=5\
IP_EXPIRATION_INTERVAL=2s\
IP_WINDOW_INTERVAL=1s\
TOKENS=abc123,def456,ghi789\
TOKENS_LIMIT=2,5,15\
TOKENS_EXPIRATION_INTERVAL=2s,2s,5s\
TOKEN_WINDOW_INTERVAL=1s\
LIMIT_BY=IP,API_KEY

##### Essa configuração de exemplo define que o rate limiter deverá limitar as requisições por token e IP (LIMIT_BY=IP,API_KEY).
##### São permitidas 5 req/s (IP_LIMIT=5 e IP_WINDOW_INTERVAL=1s) para um dado endereço IP e, uma vez atingido esse limite, o IP será bloqueado por 2 segundos (IP_EXPIRATION_INTERVAL=2s).
##### São permitidas 2 req/s para o token abc123, 5 req/s para o token def456 e 15 req/s para o token ghi789 (TOKENS=abc123,def456,ghi789 e TOKENS_LIMIT=2,5,15 e TOKEN_WINDOW_INTERVAL=1s). Os tempos de bloqueio (TOKENS_EXPIRATION_INTERVAL=2s,2s,5s) são de 2 segundos, 2 segundos e 5 segundos para os tokens abc123, def456 e ghi789, respectivamente.

## Testes
### Para realizar os testes é necessário estar com o docker rodando, pois é utilizado o testcontainers com uma instância do Redis.
