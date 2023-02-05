package domain_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/tullo/invoice-mvp/domain"
)

func TestAddPosition(t *testing.T) {
	// Setup
	var i domain.Invoice

	// Run
	i.AddPosition(1, "Programming", 20, 60)
	i.AddPosition(1, "Programming", 12, 60)
	i.AddPosition(1, "Quality control", 3, 55)
	i.AddPosition(2, "Project management", 24, 50)
	i.AddPosition(2, "Quality control", 8, 55)

	// Asserts
	expected := domain.Position{Hours: 32, Price: 1920}
	// Project 1
	assert.Equal(t, expected, i.Positions[1]["Programming"])
	expected = domain.Position{Hours: 3, Price: 165}
	assert.Equal(t, expected, i.Positions[1]["Quality control"])
	expected = domain.Position{Hours: 24, Price: 1200}
	// Project 2
	assert.Equal(t, expected, i.Positions[2]["Project management"])
	expected = domain.Position{Hours: 8, Price: 440}
	assert.Equal(t, expected, i.Positions[2]["Quality control"])
}
