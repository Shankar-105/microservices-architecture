FROM python:3.12-slim

WORKDIR /workspace
COPY catalog-service-py/requirements.txt /tmp/requirements.txt
RUN pip install --no-cache-dir -r /tmp/requirements.txt
COPY catalog-service-py /workspace/catalog-service-py
COPY generated-py /workspace/generated-py

WORKDIR /workspace/catalog-service-py

EXPOSE 50053 8001
CMD ["python", "-m", "app.main"]
