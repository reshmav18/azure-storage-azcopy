// Copyright © 2017 Microsoft <wastore@microsoft.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"time"

	chk "gopkg.in/check.v1"
)

type syncComparatorSuite struct{}

var _ = chk.Suite(&syncComparatorSuite{})

func (s *syncComparatorSuite) TestSyncSourceComparator(c *chk.C) {
	dummyCopyScheduler := dummyProcessor{}
	srcMD5 := []byte{'s'}
	destMD5 := []byte{'d'}

	// set up the indexer as well as the source comparator
	indexer := newObjectIndexer()
	sourceComparator := newSyncSourceComparator(indexer, dummyCopyScheduler.process, false)

	// create a sample destination object
	sampleDestinationObject := StoredObject{name: "test", relativePath: "/usr/test", lastModifiedTime: time.Now(), md5: destMD5}

	// test the comparator in case a given source object is not present at the destination
	// meaning no entry in the index, so the comparator should pass the given object to schedule a transfer
	compareErr := sourceComparator.processIfNecessary(StoredObject{name: "only_at_source", relativePath: "only_at_source", lastModifiedTime: time.Now(), md5: srcMD5})
	c.Assert(compareErr, chk.Equals, nil)

	// check the source object was indeed scheduled
	c.Assert(len(dummyCopyScheduler.record), chk.Equals, 1)
	c.Assert(dummyCopyScheduler.record[0].md5, chk.DeepEquals, srcMD5)

	// reset the processor so that it's empty
	dummyCopyScheduler = dummyProcessor{}

	// test the comparator in case a given source object is present at the destination
	// and it has a later modified time, so the comparator should pass the give object to schedule a transfer
	err := indexer.store(sampleDestinationObject)
	c.Assert(err, chk.IsNil)
	compareErr = sourceComparator.processIfNecessary(StoredObject{name: "test", relativePath: "/usr/test", lastModifiedTime: time.Now().Add(time.Hour), md5: srcMD5})
	c.Assert(compareErr, chk.Equals, nil)

	// check the source object was indeed scheduled
	c.Assert(len(dummyCopyScheduler.record), chk.Equals, 1)
	c.Assert(dummyCopyScheduler.record[0].md5, chk.DeepEquals, srcMD5)
	c.Assert(len(indexer.indexMap), chk.Equals, 0)

	// reset the processor so that it's empty
	dummyCopyScheduler = dummyProcessor{}

	// test the comparator in case a given source object is present at the destination
	// but is has an earlier modified time compared to the one at the destination
	// meaning that the source object is considered stale, so no transfer should be scheduled
	err = indexer.store(sampleDestinationObject)
	c.Assert(err, chk.IsNil)
	compareErr = sourceComparator.processIfNecessary(StoredObject{name: "test", relativePath: "/usr/test", lastModifiedTime: time.Now().Add(-time.Hour), md5: srcMD5})
	c.Assert(compareErr, chk.Equals, nil)

	// check no source object was scheduled
	c.Assert(len(dummyCopyScheduler.record), chk.Equals, 0)
	c.Assert(len(indexer.indexMap), chk.Equals, 0)
}

func (s *syncComparatorSuite) TestSyncSrcCompDisableComparator(c *chk.C) {
	dummyCopyScheduler := dummyProcessor{}
	srcMD5 := []byte{'s'}
	destMD5 := []byte{'d'}

	// set up the indexer as well as the source comparator
	indexer := newObjectIndexer()
	sourceComparator := newSyncSourceComparator(indexer, dummyCopyScheduler.process, true)

	// test the comparator in case a given source object is not present at the destination
	// meaning no entry in the index, so the comparator should pass the given object to schedule a transfer
	compareErr := sourceComparator.processIfNecessary(StoredObject{name: "only_at_source", relativePath: "only_at_source", lastModifiedTime: time.Now(), md5: srcMD5})
	c.Assert(compareErr, chk.Equals, nil)

	// check the source object was indeed scheduled
	c.Assert(len(dummyCopyScheduler.record), chk.Equals, 1)
	c.Assert(dummyCopyScheduler.record[0].md5, chk.DeepEquals, srcMD5)

	// reset the processor so that it's empty
	dummyCopyScheduler = dummyProcessor{}

	// create a sample source object
	currTime := time.Now()
	destinationStoredObjects := []StoredObject{
		// file whose last modified time is greater than that of source
		{name: "test1", relativePath: "/usr/test1", lastModifiedTime: currTime, md5: destMD5},
		// file whose last modified time is less than that of source
		{name: "test2", relativePath: "/usr/test2", lastModifiedTime: currTime, md5: destMD5},
	}

	sourceStoredObjects := []StoredObject{
		{name: "test1", relativePath: "/usr/test1", lastModifiedTime: currTime.Add(time.Hour), md5: srcMD5},
		{name: "test2", relativePath: "/usr/test2", lastModifiedTime: currTime.Add(-time.Hour), md5: srcMD5},
	}

	// test the comparator in case a given source object is present at the destination
	// but is has an earlier modified time compared to the one at the destination
	// meaning that the source object is considered stale, so no transfer should be scheduled
	for key, dstStoredObject := range destinationStoredObjects {
		err := indexer.store(dstStoredObject)
		c.Assert(err, chk.IsNil)
		compareErr = sourceComparator.processIfNecessary(sourceStoredObjects[key])
		c.Assert(compareErr, chk.Equals, nil)
		c.Assert(len(dummyCopyScheduler.record), chk.Equals, key+1)
		c.Assert(len(indexer.indexMap), chk.Equals, 0)
	}
}
