from fastapi import FastAPI
from app.routers import recommend, estimate, similar

app = FastAPI(
    title="NeoBarter AI Service",
    description="智能推荐、物品估值、相似物品推荐",
    version="1.0.0",
)

app.include_router(recommend.router, prefix="/recommend", tags=["推荐"])
app.include_router(estimate.router, prefix="/estimate", tags=["估值"])
app.include_router(similar.router, prefix="/similar", tags=["相似"])


@app.get("/health")
async def health():
    return {"status": "ok"}
