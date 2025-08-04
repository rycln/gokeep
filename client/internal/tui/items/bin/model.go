// Bin package handles binary file storage and display.
package bin

// BinFile represents binary file data for storage/transmission
// The 'bin' json tag indicates how data is serialized
type BinFile struct {
	Data []byte `json:"bin"`
}
