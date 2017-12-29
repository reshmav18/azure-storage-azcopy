package main

import (
	"time"
	"encoding/json"
)

//These constant defines the various types of source and destination of the transfers


type LocationType uint16
const (
	fileSystemLocation LocationType = 1
	blockBlobLocation
	pageBlobLocation
	appendBlobLocation
	azureFileLocation
)

type CommandType int
const (
	REQUEST_COMMAND_AZCOPY_UPLOAD CommandType = 1 + iota
	REQUEST_COMMAND_AZCOPY_DOWNLOAD
)

const (
	TransferEntryChunkUpdateSuccess = "updated the chunk %d of transfer %d of Job %s"
)
const (
	DirectoryListingError = "not able to list contents of directory %s"
	FileCreationError = "Error %s Occured while creating the File for JobId %s \n"
	FileAlreadyExists = "file %s already exists"
	TransferIndexOutOfBoundError = "transfer %d of JobPart %s does not exists. Transfer Index exceeds number of transfer for this JobPart"
	MemoryMapFileUnmappedAlreadyError = "memory map file already unmapped. Map it again to use further"
	ChunkIndexOutOfBoundError = "given chunk %d of transfer %d of JobPart %s does not exists. Chunk Index exceeds number of chunks for transfer %d"
	MemoryMapFileInitializationError = "error memory mapping the file for jobId %s and partNo %d with err %s"
	InvalidFileName = "invalid file name for JobId %s and partNo %d"
	TransferTaskOffsetCalculationError = "calculated offset %d and actual offset %d of Job %s part %d and transfer entry %d does not match"
	InvalidJobId = "no active job for Job Id %s"
	InvalidPartNo = "no active part %s for Job Id %s"
	TransferStatusMarshallingError = "error marshalling the transfer status for Job Id %s and part No %s"
	InvalidHttpRequestBody = "the http Request Does not have a valid body definition"
	HttpRequestBodyReadError = "error reading the HTTP Request Body"
	HttpRequestUnmarshalError = "error UnMarshalling the HTTP Request Body"
	InvalidJobIDError = ""
	DateFormat = ""
	InvalidDirectoryError = ""
	InvalidArgumentsError = ""
	ContentEncodingLengthError = ""
	ContentTypeLengthError = ""
	MetaDataLengthError = ""
)

const (
	ChunkTransferStatusInactive = 0
	ChunkTransferStatusActive = 1
	ChunkTransferStatusProgress = 2
	ChunkTransferStatusComplete = 3
)

const (
	TransferStatusInactive = 0
	TransferStatusActive = 1
	TransferStatusProgress = 2
	TransferStatusComplete = 3
)

const (
	HighJobPriority = 3
	MediumJobPriority = 2
	LowJobPriority = 1
	DefaultJobPriority = HighJobPriority
)

const (
	BLOCK_LENGTH = 1000000 // in bytes
	MAX_NUMBER_ROUTINES = 10
	MAX_SIZE_CONTENT_TYPE = 256
	MAX_SIZE_CONTENT_ENCODING = 256
	MAX_SIZE_META_DATA = 1000
	OffSetToNumTransfer = 29
	OffsetToDataLength = 33
)

const (
	SuccessfulAZCopyRequest = "Succesfully Trigger the AZCopy request"
	UnsuccessfulAZCopyRequest = "Not able to trigger the AZCopy request"
	InvalidParametersInAZCopyRequest = "invalid parameters in az copy request"
	ACCOUNT_NAME = "azcopynextgendev1"
	ACCOUNT_KEY = "iVpnmZUb1aI9vMJ+yx+NOiBVJ63E25Ejs3ZA74Uqgml9bjrRQ4CEayQdZV7KqRnw0CrYQzc456+SwyI1KpyitA=="
	CONTAINER_NAME = "mycontainer"
	CRC64BitExample ="AFC0A0012976B444"
)

type task struct{
	Src string
	SecLastModifiedTime time.Time
	Dst string
}

type SrcDstPair struct {
	Src string
	Dst string
}

type transferStatus struct {
	Src string
	Dst string
	Status uint8
}

type transfersStatus struct {
	Status []transferStatus
}

type JobPart struct {
	Version uint32
	JobID string
	PartNo string
	SrcLocationType LocationType
	DstLocationType LocationType
	Tasks []task
}

type chunkInfo struct {
	BlockId [128 / 8]byte
	Status uint8
}

type transferEntry struct {
	Offset uint64
	SrcLength uint16
	DstLength uint16
	NumChunks uint16
	ModifiedTime uint32
	Status uint8
}

type JobPartPlan struct {
	Version uint32
	JobId [128 / 8] byte
	PartNo uint32
	Priority uint8
	TTLAfterCompletion uint32
	SrcLocationType LocationType
	DstLocationType LocationType
	NumTransfers uint32
}

type destinationBlobData struct {
	ContentTypeLength uint8
	ContentType [256]byte
	ContentEncodingLength uint8
	ContentEncoding [256]byte
	MetaDataLength uint16
	MetaData [1000]byte
	BlockSize uint16
}

type blobData struct {
	ContentType string
	ContentEncoding string
	MetaData string
}
type blockBlobData struct {
	BlobData blobData
	BlockSize uint16
}

type jobPartToBlockBlob struct {
	JobPart
	Data blockBlobData
}

type pageBlobData struct {
	BlobData blobData
}

type appendBlobData struct {
	blobData
	BlockSize int
}

type jobPartToPageBlob struct {
	JobPart
	Data pageBlobData
}
type jobPartToAppendBlob struct {
	JobPart
	Data appendBlobData
}

type jobPartToUnknown struct {
	JobPart
	Data json.RawMessage
}

type statusQuery struct {
	Guid string
	PartNo string
	TransferIndex uint32
	ChunkIndex uint16
}