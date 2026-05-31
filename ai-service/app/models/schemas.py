from pydantic import BaseModel
from typing import Optional


class RecommendRequest(BaseModel):
    user_id: int
    limit: int = 10


class RecommendResponse(BaseModel):
    item_ids: list[int]
    scores: list[float]


class EstimateRequest(BaseModel):
    category_id: int
    condition: str
    title: str
    description: Optional[str] = None


class EstimateResponse(BaseModel):
    estimated_value: float
    min_value: float
    max_value: float
    confidence: float


class SimilarRequest(BaseModel):
    item_id: int
    limit: int = 10


class SimilarResponse(BaseModel):
    item_ids: list[int]
    scores: list[float]
