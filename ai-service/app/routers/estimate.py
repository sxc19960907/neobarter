from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from app.database import get_db
from app.models.schemas import EstimateRequest, EstimateResponse
from app.services.estimate import estimate_service

router = APIRouter()


@router.post("/value", response_model=EstimateResponse)
async def estimate_value(req: EstimateRequest, db: Session = Depends(get_db)):
    """物品估值"""
    estimated, min_val, max_val, confidence = estimate_service.estimate_value(
        db, req.category_id, req.condition, req.title, req.description
    )
    return EstimateResponse(
        estimated_value=estimated,
        min_value=min_val,
        max_value=max_val,
        confidence=confidence,
    )
