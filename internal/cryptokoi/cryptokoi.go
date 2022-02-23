package cryptokoi

import (
	"math/rand"
)

type CryptoKoi struct {
	// used to cache the attributes of the koi
	generatedAttributes KoiAttributes

	// the type of the koi
	// object which provides specific attributes per koi type.
	wrappedKoi  Koi
	randomizers struct {
		r1 *rand.Rand
		r2 *rand.Rand
		r3 *rand.Rand
	}
}

// only the first call to the method is valid.
func (c *CryptoKoi) generateAttributes() KoiAttributes {
	minBodyImages, maxBodyImages := c.wrappedKoi.amountBodyImages()
	minHeadImages, maxHeadImages := c.wrappedKoi.amountHeadImages()
	minFinImages, maxFinImages := c.wrappedKoi.amountFinImages()

	amountOfBodyImages := maxBodyImages
	if maxBodyImages != minBodyImages {
		// increment by 1 to include the max value into the possible values
		amountOfBodyImages = c.randomizers.r1.Intn(maxBodyImages+1-minBodyImages) + minBodyImages
	}
	amountOfFinImages := maxFinImages
	if maxFinImages != minFinImages {
		// increment by 1 to include the max value into the possible values
		amountOfFinImages = c.randomizers.r1.Intn(maxFinImages+1-minFinImages) + minFinImages
	}
	amountOfHeadImages := maxHeadImages
	if maxHeadImages != minHeadImages {
		// increment by 1 to include the max value into the possible values
		amountOfHeadImages = c.randomizers.r1.Intn(maxHeadImages+1-minHeadImages) + minHeadImages
	}
	attributes := KoiAttributes{
		BodyImages:   c.wrappedKoi.getBodyImages(amountOfBodyImages, c.randomizers.r2.Intn(255)),
		HeadImages:   c.wrappedKoi.getHeadImages(amountOfHeadImages, c.randomizers.r2.Intn(255)),
		FinImages:    c.wrappedKoi.getFinImages(amountOfFinImages, c.randomizers.r2.Intn(255)),
		BodyColor:    c.wrappedKoi.getBodyColor(c.randomizers.r3.Intn(255)),
		FinColor:     c.wrappedKoi.getFinBackgroundColor(c.randomizers.r3.Intn(255)),
		PrimaryColor: c.wrappedKoi.primaryColor(),
		KoiType:      c.wrappedKoi.getType(),
	}
	// reset the randomizers. If this method gets called again
	// it will throw an error. This is to ensure, that only valid koi images are generated.
	c.randomizers = struct {
		r1 *rand.Rand
		r2 *rand.Rand
		r3 *rand.Rand
	}{
		nil,
		nil,
		nil,
	}

	return attributes
}

func (c *CryptoKoi) GetAttributes() KoiAttributes {
	if c.randomizers.r1 != nil {
		c.generatedAttributes = c.generateAttributes()
	}
	return c.generatedAttributes
}

func (b *CryptoKoi) set(koi Koi, r1, r2, r3 *rand.Rand) {
	b.randomizers.r1 = r1
	b.randomizers.r2 = r2
	b.randomizers.r3 = r3
}
