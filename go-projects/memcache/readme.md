key -> max 255 bytes
value -> max 1MB

RAM -> donner l'option de setup le max de mémoire utilisée
Eviction -> LRU (Least Recently Used) -> si on atteint la limite de mémoire, on supprime les données les moins utilisées. Quoi d'autre? 

Flush -> Print all keys and values in a file or whatever
Monitoring -> Comment on peut monitorer le nombre de requêtes, le temps de réponse, etc
server tcp/udp
client tcp/udp

On compute un hash de la clé pour savoir sur quel serveur on doit aller.
On peut avoir plusieurs serveurs.
On recompute la clé sur le serveur pour la store.

Authentication à prévoir dans le protocole binaire
