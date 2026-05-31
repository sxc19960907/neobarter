"""相似物品服务 - 基于 TF-IDF 文本相似度"""

from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity
from sqlalchemy import text
from sqlalchemy.orm import Session


class SimilarService:
    """相似物品推荐引擎"""

    def find_similar(self, db: Session, item_id: int, limit: int = 10) -> tuple[list[int], list[float]]:
        """
        查找相似物品
        策略：
        1. 获取目标物品的标题+描述
        2. 获取同分类的候选物品
        3. 用 TF-IDF + 余弦相似度计算文本相似度
        4. 返回最相似的 N 个
        """
        # 获取目标物品
        target_query = text("""
            SELECT id, title, description, category_id, user_id
            FROM items WHERE id = :item_id
        """)
        target_row = db.execute(target_query, {"item_id": item_id}).fetchone()
        if not target_row:
            return [], []

        target_id, target_title, target_desc, category_id, target_user_id = target_row
        target_text = f"{target_title} {target_desc or ''}"

        # 获取候选物品（同分类 + 上架中 + 非自己的）
        candidates_query = text("""
            SELECT id, title, description
            FROM items
            WHERE status = 'active'
              AND id != :item_id
              AND user_id != :user_id
              AND (category_id = :category_id OR category_id IS NULL)
            ORDER BY created_at DESC
            LIMIT 200
        """)
        candidates = db.execute(candidates_query, {
            "item_id": item_id,
            "user_id": target_user_id,
            "category_id": category_id,
        }).fetchall()

        if not candidates:
            return [], []

        # 构建文本列表
        candidate_ids = [row[0] for row in candidates]
        candidate_texts = [f"{row[1]} {row[2] or ''}" for row in candidates]

        # TF-IDF 向量化
        all_texts = [target_text] + candidate_texts
        try:
            vectorizer = TfidfVectorizer(max_features=5000, analyzer="char_wb", ngram_range=(2, 4))
            tfidf_matrix = vectorizer.fit_transform(all_texts)
        except ValueError:
            return [], []

        # 计算余弦相似度
        similarities = cosine_similarity(tfidf_matrix[0:1], tfidf_matrix[1:]).flatten()

        # 排序取 top N
        top_indices = similarities.argsort()[::-1][:limit]

        item_ids = [candidate_ids[i] for i in top_indices if similarities[i] > 0]
        scores = [round(float(similarities[i]), 4) for i in top_indices if similarities[i] > 0]

        return item_ids, scores


similar_service = SimilarService()
