module almeng.com/ggc

require almeng.com/ggc/file v1.0.0

require almeng.com/ggc/general v1.0.0 // indirect

replace (
	almeng.com/ggc/file => ./file
	almeng.com/ggc/general => ./general
)

go 1.17
