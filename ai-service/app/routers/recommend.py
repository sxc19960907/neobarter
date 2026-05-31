from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from app.database import get_db
from app.models.schemas import RecommendRequest, RecommendResponse
from app.services.recommend import recommend_service

router = APIRouter()


@router.post("/items", response_model=RecommendResponse)
async def recommend_items(req: RecommendRequest, db: Session = Depends(get_db)):
    """为用户推荐物品"""
    item_ids, scores = recommend_service.recommend_for_user(db, req.user_id, req.limit)
    return RecommendResponse(item_ids=item_ids, scores=scores)
