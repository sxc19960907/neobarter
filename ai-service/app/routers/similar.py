from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from app.database import get_db
from app.models.schemas import SimilarRequest, SimilarResponse
from app.services.similar import similar_service

router = APIRouter()


@router.post("/items", response_model=SimilarResponse)
async def find_similar(req: SimilarRequest, db: Session = Depends(get_db)):
    """查找相似物品"""
    item_ids, scores = similar_service.find_similar(db, req.item_id, req.limit)
    return SimilarResponse(item_ids=item_ids, scores=scores)
