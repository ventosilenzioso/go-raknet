package systems

import "log"

// VehicleSystem manages vehicle spawning and management
type VehicleSystem struct {
	vehicles map[uint16]*VehicleData
	nextID   uint16
}

// VehicleData represents vehicle information
type VehicleData struct {
	ID       uint16
	ModelID  int
	X, Y, Z  float32
	Rotation float32
	Color1   int
	Color2   int
	Owner    uint16
}

// NewVehicleSystem creates a new vehicle system
func NewVehicleSystem() *VehicleSystem {
	return &VehicleSystem{
		vehicles: make(map[uint16]*VehicleData),
		nextID:   1,
	}
}

// SpawnVehicle spawns a new vehicle
func (vs *VehicleSystem) SpawnVehicle(modelID int, x, y, z, rotation float32, color1, color2 int, owner uint16) uint16 {
	vehicleID := vs.nextID
	vs.nextID++
	
	vehicle := &VehicleData{
		ID:       vehicleID,
		ModelID:  modelID,
		X:        x,
		Y:        y,
		Z:        z,
		Rotation: rotation,
		Color1:   color1,
		Color2:   color2,
		Owner:    owner,
	}
	
	vs.vehicles[vehicleID] = vehicle
	
	log.Printf("🚗 Vehicle %d (model %d) spawned at %.2f, %.2f, %.2f", vehicleID, modelID, x, y, z)
	
	return vehicleID
}

// DestroyVehicle destroys a vehicle
func (vs *VehicleSystem) DestroyVehicle(vehicleID uint16) bool {
	if _, exists := vs.vehicles[vehicleID]; exists {
		delete(vs.vehicles, vehicleID)
		log.Printf("🚗 Vehicle %d destroyed", vehicleID)
		return true
	}
	return false
}

// GetVehicle returns vehicle data
func (vs *VehicleSystem) GetVehicle(vehicleID uint16) (*VehicleData, bool) {
	vehicle, exists := vs.vehicles[vehicleID]
	return vehicle, exists
}

// GetVehicleCount returns the number of spawned vehicles
func (vs *VehicleSystem) GetVehicleCount() int {
	return len(vs.vehicles)
}
