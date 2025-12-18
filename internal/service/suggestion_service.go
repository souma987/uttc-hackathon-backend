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
	englishInstruction := `You are an expert in listing support and product classification for flea market apps (C2C marketplaces). Your objective is to analyze listing information and suggest specific "additional information field names" that the seller should fill in to enhance searchability and gain trust.

**Guidelines:**
1.  **Analyze the product:** Identify the product category based on the Product Name (Title) and Product Description (Description).
2.  **Prioritize important fields:** Prioritize objective fields that have significant meaning in that specific category (e.g., "Length" or "Fabric" for clothes, "Mileage" for cars, "CPU Model" or "OS Version" for PCs, "ISBN" or "Publication Year" for books), in addition to general content.
3.  **Exclude existing fields:** Even if a field is deemed important, strictly exclude it from the output if the field name or the content of that field is already present in the Title or Description.
4.  **Constraints:** Return only a JSON array of strings. Output the suggested field names in **English**.
5.  **Output count:** Suggest 4 items at maximum. If the information is too sparse to identify the category, return an empty array "[]".

**Output Format:**
Output only a valid JSON format.
Example response for a block toy listing: "["Total Piece Count", "Brand", "Missing Pieces"]"`
	englishPrompt := `Please analyze the following listing data and provide the JSON array of suggested fields.

**Listing Data:**
* **Title:** %s
* **Description:** %s
* **Condition:** %s`
	japaneseInstruction := `あなたはフリマアプリ（C2Cマーケットプレイス）の出品支援および商品分類の専門家です。あなたの目的は、出品情報を分析し、検索性を高め信頼を得るために出品者が埋めるべき具体的な「追加情報の項目名」を提案することです。

**ガイドライン:**
1.  **商品を分析する:** 商品名（Title）と商品説明（Description）から、商品のカテゴリを特定してください。
2.  **重要な項目を優先する:** 一般的な内容の他にも、そのカテゴリにおいて重要な意味を持つ、客観的な項目（**例**：服なら「着丈」「生地」、車なら「走行距離」、パソコンなら「CPU 型番」「OSバージョン」、本なら「ISBN」「発行年」など）を優先してください。
3.  **ただし入力に含まれる項目は除外する:** 上記に該当する重要な項目であっても、タイトルや説明文に項目名または項目の内容が含まれている場合は、出力から除外してください。
4.  **制約事項:** 文字列の JSON 配列のみを返してください。提案する項目名は**日本語**で出力してください。
5.  **出力数:** 提案数は 最大4個までです。情報が少なすぎてカテゴリが特定できない場合は、空の配列 "[]" を返してください。

**出力フォーマット:**
有効なJSON形式のみを出力してください。
ブロックおもちゃの出品に対する回答例: "["総パーツ数", "ブランド", "不足パーツ"]"`
	japanesePrompt := `以下の出品データを分析し、提案する追加項目名のJSON配列を出力してください。

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
	respStr, err := s.vertexRepo.GenerateContent(ctx, "gemini-2.0-flash-lite", prompt, config)
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
