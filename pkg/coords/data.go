package coords

// TileData contains a raw tile data
type TileData []byte

// MetatileData contains a set of TileData
type MetatileData [MaxMetatileSize * MaxMetatileSize]TileData
