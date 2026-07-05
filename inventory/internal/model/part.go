package model

import (
	"time"

	errs "github.com/1mpuser/inventory/internal/errors"
)

type PartFilter struct {
	UUIDs    []string
	PartType PartType
}

type Part struct {
	uuid          string
	name          string
	description   string
	price         int64 // в копейках
	partType      PartType
	stockQuantity int
	reserved      int
	properties    PartProperties
	createdAt     time.Time
}

func RestorePart(uuid, name, description string, partType PartType, price int64, stockQuantity, reserved int, properties PartProperties, createdAt time.Time) Part {
	return Part{
		uuid:          uuid,
		name:          name,
		description:   description,
		partType:      partType,
		price:         price,
		stockQuantity: stockQuantity,
		reserved:      reserved,
		properties:    properties,
		createdAt:     createdAt,
	}
}

func (p *Part) Reserve(quantity int) error {
	if p.stockQuantity-p.reserved < quantity {
		return errs.ErrOutOfStock
	}

	p.reserved += quantity
	return nil
}

func (p *Part) Release(quantity int) error {
	if p.reserved < quantity {
		return errs.ErrNothingToRelease
	}

	p.reserved -= quantity

	return nil
}

func (p *Part) UUID() string               { return p.uuid }
func (p *Part) Description() string        { return p.description }
func (p *Part) Name() string               { return p.name }
func (p *Part) PartType() PartType         { return p.partType }
func (p *Part) Price() int64               { return p.price }
func (p *Part) StockQuantity() int         { return p.stockQuantity }
func (p *Part) Reserved() int              { return p.reserved }
func (p *Part) Properties() PartProperties { return p.properties }
func (p *Part) CreatedAt() time.Time       { return p.createdAt }
