from celery import Celery

celery_app = Celery(
    "job_applier",
    broker="redis://redis:6379/0",
    backend="redis://redis:6379/1",
)

celery_app.conf.update(
    task_serializer="json",
    accept_content=["json"],
    result_serializer="json",
    timezone="UTC",
    enable_utc=True,
    task_track_started=True,
    worker_max_memory_per_child=512000,  # 512MB, restart worker if exceeded
)
