

cleanFilebeatData:
	rm meta.json&&rm -fr registry

start:
	filebeat -e  --path.config . --path.data .  -c filebeat.es.format1.yml

start-prod:
	filebeat -e  --path.config . --path.data .  -c filebeat.es.prod.yml

restart: cleanFilebeatData start

restart-prod: cleanFilebeatData start-prod
