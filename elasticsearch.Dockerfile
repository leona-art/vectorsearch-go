FROM elasticsearch:8.13.4

RUN elasticsearch-plugin install analysis-kuromoji

EXPOSE 9200