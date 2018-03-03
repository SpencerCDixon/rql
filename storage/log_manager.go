package storage

// LogManager is responsible for logging changes in our DBMS so they can be
// undone.  There is only ever one log file per DB.
type LogManager struct {
	// fm is used for creating iterators
	fm *FileManager
	// name of our log file
	filename string
	// page used to buffer log records
	page *Page
	// currentBlk that our log manager has loaded in memory
	currentBlk *Block
	// currentPos is a pointer to the next record
	currentPos int
	// Position of the last record in our log file
	LastRecordPos int
}

// NewLogManager creates a new log manager.  If a file does not already exist
// for this DB one will get created.
func NewLogManager(filename string, fm *FileManager) (*LogManager, error) {
	size, err := fm.Size(filename)
	if err != nil {
		return nil, err
	}

	lm := &LogManager{
		fm:         fm,
		filename:   filename,
		currentPos: 0,
		page:       NewPage(fm),
	}

	if size == 0 {
		if err := lm.appendNewBlock(); err != nil {
			return nil, err
		}
	} else {
		lm.currentBlk = NewBlock(filename, size-1)
		lm.page.Read(lm.currentBlk)
		lm.currentPos = lm.getLastRecordPosition() + IntSize
	}
	return lm, nil
}

// Flush takes our log record contents and writes them to disk.  Flush can be
// called for two reasons:

//   1. The page is full and needs to be written so more records can be recorded
//   2. Other parts of the system need the logs to be recorded before progressing
func (lm *LogManager) Flush() {
	lm.page.Write(lm.currentBlk)
}
func (lm *LogManager) FlushLSN(lsn int) {
	if lsn >= lm.currentLSN() {
		lm.Flush()
	}
}

// Append
func (lm *LogManager) Append(lrs []interface{}) int {
	recordSize := IntSize
	for _, lr := range lrs {
		recordSize += lm.size(lr)
	}

	// Not enough room, write to disk, and add room.
	if lm.currentPos+recordSize >= BlockSize {
		lm.Flush()
		lm.appendNewBlock()
	}

	// Add log record to buffer.
	for _, lr := range lrs {
		lm.appendValue(lr)
	}

	// Offset current values and return LSN.
	lm.finalizeRecord()

	return lm.currentLSN()
}

// Iterator returns a LogRecordIterator which can be cycled through.  Log
// records come from the front of the file to the end because that is the way
// the consumers will want to get them.  So, if the log file was:
//
// [1, 2]
// [3, 4]
//
// The first value would be [3, 4] and second would be [1, 2].  Any in memory
// records will first be flushed to disk before returning the iterator for
// accessing records.
func (lm *LogManager) Iterator() *RecordIterator {
	lm.Flush()
	iter := NewRecordIterator(lm.currentBlk, lm)
	return iter
}

func (lm *LogManager) appendNewBlock() error {
	lm.setLastRecordPosition(0)

	blk, err := lm.page.Append(lm.filename)
	if err != nil {
		return err
	}

	lm.currentBlk = blk
	lm.currentPos = IntSize

	return nil
}

func (lm *LogManager) appendValue(lr interface{}) {
	switch lr := lr.(type) {
	case int:
		lm.page.SetInt(lm.currentPos, lr)
	case string:
		lm.page.SetString(lm.currentPos, lr)
	default:
		// TODO:
	}
	lm.currentPos += lm.size(lr)
}

func (lm *LogManager) getLastRecordPosition() int {
	return lm.page.GetInt(lm.LastRecordPos)
}
func (lm *LogManager) setLastRecordPosition(pos int) error {
	return lm.page.SetInt(lm.LastRecordPos, pos)
}
func (lm *LogManager) finalizeRecord() error {
	lastPos := lm.getLastRecordPosition()
	err := lm.page.SetInt(lm.currentPos, lastPos)
	if err != nil {
		return err
	}
	err = lm.setLastRecordPosition(lm.currentPos)
	if err != nil {
		return err
	}

	lm.currentPos = lm.currentPos + IntSize
	return nil
}

// currentLSN returns the current Log Sequence Number.  Currently, the LSN is
// based on the block number but in the future this would be a good place to
// optimize.
func (lm *LogManager) currentLSN() int {
	return lm.currentBlk.BlockNum
}

// size returns how many bytes writing a log record value will be.  Currently,
// only int/strings are supported in our log files but we may add other types in
// the future.
// TODO: this should really just live in the file manager/one place.
// SizeForVal() or ByteSizeForValue() somewhere
func (lm *LogManager) size(obj interface{}) int {
	switch obj := obj.(type) {
	case int:
		return IntSize
	case string:
		return stringSize(obj) + IntSize
	default:
		return 0
	}
}

// Record/Helpers
type LogRecord struct {
	page *Page
	pos  int
}

func NewLogRecord(page *Page, pos int) *LogRecord {
	return &LogRecord{
		page: page,
		pos:  pos,
	}
}

func (lr *LogRecord) NextInt() int {
	nextInt := lr.page.GetInt(lr.pos)
	lr.pos = IntSize
	return nextInt
}

func (lr *LogRecord) NextString() string {
	nextStr := lr.page.GetString(lr.pos)
	lr.pos = stringSize(nextStr) + IntSize
	return nextStr
}

type RecordIterator struct {
	blk           *Block
	page          *Page
	lm            *LogManager
	currentRecord int
}

func NewRecordIterator(blk *Block, lm *LogManager) *RecordIterator {
	ri := &RecordIterator{
		blk:  blk,
		page: NewPage(lm.fm),
		lm:   lm,
	}

	ri.page.Read(blk)
	ri.currentRecord = ri.page.GetInt(lm.LastRecordPos)
	return ri
}

func (ri *RecordIterator) Next() bool {
	return ri.currentRecord > 0 || ri.blk.BlockNum > 0
}

func (ri *RecordIterator) Value() *LogRecord {
	// We got to the end of the current block since the linked chain goes from the
	// front to end of the file. Load in a new block and continue
	if ri.currentRecord == 0 {
		ri.moveToNextBlock()
	}
	ri.currentRecord = ri.page.GetInt(ri.currentRecord)

	lr := NewLogRecord(ri.page, ri.currentRecord+IntSize)
	return lr
}

func (ri *RecordIterator) moveToNextBlock() {
	blk := NewBlock(ri.blk.FileName, ri.blk.BlockNum-1)
	ri.page.Read(blk)
	ri.currentRecord = ri.page.GetInt(ri.lm.LastRecordPos)
}
