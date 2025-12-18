package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"uttc-hackathon-backend/internal/repository"
)

type VertexGenerativeClient interface {
	GenerateContent(ctx context.Context, modelName string, prompt string, config repository.GenerationConfig) (string, error)
}

type SuggestionService struct {
	vertexRepo VertexGenerativeClient
}

func NewSuggestionService(vertexRepo VertexGenerativeClient) *SuggestionService {
	return &SuggestionService{
		vertexRepo: vertexRepo,
	}
}

// GetListingSuggestion generates a suggestion for a listing based on the provided title and description.
func (s *SuggestionService) GetListingSuggestion(ctx context.Context, title, description, condition, language string) []string {
	englishInstruction := `You are an expert e-commerce listing assistant and taxonomist. Your goal is to analyze a P2P marketplace listing and suggest specific metadata **field names** that the seller should define to make their listing more searchable and trustworthy.

**Guidelines:**
1.  **Analyze the Item:** Determine the product category based on the Title and Description.
2.  **Contextual Awareness:** Consider the "Condition". If an item is "Poor", fields like "Defects" or "Battery Health" might be relevant. If "New", ignore wear-related fields.
3.  **Prioritize High Value:** Suggest fields that are standard for that category (e.g., "Size" for clothes, "Storage" for electronics, "ISBN" for books). Prioritize attributes that are missing or crucial for filtering.
4.  **Constraint:** Return **only** a JSON array of strings. The strings should be in **English**.
5.  **Quantity:** Return between 0 and 3 items. If the listing is too vague to categorize, return an empty array "[]".

**Output Format:**
Strictly valid JSON. Example: "["Brand", "Size", "Material"]"`
	englishPrompt := `Please analyze the following listing data and provide the JSON array of suggested fields.

**Listing Data:**
* **Title:** %s
* **Description:** %s
* **Condition:** %s`
	japaneseInstruction := `あなたはフリマアプリ（C2Cマーケットプレイス）の出品支援および商品分類の専門家です。あなたの目的は、出品情報を分析し、検索性を高め信頼を得るために出品者が埋めるべき具体的な「追加情報の項目名（フィールド名）」を提案することです。

**ガイドライン:**
1.  **商品を分析する:** 商品名（Title）と商品説明（Description）から、商品のカテゴリを特定してください。
2.  **文脈を考慮する:** 商品状態（Condition）を考慮してください。「中古」や「傷あり」の場合は「ダメージ箇所」や「バッテリー最大容量」などが重要になる場合があります。「新品」の場合は使用感に関する項目は避けてください。
3.  **重要な項目を優先する:** そのカテゴリで一般的かつ重要な項目（例：服なら「サイズ」「着丈」、家電なら「ストレージ容量」「年式」、本なら「ISBN」など）を優先してください。説明文にすでに詳しく書かれていることよりも、検索フィルタリングに使われるような構造化データを優先します。
4.  **制約事項:** 文字列の JSON 配列のみを返してください。提案する項目名は**日本語**で出力してください。
5.  **数量:** 提案数は 0個から最大3個までです。情報が少なすぎてカテゴリが特定できない場合は、空の配列 "[]" を返してください。

**出力フォーマット:**
有効なJSON形式のみを出力してください。
例: "["ブランド", "サイズ", "素材"]"`
	japanesePrompt := `以下の出品データを分析し、提案する項目のJSON配列を出力してください。

**出品データ:**
* **商品名 (Title):** %s
* **商品説明 (Description):** %s
* **商品状態 (Condition):** %s`

	var prompt, systemInstruction string
	if language == "ja" {
		prompt = fmt.Sprintf(japanesePrompt, title, description, condition)
		systemInstruction = japaneseInstruction
	} else {
		prompt = fmt.Sprintf(englishPrompt, title, description, condition)
		systemInstruction = englishInstruction
	}

	temperature := float32(0.2)
	config := repository.GenerationConfig{
		SystemInstruction: systemInstruction,
		Temperature:       &temperature,
		JsonResponse:      true,
	}
	// gemini-1.5-flash-002 is no longer available
	respStr, err := s.vertexRepo.GenerateContent(ctx, "gemini-2.5-flash", prompt, config)
	if err != nil {
		log.Printf("failed to generate suggestion: %v", err)
		return []string{}
	}

	var suggestions []string
	if err := json.Unmarshal([]byte(respStr), &suggestions); err != nil {
		log.Printf("failed to unmarshal suggestion response: %v, response: %s", err, respStr)
		return []string{}
	}

	return suggestions
}
