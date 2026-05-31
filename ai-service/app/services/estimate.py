"""估值服务 - 基于历史交易数据统计"""

import numpy as np
from sqlalchemy import text
from sqlalchemy.orm import Session


class EstimateService:
    """物品估值引擎"""

    # 成色折扣系数
    CONDITION_FACTOR = {
        "new": 1.0,
        "like_new": 0.85,
        "good": 0.7,
        "fair": 0.5,
    }

    def estimate_value(
        self, db: Session, category_id: int, condition: str, title: str, description: str | None = None
    ) -> tuple[float, float, float, float]:
        """
        估算物品价值
        返回: (estimated_value, min_value, max_value, confidence)
        策略：
        1. 查询同分类已完成交易的物品估值
        2. 按成色调整
        3. 返回中位数作为推荐值，P25/P75 作为范围
        """
        # 查询同分类历史估值
        query = text("""
            SELECT i.estimated_value
            FROM items i
            WHERE i.category_id = :category_id
              AND i.estimated_value > 0
              AND i.status IN ('active', 'traded')
            ORDER BY i.created_at DESC
            LIMIT 100
        """)

        result = db.execute(query, {"category_id": category_id})
        values = [float(row[0]) for row in result if row[0]]

        if not values:
            # 没有历史数据，返回默认值
            return 100.0, 50.0, 200.0, 0.1

        values_array = np.array(values)
        condition_factor = self.CONDITION_FACTOR.get(condition, 0.7)

        median_value = float(np.median(values_array)) * condition_factor
        p25 = float(np.percentile(values_array, 25)) * condition_factor
        p75 = float(np.percentile(values_array, 75)) * condition_factor

        # 置信度基于样本量
        confidence = min(len(values) / 50.0, 1.0)

        return (
            round(median_value, 2),
            round(p25, 2),
            round(p75, 2),
            round(confidence, 2),
        )


estimate_service = EstimateService()
