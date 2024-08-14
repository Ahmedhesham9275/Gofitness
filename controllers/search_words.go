package controllers

import (
	"encoding/json"
	"myblog/config"
	"myblog/models"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func SearchWords(c *gin.Context) {
	var input struct {
		Words []string `json:"words" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var posts []models.Post
	if err := config.DB.Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	var wg sync.WaitGroup
	resultChannel := make(chan map[string]interface{}, len(input.Words))

	for _, word := range input.Words {
		wg.Add(1)
		go func(word string) {
			defer wg.Done()
			word = strings.ToLower(word)
			tf := 0
			df := 0

			// Calculate TF and DF
			for _, post := range posts {
				if strings.Contains(strings.ToLower(post.Content), word) {
					df++
					tf += strings.Count(strings.ToLower(post.Content), word)
				}
			}

			// Update statistics
			var stat models.SearchStatistic
			var err error

			// Lock mutex to handle concurrent access
			models.Mutex.Lock()
			defer models.Mutex.Unlock()

			if err = config.DB.Where("word = ?", word).First(&stat).Error; err != nil {
				if err.Error() == "record not found" {
					stat = models.SearchStatistic{Word: word}
				} else {
					resultChannel <- map[string]interface{}{
						"word":  word,
						"error": "Failed to fetch or create word statistics",
					}
					return
				}
			}

			stat.SearchCount++
			stat.LastTF = uint(tf)
			stat.LastDF = uint(df)

			// Update history
			var history []models.TFDFHistory
			if err := json.Unmarshal(stat.History, &history); err != nil {
				history = []models.TFDFHistory{}
			}
			history = append(history, models.TFDFHistory{
				TF:        tf,
				DF:        df,
				Timestamp: time.Now(),
			})
			historyJSON, err := json.Marshal(history)
			if err != nil {
				resultChannel <- map[string]interface{}{
					"word":  word,
					"error": "Failed to marshal history",
				}
				return
			}
			stat.History = historyJSON

			if err := config.DB.Save(&stat).Error; err != nil {
				resultChannel <- map[string]interface{}{
					"word":  word,
					"error": "Failed to save word statistics",
				}
				return
			}

			// Send the result back
			resultChannel <- map[string]interface{}{
				word: stat,
			}
		}(word)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	// Collect results
	results := []map[string]interface{}{}
	for result := range resultChannel {
		results = append(results, result)
	}

	c.JSON(http.StatusOK, results)
}
