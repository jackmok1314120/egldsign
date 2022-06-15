package egld

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
)

type shardCoordinator struct {
	coordinator Coordinator
}
type multiShardCoordinator struct {
	maskHigh       uint32
	maskLow        uint32
	selfId         uint32
	numberOfShards uint32
}

var _ Coordinator = (*multiShardCoordinator)(nil)

// ComputeId calculates the shard for a given address container
func (msc *multiShardCoordinator) ComputeId(address []byte) uint32 {
	return msc.ComputeIdFromBytes(address)
}

// SystemAccountAddress is the hard-coded address in which we save global settings on all shards
var SystemAccountAddress = bytes.Repeat([]byte{255}, 32)

// NumInitCharactersForScAddress numbers of characters for smart contract address identifier
const NumInitCharactersForScAddress = 10

// VMTypeLen number of characters with VMType identifier in an address, these are the last 2 characters from the
// initial identifier
const VMTypeLen = 2

// ShardIdentiferLen number of characters for shard identifier in an address
const ShardIdentiferLen = 2

const metaChainShardIdentifier uint8 = 255
const numInitCharactersForOnMetachainSC = 15

const numInitCharactersForSystemAccountAddress = 30

// IsSystemAccountAddress returns true if given address is system account address
func IsSystemAccountAddress(address []byte) bool {
	if len(address) < numInitCharactersForSystemAccountAddress {
		return false
	}
	return bytes.Equal(address[:numInitCharactersForSystemAccountAddress], SystemAccountAddress[:numInitCharactersForSystemAccountAddress])
}

// IsSmartContractAddress verifies if a set address is of type smart contract
func IsSmartContractAddress(rcvAddress []byte) bool {
	if len(rcvAddress) <= NumInitCharactersForScAddress {
		return false
	}

	if IsEmptyAddress(rcvAddress) {
		return true
	}

	numOfZeros := NumInitCharactersForScAddress - VMTypeLen
	isSCAddress := bytes.Equal(rcvAddress[:numOfZeros], make([]byte, numOfZeros))
	return isSCAddress
}

// IsEmptyAddress returns whether an address is empty
func IsEmptyAddress(address []byte) bool {
	isEmptyAddress := bytes.Equal(address, make([]byte, len(address)))
	return isEmptyAddress
}

// IsMetachainIdentifier verifies if the identifier is of type metachain
func IsMetachainIdentifier(identifier []byte) bool {
	if len(identifier) == 0 {
		return false
	}

	for i := 0; i < len(identifier); i++ {
		if identifier[i] != metaChainShardIdentifier {
			return false
		}
	}

	return true
}

// IsSmartContractOnMetachain verifies if an address is smart contract on metachain
func IsSmartContractOnMetachain(identifier []byte, rcvAddress []byte) bool {
	if len(rcvAddress) <= NumInitCharactersForScAddress+numInitCharactersForOnMetachainSC {
		return false
	}

	if !IsMetachainIdentifier(identifier) {
		return false
	}

	if !IsSmartContractAddress(rcvAddress) {
		return false
	}

	leftSide := rcvAddress[NumInitCharactersForScAddress:(NumInitCharactersForScAddress + numInitCharactersForOnMetachainSC)]
	isOnMetaChainSCAddress := bytes.Equal(leftSide,
		make([]byte, numInitCharactersForOnMetachainSC))
	return isOnMetaChainSCAddress
}

// ComputeIdFromBytes calculates the shard for a given address
func (msc *multiShardCoordinator) ComputeIdFromBytes(address []byte) uint32 {

	var bytesNeed int
	if msc.numberOfShards <= 256 {
		bytesNeed = 1
	} else if msc.numberOfShards <= 65536 {
		bytesNeed = 2
	} else if msc.numberOfShards <= 16777216 {
		bytesNeed = 3
	} else {
		bytesNeed = 4
	}

	startingIndex := 0
	if len(address) > bytesNeed {
		startingIndex = len(address) - bytesNeed
	}

	buffNeeded := address[startingIndex:]
	if IsSmartContractOnMetachain(buffNeeded, address) {
		return MetachainShardId
	}

	addr := uint32(0)
	for i := 0; i < len(buffNeeded); i++ {
		addr = addr<<8 + uint32(buffNeeded[i])
	}

	shard := addr & msc.maskHigh
	if shard > msc.numberOfShards-1 {
		shard = addr & msc.maskLow
	}

	return shard
}

// NumberOfShards returns the number of shards
func (msc *multiShardCoordinator) NumberOfShards() uint32 {
	return msc.numberOfShards
}

// SelfId gets the shard id of the current node
func (msc *multiShardCoordinator) SelfId() uint32 {
	return msc.selfId
}

// SameShard returns weather two addresses belong to the same shard
func (msc *multiShardCoordinator) SameShard(firstAddress, secondAddress []byte) bool {
	if bytes.Equal(firstAddress, secondAddress) {
		return true
	}

	return msc.ComputeId(firstAddress) == msc.ComputeId(secondAddress)
}

// CommunicationIdentifier returns the identifier between current shard ID and destination shard ID
// identifier is generated such as the first shard from identifier is always smaller or equal than the last
func (msc *multiShardCoordinator) CommunicationIdentifier(destShardID uint32) string {
	return CommunicationIdentifierBetweenShards(msc.selfId, destShardID)
}

// ConvertBytes converts the input bytes in a readable string using multipliers (k, M, G)
func ConvertBytes(bytes uint64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(bytes)/1024.0)
	}
	if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(bytes)/1024.0/1024.0)
	}
	return fmt.Sprintf("%.2f GB", float64(bytes)/1024.0/1024.0/1024.0)
}

