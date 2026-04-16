# Deployment Log

## Project Profile
- **Name**: job-applier-saas
- **Stack**: Go (backend), Python/FastAPI (AI service), SvelteKit/Bun (frontend), PostgreSQL, Redis, Celery
- **Domain**: jobs.murgana.online
- **VPS IP**: 72.167.141.90
- **Services**: backend (Go:8080), frontend (SvelteKit:3000), python-service (FastAPI:8001), celery-worker, postgres, redis

## Progress
- [x] Phase 1: Dockerfiles verified/fixed for production
- [x] Phase 2: docker-compose.coolify.yml created with Traefik labels
- [x] Phase 3: Gitea repo created and code pushed
- [ ] Phase 4: DNS verified (jobs.murgana.online -> 72.167.141.90) -- WAITING for propagation
- [ ] Phase 5: Deployed on VPS via Coolify
- [ ] Phase 6: HTTPS verified
- [ ] Phase 7: All health checks passing
- [ ] Phase 8: All pages tested from browser

## Notes
- Existing project on VPS: murgana-scrapers on murgana.online (DO NOT touch)
- Traefik slug: jobapplier (globally unique)
- Memory budget: ~3.5GB (within 8-10GB available)
