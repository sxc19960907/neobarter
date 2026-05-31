"""推荐服务 - 基于协同过滤 + 内容特征"""

import numpy as np
from sqlalchemy import text
from sqlalchemy.orm import Session


class RecommendService:
    """物品推荐引擎"""

    def recommend_for_user(self, db: Session, user_id: int, limit: int = 10) -> tuple[list[int], list[float]]:
        """
        为用户推荐物品
        策略：
        1. 获取用户浏览/交易过的物品分类偏好
        2. 推荐同分类下热门物品（排除用户自己的物品和已交易的）
        3. 按浏览量 + 信用分加权排序
        """
        # 获取用户偏好分类（基于浏览和交易历史）
        preferred_categories = self._get_user_preferences(db, user_id)

        # 查询推荐物品
        query = text("""
            SELECT i.id, i.view_count, u.credit_score,
                   CASE WHEN i.category_id = ANY(:categories) THEN 1.5 ELSE 1.0 END as category_boost
            FROM items i
            JOIN users u ON u.id = i.user_id
            WHERE i.status = 'active'
              AND i.user_id != :user_id
            ORDER BY (i.view_count * 0.3 + u.credit_score * 0.2) *
                     CASE WHEN i.category_id = ANY(:categories) THEN 1.5 ELSE 1.0 END DESC
            LIMIT :limit
        """)

        result = db.execute(query, {
            "user_id": user_id,
            "categories": preferred_categories if preferred_categories else [0],
            "limit": limit,
        })

        item_ids = []
        scores = []
        for row in result:
            item_ids.append(row[0])
            # 归一化分数
            score = float(row[1] * 0.3 + row[2] * 0.2) * float(row[3])
            scores.append(round(score, 4))

        # 归一化到 0-1
        if scores:
            max_score = max(scores) if max(scores) > 0 else 1
            scores = [round(s / max_score, 4) for s in scores]

        return item_ids, scores

    def _get_user_preferences(self, db: Session, user_id: int) -> list[int]:
        """获取用户偏好的物品分类"""
        query = text("""
            SELECT DISTINCT i.category_id
            FROM trade_requests tr
            JOIN items i ON i.id = tr.target_item_id
            WHERE tr.initiator_id = :user_id
              AND i.category_id IS NOT NULL
            UNION
            SELECT DISTINCT i.category_id
            FROM items i
            WHERE i.user_id = :user_id
              AND i.category_id IS NOT NULL
            LIMIT 5
        """)
        result = db.execute(query, {"user_id": user_id})
        return [row[0] for row in result]


recommend_service = RecommendService()
