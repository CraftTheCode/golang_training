// Exercise 04: In-Memory Contact Manager
// Objective: Practice structs, pointers, slices, maps, validation, and error handling.
//
// Task:
// Build an in-memory Contact Book supporting Add, Find, Update, Delete, and List.
package main

import (
	"errors"
	"fmt"
	"strings"
)

type Contact struct {
	ID    int
	Name  string
	Email string
	Phone string
}

type ContactManager struct {
	contacts map[int]*Contact
	nextID   int
}

func NewContactManager() *ContactManager {
	return &ContactManager{
		contacts: make(map[int]*Contact),
		nextID:   1,
	}
}

func (cm *ContactManager) Add(name, email, phone string) (*Contact, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("contact name cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return nil, errors.New("invalid email address")
	}

	contact := &Contact{
		ID:    cm.nextID,
		Name:  name,
		Email: email,
		Phone: phone,
	}
	cm.contacts[contact.ID] = contact
	cm.nextID++
	return contact, nil
}

func (cm *ContactManager) FindByName(query string) []*Contact {
	query = strings.ToLower(query)
	var matches []*Contact
	for _, c := range cm.contacts {
		if strings.Contains(strings.ToLower(c.Name), query) {
			matches = append(matches, c)
		}
	}
	return matches
}

func (cm *ContactManager) Update(id int, newEmail, newPhone string) error {
	contact, exists := cm.contacts[id]
	if !exists {
		return fmt.Errorf("contact with ID %d not found", id)
	}
	if newEmail != "" {
		if !strings.Contains(newEmail, "@") {
			return errors.New("invalid email address")
		}
		contact.Email = newEmail
	}
	if newPhone != "" {
		contact.Phone = newPhone
	}
	return nil
}

func (cm *ContactManager) Delete(id int) error {
	if _, exists := cm.contacts[id]; !exists {
		return fmt.Errorf("contact with ID %d not found", id)
	}
	delete(cm.contacts, id)
	return nil
}

func (cm *ContactManager) ListAll() []*Contact {
	all := make([]*Contact, 0, len(cm.contacts))
	for _, c := range cm.contacts {
		all = append(all, c)
	}
	return all
}

func main() {
	cm := NewContactManager()

	fmt.Println("=== 1. Adding Contacts ===")
	c1, _ := cm.Add("Alice Smith", "alice@example.com", "+1-555-0101")
	c2, _ := cm.Add("Bob Jones", "bob@example.com", "+1-555-0102")
	fmt.Printf("Added: %+v\n", c1)
	fmt.Printf("Added: %+v\n", c2)

	// Validation check
	_, err := cm.Add("", "bad-email", "123")
	fmt.Println("Validation error handled:", err)

	fmt.Println("\n=== 2. Searching Contacts ===")
	results := cm.FindByName("alice")
	for _, r := range results {
		fmt.Printf("Found: %s (Email: %s)\n", r.Name, r.Email)
	}

	fmt.Println("\n=== 3. Updating and Deleting ===")
	cm.Update(c2.ID, "bob.jones@newcorp.com", "")
	fmt.Printf("Updated Bob: %+v\n", cm.contacts[c2.ID])

	cm.Delete(c1.ID)
	fmt.Printf("Total contacts after deleting Alice: %d\n", len(cm.ListAll()))
}
