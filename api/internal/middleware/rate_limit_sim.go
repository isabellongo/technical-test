package middleware

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitSim simula respostas 429 (Too Many Requests) aleatórias.
// DECISÃO: 5% de chance por requisição no endpoint /people/v1/enrichments,
// para testar retry no n8n. Lógica documentada no README da API.
const rate429Chance = 0.05

var (
	simRand  = rand.New(rand.NewSource(time.Now().UnixNano()))
	simMutex sync.Mutex
)

// Simulate429 retorna 429 com ~5% de probabilidade, senão passa adiante.
func Simulate429() gin.HandlerFunc {
	return func(c *gin.Context) {
		simMutex.Lock()
		chance := simRand.Float64()
		simMutex.Unlock()

		if chance < rate429Chance {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate limit exceeded",
				"message": "simulated 429 for pipeline retry testing",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