func plural(count int, singular string) (result string) {
	if count < 2 {
		result = strconv.Itoa(count) + " " + singular + " "
	} else {
		result = strconv.Itoa(count) + " " + singular + "s "
	}
	return
}

// SecondsToHourMinSec transform seconds input in a human friendly format
func SecondsToHourMinSec(input int) string {
	numSecondsInAMinute := 60
	numMinutesInAHour := 60
	numSecondsInAHour := numSecondsInAMinute * numMinutesInAHour
	result := ""

	hours := math.Floor(float64(input) / float64(numSecondsInAMinute) / float64(numMinutesInAHour))
	seconds := input % (numSecondsInAHour)
	minutes := math.Floor(float64(seconds) / float64(numSecondsInAMinute))
	seconds = input % numSecondsInAMinute

	if hours > 0 {
		result = plural(int(hours), "hour")
	}
	if minutes > 0 {
		result += plural(int(minutes), "minute")
	}
	if seconds > 0 {
		result += plural(seconds, "second")
	}

	return result
}

// GetShardIDString will return the string representation of the shard id
func GetShardIDString(shardID uint32) string {
	if shardID == math.MaxUint32 {
		return "metachain"
	}

	return fmt.Sprintf("%d", shardID)
}

// ConvertShardIDToUint32 converts shard id from string to uint32
func ConvertShardIDToUint32(shardIDStr string) (uint32, error) {
	if shardIDStr == "metachain" {
		return MetachainShardId, nil
	}

	shardID, err := strconv.ParseInt(shardIDStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint32(shardID), nil
}

// EpochStartIdentifier returns the string for the epoch start identifier
func EpochStartIdentifier(epoch uint32) string {
	return fmt.Sprintf("epochStartBlock_%d", epoch)
}

// AllShardId will be used to identify that a message is for all shards
const AllShardId = uint32(0xFFFFFFF0)

// CommunicationIdentifierBetweenShards is used to generate the identifier between shardID1 and shardID2
// identifier is generated such as the first shard from identifier is always smaller or equal than the last
func CommunicationIdentifierBetweenShards(shardId1 uint32, shardId2 uint32) string {
	if shardId1 == AllShardId || shardId2 == AllShardId {
		return ShardIdToString(AllShardId)
	}

	if shardId1 == shardId2 {
		return ShardIdToString(shardId1)
	}

	if shardId1 < shardId2 {
		return ShardIdToString(shardId1) + ShardIdToString(shardId2)
	}

	return ShardIdToString(shardId2) + ShardIdToString(shardId1)
}

// ShardIdToString returns the string according to the shard id
func ShardIdToString(shardId uint32) string {
	if shardId == MetachainShardId {
		return "_META"
	}
	if shardId == AllShardId {
		return "_ALL"
	}
	return fmt.Sprintf("_%d", shardId)
}

// IsInterfaceNil returns true if there is no value under the interface
func (msc *multiShardCoordinator) IsInterfaceNil() bool {
	return msc == nil
}

// ErrInvalidShardId signals that an invalid shard is was passed
var ErrInvalidShardId = errors.New("shard id must be smaller than the total number of shards")

// ErrInvalidNumberOfShards signals that an invalid number of shards was passed to the sharding registry
var ErrInvalidNumberOfShards = errors.New("the number of shards must be greater than zero")

// NewMultiShardCoordinator returns a new multiShardCoordinator and initializes the masks

// calculateMasks will create two numbers who's binary form is composed from as many
// ones needed to be taken into consideration for the shard assignment. The result
// of a bitwise AND operation of an address with this mask will result in the
// shard id where a transaction from that address will be dispatched
func (msc *multiShardCoordinator) calculateMasks() (uint32, uint32) {
	n := math.Ceil(math.Log2(float64(msc.numberOfShards)))
	return (1 << uint(n)) - 1, (1 << uint(n-1)) - 1
}

func NewMultiShardCoordinator(numberOfShards, selfId uint32) (*multiShardCoordinator, error) {
	if numberOfShards < 1 {
		return nil, ErrInvalidNumberOfShards
	}
	if selfId >= numberOfShards && selfId != MetachainShardId {
		return nil, ErrInvalidShardId
	}

	sr := &multiShardCoordinator{}
	sr.selfId = selfId
	sr.numberOfShards = numberOfShards
	sr.maskHigh, sr.maskLow = sr.calculateMasks()

	return sr, nil
}

// NewShardCoordinator returns a shard coordinator instance that is able to execute sharding-related operations
func NewShardCoordinator(numOfShardsWithoutMeta uint32, currentShard uint32) (*shardCoordinator, error) {
	coord, err := NewMultiShardCoordinator(numOfShardsWithoutMeta, currentShard)
	if err != nil {
		return nil, err
	}

	return &shardCoordinator{
		coordinator: coord,
	}, nil
}

// ComputeShardId computes the shard ID of a provided address
func (sc *shardCoordinator) ComputeShardId(address AddressHandler) (uint32, error) {
	if IfNil(address) {
		return 0, ErrNilAddress
	}
	if len(address.AddressBytes()) == 0 {
		return 0, ErrInvalidAddress
	}

	return sc.coordinator.ComputeId(address.AddressBytes()), nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (sc *shardCoordinator) IsInterfaceNil() bool {
	return sc == nil
}

//
//// NewAddressFromBech32String returns a new address from provided bech32 string
//func NewAddressFromBech32String(bech32 string) (*address, error) {
//	buff, err := AddressPublicKeyConverter.Decode(bech32)
//	if err != nil {
//		return nil, err
//	}
//
//	return &address{
//		bytes: buff,
//	}, err
//}
//
//// AddressAsBech32String returns the address as a bech32 string
//func (a *address) AddressAsBech32String() string {
//	return AddressPublicKeyConverter.Encode(a.bytes)
//}
